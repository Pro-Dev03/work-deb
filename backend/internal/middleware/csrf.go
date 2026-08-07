package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware يوفر حماية من هجمات CSRF
// بما أننا نستخدم httpOnly cookies مع SameSite=Strict، هذا middleware إضافي
type CSRFMiddleware struct {
	tokens map[string]bool
	mu     sync.RWMutex
}

func NewCSRFMiddleware() *CSRFMiddleware {
	return &CSRFMiddleware{
		tokens: make(map[string]bool),
	}
}

// GenerateToken ينشئ CSRF token جديد
func (m *CSRFMiddleware) GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)

	m.mu.Lock()
	m.tokens[token] = true
	m.mu.Unlock()

	return token
}

// ValidateToken يتحقق من صحة CSRF token
func (m *CSRFMiddleware) ValidateToken(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tokens[token]
}

// RevokeToken يلغي CSRF token
func (m *CSRFMiddleware) RevokeToken(token string) {
	m.mu.Lock()
	delete(m.tokens, token)
	m.mu.Unlock()
}

// CSRF middleware للتحقق من CSRF token في الطلبات المتغيرة
func (m *CSRFMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// تجاهل الطلبات الآمنة (GET, HEAD, OPTIONS)
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// تجاهل مسارات المصادقة
		path := c.Request.URL.Path
		if path == "/api/v1/auth/login" || path == "/api/v1/auth/phone-login" || 
		   path == "/api/v1/auth/refresh" || path == "/api/v1/auth/logout" {
			c.Next()
			return
		}

		// الحصول على CSRF token من header
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token مطلوب"})
			return
		}

		// التحقق من التوكن
		if !m.ValidateToken(csrfToken) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token غير صالح"})
			return
		}

		c.Next()
	}
}

// SetCSRFToken يضيف CSRF token إلى الاستجابة
func (m *CSRFMiddleware) SetCSRFToken(c *gin.Context) {
	token := m.GenerateToken()
	c.Header("X-CSRF-Token", token)
	c.Set("csrf_token", token)
}
