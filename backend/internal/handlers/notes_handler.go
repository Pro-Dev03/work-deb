package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	DB *sql.DB
}

func NewNotesHandler(db *sql.DB) *NotesHandler {
	return &NotesHandler{DB: db}
}

type createNoteRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	WorksiteID string `json:"worksite_id"`
}

type createBulkNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

type updateNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateNote إنشاء ملاحظة جديدة من الأدمن للموظف
func (h *NotesHandler) CreateNote(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	// التحقق من أن المستخدم أدمن
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 إنشاء ملاحظة: admin_id=%s, employee_id=%s, worksite_id=%s", adminID, req.EmployeeID, req.WorksiteID)

	var id string
	var err error
	if req.WorksiteID != "" {
		err = h.DB.QueryRow(`
			INSERT INTO notes (admin_id, employee_id, content, worksite_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, adminID, req.EmployeeID, req.Content, req.WorksiteID).Scan(&id)
	} else {
		err = h.DB.QueryRow(`
			INSERT INTO notes (admin_id, employee_id, content)
			VALUES ($1, $2, $3)
			RETURNING id
		`, adminID, req.EmployeeID, req.Content).Scan(&id)
	}

	if err != nil {
		log.Printf("❌ فشل إنشاء الملاحظة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنشاء الملاحظة"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ تم إنشاء الملاحظة بنجاح",
		"id":      id,
	})
}

// UpdateNote تعديل ملاحظة موجودة
func (h *NotesHandler) UpdateNote(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	// التحقق من أن المستخدم أدمن
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	noteID := c.Param("id")
	var req updateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 تعديل ملاحظة: note_id=%s, admin_id=%s", noteID, adminID)

	result, err := h.DB.Exec(`
		UPDATE notes
		SET content = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND admin_id = $3
	`, req.Content, noteID, adminID)

	if err != nil {
		log.Printf("❌ فشل تعديل الملاحظة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تعديل الملاحظة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "الملاحظة غير موجودة أو ليس لديك صلاحية تعديلها"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ تم تعديل الملاحظة بنجاح",
	})
}

// DeleteNote حذف ملاحظة
func (h *NotesHandler) DeleteNote(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	// التحقق من أن المستخدم أدمن
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	noteID := c.Param("id")
	log.Printf("📝 حذف ملاحظة: note_id=%s, admin_id=%s", noteID, adminID)

	result, err := h.DB.Exec(`
		DELETE FROM notes
		WHERE id = $1 AND admin_id = $2
	`, noteID, adminID)

	if err != nil {
		log.Printf("❌ فشل حذف الملاحظة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف الملاحظة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "الملاحظة غير موجودة أو ليس لديك صلاحية حذفها"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ تم حذف الملاحظة بنجاح",
	})
}

// GetAdminNotes جلب جميع ملاحظات الأدمن
func (h *NotesHandler) GetAdminNotes(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	// التحقق من أن المستخدم أدمن
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	employeeID := c.DefaultQuery("employee_id", "")

	query := `
		SELECT
			n.id,
			n.admin_id,
			a.full_name as admin_name,
			n.employee_id,
			e.full_name as employee_name,
			e.phone as employee_phone,
			n.content,
			n.created_at,
			n.updated_at,
			n.is_read,
			n.worksite_id,
			w.name as worksite_name
		FROM notes n
		LEFT JOIN users a ON n.admin_id = a.id
		LEFT JOIN users e ON n.employee_id = e.id
		LEFT JOIN worksites w ON n.worksite_id = w.id
		WHERE n.admin_id = $1
	`
	args := []interface{}{adminID}
	argCount := 1

	if employeeID != "" {
		argCount++
		query += fmt.Sprintf(" AND n.employee_id = $%d", argCount)
		args = append(args, employeeID)
	}

	query += " ORDER BY n.created_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ فشل جلب الملاحظات: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملاحظات"})
		return
	}
	defer rows.Close()

	var notes []gin.H
	for rows.Next() {
		var id, adminID, adminName, employeeID, employeeName, employeePhone, content string
		var createdAt, updatedAt time.Time
		var isRead bool
		var worksiteID, worksiteName sql.NullString

		err := rows.Scan(
			&id, &adminID, &adminName, &employeeID, &employeeName, &employeePhone,
			&content, &createdAt, &updatedAt, &isRead, &worksiteID, &worksiteName,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		note := gin.H{
			"id":             id,
			"admin_id":       adminID,
			"admin_name":     adminName,
			"employee_id":    employeeID,
			"employee_name":  employeeName,
			"employee_phone": employeePhone,
			"content":        content,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
			"is_read":        isRead,
		}

		if worksiteID.Valid && worksiteID.String != "" {
			note["worksite_id"] = worksiteID.String
			note["worksite_name"] = worksiteName.String
		}

		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

// GetEmployeeNotes جلب ملاحظات الموظف من الأدمن
func (h *NotesHandler) GetEmployeeNotes(c *gin.Context) {
	employeeID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	query := `
		SELECT
			n.id,
			n.admin_id,
			a.full_name as admin_name,
			n.employee_id,
			n.content,
			n.created_at,
			n.updated_at,
			n.is_read,
			n.worksite_id,
			w.name as worksite_name
		FROM notes n
		LEFT JOIN users a ON n.admin_id = a.id
		LEFT JOIN worksites w ON n.worksite_id = w.id
		WHERE n.employee_id = $1
		ORDER BY n.created_at DESC
	`

	rows, err := h.DB.Query(query, employeeID)
	if err != nil {
		log.Printf("❌ فشل جلب الملاحظات: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملاحظات"})
		return
	}
	defer rows.Close()

	var notes []gin.H
	for rows.Next() {
		var id, adminID, adminName, employeeID, content string
		var createdAt, updatedAt time.Time
		var isRead bool
		var worksiteID, worksiteName sql.NullString

		err := rows.Scan(
			&id, &adminID, &adminName, &employeeID,
			&content, &createdAt, &updatedAt, &isRead, &worksiteID, &worksiteName,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		note := gin.H{
			"id":         id,
			"admin_id":   adminID,
			"admin_name": adminName,
			"content":    content,
			"created_at": createdAt,
			"updated_at": updatedAt,
			"is_read":    isRead,
		}

		if worksiteID.Valid && worksiteID.String != "" {
			note["worksite_id"] = worksiteID.String
			note["worksite_name"] = worksiteName.String
		}

		notes = append(notes, note)
	}

	// تحديث حالة القراءة للملاحظات غير المقروءة
	if len(notes) > 0 {
		_, err = h.DB.Exec(`
			UPDATE notes
			SET is_read = true
			WHERE employee_id = $1 AND is_read = false
		`, employeeID)
		if err != nil {
			log.Printf("❌ فشل تحديث حالة القراءة: %v", err)
		}
	}

	c.JSON(http.StatusOK, notes)
}

// MarkNoteAsRead تحديد ملاحظة كمقروءة
func (h *NotesHandler) MarkNoteAsRead(c *gin.Context) {
	employeeID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	noteID := c.Param("id")

	result, err := h.DB.Exec(`
		UPDATE notes
		SET is_read = true
		WHERE id = $1 AND employee_id = $2
	`, noteID, employeeID)

	if err != nil {
		log.Printf("❌ فشل تحديث حالة القراءة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تحديث حالة القراءة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "الملاحظة غير موجودة"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ تم تحديث حالة القراءة",
	})
}
