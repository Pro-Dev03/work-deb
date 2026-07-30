package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

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

// List - جلب جميع نقاط العمل مع الموظفين العاملين حالياً
func (h *WorksiteHandler) List(c *gin.Context) {
	// التحقق من وجود عمود assigned_employee_id للتوافق مع الإصدارات القديمة
	var columnExists bool
	h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'worksites' AND column_name = 'assigned_employee_id'
		)
	`).Scan(&columnExists)

	var rows *sql.Rows
	var err error

	if columnExists {
		// استخدام الاستعلام مع assigned_employee_id (الإصدار الحديث)
		rows, err = h.DB.Query(`
			SELECT w.id, w.name, w.address, w.latitude, w.longitude, w.radius_meters, w.is_active, w.created_at,
				u.id as assigned_employee_id, u.full_name as assigned_employee_name
			FROM worksites w
			LEFT JOIN users u ON w.assigned_employee_id = u.id
			WHERE w.is_deleted = FALSE
			ORDER BY w.created_at DESC`)
	} else {
		// استخدام الاستعلام بدون assigned_employee_id (للتوافق مع الإصدارات القديمة)
		rows, err = h.DB.Query(`
			SELECT w.id, w.name, w.address, w.latitude, w.longitude, w.radius_meters, w.is_active, w.created_at
			FROM worksites w
			WHERE w.is_deleted = FALSE
			ORDER BY w.created_at DESC`)
	}

	if err != nil {
		log.Printf("❌ فشل جلب نقاط العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب نقاط العمل"})
		return
	}
	defer rows.Close()

	var worksites []gin.H
	for rows.Next() {
		var w models.Worksite
		var assignedEmployeeID, assignedEmployeeName sql.NullString

		if columnExists {
			if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
				&w.RadiusMeters, &w.IsActive, &w.CreatedAt, &assignedEmployeeID, &assignedEmployeeName); err != nil {
				log.Printf("⚠️ خطأ في قراءة صف نقاط العمل: %v", err)
				continue
			}
		} else {
			if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
				&w.RadiusMeters, &w.IsActive, &w.CreatedAt); err != nil {
				log.Printf("⚠️ خطأ في قراءة صف نقاط العمل: %v", err)
				continue
			}
		}

		worksite := gin.H{
			"id":            w.ID,
			"name":          w.Name,
			"address":       w.Address,
			"latitude":      w.Latitude,
			"longitude":     w.Longitude,
			"radius_meters": w.RadiusMeters,
			"is_active":     w.IsActive,
			"created_at":    w.CreatedAt,
		}

		if columnExists && assignedEmployeeID.Valid && assignedEmployeeName.Valid {
			worksite["assigned_to"] = gin.H{
				"id":   assignedEmployeeID.String,
				"name": assignedEmployeeName.String,
			}
		}

		worksites = append(worksites, worksite)
	}

	// جلب الموظفين العاملين حالياً في كل نقطة عمل
	for i, ws := range worksites {
		if ws["id"] == "unassigned" {
			continue
		}
		
		workingRows, err := h.DB.Query(`
			SELECT DISTINCT u.id, u.full_name, a.id as attendance_id
			FROM users u
			JOIN attendance a ON u.id = a.user_id
			WHERE a.status = 'in_progress' 
			AND a.worksite_id = $1
			AND u.role = 'employee'
			AND a.worksite_id IS NOT NULL
		`, ws["id"])
		
		if err == nil {
			var workingEmployees []gin.H
			for workingRows.Next() {
				var id, name, attendanceID string
				if err := workingRows.Scan(&id, &name, &attendanceID); err == nil {
					workingEmployees = append(workingEmployees, gin.H{
						"id":             id,
						"name":           name,
						"attendance_id":  attendanceID,
					})
				}
			}
			workingRows.Close()
			worksites[i]["working_employees"] = workingEmployees
		}
	}

	// إضافة خيار "غير معين" للموظفين الذين يعملون بدون نقطة عمل
	unassignedRows, err := h.DB.Query(`
		SELECT DISTINCT u.id, u.full_name, a.id as attendance_id
		FROM users u
		JOIN attendance a ON u.id = a.user_id
		WHERE a.status = 'in_progress'
		AND a.worksite_id IS NULL
		AND u.role = 'employee'
	`)
	if err == nil {
		var unassignedEmployees []gin.H
		for unassignedRows.Next() {
			var id, name, attendanceID string
			if err := unassignedRows.Scan(&id, &name, &attendanceID); err == nil {
				unassignedEmployees = append(unassignedEmployees, gin.H{
					"id":             id,
					"name":           name,
					"attendance_id":  attendanceID,
				})
			}
		}
		unassignedRows.Close()

		if len(unassignedEmployees) > 0 {
			worksites = append([]gin.H{
				{
					"id":                 "unassigned",
					"name":               "غير معين",
					"address":            "الموظفون الذين يعملون بدون نقطة عمل محددة",
					"latitude":           0,
					"longitude":          0,
					"radius_meters":      0,
					"is_active":          true,
					"created_at":         "now()",
					"is_unassigned":      true,
					"working_employees":  unassignedEmployees,
				},
			}, worksites...)
		}
	}

	c.JSON(http.StatusOK, worksites)
}

// GetAvailableWorksites - جلب نقاط العمل المتاحة للموظفين
func (h *WorksiteHandler) GetAvailableWorksites(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, name, address, latitude, longitude, radius_meters
		FROM worksites WHERE is_active = TRUE AND is_deleted = FALSE ORDER BY name`)
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

	// التحقق من وجود الاسم في النقاط النشطة فقط (تجاهل المحذوفة)
	var exists bool
	errCheck := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM worksites 
			WHERE name = $1 AND is_deleted = FALSE
		)
	`, req.Name).Scan(&exists)
	
	if errCheck != nil {
		log.Printf("❌ فشل التحقق من الاسم: %v", errCheck)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من الاسم"})
		return
	}
	
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "اسم نقطة العمل موجود بالفعل",
			"suggestion": "يمكنك استخدام اسم مختلف أو استعادة النقطة المحذوفة",
		})
		return
	}

	id := uuid.NewString()
	_, errInsert := h.DB.Exec(`
		INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, is_deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE, FALSE, now(), now())`,
		id, req.Name, req.Address, req.Latitude, req.Longitude, req.RadiusMeters,
	)
	if errInsert != nil {
		log.Printf("❌ فشل الإضافة: %v", errInsert)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الإضافة"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "تم الإضافة"})
}

// Update - تعديل نقطة العمل
func (h *WorksiteHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
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

	// التحقق من وجود نقطة العمل
	var exists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM worksites 
			WHERE id = $1 AND is_deleted = FALSE
		)
	`, id).Scan(&exists)
	
	if err != nil {
		log.Printf("❌ فشل التحقق من نقطة العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من نقطة العمل"})
		return
	}
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "نقطة العمل غير موجودة"})
		return
	}

	// التحقق من وجود الاسم في نقاط العمل الأخرى (تجاهل المحذوفة ونفس النقطة)
	var nameExists bool
	errCheck := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM worksites 
			WHERE name = $1 AND id != $2 AND is_deleted = FALSE
		)
	`, req.Name, id).Scan(&nameExists)
	
	if errCheck != nil {
		log.Printf("❌ فشل التحقق من الاسم: %v", errCheck)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من الاسم"})
		return
	}
	
	if nameExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "اسم نقطة العمل موجود بالفعل",
			"suggestion": "يمكنك استخدام اسم مختلف",
		})
		return
	}

	// تحديث نقطة العمل
	_, err = h.DB.Exec(`
		UPDATE worksites 
		SET name = $1, address = $2, latitude = $3, longitude = $4, 
		    radius_meters = $5, updated_at = now()
		WHERE id = $6
	`, req.Name, req.Address, req.Latitude, req.Longitude, req.RadiusMeters, id)
	
	if err != nil {
		log.Printf("❌ فشل التعديل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التعديل"})
		return
	}

	log.Printf("✅ تم تعديل نقطة العمل: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "تم التعديل بنجاح"})
}

// Delete - حذف نقطة العمل (Soft Delete)
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

	// 1. حذف المهام المرتبطة
	result, err := tx.Exec(`
		DELETE FROM tasks WHERE worksite_id = $1
	`, id)
	if err != nil {
		log.Printf("❌ فشل حذف المهام: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف المهام"})
		return
	}
	tasksRows, _ := result.RowsAffected()
	log.Printf("🗑️ تم حذف %d مهمة", tasksRows)

	// 2. Soft Delete لنقطة العمل (بدلاً من الحذف الفعلي)
	result, err = tx.Exec(`
		UPDATE worksites SET is_deleted = TRUE, updated_at = now() WHERE id = $1
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

	log.Printf("✅ تم حذف نقطة العمل: %s (Soft Delete، حذف %d مهام)", id, tasksRows)
	c.JSON(http.StatusOK, gin.H{
		"message":       "تم حذف نقطة العمل بنجاح",
		"tasks_deleted": tasksRows,
		"soft_delete":   true,
	})
}

