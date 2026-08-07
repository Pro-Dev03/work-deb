package handlers

import (
	"database/sql"
	"runtime"
	"time"

	"worktrack/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	DB *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

// Health فحص صحة شامل للخدمة
func (h *HealthHandler) Health(c *gin.Context) {
	startTime := time.Now()
	
	healthStatus := gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}

	// فحص الاتصال بقاعدة البيانات
	dbStatus := h.checkDatabase()
	healthStatus["database"] = dbStatus

	// فحص موارد النظام
	systemStatus := h.checkSystemResources()
	healthStatus["system"] = systemStatus

	// حساب وقت الاستجابة
	responseTime := time.Since(startTime).Milliseconds()
	healthStatus["response_time_ms"] = responseTime

	// تحديد الحالة العامة
	if dbStatus["status"] != "ok" {
		healthStatus["status"] = "degraded"
	}

	if responseTime > 1000 {
		healthStatus["status"] = "slow"
	}

	statusCode := 200
	if healthStatus["status"] != "ok" {
		statusCode = 503
	}

	utils.LogInfo("Health check performed", map[string]interface{}{
		"status": healthStatus["status"],
		"response_time_ms": responseTime,
		"database_status": dbStatus["status"],
	})

	c.JSON(statusCode, healthStatus)
}

// checkDatabase فحص الاتصال بقاعدة البيانات
func (h *HealthHandler) checkDatabase() gin.H {
	startTime := time.Now()
	
	err := h.DB.Ping()
	responseTime := time.Since(startTime).Milliseconds()

	if err != nil {
		utils.LogError("Database health check failed", map[string]interface{}{
			"error": err,
			"response_time_ms": responseTime,
		})
		return gin.H{
			"status":          "error",
			"error":           err.Error(),
			"response_time_ms": responseTime,
		}
	}

	// فحص إحصائيات الاتصال
	var maxOpenConns, openConns, inUse int
	h.DB.Stats()
	stats := h.DB.Stats()
	maxOpenConns = stats.MaxOpenConnections
	openConns = stats.OpenConnections
	inUse = stats.InUse

	return gin.H{
		"status":           "ok",
		"response_time_ms": responseTime,
		"connections": gin.H{
			"max_open":  maxOpenConns,
			"open":      openConns,
			"in_use":    inUse,
			"idle":      stats.Idle,
		},
	}
}

// checkSystemResources فحص موارد النظام
func (h *HealthHandler) checkSystemResources() gin.H {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// تحويل إلى ميجابايت
	allocMB := m.Alloc / 1024 / 1024
	totalAllocMB := m.TotalAlloc / 1024 / 1024
	sysMB := m.Sys / 1024 / 1024

	// عدد goroutines
	numGoroutines := runtime.NumGoroutine()

	return gin.H{
		"memory": gin.H{
			"alloc_mb":       allocMB,
			"total_alloc_mb": totalAllocMB,
			"sys_mb":         sysMB,
			"num_gc":         m.NumGC,
		},
		"goroutines": numGoroutines,
		"cpu": gin.H{
			"num_cpu": runtime.NumCPU(),
		},
	}
}

// Ready فحص جاهزية الخدمة
func (h *HealthHandler) Ready(c *gin.Context) {
	// فحص قاعدة البيانات فقط
	err := h.DB.Ping()
	if err != nil {
		utils.LogError("Readiness check failed", map[string]interface{}{
			"error": err,
		})
		c.JSON(503, gin.H{
			"status": "not_ready",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
	})
}

// Live فحص حياة الخدمة (بسيط وسريع)
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}
