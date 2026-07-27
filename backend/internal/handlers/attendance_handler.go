package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service             *services.AttendanceService
	NotificationService *services.NotificationService
}

func NewAttendanceHandler(s *services.AttendanceService, n *services.NotificationService) *AttendanceHandler {
	return &AttendanceHandler{Service: s, NotificationService: n}
}

type checkInRequest struct {
	WorksiteID string  `json:"worksite_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req checkInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات غير صحيحة. تأكد من إرسال worksite_id, latitude, longitude",
		})
		return
	}

	log.Printf("📝 طلب CheckIn: worksite_id=%s, lat=%f, lng=%f", req.WorksiteID, req.Latitude, req.Longitude)

	attendance, geofenceResult, err := h.Service.CheckIn(
		userID.(string),
		req.WorksiteID,
		req.Latitude,
		req.Longitude,
	)

	if errors.Is(err, services.ErrOutsideGeofence) {
		log.Printf("❌ المستخدم خارج النطاق: %.2f > %d",
			geofenceResult.DistanceMeters, geofenceResult.AllowedRadius)

		if h.NotificationService != nil {
			_ = h.NotificationService.Send(
				userID.(string),
				"🚨 خارج النطاق!",
				"أنت خارج نطاق موقع العمل المسموح. المسافة: "+formatDistance(geofenceResult.DistanceMeters),
			)
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "❌ أنت خارج نطاق موقع العمل",
			"geofence": geofenceResult,
			"details": gin.H{
				"distance":    geofenceResult.DistanceMeters,
				"max_allowed": geofenceResult.AllowedRadius,
				"is_inside":   false,
			},
		})
		return
	}

	if err != nil {
		log.Printf("❌ فشل CheckIn: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "✅ تم تسجيل بدء الدوام بنجاح",
		"attendance": attendance,
		"geofence":   geofenceResult,
	})
}

type checkOutRequest struct {
	AttendanceID string  `json:"attendance_id" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req checkOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	// =============================================
	// 🔒 التحقق من وجود إحداثيات الموقع
	// =============================================
	if req.Latitude == 0 && req.Longitude == 0 {
		log.Printf("❌ الموقع غير محدد (0,0)")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "📍 يرجى تحديد موقعك أولاً قبل إنهاء الدوام",
			"code":  "location_required",
		})
		return
	}

	log.Printf("📝 طلب CheckOut: attendance_id=%s, lat=%f, lng=%f", req.AttendanceID, req.Latitude, req.Longitude)

	geofenceResult, workedHours, err := h.Service.CheckOut(
		userID.(string),
		req.AttendanceID,
		req.Latitude,
		req.Longitude,
	)

	if errors.Is(err, services.ErrOutsideGeofence) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":    "❌ أنت خارج نطاق موقع العمل عند الخروج",
			"geofence": geofenceResult,
		})
		return
	}

	if err != nil {
		log.Printf("❌ فشل CheckOut: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "✅ تم تسجيل إنهاء الدوام بنجاح",
		"geofence":     geofenceResult,
		"worked_hours": workedHours,
	})
}

func (h *AttendanceHandler) GetCurrentAttendance(c *gin.Context) {
	userID, _ := c.Get("user_id")

	attendance, err := h.Service.GetCurrentAttendance(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب حالة الحضور"})
		return
	}

	if attendance == nil {
		c.JSON(http.StatusOK, gin.H{
			"has_active": false,
			"message":    "لا يوجد سجل حضور نشط",
		})
		return
	}

	var elapsedSeconds int64
	if attendance.CheckInTime != nil {
		elapsedSeconds = int64(time.Since(*attendance.CheckInTime).Seconds())
	}

	c.JSON(http.StatusOK, gin.H{
		"has_active":      true,
		"attendance_id":   attendance.ID,
		"worksite_id":     attendance.WorksiteID,
		"check_in_time":   attendance.CheckInTime,
		"elapsed_seconds": elapsedSeconds,
		"status":          attendance.Status,
	})
}

