package database

import (
	"database/sql"
	"fmt"
	"time"

	"worktrack/backend/pkg/utils"

	_ "github.com/lib/pq"
)

// Connect يفتح اتصالاً بقاعدة بيانات Postgres (Supabase أو أي مزوّد آخر)
func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL غير معرّف في متغيرات البيئة")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		utils.LogError("فشل فتح الاتصال بقاعدة البيانات", map[string]interface{}{
			"error": err,
		})
		return nil, fmt.Errorf("فشل فتح الاتصال بقاعدة البيانات: %w", err)
	}

	// إعدادات Connection Pooling المحسّنة
	// هذه الإعدادات مُحسّنة لـ Supabase PostgreSQL
	db.SetMaxOpenConns(20)        // الحد الأقصى للاتصالات المفتوحة (Supabase يفضل 20-30)
	db.SetMaxIdleConns(10)        // الحد الأقصى للاتصالات الخاملة (زيادة لتحسين الأداء)
	db.SetConnMaxLifetime(30 * time.Minute) // أقصى عمر للاتصال (30 دقيقة أفضل لإعادة الاستخدام)
	db.SetConnMaxIdleTime(5 * time.Minute)  // أقصى وقت للاتصال الخامل (5 دقائق لمنع الاتصالات القديمة)

	if err := db.Ping(); err != nil {
		utils.LogError("فشل الاتصال بقاعدة البيانات", map[string]interface{}{
			"error": err,
		})
		return nil, fmt.Errorf("فشل الاتصال بقاعدة البيانات: %w", err)
	}

	utils.LogInfo("تم الاتصال بقاعدة البيانات بنجاح", map[string]interface{}{
		"max_open_conns": 20,
		"max_idle_conns": 10,
		"max_lifetime":   "30m",
		"max_idle_time":  "5m",
	})

	return db, nil
}
