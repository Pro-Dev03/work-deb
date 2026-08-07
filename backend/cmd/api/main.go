package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/database"
	"worktrack/backend/internal/router"
	"worktrack/backend/pkg/utils"
)

func main() {
	// Set memory limit to prevent OOM crashes
	debug.SetMemoryLimit(512 * 1024 * 1024) // 512MB

	cfg := config.Load()

	// Initialize logger
	utils.InitLoggerWithConfig(cfg.LogLevel, cfg.AppEnv)

	if err := cfg.ValidateProduction(); err != nil {
		utils.LogError("invalid production configuration", map[string]interface{}{
			"error": err,
		})
		log.Fatalf("❌ invalid production configuration: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		utils.LogError("failed to connect to database", map[string]interface{}{
			"error": err,
		})
		log.Fatalf("❌ %v", err)
	}
	defer db.Close()

	// تشغيل الترحيلات
	migrationsDir := "./internal/database/migrations"
	if err := database.Migrate(db, migrationsDir); err != nil {
		utils.LogWarning("failed to run migrations", map[string]interface{}{
			"error": err,
		})
		// لا نوقف الخادم إذا فشلت الترحيلات، قد تكون قاعدة البيانات موجودة بالفعل
	}

	r := router.Setup(db, cfg)

	// إنشاء مجلد uploads
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.MkdirAll(uploadsDir, 0755)
		utils.LogInfo("تم إنشاء مجلد uploads", nil)
	}

	// خدمة الملفات الثابتة
	r.Static("/uploads", "./uploads")

	// إنشاء HTTP server مع إعدادات مناسبة
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		// إعدادات timeout لمنع هجمات Slowloris
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("⚠️ Recovered from panic in server goroutine: %v", r)
			}
		}()

		utils.LogInfo("WorkTrack API يعمل", map[string]interface{}{
			"port": cfg.Port,
			"uploads_url": "http://localhost:" + cfg.Port + "/uploads/",
			"memory_limit": "512MB",
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogError("فشل تشغيل السيرفر", map[string]interface{}{
				"error": err,
			})
			log.Fatalf("❌ فشل تشغيل السيرفر: %v", err)
		}
	}()

	// انتظار إشارة الإنهاء
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.LogInfo("جاري إيقاف السيرفر بشكل آمن", nil)

	// إيقاف السيرفر مع timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError("فشل إيقاف السيرفر", map[string]interface{}{
			"error": err,
		})
	}

	utils.LogInfo("تم إيقاف السيرفر بنجاح", nil)
}
