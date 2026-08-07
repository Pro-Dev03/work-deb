package handlers

import (
	"net/http"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// RefreshToken يجدد access token باستخدام refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	lang := i18n.Detect(c)

	// الحصول على refresh token من cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_refresh_token")})
		return
	}

	// التحقق من refresh token
	userID, err := h.RefreshTokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_refresh_token")})
		return
	}

	// جلب بيانات المستخدم
	var role string
	err = h.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// إنشاء access token جديد
	newToken, err := h.AuthService.GenerateToken(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_session_create_failed")})
		return
	}

	// تدوير refresh token (إنشاء جديد وإلغاء القديم)
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	newRefreshToken, err := h.RefreshTokenService.RotateRefreshToken(refreshToken, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_token_rotation_failed")})
		return
	}

	// إرسال access token جديد كـ httpOnly cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"access_token",
		newToken,
		3600, // 1 hour
		"/",
		"",
		h.Config.ShouldUseSecureCookies(), // secure - true in production, false in development
		true, // httpOnly
	)

	// إرسال refresh token جديد كـ httpOnly cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*3600, // 7 days
		"/",
		"",
		h.Config.ShouldUseSecureCookies(), // secure - true in production, false in development
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(lang, "msg_session_renewed")})
}

// Logout يخرج المستخدم ويلغي refresh token
func (h *AuthHandler) Logout(c *gin.Context) {
	lang := i18n.Detect(c)

	// الحصول على refresh token من cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		// إلغاء refresh token
		h.RefreshTokenService.RevokeRefreshToken(refreshToken)
	}

	// حذف cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", h.Config.ShouldUseSecureCookies(), true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.Config.ShouldUseSecureCookies(), true)

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(lang, "msg_logged_out")})
}

// LogoutAll يخرج المستخدم من جميع الأجهزة
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	lang := i18n.Detect(c)

	// الحصول على user_id من context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_please_login")})
		return
	}

	// إلغاء جميع refresh tokens للمستخدم
	if err := h.RefreshTokenService.RevokeAllUserTokens(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// حذف cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", h.Config.ShouldUseSecureCookies(), true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.Config.ShouldUseSecureCookies(), true)

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(lang, "msg_logged_out_all_devices")})
}
