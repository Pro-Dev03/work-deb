package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware يسمح فقط لروابط الواجهات الأمامية (Vercel) بالوصول للـ API
// نسمح أيضاً بهيدر X-Lang المخصص الذي تستخدمه الواجهات لتحديد لغة الرد (ar/he/en)
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	// تحذير إذا كان يوجد * في الـ origins
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "*" {
			log.Println("⚠️ تحذير أمني: استخدام '*' في CORS مسموح فقط في التطوير")
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Lang", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	})
}
