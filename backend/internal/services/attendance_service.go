package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"worktrack/backend/internal/models"
	"worktrack/backend/pkg/utils"

	"github.com/google/uuid"
)

type AttendanceService struct {
	DB *sql.DB
}

func NewAttendanceService(db *sql.DB) *AttendanceService {
	return &AttendanceService{DB: db}
}

var ErrOutsideGeofence = errors.New("خارج نطاق الموقع المسموح")

func (s *AttendanceService) GetCurrentAttendance(userID string) (*models.Attendance, error) {
	var attendance models.Attendance
	var checkInTime time.Time

	err := s.DB.QueryRow(`
		SELECT id, worksite_id, check_in_time, status
		FROM attendance 
		WHERE user_id = $1 AND status = 'in_progress'
		ORDER BY check_in_time DESC
		LIMIT 1
	`, userID).Scan(&attendance.ID, &attendance.WorksiteID, &checkInTime, &attendance.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	attendance.UserID = userID
	attendance.CheckInTime = &checkInTime
	return &attendance, nil
}

func (s *AttendanceService) CheckIn(userID, worksiteID string, lat, lng float64) (*models.Attendance, *GeofenceCheckResult, error) {
	log.Printf("📌 CheckIn - المستخدم: %s, الموقع: %s", userID, worksiteID)
	log.Printf("📍 إحداثيات المستخدم: %f, %f", lat, lng)

	// التحقق من وجود وردية نشطة
	existing, _ := s.GetCurrentAttendance(userID)
	if existing != nil {
		return nil, nil, errors.New("يوجد وردية نشطة بالفعل")
	}

	if !utils.IsValidCoordinates(lat, lng) {
		return nil, nil, errors.New("إحداثيات غير صالحة")
	}

	// جلب نقطة العمل
	var worksite models.Worksite
	err := s.DB.QueryRow(`
		SELECT id, name, address, latitude, longitude, radius_meters, is_active
		FROM worksites WHERE id = $1 AND is_active = TRUE
	`, worksiteID).Scan(&worksite.ID, &worksite.Name, &worksite.Address,
		&worksite.Latitude, &worksite.Longitude, &worksite.RadiusMeters,
		&worksite.IsActive)

	if err != nil {
		log.Printf("❌ نقطة العمل غير موجودة: %v", err)
		return nil, nil, fmt.Errorf("نقطة العمل غير موجودة: %w", err)
	}

	log.Printf("📍 نقطة العمل: %s (%.6f, %.6f) - النطاق: %d متر",
		worksite.Name, worksite.Latitude, worksite.Longitude, worksite.RadiusMeters)

	// حساب المسافة بين المستخدم ونقطة العمل
	distance := utils.HaversineDistance(lat, lng, worksite.Latitude, worksite.Longitude)
	log.Printf("📏 المسافة المحسوبة: %.2f متر", distance)

	// التحقق من النطاق
	if distance > float64(worksite.RadiusMeters) {
		log.Printf("❌ المستخدم خارج النطاق: %.2f > %d", distance, worksite.RadiusMeters)
		result := &GeofenceCheckResult{
			IsWithinRange:  false,
			DistanceMeters: distance,
			AllowedRadius:  worksite.RadiusMeters,
		}
		return nil, result, ErrOutsideGeofence
	}

	log.Printf("✅ المستخدم داخل النطاق: %.2f <= %d", distance, worksite.RadiusMeters)

	// إنشاء سجل الحضور
	id := uuid.NewString()
	now := utils.NowInJerusalem()

	_, err = s.DB.Exec(`
		INSERT INTO attendance (id, user_id, worksite_id, check_in_time, 
			check_in_lat, check_in_lng, check_in_distance_meters, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'in_progress')`,
		id, userID, worksite.ID, now, lat, lng, distance,
	)
	if err != nil {
		log.Printf("❌ فشل حفظ بداية الدوام: %v", err)
		return nil, nil, fmt.Errorf("فشل حفظ بداية الدوام: %w", err)
	}

	result := &GeofenceCheckResult{
		IsWithinRange:  true,
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}

	attendance := &models.Attendance{
		ID:                    id,
		UserID:                userID,
		WorksiteID:            &worksite.ID,
		CheckInTime:           &now,
		CheckInLat:            &lat,
		CheckInLng:            &lng,
		CheckInDistanceMeters: &distance,
		Status:                "in_progress",
	}

	log.Printf("✅ تم تسجيل بدء الدوام بنجاح: %s", id)
	return attendance, result, nil
}

func (s *AttendanceService) CheckOut(userID, attendanceID string, lat, lng float64) (*GeofenceCheckResult, float64, error) {
	log.Printf("📌 CheckOut - المستخدم: %s", userID)
	log.Printf("📍 إحداثيات المستخدم: %f, %f", lat, lng)

	if !utils.IsValidCoordinates(lat, lng) {
		return nil, 0, errors.New("إحداثيات غير صالحة")
	}

	var worksiteID string
	var checkInTime time.Time

	err := s.DB.QueryRow(`
		SELECT worksite_id, check_in_time FROM attendance 
		WHERE id = $1 AND user_id = $2 AND status = 'in_progress'
	`, attendanceID, userID).Scan(&worksiteID, &checkInTime)
	if err != nil {
		log.Printf("❌ لا يوجد وردية نشطة: %v", err)
		return nil, 0, errors.New("لا يوجد وردية نشطة لهذا المستخدم")
	}

	var worksite models.Worksite
	err = s.DB.QueryRow(`
		SELECT id, name, latitude, longitude, radius_meters
		FROM worksites WHERE id = $1
	`, worksiteID).Scan(&worksite.ID, &worksite.Name,
		&worksite.Latitude, &worksite.Longitude, &worksite.RadiusMeters)
	if err != nil {
		log.Printf("❌ نقطة العمل غير موجودة: %v", err)
		return nil, 0, fmt.Errorf("نقطة العمل غير موجودة: %w", err)
	}

	// حساب المسافة عند الخروج
	distance := utils.HaversineDistance(lat, lng, worksite.Latitude, worksite.Longitude)
	log.Printf("📏 المسافة عند الخروج: %.2f متر", distance)

	if distance > float64(worksite.RadiusMeters) {
		log.Printf("❌ المستخدم خارج النطاق عند الخروج: %.2f > %d", distance, worksite.RadiusMeters)
		result := &GeofenceCheckResult{
			IsWithinRange:  false,
			DistanceMeters: distance,
			AllowedRadius:  worksite.RadiusMeters,
		}
		return result, 0, ErrOutsideGeofence
	}

	now := utils.NowInJerusalem()
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, check_out_lat = $2, check_out_lng = $3,
		    check_out_distance_meters = $4, status = 'completed'
		WHERE id = $5`,
		now, lat, lng, distance, attendanceID,
	)
	if err != nil {
		log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
		return nil, 0, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
	}

	workedHours := now.Sub(checkInTime).Hours()
	log.Printf("✅ تم تسجيل إنهاء الدوام: %.2f ساعة", workedHours)

	result := &GeofenceCheckResult{
		IsWithinRange:  true,
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}

	return result, workedHours, nil
}

// ForceCheckOut إنهاء دوام الموظف من قبل المدير (بدون التحقق من الموقع)
func (s *AttendanceService) ForceCheckOut(attendanceID string, adminID string) (map[string]interface{}, error) {
	log.Printf("📌 ForceCheckOut - الـ ID: %s, المدير: %s", attendanceID, adminID)

	var checkInTime time.Time
	var userID string
	var worksiteID string
	var worksiteName string
	var employeeName string

	// جلب معلومات الوردية
	err := s.DB.QueryRow(`
		SELECT a.user_id, a.check_in_time, a.worksite_id, u.full_name, w.name
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.id = $1 AND a.status = 'in_progress'
	`, attendanceID).Scan(&userID, &checkInTime, &worksiteID, &employeeName, &worksiteName)
	if err != nil {
		log.Printf("❌ لا يوجد وردية نشطة: %v", err)
		return nil, errors.New("لا يوجد وردية نشطة")
	}

	now := utils.NowInJerusalem()

	// محاولة التحديث مع check_out_notes أولاً
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, status = 'completed', check_out_notes = 'تم إنهاء الدوام من قبل المدير'
		WHERE id = $2
	`, now, attendanceID)

	// إذا فشل بسبب عدم وجود عمود check_out_notes، حاول بدونه
	if err != nil {
		log.Printf("⚠️ فشل التحديث مع check_out_notes، محاولة بدونه: %v", err)
		_, err = s.DB.Exec(`
			UPDATE attendance
			SET check_out_time = $1, status = 'completed'
			WHERE id = $2
		`, now, attendanceID)
		if err != nil {
			log.Printf("❌ فشل إنهاء الدوام: %v", err)
			return nil, fmt.Errorf("فشل إنهاء الدوام: %w", err)
		}
	}

	workedHours := now.Sub(checkInTime).Hours()
	log.Printf("✅ تم إنهاء الدوام: %.2f ساعة", workedHours)

	// إرجاع جميع المعلومات المطلوبة
	result := map[string]interface{}{
		"employee_name": employeeName,
		"worksite_name": worksiteName,
		"check_in_time": checkInTime,
		"check_out_time": now,
		"worked_hours": workedHours,
		"message": "تم إنهاء الدوام بنجاح",
	}

	return result, nil
}