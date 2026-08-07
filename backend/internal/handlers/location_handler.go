package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationHandler struct {
	DB *sql.DB
	WS *WSHandler
}

func NewLocationHandler(db *sql.DB, wsHandler *WSHandler) *LocationHandler {
	return &LocationHandler{DB: db, WS: wsHandler}
}

// UpdateLocation - تحديث موقع المستخدم
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Accuracy  float64 `json:"accuracy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	_, err := h.DB.Exec(`
		INSERT INTO location_tracking (id, user_id, latitude, longitude, accuracy, recorded_at)
		VALUES ($1, $2, $3, $4, $5, now())
	`, uuid.NewString(), userID, req.Latitude, req.Longitude, req.Accuracy)
	if err != nil {
		log.Printf("❌ فشل حفظ الموقع: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// التحقق من الخروج عن النطاق
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("⚠️ Recovered from panic in geofence check: %v", r)
			}
		}()
		h.checkGeofenceViolation(userID.(string), req.Latitude, req.Longitude)
	}()

	// إرسال تحديث فوري عبر WebSocket
	if h.WS != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("⚠️ Recovered from panic in WebSocket broadcast: %v", r)
				}
			}()

			// جلب معلومات الموظف
			var fullName, email string
			err := h.DB.QueryRow(`
				SELECT full_name, email FROM users WHERE id = $1
			`, userID).Scan(&fullName, &email)
			if err != nil {
				log.Printf("⚠️ Failed to fetch user info for WebSocket broadcast: %v", err)
				return
			}

			// إرسال التحديث
			h.WS.BroadcastLocationUpdate(map[string]interface{}{
				"user_id":    userID,
				"full_name":  fullName,
				"email":      email,
				"latitude":   req.Latitude,
				"longitude":  req.Longitude,
				"accuracy":   req.Accuracy,
				"immediate":  true, // تحديث فوري من المستخدم
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم تحديث الموقع"})
}

// checkGeofenceViolation - التحقق من خروج الموظف عن النطاق
func (h *LocationHandler) checkGeofenceViolation(userID string, lat, lng float64) {
	var attendanceID string
	var worksiteID *string
	var worksiteLat, worksiteLng float64
	var radiusMeters int

	err := h.DB.QueryRow(`
		SELECT a.id, a.worksite_id, w.latitude, w.longitude, w.radius_meters
		FROM attendance a
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1 AND a.status = 'in_progress'
	`, userID).Scan(&attendanceID, &worksiteID, &worksiteLat, &worksiteLng, &radiusMeters)

	if err != nil {
		return
	}

	// إذا كانت نقطة العمل محذوفة، لا يمكن التحقق من المسافة
	if worksiteID == nil {
		return
	}

	distance := utils.HaversineDistance(lat, lng, worksiteLat, worksiteLng)
	if distance > float64(radiusMeters) {
		var count int
		_ = h.DB.QueryRow(`
			SELECT COUNT(*) FROM notifications 
			WHERE user_id = $1 AND type = 'geofence_alert' 
			AND created_at > now() - interval '5 minutes'
		`, userID).Scan(&count)

		if count == 0 {
			_, _ = h.DB.Exec(`
				INSERT INTO notifications (id, user_id, title, body, type, is_read, created_at)
				VALUES ($1, $2, $3, $4, 'geofence_alert', FALSE, now())
			`, uuid.NewString(), userID,
				"🚨 خروج عن النطاق!",
				fmt.Sprintf("خرج الموظف عن النطاق المسموح به. المسافة: %.0f متر", distance))
		}
	}
}

// GetActiveEmployees - جلب الموظفين النشطين مع مواقعهم
func (h *LocationHandler) GetActiveEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT DISTINCT ON (u.id)
			u.id,
			u.full_name,
			u.email,
			u.phone,
			lt.latitude,
			lt.longitude,
			lt.recorded_at,
			a.id as attendance_id,
			a.check_in_time,
			w.id as worksite_id,
			COALESCE(w.name, a.worksite_name_for_history) as worksite_name,
			w.latitude as worksite_lat,
			w.longitude as worksite_lng,
			w.radius_meters,
			EXTRACT(EPOCH FROM (now() - a.check_in_time)) / 3600 as hours_worked
		FROM users u
		JOIN attendance a ON a.user_id = u.id AND a.status = 'in_progress'
		JOIN location_tracking lt ON lt.user_id = u.id
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee' AND u.is_active = TRUE
		ORDER BY u.id, lt.recorded_at DESC
	`)

	if err != nil {
		log.Printf("❌ فشل جلب الموظفين النشطين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب البيانات"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, attendanceID, worksiteID, worksiteName string
		var lat, lng, worksiteLat, worksiteLng float64
		var radiusMeters int
		var recordedAt, checkInTime time.Time
		var hoursWorked float64

		if err := rows.Scan(&id, &fullName, &email, &phone, &lat, &lng, &recordedAt,
			&attendanceID, &checkInTime, &worksiteID, &worksiteName, 
			&worksiteLat, &worksiteLng, &radiusMeters, &hoursWorked); err != nil {
			continue
		}

		distance := utils.HaversineDistance(lat, lng, worksiteLat, worksiteLng)
		isInside := distance <= float64(radiusMeters)

		// تحديد حالة الموظف
		status := "inside"
		statusText := "✅ داخل النطاق"
		if !isInside {
			status = "outside"
			statusText = "❌ خارج النطاق"
		}

		// حساب وقت الانقطاع (آخر تحديث للموقع)
		timeSinceLastUpdate := time.Since(recordedAt)
		isActive := timeSinceLastUpdate < 5*time.Minute

		employees = append(employees, gin.H{
			"id":              id,
			"full_name":       fullName,
			"email":           email,
			"phone":           phone,
			"latitude":        lat,
			"longitude":       lng,
			"last_update":     recordedAt,
			"is_active":       isActive,
			"status":          status,
			"status_text":     statusText,
			"attendance_id":   attendanceID,
			"check_in_time":   checkInTime,
			"hours_worked":    hoursWorked,
			"worksite": gin.H{
				"id":        worksiteID,
				"name":      worksiteName,
				"latitude":  worksiteLat,
				"longitude": worksiteLng,
				"radius":    radiusMeters,
				"is_inside": isInside,
				"distance":  distance,
			},
		})
	}

	c.JSON(http.StatusOK, employees)
}

