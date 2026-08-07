package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DB          *sql.DB
	AuthService *services.AuthService
}

func NewAuthHandler(db *sql.DB, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{DB: db, AuthService: authService}
}

func (h *AuthHandler) checkSubscriptionStatus(userID string, status string, expiresAt sql.NullTime) error {
	if status == "canceled" {
		return errors.New("subscription canceled")
	}

	if status == "expired" {
		return errors.New("subscription expired")
	}

	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		_, err := h.DB.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, userID)
		if err != nil {
			log.Printf("⚠️ failed to mark subscription expired for user %s: %v", userID, err)
		}
		return errors.New("subscription expired")
	}

	return nil
}

func (h *AuthHandler) validatePassword(password, storedHash, email string) bool {
	if h.AuthService.CheckPassword(password, storedHash) {
		return true
	}

	var passwordMatches bool
	err := h.DB.QueryRow(`SELECT crypt($1, password_hash) = password_hash FROM users WHERE email = $2`, password, email).Scan(&passwordMatches)
	if err != nil {
		log.Printf("⚠️ password fallback check failed for %s: %v", email, err)
		return false
	}

	return passwordMatches
}

// =============================================
// 1. تسجيل دخول المدير
// =============================================
func (h *AuthHandler) Login(c *gin.Context) {
	lang := i18n.Detect(c)

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_login_fields")})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	var userID, fullName, role, passwordHash, subscriptionStatus string
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT id, full_name, role, password_hash, subscription_status, subscription_expires_at
		FROM users 
		WHERE email = $1 AND is_active = TRUE
	`, req.Email).Scan(&userID, &fullName, &role, &passwordHash, &subscriptionStatus, &subscriptionExpiresAt)

	if err != nil {
		log.Printf("❌ مستخدم غير موجود: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_credentials")})
		return
	}

	if !h.validatePassword(req.Password, passwordHash, req.Email) {
		log.Printf("❌ كلمة مرور خاطئة")
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_credentials")})
		return
	}

	if err := h.checkSubscriptionStatus(userID, subscriptionStatus, subscriptionExpiresAt); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
		return
	}

	token, err := h.AuthService.GenerateToken(userID, role)
	if err != nil {
		log.Printf("❌ فشل إنشاء التوكن: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_session_create_failed")})
		return
	}

	log.Printf("✅ تسجيل دخول ناجح: %s", fullName)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":                  userID,
			"full_name":           fullName,
			"role":                role,
			"subscription_status": subscriptionStatus,
			"subscription_expires_at": func() interface{} {
				if subscriptionExpiresAt.Valid {
					return subscriptionExpiresAt.Time
				}
				return nil
			}(),
		},
	})
}

// =============================================
// 2. جلب جميع الموظفين
// =============================================
func (h *AuthHandler) ListEmployees(c *gin.Context) {
	log.Println("📋 جلب قائمة الموظفين...")

	rows, err := h.DB.Query(`
		SELECT
			u.id, u.full_name, u.email, u.phone, u.role, u.is_active, u.created_at,
			COALESCE(u.device_model, '') as device_model,
			COALESCE(ws.name, a.worksite_name_for_history) as current_worksite,
			ws.id as current_worksite_id
		FROM users u
		LEFT JOIN attendance a ON u.id = a.user_id AND a.status = 'in_progress'
		LEFT JOIN worksites ws ON a.worksite_id = ws.id
		ORDER BY u.created_at DESC
	`)

	if err != nil {
		log.Printf("❌ خطأ في جلب الموظفين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, role, createdAt, deviceModel string
		var isActive bool
		var currentWorksite, currentWorksiteID sql.NullString

		if err := rows.Scan(&id, &fullName, &email, &phone, &role, &isActive, &createdAt, &deviceModel, &currentWorksite, &currentWorksiteID); err != nil {
			log.Printf("⚠️ خطأ في القراءة: %v", err)
			continue
		}

		employee := gin.H{
			"id":                id,
			"full_name":         fullName,
			"email":             email,
			"phone":             phone,
			"role":              role,
			"is_active":         isActive,
			"created_at":        createdAt,
			"device_model":      deviceModel,
			"is_registered":     deviceModel != "",
			"current_worksite":  currentWorksite.String,
			"current_worksite_id": currentWorksiteID.String,
		}

		employees = append(employees, employee)
	}

	log.Printf("✅ تم جلب %d موظف", len(employees))
	c.JSON(http.StatusOK, employees)
}

// =============================================
// 3. إنشاء موظف جديد
// =============================================
func (h *AuthHandler) CreateEmployeePhone(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 إنشاء موظف: %s, %s", req.FullName, req.Phone)

	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`, req.Phone).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "حدث خطأ"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "رقم الهاتف مستخدم"})
		return
	}

	id := uuid.NewString()
	role := "employee"
	if req.Role != "" {
		role = req.Role
	}

	bytes := make([]byte, 8)
	rand.Read(bytes)
	password := hex.EncodeToString(bytes)
	hash, _ := h.AuthService.HashPassword(password)
	email := req.Phone + "@worktrack.com"

	_, err = h.DB.Exec(`
		INSERT INTO users (
			id, full_name, email, phone, password_hash, 
			role, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, TRUE, now(), now())
	`, id, req.FullName, email, req.Phone, hash, role)

	if err != nil {
		log.Printf("❌ فشل إنشاء الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنشاء الموظف"})
		return
	}

	log.Printf("✅ تم إنشاء موظف: %s", req.FullName)

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "تم إنشاء الموظف بنجاح",
		"user": gin.H{
			"id":        id,
			"full_name": req.FullName,
			"phone":     req.Phone,
			"email":     email,
			"role":      role,
		},
	})
}

