package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"worktrack/backend/internal/database"
	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// حدود الطلبات - محسنة للأمان
const (
	generalLimit = 100  // حد عام محسد للإنتاج (كان 1000)
	authLimit    = 10   // حد صارم للمصادقة (كان 100)
	strictLimit  = 5    // حد للعمليات الحساسة
	blockDuration = 10 * time.Minute // مدة الحظر عند الانتهاك (كان 5 دقائق)
	rateWindow   = 1 * time.Minute // نافذة الوقت
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// إذا لم يكن Redis متصلاً، تجاوز الـ rate limiter
		if database.RedisClient == nil {
			c.Next()
			return
		}

		lang := i18n.Detect(c)
		ip := c.ClientIP()
		ctx := context.Background()

		// تحديد نوع الطلب
		isAuthEndpoint := c.Request.URL.Path == "/api/v1/auth/login" ||
			c.Request.URL.Path == "/api/v1/auth/phone-login"

		isStrictEndpoint := c.Request.URL.Path == "/api/v1/admin/reset-device" ||
			c.Request.URL.Path == "/api/v1/admin/reset-customer-password" ||
			c.Request.URL.Path == "/api/v1/auth/employee-phone"

		limit := generalLimit
		if isAuthEndpoint {
			limit = authLimit
		} else if isStrictEndpoint {
			limit = strictLimit
		}

		// مفتاح Redis
		key := fmt.Sprintf("ratelimit:%s:%s", ip, c.Request.URL.Path)

		// التحقق من الحظر
		blockedKey := fmt.Sprintf("blocked:%s", ip)
		blocked, err := database.RedisClient.Exists(ctx, blockedKey).Result()
		if err == nil && blocked > 0 {
			ttl, _ := database.RedisClient.TTL(ctx, blockedKey).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":         i18n.T(lang, "err_too_many_requests"),
				"blocked_until": time.Now().Add(ttl).Format(time.RFC3339),
			})
			return
		}

		// زيادة العداد
		count, err := database.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			// في حالة الخطأ، تجاوز الـ rate limiter
			c.Next()
			return
		}

		// إذا كان أول طلب، اضبط وقت انتهاء الصلاحية
		if count == 1 {
			database.RedisClient.Expire(ctx, key, rateWindow)
		}

		// التحقق من الحد
		if count > int64(limit) {
			// حظر العنوان
			database.RedisClient.Set(ctx, blockedKey, "1", blockDuration)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":         i18n.T(lang, "err_too_many_requests"),
				"blocked_until": time.Now().Add(blockDuration).Format(time.RFC3339),
			})
			return
		}

		// إضافة معلومات عن الـ rate limit في الرأس
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(limit)-int(count)))
		c.Header("X-RateLimit-Reset", time.Now().Add(rateWindow).Format(time.RFC3339))

		c.Next()
	}
}
