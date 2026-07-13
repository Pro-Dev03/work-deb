package middleware

import (
	"net/http"
	"sync"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// RateLimiter بسيط يمنع أكثر من 60 طلب في الدقيقة لكل IP (حماية أساسية للـ API المجاني)
type visitor struct {
	count     int
	resetTime time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Now().After(v.resetTime) {
			v = &visitor{count: 0, resetTime: time.Now().Add(time.Minute)}
			visitors[ip] = v
		}
		v.count++
		count := v.count
		mu.Unlock()

		if count > 60 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": i18n.T(lang, "err_too_many_requests")})
			return
		}

		c.Next()
	}
}