func (h *AttendanceHandler) GetAttendanceSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var todayHours, weekHours, monthHours float64

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(COALESCE(night_hours, 0) + COALESCE(day_hours, 0)), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) = CURRENT_DATE AND check_out_time IS NOT NULL
	`, userID).Scan(&todayHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(COALESCE(night_hours, 0) + COALESCE(day_hours, 0)), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) >= DATE_TRUNC('week', CURRENT_DATE) AND check_out_time IS NOT NULL
	`, userID).Scan(&weekHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(COALESCE(night_hours, 0) + COALESCE(day_hours, 0)), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) >= DATE_TRUNC('month', CURRENT_DATE) AND check_out_time IS NOT NULL
	`, userID).Scan(&monthHours)

	c.JSON(http.StatusOK, gin.H{
		"today_hours": todayHours,
		"week_hours":  weekHours,
		"month_hours": monthHours,
	})
}

func formatDistance(d float64) string {
	if d < 1000 {
		return fmt.Sprintf("%.0f متر", d)
	}
	return fmt.Sprintf("%.2f كيلومتر", d/1000)
}

func (h *AttendanceHandler) GetAllAttendanceSummary(c *gin.Context) {
	rows, err := h.Service.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.email,
			COALESCE(SUM(EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600), 0) as total_hours,
			COUNT(CASE WHEN a.check_out_time IS NOT NULL THEN 1 END) as completed_shifts,
			COUNT(CASE WHEN a.check_out_time IS NULL THEN 1 END) as active_shifts
		FROM users u
		LEFT JOIN attendance a ON a.user_id = u.id
		WHERE u.role = 'employee'
		AND (a.check_in_time IS NULL OR DATE(a.check_in_time) >= DATE_TRUNC('week', CURRENT_DATE))
		GROUP BY u.id, u.full_name, u.email
		ORDER BY total_hours DESC
	`)

	if err != nil {
		log.Printf("❌ فشل جلب ملخص ساعات العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب ملخص ساعات العمل"})
		return
	}
	defer rows.Close()

	var summary []gin.H
	for rows.Next() {
		var id, fullName, email string
		var totalHours float64
		var completedShifts, activeShifts int

		if err := rows.Scan(&id, &fullName, &email, &totalHours, &completedShifts, &activeShifts); err == nil {
			summary = append(summary, gin.H{
				"id":               id,
				"full_name":        fullName,
				"email":            email,
				"total_hours":      totalHours,
				"completed_shifts": completedShifts,
				"active_shifts":    activeShifts,
			})
		}
	}

	c.JSON(http.StatusOK, summary)
}

