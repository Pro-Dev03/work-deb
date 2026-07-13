package middleware

import (
	"net/http"
	"strings"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware يتحقق من وجود JWT صالح في هيدر Authorization
// ويضع user_id و role في السياق (Context) لاستخدامها في الـ Handlers
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)

		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_please_login")})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_session")})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
