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
		SELECT id, worksite_id, worksite_name_for_history, check_in_time, status
		FROM attendance
		WHERE user_id = $1 AND status = 'in_progress'
		ORDER BY check_in_time DESC
		LIMIT 1
	`, userID).Scan(&attendance.ID, &attendance.WorksiteID, &attendance.WorksiteNameForHistory, &checkInTime, &attendance.Status)

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

	// حفظ اسم نقطة العمل في سجل الحضور للرجوع إليه لاحقاً
	// هذا مهم في حال تم حذف نقطة العمل
	var worksiteNameForHistory string = worksite.Name

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
		INSERT INTO attendance (id, user_id, worksite_id, worksite_name_for_history, check_in_time, 
			check_in_lat, check_in_lng, check_in_distance_meters, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'in_progress')`,
		id, userID, worksite.ID, worksiteNameForHistory, now, lat, lng, distance,
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
		WorksiteNameForHistory: &worksiteNameForHistory,
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

	var worksiteID *string
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
	// إذا كانت نقطة العمل لا تزال موجودة، جلب بياناتها
	if worksiteID != nil {
		err = s.DB.QueryRow(`
			SELECT id, name, latitude, longitude, radius_meters
			FROM worksites WHERE id = $1
		`, *worksiteID).Scan(&worksite.ID, &worksite.Name,
			&worksite.Latitude, &worksite.Longitude, &worksite.RadiusMeters)
		if err != nil {
			log.Printf("❌ نقطة العمل غير موجودة: %v", err)
			return nil, 0, fmt.Errorf("نقطة العمل غير موجودة: %w", err)
		}
	} else {
		// إذا تم حذف نقطة العمل، لا يمكن التحقق من المسافة
		// سنسمح بالخروج بدون التحقق من الموقع
		log.Printf("⚠️ نقطة العمل محذوفة، السماح بالخروج بدون التحقق من المسافة")

		now := utils.NowInJerusalem()

		// حساب التقسيم عبر الأيام
		dayOneDate, dayTwoDate, dayOneHours, dayTwoHours := utils.SplitShiftAcrossDays(checkInTime, now)
		spansMultipleDays := dayTwoDate != nil

		// حساب ساعات العمل الليلية والنهارية
		dayNightPeriods := utils.CalculateDayNightHours(checkInTime, now)
		totalHours := now.Sub(checkInTime).Hours()
		isNightShift := utils.IsNightShift(dayNightPeriods.NightHours, totalHours)

		// تحديث سجل الحضور بدون التحقق من المسافة
		_, err = s.DB.Exec(`
			UPDATE attendance
			SET check_out_time = $1, status = 'completed',
			    spans_multiple_days = $2, day_one_date = $3, day_two_date = $4,
			    day_one_hours = $5, day_two_hours = $6,
			    night_hours = $7, day_hours = $8, is_night_shift = $9
			WHERE id = $10
		`, now, spansMultipleDays, dayOneDate, dayTwoDate,
			dayOneHours, dayTwoHours,
			dayNightPeriods.NightHours, dayNightPeriods.DayHours, isNightShift,
			attendanceID)

		if err != nil {
			log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
			return nil, 0, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
		}

		workedHours := totalHours
		log.Printf("✅ تم تسجيل إنهاء الدوام (نقطة عمل محذوفة): %.2f ساعة", workedHours)

		return nil, workedHours, nil
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

	// حساب التقسيم عبر الأيام
	dayOneDate, dayTwoDate, dayOneHours, dayTwoHours := utils.SplitShiftAcrossDays(checkInTime, now)
	spansMultipleDays := dayTwoDate != nil

	// حساب ساعات العمل الليلية والنهارية
	dayNightPeriods := utils.CalculateDayNightHours(checkInTime, now)
	totalHours := now.Sub(checkInTime).Hours()
	isNightShift := utils.IsNightShift(dayNightPeriods.NightHours, totalHours)

	// تحديث سجل الحضور مع الحقول الجديدة
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, check_out_lat = $2, check_out_lng = $3,
		    check_out_distance_meters = $4, status = 'completed',
		    spans_multiple_days = $5, day_one_date = $6, day_two_date = $7,
		    day_one_hours = $8, day_two_hours = $9,
		    night_hours = $10, day_hours = $11, is_night_shift = $12
		WHERE id = $13`,
		now, lat, lng, distance,
		spansMultipleDays, dayOneDate, dayTwoDate,
		dayOneHours, dayTwoHours,
		dayNightPeriods.NightHours, dayNightPeriods.DayHours, isNightShift,
		attendanceID,
	)
	if err != nil {
		log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
		return nil, 0, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
	}

	workedHours := totalHours
	log.Printf("✅ تم تسجيل إنهاء الدوام: %.2f ساعة (ليلي: %.2f، نهاري: %.2f)", workedHours, dayNightPeriods.NightHours, dayNightPeriods.DayHours)

	result := &GeofenceCheckResult{
		IsWithinRange:  true,
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}

	return result, workedHours, nil
}

// ForceCheckOut إنهاء دوام الموظف من قبل المدير (بدون التحقق من الموقع)
func (s *AttendanceService) ForceCheckOut(attendanceID string, adminID string) (float64, error) {
	log.Printf("📌 ForceCheckOut - الـ ID: %s, المدير: %s", attendanceID, adminID)

	var checkInTime time.Time
	var userID string

	// جلب معلومات الوردية
	err := s.DB.QueryRow(`
		SELECT user_id, check_in_time FROM attendance 
		WHERE id = $1 AND status = 'in_progress'
	`, attendanceID).Scan(&userID, &checkInTime)
	if err != nil {
		log.Printf("❌ لا يوجد وردية نشطة: %v", err)
		return 0, errors.New("لا يوجد وردية نشطة")
	}

	now := utils.NowInJerusalem()

	// حساب التقسيم عبر الأيام
	dayOneDate, dayTwoDate, dayOneHours, dayTwoHours := utils.SplitShiftAcrossDays(checkInTime, now)
	spansMultipleDays := dayTwoDate != nil

	// حساب ساعات العمل الليلية والنهارية
	dayNightPeriods := utils.CalculateDayNightHours(checkInTime, now)
	totalHours := now.Sub(checkInTime).Hours()
	isNightShift := utils.IsNightShift(dayNightPeriods.NightHours, totalHours)

	// محاولة التحديث مع check_out_notes أولاً
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, status = 'completed', check_out_notes = 'تم إنهاء الدوام من قبل المدير',
		    spans_multiple_days = $2, day_one_date = $3, day_two_date = $4,
		    day_one_hours = $5, day_two_hours = $6,
		    night_hours = $7, day_hours = $8, is_night_shift = $9
		WHERE id = $10
	`, now, spansMultipleDays, dayOneDate, dayTwoDate,
		dayOneHours, dayTwoHours,
		dayNightPeriods.NightHours, dayNightPeriods.DayHours, isNightShift,
		attendanceID)

	// إذا فشل بسبب عدم وجود عمود check_out_notes، حاول بدونه
	if err != nil {
		log.Printf("⚠️ فشل التحديث مع check_out_notes، محاولة بدونه: %v", err)
		_, err = s.DB.Exec(`
			UPDATE attendance
			SET check_out_time = $1, status = 'completed',
			    spans_multiple_days = $2, day_one_date = $3, day_two_date = $4,
			    day_one_hours = $5, day_two_hours = $6,
			    night_hours = $7, day_hours = $8, is_night_shift = $9
			WHERE id = $10
		`, now, spansMultipleDays, dayOneDate, dayTwoDate,
			dayOneHours, dayTwoHours,
			dayNightPeriods.NightHours, dayNightPeriods.DayHours, isNightShift,
			attendanceID)
		if err != nil {
			log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
			return 0, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
		}
	}

	workedHours := totalHours
	log.Printf("✅ تم إنهاء الدوام من قبل المدير: %.2f ساعة (ليلي: %.2f، نهاري: %.2f)", workedHours, dayNightPeriods.NightHours, dayNightPeriods.DayHours)

	return workedHours, nil
}
