package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware يضيف رؤوس الأمان المهمة لجميع الاستجابات
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy - يمنع XSS من خلال السماح بمصادر محددة
		c.Header("Content-Security-Policy", 
			"default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; " +
			"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
			"img-src 'self' data: https: https://cdn.jsdelivr.net; " +
			"font-src 'self' data: https://cdn.jsdelivr.net; " +
			"connect-src 'self' https://worktrack-admin.onrender.com wss://worktrack-admin.onrender.com https://*.onrender.com wss://*.onrender.com; " +
			"frame-ancestors 'none'; " +
			"form-action 'self'; " +
			"base-uri 'self';");

		// X-Content-Type-Options - يمنع sniffing من نوع المحتوى
		c.Header("X-Content-Type-Options", "nosniff");

		// X-Frame-Options - يمنع clickjacking (DENY أكثر أماناً من SAMEORIGIN)
		c.Header("X-Frame-Options", "DENY");

		// Referrer-Policy - يتحكم في معلومات referrer
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin");

		// Permissions-Policy - يتحكم في استخدام ميزات المتصفح
		c.Header("Permissions-Policy", 
			"geolocation=(self), " +
			"camera=(self), " +
			"microphone=(self), " +
			"payment=(self)");

		// X-XSS-Protection - حماية إضافية ضد XSS
		c.Header("X-XSS-Protection", "1; mode=block");

		// Strict-Transport-Security - يفرض HTTPS (بدون preload لعدم وجود شهادة ثابتة)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains");

		c.Next()
	}
}

// HTTPSEnforcement يفرض HTTPS في الإنتاج بشكل ذكي
func HTTPSEnforcement() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only enforce in production
		if os.Getenv("APP_ENV") != "production" {
			c.Next()
			return
		}

		// Check if the request is secure
		if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
			// Redirect to HTTPS
			target := "https://" + c.Request.Host + c.Request.URL.Path
			if c.Request.URL.RawQuery != "" {
				target += "?" + c.Request.URL.RawQuery
			}
			c.Redirect(http.StatusMovedPermanently, target)
			c.Abort()
			return
		}

		c.Next()
	}
}
