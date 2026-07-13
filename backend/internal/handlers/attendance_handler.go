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
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) = CURRENT_DATE AND check_out_time IS NOT NULL
	`, userID).Scan(&todayHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) >= DATE_TRUNC('week', CURRENT_DATE) AND check_out_time IS NOT NULL
	`, userID).Scan(&weekHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600), 0)
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
