package middleware

import (
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
			"frame-ancestors 'self'; " +
			"form-action 'self'; " +
			"base-uri 'self';");

		// X-Content-Type-Options - يمنع sniffing من نوع المحتوى
		c.Header("X-Content-Type-Options", "nosniff");

		// X-Frame-Options - يمنع clickjacking
		c.Header("X-Frame-Options", "SAMEORIGIN");

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
