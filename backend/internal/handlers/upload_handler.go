package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	DB *sql.DB
}

func NewUploadHandler(db *sql.DB) *UploadHandler {
	return &UploadHandler{DB: db}
}

func (h *UploadHandler) UploadTaskPhoto(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	log.Printf("📸 محاولة رفع صورة من المستخدم: %s", userID)

	attendanceID := c.PostForm("attendance_id")
	log.Printf("📝 attendance_id: %s", attendanceID)

	if attendanceID == "" {
		log.Println("❌ attendance_id مطلوب")
		c.JSON(http.StatusBadRequest, gin.H{"error": "attendance_id مطلوب"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		log.Printf("❌ فشل الحصول على الصورة: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_no_photo")})
		return
	}

	log.Printf("📄 اسم الملف: %s, الحجم: %d بايت", file.Filename, file.Size)

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "نوع الملف غير مدعوم"})
		return
	}

	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "حجم الملف كبير جداً"})
		return
	}

	var employeeID string
	err = h.DB.QueryRow(`
		SELECT user_id FROM attendance WHERE id = $1
	`, attendanceID).Scan(&employeeID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ سجل الحضور غير موجود: %s", attendanceID)
			c.JSON(http.StatusNotFound, gin.H{"error": "سجل الحضور غير موجود"})
			return
		}
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}

	if employeeID != userID {
		log.Printf("❌ المستخدم ليس صاحب السجل")
		c.JSON(http.StatusForbidden, gin.H{"error": "ليس لديك صلاحية"})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	filename := fmt.Sprintf("%d_%s_%s%s", 
		time.Now().Unix(), 
		userID.(string)[:8], 
		uuid.NewString()[:8],
		ext,
	)

	filePath := filepath.Join(uploadDir, filename)
	
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("❌ فشل حفظ الصورة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حفظ الصورة"})
		return
	}

	photoURL := fmt.Sprintf("/uploads/%s", filename)
	
	log.Printf("📸 تم حفظ الصورة: %s", filePath)
	log.Printf("📸 رابط الصورة: %s", photoURL)

	now := time.Now()
	
	result, err := h.DB.Exec(`
		UPDATE attendance 
		SET photo_url = $1, photo_uploaded_at = $2, photo_notes = $3
		WHERE id = $4
	`, photoURL, now, "تم رفع صورة إثبات الإنجاز", attendanceID)

	if err != nil {
		log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حفظ رابط الصورة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ تم تحديث %d سجل", rowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message":       "✅ تم رفع الصورة بنجاح",
		"photo_url":     photoURL,
		"attendance_id": attendanceID,
		"filename":      filename,
	})
}

func (h *UploadHandler) GetAttendancePhotos(c *gin.Context) {
	log.Println("📸 جلب صور إثبات الإنجاز...")

	rows, err := h.DB.Query(`
		SELECT 
			a.id,
			a.photo_url,
			a.photo_uploaded_at,
			a.photo_notes,
			u.full_name as employee_name,
			u.email as employee_email,
			COALESCE(t.title, 'بدون مهمة') as task_title,
			COALESCE(w.name, a.worksite_name_for_history, 'بدون موقع') as worksite_name
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN tasks t ON a.task_id = t.id
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.photo_url IS NOT NULL AND a.photo_url != ''
		ORDER BY a.photo_uploaded_at DESC
		LIMIT 50
	`)

	if err != nil {
		log.Printf("❌ فشل جلب الصور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الصور"})
		return
	}
	defer rows.Close()

	var photos []gin.H
	for rows.Next() {
		var id, photoURL, photoNotes, employeeName, employeeEmail, taskTitle, worksiteName string
		var photoUploadedAt sql.NullTime

		if err := rows.Scan(&id, &photoURL, &photoUploadedAt, &photoNotes,
			&employeeName, &employeeEmail, &taskTitle, &worksiteName); err != nil {
			log.Printf("⚠️ خطأ في قراءة البيانات: %v", err)
			continue
		}

		photo := gin.H{
			"id":              id,
			"photo_url":       photoURL,
			"photo_notes":     photoNotes,
			"employee_name":   employeeName,
			"employee_email":  employeeEmail,
			"task_title":      taskTitle,
			"worksite_name":   worksiteName,
			"uploaded_at":     photoUploadedAt.Time,
		}

		photos = append(photos, photo)
	}

	log.Printf("✅ تم جلب %d صورة", len(photos))
	c.JSON(http.StatusOK, photos)
}

func (h *UploadHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	log.Printf("🗑️ محاولة حذف الصورة: %s", id)

	var photoURL string
	err := h.DB.QueryRow(`
		SELECT photo_url FROM attendance WHERE id = $1
	`, id).Scan(&photoURL)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ الصورة غير موجودة: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "الصورة غير موجودة"})
			return
		}
		log.Printf("❌ فشل جلب الصورة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الصورة"})
		return
	}

	if photoURL != "" {
		filename := filepath.Base(photoURL)
		filePath := filepath.Join("./uploads", filename)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				log.Printf("⚠️ فشل حذف الملف: %v", err)
			} else {
				log.Printf("✅ تم حذف الملف: %s", filePath)
			}
		}
	}

	_, err = h.DB.Exec(`
		UPDATE attendance 
		SET photo_url = NULL, photo_uploaded_at = NULL, photo_notes = NULL
		WHERE id = $1
	`, id)

	if err != nil {
		log.Printf("❌ فشل حذف الصورة من قاعدة البيانات: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف الصورة"})
		return
	}

	log.Printf("✅ تم حذف الصورة: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "✅ تم حذف الصورة بنجاح"})
}
