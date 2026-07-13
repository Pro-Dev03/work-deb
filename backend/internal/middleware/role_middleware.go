package middleware

import (
	"net/http"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// RequireRole يقيّد وصول Endpoint معين على دور محدد فقط (مثلاً "admin")
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)

		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_forbidden_role")})
			return
		}
		c.Next()
	}
}
