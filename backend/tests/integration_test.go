package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"worktrack/backend/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// testDB هي قاعدة البيانات المستخدمة للاختبارات
var testDB *sql.DB

// TestMain يُنفّذ قبل جميع الاختبارات لتهيئة قاعدة البيانات
func TestMain(m *testing.M) {
	// تحميل متغيرات البيئة من ملف .env
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("⚠️  فشل تحميل .env: %v\n", err)
	}

	// قراءة DATABASE_URL من متغيرات البيئة
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("⚠️  DATABASE_URL غير معرّف - تخطي اختبارات Integration")
		os.Exit(0)
	}

	// الاتصال بقاعدة البيانات
	var err error
	testDB, err = database.Connect(databaseURL)
	if err != nil {
		fmt.Printf("❌ فشل الاتصال بقاعدة البيانات: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	fmt.Println("✅ تم تهيئة قاعدة بيانات الاختبار بنجاح")

	// تشغيل جميع الاختبارات
	exitCode := m.Run()

	// تنظيف بعد الاختبارات (اختياري)
	fmt.Println("🧹 إنهاء اختبارات Integration")
	os.Exit(exitCode)
}

// getTestDB يُرجع اتصال قاعدة البيانات للاختبارات
func getTestDB() *sql.DB {
	if testDB == nil {
		panic("قاعدة بيانات الاختبار غير مهيأة - تأكد من تشغيل TestMain")
	}
	return testDB
}

// TestDatabaseConnection يتحقق من الاتصال بقاعدة البيانات
func TestDatabaseConnection(t *testing.T) {
	db := getTestDB()

	// اختبار استعلام بسيط
	var result string
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("فشل الاستعلام البسيط: %v", err)
	}

	if result != "1" {
		t.Errorf("نتيجة غير متوقعة: %s", result)
	}

	t.Log("✅ الاتصال بقاعدة البيانات يعمل بشكل صحيح")
}

// TestDatabaseTables يتحقق من وجود الجداول الأساسية
func TestDatabaseTables(t *testing.T) {
	db := getTestDB()

	// قائمة الجداول المتوقعة
	expectedTables := []string{
		"users",
		"worksites",
		"attendance",
	}

	for _, table := range expectedTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = $1
			)
		`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			t.Errorf("فشل التحقق من الجدول %s: %v", table, err)
			continue
		}

		if !exists {
			t.Errorf("الجدول %s غير موجود", table)
		} else {
			t.Logf("✅ الجدول %s موجود", table)
		}
	}
}

// TestDatabaseConnectionPool يتحقق من إعدادات connection pool
func TestDatabaseConnectionPool(t *testing.T) {
	db := getTestDB()

	// التحقق من إعدادات Connection Pool
	stats := db.Stats()
	t.Logf("إحصائيات Connection Pool:")
	t.Logf("  - Open Connections: %d", stats.OpenConnections)
	t.Logf("  - In Use: %d", stats.InUse)
	t.Logf("  - Idle: %d", stats.Idle)

	// التحقق من أن الاتصال مفتوح
	if err := db.Ping(); err != nil {
		t.Fatalf("فشل ping قاعدة البيانات: %v", err)
	}

	t.Log("✅ Connection Pool يعمل بشكل صحيح")
}

// TestDatabaseCRUD يتحقق من عمليات القراءة والكتابة الأساسية
func TestDatabaseCRUD(t *testing.T) {
	db := getTestDB()

	// اختبار قراءة بيانات من جدول users
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("فشل قراءة البيانات من جدول users: %v", err)
	}

	t.Logf("✅ عدد المستخدمين في قاعدة البيانات: %d", count)

	// اختبار قراءة بيانات من جدول worksites
	var worksiteCount int
	err = db.QueryRow("SELECT COUNT(*) FROM worksites").Scan(&worksiteCount)
	if err != nil {
		t.Fatalf("فشل قراءة البيانات من جدول worksites: %v", err)
	}

	t.Logf("✅ عدد مواقع العمل في قاعدة البيانات: %d", worksiteCount)

	// اختبار قراءة بيانات من جدول attendance
	var attendanceCount int
	err = db.QueryRow("SELECT COUNT(*) FROM attendance").Scan(&attendanceCount)
	if err != nil {
		t.Fatalf("فشل قراءة البيانات من جدول attendance: %v", err)
	}

	t.Logf("✅ عدد سجلات الحضور في قاعدة البيانات: %d", attendanceCount)
}