// =============================================
// 4. تسجيل دخول الموظف (مع التحقق من نوع الجهاز)
// =============================================
func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	lang := i18n.Detect(c)
	log.Println("========================================")
	log.Println("📱 تسجيل دخول الموظف...")

	var req struct {
		Phone       string `json:"phone" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
		DeviceModel string `json:"device_model" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📱 رقم الهاتف: %s", req.Phone)
	log.Printf("🆔 معرف الجهاز: %s", req.DeviceID)
	log.Printf("📱 نوع الجهاز: %s", req.DeviceModel)

	var userID, fullName, role, storedDeviceID, storedDeviceModel, subscriptionStatus string
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT 
			id, 
			full_name, 
			role, 
			COALESCE(device_id, ''), 
			COALESCE(device_model, ''),
			subscription_status,
			subscription_expires_at
		FROM users 
		WHERE phone = $1 AND is_active = TRUE
	`, req.Phone).Scan(&userID, &fullName, &role, &storedDeviceID, &storedDeviceModel, &subscriptionStatus, &subscriptionExpiresAt)

	if err != nil {
		log.Printf("❌ المستخدم غير موجود: %s", req.Phone)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "رقم الهاتف غير مسجل"})
		return
	}

	log.Printf("✅ تم العثور على المستخدم: %s", fullName)

	// =============================================
	// التحقق من الاشتراك للمديرين فقط
	// =============================================
	if role == "admin" {
		if err := h.checkSubscriptionStatus(userID, subscriptionStatus, subscriptionExpiresAt); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}
	}

	// =============================================
	// 🔒 طبقة أمان 1: التحقق من Device ID
	// =============================================
	if storedDeviceID != "" && storedDeviceID != req.DeviceID {
		log.Printf("🚨 Device ID غير مطابق!")
		log.Printf("   المسجل: %s", storedDeviceID)
		log.Printf("   الحالي: %s", req.DeviceID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":           "🚫 هذا الجهاز غير مصرح به",
			"device_mismatch": true,
		})
		return
	}

	// =============================================
	// 🔒 طبقة أمان 2: التحقق من Device Model (نوع الجهاز)
	// =============================================
	if storedDeviceModel != "" && storedDeviceModel != req.DeviceModel {
		log.Printf("🚨 نوع الجهاز غير مطابق!")
		log.Printf("   المسجل: %s", storedDeviceModel)
		log.Printf("   الحالي: %s", req.DeviceModel)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":          "🚫 نوع الجهاز غير مصرح به",
			"model_mismatch": true,
		})
		return
	}

	// =============================================
	// تسجيل الجهاز إذا كان جديداً
	// =============================================
	if storedDeviceID == "" || storedDeviceModel == "" {
		log.Println("📱 أول تسجيل دخول - تسجيل الجهاز...")

		_, err = h.DB.Exec(`
			UPDATE users 
			SET 
				device_id = $1,
				device_model = $2,
				phone_verified = TRUE,
				updated_at = now()
			WHERE id = $3
		`, req.DeviceID, req.DeviceModel, userID)

		if err != nil {
			log.Printf("❌ فشل تسجيل الجهاز: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تسجيل الجهاز"})
			return
		}

		log.Printf("✅ تم تسجيل الجهاز بنجاح!")
		log.Printf("   Device ID: %s", req.DeviceID)
		log.Printf("   Device Model: %s", req.DeviceModel)
	}

	// تحديث وقت آخر دخول
	_, err = h.DB.Exec(`
		UPDATE users 
		SET updated_at = now()
		WHERE id = $1
	`, userID)
	if err != nil {
		log.Printf("⚠️ فشل تحديث وقت الدخول: %v", err)
	}

	// إنشاء التوكن
	token, err := h.AuthService.GenerateToken(userID, role)
	if err != nil {
		log.Printf("❌ فشل إنشاء التوكن: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تسجيل الدخول"})
		return
	}

	log.Printf("✅ تسجيل دخول ناجح: %s", fullName)
	log.Printf("📱 نوع الجهاز المسجل: %s", req.DeviceModel)
	log.Println("========================================")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        userID,
			"full_name": fullName,
			"role":      role,
			"phone":     req.Phone,
		},
		"device_registered": storedDeviceID != "",
	})
}

// =============================================
// 5. بيانات المستخدم الحالي
// =============================================
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var fullName, email, role, subscriptionStatus string
	var phone sql.NullString
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT full_name, email, phone, role, subscription_status, subscription_expires_at
		FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email, &phone, &role, &subscriptionStatus, &subscriptionExpiresAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "المستخدم غير موجود"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                  userID,
		"full_name":           fullName,
		"email":               email,
		"phone":               func() interface{} {
			if phone.Valid {
				return phone.String
			}
			return ""
		}(),
		"role":                role,
		"subscription_status": subscriptionStatus,
		"subscription_expires_at": func() interface{} {
			if subscriptionExpiresAt.Valid {
				return subscriptionExpiresAt.Time
			}
			return nil
		}(),
	})
}

// =============================================
// 6. حذف موظف
// =============================================
func (h *AuthHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if id == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "لا يمكن حذف المدير الرئيسي"})
		return
	}

	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "المستخدم غير موجود"})
		return
	}

	_, err = h.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		log.Printf("❌ فشل حذف الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الحذف"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم حذف الموظف"})
}

// =============================================
// 7. إعادة تعيين جهاز الموظف
// =============================================
func (h *AuthHandler) ResetDevice(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	_, err := h.DB.Exec(`
		UPDATE users 
		SET 
			device_id = NULL,
			device_model = NULL,
			phone_verified = FALSE,
			updated_at = now()
		WHERE id = $1
	`, req.UserID)

	if err != nil {
		log.Printf("❌ فشل إعادة تعيين الجهاز: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إعادة تعيين الجهاز"})
		return
	}

	log.Printf("✅ تم إعادة تعيين جهاز المستخدم: %s", req.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "تم إعادة تعيين الجهاز"})
}

// =============================================
// 8. معلومات الجهاز
// =============================================
func (h *AuthHandler) GetDeviceInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var deviceID, deviceModel string
	err := h.DB.QueryRow(`
		SELECT 
			COALESCE(device_id, ''), 
			COALESCE(device_model, '') 
		FROM users WHERE id = $1
	`, userID).Scan(&deviceID, &deviceModel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب معلومات الجهاز"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id":     deviceID,
		"device_model":  deviceModel,
		"is_registered": deviceID != "" && deviceModel != "",
	})
}

// =============================================
// 9. تعديل اسم الموظف
// =============================================
func (h *AuthHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "المستخدم غير موجود"})
		return
	}

	_, err = h.DB.Exec(`
		UPDATE users 
		SET full_name = $1, updated_at = now()
		WHERE id = $2
	`, req.FullName, id)

	if err != nil {
		log.Printf("❌ فشل تعديل الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تعديل الموظف"})
		return
	}

	log.Printf("✅ تم تعديل اسم الموظف: %s -> %s", id, req.FullName)

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تعديل الموظف بنجاح",
		"user": gin.H{
			"id":        id,
			"full_name": req.FullName,
		},
	})
}
