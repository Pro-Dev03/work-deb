package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate يقوم بتشغيل الترحيلات من مجلد migrations
func Migrate(db *sql.DB, migrationsDir string) error {
	// إنشاء جدول الترحيلات إذا لم يكن موجوداً
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(14) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return fmt.Errorf("فشل إنشاء جدول الترحيلات: %w", err)
	}

	// جلب الترحيلات المطبقة بالفعل
	appliedMigrations := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("فشل جلب الترحيلات المطبقة: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("فشل قراءة إصدار الترحيل: %w", err)
		}
		appliedMigrations[version] = true
	}

	// قراءة ملفات الترحيل
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("فشل قراءة مجلد الترحيلات: %w", err)
	}

	// ترتيب الملفات حسب الرقم
	var upMigrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	sort.Strings(upMigrations)

	// تشغيل الترحيلات
	for _, filename := range upMigrations {
		version := strings.TrimSuffix(filename, ".up.sql")
		if appliedMigrations[version] {
			log.Printf("⏭️  الترحيل %s مطبق بالفعل", version)
			continue
		}

		log.Printf("🔄 تشغيل الترحيل: %s", version)

		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("فشل قراءة ملف الترحيل %s: %w", filename, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("فشل بدء المعاملة: %w", err)
		}

		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("فشل تنفيذ الترحيل %s: %w", version, err)
		}

		_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("فشل تسجيل الترحيل %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("فشل إتمام المعاملة: %w", err)
		}

		log.Printf("✅ تم تطبيق الترحيل: %s", version)
	}

	log.Println("✅ تم الانتهاء من جميع الترحيلات")
	return nil
}
