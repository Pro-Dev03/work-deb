package middleware

import (
	"net/http"
	"sync"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// RateLimiter محسّن مع حدود مختلفة لأنواع مختلفة من الطلبات
type visitor struct {
	requests    []time.Time // sliding window
	limit       int
	window      time.Duration
	lastCleanup time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex
var cleanupOnce sync.Once

func RateLimiter() gin.HandlerFunc {
	// بدء cleanup coroutine مرة واحدة فقط
	cleanupOnce.Do(func() {
		go cleanupExpiredVisitors()
	})

	return func(c *gin.Context) {
		lang := i18n.Detect(c)
		ip := c.ClientIP()

		// تحديد الحد حسب نوع الطلب
		limit := 1000 // default: 1000 requests/minute (زيادة للإنتاج)
		window := time.Minute

		// حدود أعلى للطلبات المتغيرة (POST, PUT, DELETE)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			limit = 500 // 500 requests/minute for write operations (زيادة للإنتاج)
		}

		// حدود أعلى لتسجيل الدخول
		if c.Request.URL.Path == "/api/v1/auth/login" || c.Request.URL.Path == "/api/v1/auth/phone-login" {
			limit = 50 // 50 requests/minute for login (زيادة للإنتاج)
			window = time.Minute
		}

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{
				requests:    []time.Time{},
				limit:       limit,
				window:      window,
				lastCleanup: time.Now(),
			}
			visitors[ip] = v
		}

		// إزالة الطلبات القديمة (sliding window)
		now := time.Now()
		cutoff := now.Add(-window)
		validRequests := []time.Time{}
		for _, reqTime := range v.requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		v.requests = validRequests

		// التحقق من الحد
		if len(v.requests) >= v.limit {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": i18n.T(lang, "err_too_many_requests"),
				"retry_after": window.Seconds(),
			})
			return
		}

		// إضافة الطلب الحالي
		v.requests = append(v.requests, now)
		v.lastCleanup = now
		mu.Unlock()

		c.Next()
	}
}

// cleanupExpiredVisitors ينظف الزوار غير النشطين
func cleanupExpiredVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		now := time.Now()
		for ip, v := range visitors {
			// إزالة الزوار الذين لم يطلبوا لمدة 10 دقائق
			if now.Sub(v.lastCleanup) > 10*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
