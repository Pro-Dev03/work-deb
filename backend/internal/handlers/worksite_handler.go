package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"worktrack/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorksiteHandler struct {
	DB *sql.DB
}

func NewWorksiteHandler(db *sql.DB) *WorksiteHandler {
	return &WorksiteHandler{DB: db}
}

// List - جلب جميع نقاط العمل
func (h *WorksiteHandler) List(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, name, address, latitude, longitude, radius_meters, is_active, created_at
		FROM worksites ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("❌ فشل جلب نقاط العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب نقاط العمل"})
		return
	}
	defer rows.Close()

	var worksites []gin.H
	for rows.Next() {
		var w models.Worksite
		if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
			&w.RadiusMeters, &w.IsActive, &w.CreatedAt); err != nil {
			continue
		}
		worksites = append(worksites, gin.H{
			"id":            w.ID,
			"name":          w.Name,
			"address":       w.Address,
			"latitude":      w.Latitude,
			"longitude":     w.Longitude,
			"radius_meters": w.RadiusMeters,
			"is_active":     w.IsActive,
			"created_at":    w.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, worksites)
}

// GetAvailableWorksites - جلب نقاط العمل المتاحة للموظفين
func (h *WorksiteHandler) GetAvailableWorksites(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, name, address, latitude, longitude, radius_meters
		FROM worksites WHERE is_active = TRUE ORDER BY name`)
	if err != nil {
		log.Printf("❌ فشل جلب نقاط العمل المتاحة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب نقاط العمل"})
		return
	}
	defer rows.Close()

	var worksites []gin.H
	for rows.Next() {
		var id, name, address string
		var lat, lng float64
		var radius int
		if err := rows.Scan(&id, &name, &address, &lat, &lng, &radius); err != nil {
			continue
		}
		worksites = append(worksites, gin.H{
			"id":            id,
			"name":          name,
			"address":       address,
			"latitude":      lat,
			"longitude":     lng,
			"radius_meters": radius,
		})
	}
	c.JSON(http.StatusOK, worksites)
}

// Create - إضافة نقطة عمل جديدة
func (h *WorksiteHandler) Create(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Address      string  `json:"address"`
		Latitude     float64 `json:"latitude" binding:"required"`
		Longitude    float64 `json:"longitude" binding:"required"`
		RadiusMeters int     `json:"radius_meters" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}
	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE, now(), now())`,
		id, req.Name, req.Address, req.Latitude, req.Longitude, req.RadiusMeters,
	)
	if err != nil {
		log.Printf("❌ فشل الإضافة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الإضافة"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "تم الإضافة"})
}

// Delete - حذف نقطة العمل مع كل ما يرتبط بها
func (h *WorksiteHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	log.Printf("🗑️ محاولة حذف نقطة العمل: %s", id)

	// بدء معاملة
	tx, err := h.DB.Begin()
	if err != nil {
		log.Printf("❌ فشل بدء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل بدء عملية الحذف"})
		return
	}
	defer tx.Rollback()

	// 1. حذف سجلات الحضور المرتبطة
	result, err := tx.Exec(`
		DELETE FROM attendance WHERE worksite_id = $1
	`, id)
	if err != nil {
		log.Printf("❌ فشل حذف سجلات الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف سجلات الحضور"})
		return
	}
	attendanceRows, _ := result.RowsAffected()
	log.Printf("🗑️ تم حذف %d سجل حضور", attendanceRows)

	// 2. حذف المهام المرتبطة
	result, err = tx.Exec(`
		DELETE FROM tasks WHERE worksite_id = $1
	`, id)
	if err != nil {
		log.Printf("❌ فشل حذف المهام: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف المهام"})
		return
	}
	tasksRows, _ := result.RowsAffected()
	log.Printf("🗑️ تم حذف %d مهمة", tasksRows)

	// 3. حذف نقاط العمل
	result, err = tx.Exec(`
		DELETE FROM worksites WHERE id = $1
	`, id)
	if err != nil {
		log.Printf("❌ فشل حذف نقطة العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف نقطة العمل"})
		return
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		log.Printf("❌ نقطة العمل غير موجودة: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "نقطة العمل غير موجودة"})
		return
	}

	// إنهاء المعاملة
	if err := tx.Commit(); err != nil {
		log.Printf("❌ فشل إنهاء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنهاء عملية الحذف"})
		return
	}

	log.Printf("✅ تم حذف نقطة العمل: %s (حذف %d حضور، %d مهام)", id, attendanceRows, tasksRows)
	c.JSON(http.StatusOK, gin.H{
		"message":            "تم حذف نقطة العمل بنجاح",
		"attendance_deleted": attendanceRows,
		"tasks_deleted":      tasksRows,
	})
}

// AssignEmployee - تعيين موظف لنقطة عمل
func (h *WorksiteHandler) AssignEmployee(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
		WorksiteID string `json:"worksite_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "تم التعيين"})
}

// GetAvailableEmployees - جلب الموظفين المتاحين
func (h *WorksiteHandler) GetAvailableEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, full_name FROM users WHERE role = 'employee' AND is_active = TRUE`)
	if err != nil {
		log.Printf("❌ فشل جلب الموظفين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين"})
		return
	}
	defer rows.Close()
	var employees []gin.H
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err == nil {
			employees = append(employees, gin.H{"id": id, "full_name": name})
		}
	}
	c.JSON(http.StatusOK, employees)
}