// GetEmployeeTrack - جلب مسار موظف معين
func (h *LocationHandler) GetEmployeeTrack(c *gin.Context) {
	employeeID := c.Param("id")
	limit := c.DefaultQuery("limit", "50")

	rows, err := h.DB.Query(`
		SELECT latitude, longitude, recorded_at
		FROM location_tracking
		WHERE user_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`, employeeID, limit)

	if err != nil {
		log.Printf("❌ فشل جلب مسار الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب المسار"})
		return
	}
	defer rows.Close()

	var track []gin.H
	for rows.Next() {
		var lat, lng float64
		var recordedAt time.Time

		if err := rows.Scan(&lat, &lng, &recordedAt); err != nil {
			continue
		}

		track = append(track, gin.H{
			"latitude":  lat,
			"longitude": lng,
			"timestamp": recordedAt,
		})
	}

	c.JSON(http.StatusOK, track)
}

// GetEmployeeSecurityNotes - جلب ملاحظات أمنية عن موظف
func (h *LocationHandler) GetEmployeeSecurityNotes(c *gin.Context) {
	employeeID := c.Param("id")

	rows, err := h.DB.Query(`
		SELECT 
			title, 
			body, 
			created_at,
			type
		FROM notifications
		WHERE user_id = $1 AND type = 'geofence_alert'
		ORDER BY created_at DESC
		LIMIT 20
	`, employeeID)

	if err != nil {
		log.Printf("❌ فشل جلب الملاحظات الأمنية: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملاحظات"})
		return
	}
	defer rows.Close()

	var notes []gin.H
	for rows.Next() {
		var title, body, notificationType string
		var createdAt time.Time

		if err := rows.Scan(&title, &body, &createdAt, &notificationType); err != nil {
			continue
		}

		notes = append(notes, gin.H{
			"title":      title,
			"body":       body,
			"created_at": createdAt,
			"type":       notificationType,
		})
	}

	c.JSON(http.StatusOK, notes)
}

// GetLocationLogs - جلب سجل المواقع للمدير
func (h *LocationHandler) GetLocationLogs(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT
			lt.latitude,
			lt.longitude,
			lt.recorded_at,
			u.full_name,
			u.email,
			COALESCE(w.name, a.worksite_name_for_history) as worksite_name,
			a.status as attendance_status
		FROM location_tracking lt
		JOIN users u ON lt.user_id = u.id
		LEFT JOIN attendance a ON a.user_id = u.id AND a.status = 'in_progress'
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee'
		ORDER BY lt.recorded_at DESC
		LIMIT 100
	`)

	if err != nil {
		log.Printf("❌ فشل جلب سجل المواقع: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب السجل"})
		return
	}
	defer rows.Close()

	var logs []gin.H
	for rows.Next() {
		var lat, lng float64
		var recordedAt time.Time
		var fullName, email, worksiteName, attendanceStatus string

		if err := rows.Scan(&lat, &lng, &recordedAt, &fullName, &email, &worksiteName, &attendanceStatus); err != nil {
			continue
		}

		logs = append(logs, gin.H{
			"employee":          fullName,
			"email":             email,
			"latitude":          lat,
			"longitude":         lng,
			"recorded_at":       recordedAt,
			"worksite":          worksiteName,
			"attendance_status": attendanceStatus,
		})
	}

	c.JSON(http.StatusOK, logs)
}
