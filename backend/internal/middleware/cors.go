package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware يسمح فقط لروابط الواجهات الأمامية (Vercel) بالوصول للـ API
// نسمح أيضاً بهيدر X-Lang المخصص الذي تستخدمه الواجهات لتحديد لغة الرد (ar/he/en)
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Lang", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