// GetEmployeeAttendanceHistory جلب سجل الحضور لموظف معين
func (h *AttendanceHandler) GetEmployeeAttendanceHistory(c *gin.Context) {
	userID := c.Param("id")
	
	// الحصول على معاملات التصفية
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")
	
	// بناء الاستعلام حسب الفلتر
	query := `
		SELECT
			a.id,
			a.worksite_id,
			COALESCE(w.name, a.worksite_name_for_history) as worksite_name,
			a.check_in_time,
			a.check_in_lat,
			a.check_in_lng,
			a.check_in_distance_meters,
			a.check_out_time,
			a.check_out_lat,
			a.check_out_lng,
			a.check_out_distance_meters,
			a.status,
			a.created_at,
			a.spans_multiple_days,
			a.day_one_date,
			a.day_two_date,
			a.day_one_hours,
			a.day_two_hours,
			a.night_hours,
			a.day_hours,
			a.is_night_shift
		FROM attendance a
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1
	
	if year != "" && month != "" {
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM check_in_time) = $%d", argCount)
		args = append(args, year)
		
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM check_in_time) = $%d", argCount)
		args = append(args, month)
	}
	
	query += " ORDER BY check_in_time DESC"
	
	rows, err := h.Service.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}
	defer rows.Close()
	
	var history []gin.H
	for rows.Next() {
		var id, worksiteID, worksiteName string
		var checkInTime, checkOutTime, createdAt time.Time
		var checkInLat, checkInLng, checkInDistance, checkOutLat, checkOutLng, checkOutDistance *float64
		var status string
		var spansMultipleDays bool
		var dayOneDate, dayTwoDate *time.Time
		var dayOneHours, dayTwoHours, nightHours, dayHours *float64
		var isNightShift *bool

		err := rows.Scan(
			&id, &worksiteID, &worksiteName, &checkInTime, &checkInLat, &checkInLng, &checkInDistance,
			&checkOutTime, &checkOutLat, &checkOutLng, &checkOutDistance, &status, &createdAt,
			&spansMultipleDays, &dayOneDate, &dayTwoDate, &dayOneHours, &dayTwoHours,
			&nightHours, &dayHours, &isNightShift,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		// حساب worked_hours كمجموع الساعات الليلية والنهارية لضمان الاتساق
		var workedHours *float64
		if nightHours != nil && dayHours != nil {
			total := *nightHours + *dayHours
			workedHours = &total
		} else if checkOutTime.After(checkInTime) {
			hours := checkOutTime.Sub(checkInTime).Hours()
			workedHours = &hours
		}

		record := gin.H{
			"id":                        id,
			"worksite_id":               worksiteID,
			"worksite_name":             worksiteName,
			"check_in_time":             checkInTime,
			"check_in_lat":              checkInLat,
			"check_in_lng":              checkInLng,
			"check_in_distance_meters":  checkInDistance,
			"check_out_time":            checkOutTime,
			"check_out_lat":             checkOutLat,
			"check_out_lng":             checkOutLng,
			"check_out_distance_meters": checkOutDistance,
			"status":                    status,
			"worked_hours":              workedHours,
			"created_at":                createdAt,
			"spans_multiple_days":       spansMultipleDays,
			"day_one_date":              dayOneDate,
			"day_two_date":              dayTwoDate,
			"day_one_hours":             dayOneHours,
			"day_two_hours":             dayTwoHours,
			"night_hours":               nightHours,
			"day_hours":                 dayHours,
			"is_night_shift":            isNightShift,
		}

		history = append(history, record)
	}
	
	c.JSON(http.StatusOK, history)
}

// GetEmployeeMonthlySummary جلب ملخص شهري لموظف معين
func (h *AttendanceHandler) GetEmployeeMonthlySummary(c *gin.Context) {
	userID := c.Param("id")
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")
	
	if year == "" || month == "" {
		// استخدام الشهر الحالي إذا لم يتم تحديده
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%d", int(now.Month()))
	}
	
	// جلب إجمالي الساعات للشهر
	var totalHours float64
	err := h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(COALESCE(night_hours, 0) + COALESCE(day_hours, 0)), 0)
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&totalHours)
	
	if err != nil {
		log.Printf("❌ فشل جلب الملخص الشهري: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملخص الشهري"})
		return
	}
	
	// جلب عدد أيام العمل
	var workDays int
	err = h.Service.DB.QueryRow(`
		SELECT COUNT(DISTINCT DATE(check_in_time))
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&workDays)
	
	if err != nil {
		log.Printf("❌ فشل جلب عدد أيام العمل: %v", err)
		workDays = 0
	}
	
	// جلب معلومات الموظف
	var fullName, email string
	err = h.Service.DB.QueryRow(`
		SELECT full_name, email FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email)
	
	if err != nil {
		log.Printf("❌ فشل جلب معلومات الموظف: %v", err)
		fullName = "غير معروف"
		email = ""
	}
	
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"id":         userID,
			"full_name":  fullName,
			"email":      email,
		},
		"period": gin.H{
			"year":  year,
			"month": month,
		},
		"summary": gin.H{
			"total_hours": totalHours,
			"work_days":   workDays,
		},
	})
}

// GetMyAttendanceHistory جلب سجل الحضور للمستخدم الحالي
func (h *AttendanceHandler) GetMyAttendanceHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// الحصول على معاملات التصفية
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")

	// بناء الاستعلام حسب الفلتر
	query := `
		SELECT
			a.id,
			a.worksite_id,
			COALESCE(w.name, a.worksite_name_for_history) as worksite_name,
			a.check_in_time,
			a.check_in_lat,
			a.check_in_lng,
			a.check_in_distance_meters,
			a.check_out_time,
			a.check_out_lat,
			a.check_out_lng,
			a.check_out_distance_meters,
			a.status,
			a.created_at,
			a.spans_multiple_days,
			a.day_one_date,
			a.day_two_date,
			a.day_one_hours,
			a.day_two_hours,
			a.night_hours,
			a.day_hours,
			a.is_night_shift
		FROM attendance a
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1
	`
	args := []interface{}{userID.(string)}
	argCount := 1

	// إضافة الفلتر فقط إذا تم تحديد السنة والشهر معاً
	if year != "" && month != "" {
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM check_in_time) = $%d", argCount)
		args = append(args, year)

		argCount++
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM check_in_time) = $%d", argCount)
		args = append(args, month)
	}

	query += " ORDER BY check_in_time DESC"

	log.Printf("📝 جلب سجل الحضور للمستخدم: user_id=%s, year=%s, month=%s", userID.(string), year, month)

	rows, err := h.Service.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}
	defer rows.Close()

	var history []gin.H
	for rows.Next() {
		var id, worksiteID, worksiteName string
		var checkInTime, checkOutTime, createdAt time.Time
		var checkInLat, checkInLng, checkInDistance, checkOutLat, checkOutLng, checkOutDistance *float64
		var status string
		var spansMultipleDays bool
		var dayOneDate, dayTwoDate *time.Time
		var dayOneHours, dayTwoHours, nightHours, dayHours *float64
		var isNightShift *bool

		err := rows.Scan(
			&id, &worksiteID, &worksiteName, &checkInTime, &checkInLat, &checkInLng, &checkInDistance,
			&checkOutTime, &checkOutLat, &checkOutLng, &checkOutDistance, &status, &createdAt,
			&spansMultipleDays, &dayOneDate, &dayTwoDate, &dayOneHours, &dayTwoHours,
			&nightHours, &dayHours, &isNightShift,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		// حساب worked_hours كمجموع الساعات الليلية والنهارية لضمان الاتساق
		var workedHours *float64
		if nightHours != nil && dayHours != nil {
			total := *nightHours + *dayHours
			workedHours = &total
		} else if checkOutTime.After(checkInTime) {
			hours := checkOutTime.Sub(checkInTime).Hours()
			workedHours = &hours
		}

		record := gin.H{
			"id":                        id,
			"worksite_id":               worksiteID,
			"worksite_name":             worksiteName,
			"check_in_time":             checkInTime,
			"check_in_lat":              checkInLat,
			"check_in_lng":              checkInLng,
			"check_in_distance_meters":  checkInDistance,
			"check_out_time":            checkOutTime,
			"check_out_lat":             checkOutLat,
			"check_out_lng":             checkOutLng,
			"check_out_distance_meters": checkOutDistance,
			"status":                    status,
			"worked_hours":              workedHours,
			"created_at":                createdAt,
			"spans_multiple_days":       spansMultipleDays,
			"day_one_date":              dayOneDate,
			"day_two_date":              dayTwoDate,
			"day_one_hours":             dayOneHours,
			"day_two_hours":             dayTwoHours,
			"night_hours":               nightHours,
			"day_hours":                 dayHours,
			"is_night_shift":            isNightShift,
		}

		history = append(history, record)
	}

	log.Printf("✅ تم جلب %d سجل حضور", len(history))
	c.JSON(http.StatusOK, history)
}

// GetMyMonthlySummary جلب الملخص الشهري للمستخدم الحالي
func (h *AttendanceHandler) GetMyMonthlySummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")
	
	if year == "" || month == "" {
		// استخدام الشهر الحالي إذا لم يتم تحديده
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%d", int(now.Month()))
	}
	
	// جلب إجمالي الساعات للشهر
	var totalHours float64
	err := h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(COALESCE(night_hours, 0) + COALESCE(day_hours, 0)), 0)
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&totalHours)
	
	if err != nil {
		log.Printf("❌ فشل جلب الملخص الشهري: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملخص الشهري"})
		return
	}
	
	// جلب عدد أيام العمل
	var workDays int
	err = h.Service.DB.QueryRow(`
		SELECT COUNT(DISTINCT DATE(check_in_time))
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&workDays)
	
	if err != nil {
		log.Printf("❌ فشل جلب عدد أيام العمل: %v", err)
		workDays = 0
	}
	
	// جلب معلومات المستخدم
	var fullName, email string
	err = h.Service.DB.QueryRow(`
		SELECT full_name, email FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email)
	
	if err != nil {
		log.Printf("❌ فشل جلب معلومات المستخدم: %v", err)
		fullName = "غير معروف"
		email = ""
	}
	
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"id":         userID,
			"full_name":  fullName,
			"email":      email,
		},
		"period": gin.H{
			"year":  year,
			"month": month,
		},
		"summary": gin.H{
			"total_hours": totalHours,
			"work_days":   workDays,
		},
	})
}

// CleanupOldRecords حذف السجلات القديمة (أكثر من 3 أشهر)
func (h *AttendanceHandler) CleanupOldRecords(c *gin.Context) {
	// حذف السجلات الأقدم من 3 أشهر
	result, err := h.Service.DB.Exec(`
		DELETE FROM attendance 
		WHERE check_in_time < NOW() - INTERVAL '3 months'
		AND status = 'completed'
	`)
	if err != nil {
		log.Printf("❌ فشل حذف السجلات القديمة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف السجلات القديمة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ تم حذف %d سجل قديم", rowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("تم حذف %d سجل قديم بنجاح", rowsAffected),
		"deleted_count": rowsAffected,
	})
}

// ForceCheckOut إنهاء دوام الموظف من قبل المدير
func (h *AttendanceHandler) ForceCheckOut(c *gin.Context) {
	// التحقق من أن المستخدم مدير
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	adminIDStr := adminID.(string)

	// التحقق من حالة الاشتراك (مع التسامح مع عدم وجود الحقول)
	var subscriptionStatus string
	var subscriptionExpiresAt time.Time
	err := h.Service.DB.QueryRow(`
		SELECT subscription_status, subscription_expires_at
		FROM users
		WHERE id = $1
	`, adminIDStr).Scan(&subscriptionStatus, &subscriptionExpiresAt)

	// إذا فشل جلب حالة الاشتراك (مثلاً الحقول غير موجودة)، نتجاوز التحقق
	if err != nil {
		log.Printf("⚠️ فشل جلب حالة الاشتراك، تجاوز التحقق: %v", err)
		// نستمر بدون التحقق من الاشتراك
	} else {
		// التحقق من أن الاشتراك نشط
		if subscriptionStatus == "canceled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "اشتراكك ملغي، الرجاء التواصل مع الدعم"})
			return
		}

		if subscriptionStatus == "expired" || (!subscriptionExpiresAt.IsZero() && time.Now().After(subscriptionExpiresAt)) {
			// تحديث الحالة إلى منتهي إذا لزم الأمر
			if subscriptionStatus != "expired" {
				_, _ = h.Service.DB.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, adminIDStr)
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "اشتراكك منتهي، الرجاء تجديده للتمكن من استخدام هذه الميزة"})
			return
		}
	}

	var req struct {
		AttendanceID string `json:"attendance_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 طلب ForceCheckOut: attendance_id=%s, admin_id=%s", req.AttendanceID, adminIDStr)

	workedHours, err := h.Service.ForceCheckOut(req.AttendanceID, adminIDStr)
	if err != nil {
		log.Printf("❌ فشل ForceCheckOut: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "✅ تم إنهاء الدوام بنجاح",
		"worked_hours": workedHours,
	})
}
