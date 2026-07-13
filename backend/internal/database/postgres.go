package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Connect يفتح اتصالاً بقاعدة بيانات Postgres (Supabase أو أي مزوّد آخر)
func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL غير معرّف في متغيرات البيئة")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("فشل فتح الاتصال بقاعدة البيانات: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("فشل الاتصال بقاعدة البيانات: %w", err)
	}

	log.Println("✅ تم الاتصال بقاعدة البيانات بنجاح")
	return db, nil
}
