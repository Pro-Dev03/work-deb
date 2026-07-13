package main

import (
	"log"
	"os"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/database"
	"worktrack/backend/internal/router"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateProduction(); err != nil {
		log.Fatalf("❌ invalid production configuration: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	defer db.Close()

	r := router.Setup(db, cfg)

	// إنشاء مجلد uploads
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.MkdirAll(uploadsDir, 0755)
		log.Println("📁 تم إنشاء مجلد uploads")
	}

	// خدمة الملفات الثابتة
	r.Static("/uploads", "./uploads")

	log.Printf("🚀 WorkTrack API يعمل على المنفذ %s", cfg.Port)
	log.Println("📁 الصور متاحة على: http://localhost:" + cfg.Port + "/uploads/")
	
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ فشل تشغيل السيرفر: %v", err)
	}
}