// AssignEmployee - تعيين موظف لنقطة عمل وتسجيل دخوله تلقائياً
func (h *WorksiteHandler) AssignEmployee(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
		WorksiteID string `json:"worksite_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	// التحقق من وجود الموظف وأنه موظف (employee)
	var employeeExists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users 
			WHERE id = $1 AND role = 'employee' AND is_active = TRUE
		)
	`, req.EmployeeID).Scan(&employeeExists)
	
	if err != nil {
		log.Printf("❌ فشل التحقق من الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من الموظف"})
		return
	}
	
	if !employeeExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "الموظف غير موجود أو غير نشط"})
		return
	}

	// التحقق من وجود نقطة العمل
	var worksiteExists bool
	var worksiteLat, worksiteLng float64
	var worksiteRadius int
	
	err = h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM worksites 
			WHERE id = $1 AND is_deleted = FALSE
		), latitude, longitude, radius_meters
	`, req.WorksiteID).Scan(&worksiteExists, &worksiteLat, &worksiteLng, &worksiteRadius)
	
	if err != nil {
		log.Printf("❌ فشل التحقق من نقطة العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من نقطة العمل"})
		return
	}
	
	if !worksiteExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "نقطة العمل غير موجودة"})
		return
	}

	// التحقق من أن الموظف ليس لديه وردية نشطة حالياً
	var hasActiveShift bool
	err = h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM attendance 
			WHERE user_id = $1 AND status = 'in_progress'
		)
	`, req.EmployeeID).Scan(&hasActiveShift)
	
	if err != nil {
		log.Printf("❌ فشل التحقق من الوردية النشطة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل التحقق من الوردية النشطة"})
		return
	}
	
	if hasActiveShift {
		c.JSON(http.StatusBadRequest, gin.H{"error": "الموظف لديه وردية نشطة حالياً"})
		return
	}

	// بدء معاملة
	tx, err := h.DB.Begin()
	if err != nil {
		log.Printf("❌ فشل بدء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل بدء عملية التعيين"})
		return
	}
	defer tx.Rollback()

	// 1. تحديث assigned_employee_id في جدول worksites
	_, err = tx.Exec(`
		UPDATE worksites 
		SET assigned_employee_id = $1, updated_at = now() 
		WHERE id = $2
	`, req.EmployeeID, req.WorksiteID)
	
	if err != nil {
		log.Printf("❌ فشل تعيين الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تعيين الموظف"})
		return
	}

	// 2. إنشاء سجل حضور جديد (تسجيل دخول تلقائي)
	attendanceID := uuid.NewString()
	now := time.Now()
	
	// استخدام إحداثيات نقطة العمل كموقع للموظف (بما أنه تم تعيينه من قبل المدير)
	distance := 0.0 // المسافة 0 لأنه في موقع نقطة العمل
	
	_, err = tx.Exec(`
		INSERT INTO attendance (id, user_id, worksite_id, check_in_time, 
			check_in_lat, check_in_lng, check_in_distance_meters, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'in_progress')
	`, attendanceID, req.EmployeeID, req.WorksiteID, now, worksiteLat, worksiteLng, distance)
	
	if err != nil {
		log.Printf("❌ فشل تسجيل الدخول التلقائي: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تسجيل الدخول التلقائي"})
		return
	}

	// إنهاء المعاملة
	if err := tx.Commit(); err != nil {
		log.Printf("❌ فشل إنهاء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنهاء عملية التعيين"})
		return
	}

	log.Printf("✅ تم تعيين الموظف %s لنقطة العمل %s وتسجيل دخوله تلقائياً", req.EmployeeID, req.WorksiteID)
	c.JSON(http.StatusOK, gin.H{
		"message": "تم التعيين وتسجيل الدخول بنجاح",
		"attendance_id": attendanceID,
		"check_in_time": now,
	})
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
