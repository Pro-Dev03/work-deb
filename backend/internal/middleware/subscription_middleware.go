package middleware

import (
	"database/sql"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

func SubscriptionMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_please_login")})
			return
		}

		var status string
		var expiresAt sql.NullTime
		err := db.QueryRow(`
			SELECT subscription_status, subscription_expires_at
			FROM users
			WHERE id = $1
		`, userID).Scan(&status, &expiresAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_session")})
			return
		}

		if status == "canceled" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}

		if status == "expired" || (expiresAt.Valid && time.Now().After(expiresAt.Time)) {
			if status != "expired" {
				_, _ = db.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, userID)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}

		c.Next()
	}
}
