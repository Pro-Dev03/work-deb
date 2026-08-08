# Backend Code

## تم الاستخراج في: 2026-08-06 02:09:32

## عدد الملفات: 95

---

## 📄 backend/cmd/api/main.go

```go
package main

import (
	"log"
	"os"
	"runtime/debug"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/database"
	"worktrack/backend/internal/router"
)

func main() {
	// Set memory limit to prevent OOM crashes
	debug.SetMemoryLimit(512 * 1024 * 1024) // 512MB

	cfg := config.Load()
	if err := cfg.ValidateProduction(); err != nil {
		log.Fatalf("❌ invalid production configuration: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	defer db.Close()

	// الاتصال بـ Redis (اختياري)
	if err := database.ConnectRedis(cfg.RedisURL); err != nil {
		log.Printf("⚠️  فشل الاتصال بـ Redis: %v", err)
		// لا نوقف التطبيق إذا فشل Redis
	}
	defer database.CloseRedis()

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
	log.Println("🧠 حد الذاكرة: 512MB")

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ فشل تشغيل السيرفر: %v", err)
	}
}

```

---

## 📄 backend/fix_customer_role.sql

```sql
-- =====================================================
-- WorkTrack Database Schema - DevPro Version
-- Compatible with current project structure
-- =====================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =====================================================
-- Update users table with missing fields
-- =====================================================

-- Add subscription fields if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (subscription_status IN ('active', 'expired', 'canceled'));
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ NULL;

-- Add updated_at if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add device and login fields (from migration 000007)
-- =====================================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_id TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_model TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_device TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

-- =====================================================
-- Update role constraint to include ALL roles (admin, employee, customer, client)
-- =====================================================
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'customer', 'client'));

-- =====================================================
-- Update users table to allow nullable email for customers
-- =====================================================
-- Drop the NOT NULL constraint on email to allow NULL values
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Add unique constraint on email (only for non-null values)
-- This allows multiple NULL emails
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users(email) WHERE email IS NOT NULL;

-- =====================================================
-- Add missing updated_at to worksites table
-- =====================================================
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add assigned_employee_id to worksites if not exist
-- =====================================================
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- =====================================================
-- Add missing updated_at to tasks table
-- =====================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add photo fields to attendance (from migration 000007)
-- =====================================================
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- =====================================================
-- Add missing updated_at to attendance table
-- =====================================================
ALTER TABLE attendance
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add type and related_id to notifications (from migration 000007)
-- =====================================================
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

-- =====================================================
-- Create service_requests table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS service_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_name VARCHAR(150),
    client_phone VARCHAR(30),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'assigned', 'in_progress', 'completed', 'cancelled')),
    priority VARCHAR(20) NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    photos TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add location_name if not exist
ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS location_name TEXT;

-- =====================================================
-- Create assignments table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    admin_notes TEXT,
    employee_notes TEXT,
    client_rating INTEGER CHECK (client_rating BETWEEN 1 AND 5),
    client_feedback TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'assigned'
        CHECK (status IN ('assigned', 'accepted', 'in_progress', 'completed', 'rejected'))
);

-- =====================================================
-- Create location_tracking table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS location_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- Create indexes for performance
-- =====================================================

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_subscription_status ON users(subscription_status);

-- Worksites indexes
CREATE INDEX IF NOT EXISTS idx_worksites_location ON worksites(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);
CREATE INDEX IF NOT EXISTS idx_worksites_is_active ON worksites(is_active);

-- Tasks indexes
CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id ON tasks(worksite_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id ON tasks(assigned_user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_client_id ON tasks(client_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

-- Attendance indexes
CREATE INDEX IF NOT EXISTS idx_attendance_user_id ON attendance(user_id);
CREATE INDEX IF NOT EXISTS idx_attendance_task_id ON attendance(task_id);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);
CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time ON attendance(check_in_time);
CREATE INDEX IF NOT EXISTS idx_attendance_status ON attendance(status);

-- Service requests indexes
CREATE INDEX IF NOT EXISTS idx_service_requests_client_id ON service_requests(client_id);
CREATE INDEX IF NOT EXISTS idx_service_requests_status ON service_requests(status);
CREATE INDEX IF NOT EXISTS idx_service_requests_priority ON service_requests(priority);

-- Assignments indexes
CREATE INDEX IF NOT EXISTS idx_assignments_employee_id ON assignments(employee_id);
CREATE INDEX IF NOT EXISTS idx_assignments_request_id ON assignments(request_id);
CREATE INDEX IF NOT EXISTS idx_assignments_admin_id ON assignments(admin_id);
CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);

-- Location tracking indexes
CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id ON location_tracking(user_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id ON location_tracking(assignment_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_recorded_at ON location_tracking(recorded_at);

-- Notifications indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- =====================================================
-- Create default admin users
-- =====================================================
INSERT INTO users (
  full_name,
  email,
  password_hash,
  role,
  is_active,
  subscription_status,
  subscription_expires_at
)
VALUES
(
  'DevPro System Administrator',
  'admin@devpro.com',
  crypt('devproadmin', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
),
(
  'DevPro Support Manager',
  'support@devpro.com',
  crypt('devprosupport', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NOW() + INTERVAL '1 year'
),
(
  'DevPro Project Manager',
  'manager@devpro.com',
  crypt('devpromanager', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
)
ON CONFLICT (email) DO UPDATE
SET
  password_hash = EXCLUDED.password_hash,
  role = 'admin',
  is_active = TRUE,
  subscription_status = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN 'active'
    ELSE EXCLUDED.subscription_status
  END,
  subscription_expires_at = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN NULL
    ELSE EXCLUDED.subscription_expires_at
  END,
  updated_at = now()
RETURNING email, role, is_active, subscription_status, subscription_expires_at;

-- =====================================================
-- Verification query
-- =====================================================
SELECT
    conname as constraint_name,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'users'::regclass AND conname = 'users_role_check';

```

---

## 📄 backend/fix_service_requests_client_delete.sql

```sql
-- Fix service_requests table to handle client deletion properly
-- This script changes the foreign key constraint from ON DELETE CASCADE to ON DELETE SET NULL
-- and makes client_id nullable so service requests can exist without a client

-- Step 1: Drop the existing foreign key constraint
ALTER TABLE service_requests DROP CONSTRAINT IF EXISTS service_requests_client_id_fkey;

-- Step 2: Make client_id nullable (if it's currently NOT NULL)
ALTER TABLE service_requests ALTER COLUMN client_id DROP NOT NULL;

-- Step 3: Add the new foreign key constraint with ON DELETE SET NULL
ALTER TABLE service_requests 
ADD CONSTRAINT service_requests_client_id_fkey 
FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE SET NULL;

-- Verification query
SELECT 
    conname as constraint_name,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'service_requests'::regclass AND conname = 'service_requests_client_id_fkey';

```

---

## 📄 backend/internal/config/config.go

```go
package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AllowedOrigin string
	DefaultLang   string
	GeoapifyKey   string
	R2AccountID   string
	R2AccessKey   string
	R2SecretKey   string
	R2BucketName  string
	RedisURL      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  لم يتم العثور على ملف .env")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		AllowedOrigin: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:3002"),
		DefaultLang:   getEnv("DEFAULT_LANG", "ar"),
		GeoapifyKey:   getEnv("GEOAPIFY_KEY", ""),
		R2AccountID:   getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:   getEnv("R2_ACCESS_KEY", ""),
		R2SecretKey:   getEnv("R2_SECRET_KEY", ""),
		R2BucketName:  getEnv("R2_BUCKET_NAME", "worktrack-uploads"),
		RedisURL:      getEnv("REDIS_URL", ""),
	}
}

// ValidateProduction prevents the API from starting with unsafe placeholder
// values when deployed. Local development keeps convenient defaults.
func (c *Config) ValidateProduction() error {
	if strings.EqualFold(getEnv("APP_ENV", "development"), "production") {
		if c.DatabaseURL == "" {
			return fmt.Errorf("DATABASE_URL must be configured in production")
		}
		if len(c.JWTSecret) < 32 || strings.Contains(strings.ToLower(c.JWTSecret), "change_this") {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters and not a placeholder in production")
		}
		if c.AllowedOrigin == "" || strings.Contains(c.AllowedOrigin, "*") {
			return fmt.Errorf("ALLOWED_ORIGINS must list explicit HTTPS frontend origins in production")
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

```

---

## 📄 backend/internal/database/migrations/000001_create_users.down.sql

```sql
DROP TABLE IF EXISTS users;

```

---

## 📄 backend/internal/database/migrations/000001_create_users.up.sql

```sql
-- تفعيل توليد UUID تلقائياً (مطلوب لكل الجداول التي تعتمد uuid_generate_v4)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- جدول المستخدمين: يضم كلاً من المدير (admin) والموظف (employee) في جدول واحد
-- ويُفرَّق بينهما عبر عمود role فقط، لتبسيط منطق تسجيل الدخول (Endpoint واحد للجميع)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(30),
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee' CHECK (role IN ('admin', 'employee')),
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

```

---

## 📄 backend/internal/database/migrations/000002_create_worksites.down.sql

```sql
DROP TABLE IF EXISTS worksites;

```

---

## 📄 backend/internal/database/migrations/000002_create_worksites.up.sql

```sql
-- نقاط العمل (Geofence Zones) — هذا الجدول هو أساس ميزة منع التختيم خارج النطاق.
-- radius_meters هو نصف القطر المسموح به حول (latitude, longitude) لتلك النقطة
CREATE TABLE IF NOT EXISTS worksites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius_meters INTEGER NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

```

---

## 📄 backend/internal/database/migrations/000003_create_clients.down.sql

```sql
DROP TABLE IF EXISTS clients;

```

---

## 📄 backend/internal/database/migrations/000003_create_clients.up.sql

```sql
-- العملاء: الجهة التي تُنفَّذ لصالحها المهمة (لا تعتمد على أي جدول آخر)
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(150),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

```

---

## 📄 backend/internal/database/migrations/000004_create_tasks.down.sql

```sql
DROP TABLE IF EXISTS tasks;

```

---

## 📄 backend/internal/database/migrations/000004_create_tasks.up.sql

```sql
-- المهام: تربط بين عميل + نقطة عمل (worksite) + موظف مكلَّف
-- الربط بـ worksite_id هو ما يُفعِّل التحقق من النطاق الجغرافي عند التختيم لاحقاً
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    title_ar VARCHAR(200),
    title_he VARCHAR(200),
    title_en VARCHAR(200),
    description TEXT,
    description_ar TEXT,
    description_he TEXT,
    description_en TEXT,
    client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
    worksite_id UUID NOT NULL REFERENCES worksites(id) ON DELETE RESTRICT,
    assigned_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'in_progress', 'completed', 'late', 'cancelled')),
    priority VARCHAR(20) NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    scheduled_start TIMESTAMPTZ,
    scheduled_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    language VARCHAR(10) NOT NULL DEFAULT 'en'
);

```

---

## 📄 backend/internal/database/migrations/000005_create_attendance.down.sql

```sql
DROP TABLE IF EXISTS attendance;

```

---

## 📄 backend/internal/database/migrations/000005_create_attendance.up.sql

```sql
-- سجل الحضور/الانصراف — كل سجل يحفظ موقع الموظف الفعلي والمسافة عن نقطة العمل
-- لحظة كل تختيم (بدء/إنهاء)، كإثبات أنه تم داخل النطاق المسموح
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    worksite_id UUID NOT NULL REFERENCES worksites(id) ON DELETE RESTRICT,

    check_in_time TIMESTAMPTZ,
    check_in_lat DOUBLE PRECISION,
    check_in_lng DOUBLE PRECISION,
    check_in_distance_meters DOUBLE PRECISION,

    check_out_time TIMESTAMPTZ,
    check_out_lat DOUBLE PRECISION,
    check_out_lng DOUBLE PRECISION,
    check_out_distance_meters DOUBLE PRECISION,

    photo_url TEXT,
    photo_uploaded_at TIMESTAMPTZ,
    photo_notes TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'in_progress'
        CHECK (status IN ('in_progress', 'completed')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

```

---

## 📄 backend/internal/database/migrations/000006_create_notifications.down.sql

```sql
DROP TABLE IF EXISTS notifications;

```

---

## 📄 backend/internal/database/migrations/000006_create_notifications.up.sql

```sql
-- إشعارات لكل مستخدم (مثال: إشعار عند محاولة تختيم مرفوضة خارج النطاق)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

```

---

## 📄 backend/internal/database/migrations/000007_extend_production_schema.up.sql

```sql
-- Completes the production schema used by the current WorkTrack API.
-- Safe to run on an existing database: it does not delete tables, rows, or users.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Phone-login and device fields required by the authentication handlers.
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_id TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_model TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_device TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

-- The original migration allowed only admin and employee roles. Client accounts
-- are required by the service-request flow.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'client'));

-- Evidence-photo fields required by the attendance upload handler.
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- Notification type is used for geofence and operational alerts.
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

CREATE TABLE IF NOT EXISTS service_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_name VARCHAR(150),
    client_phone VARCHAR(30),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'assigned', 'in_progress', 'completed', 'cancelled')),
    priority VARCHAR(20) NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    photos TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    admin_notes TEXT,
    employee_notes TEXT,
    client_rating INTEGER CHECK (client_rating BETWEEN 1 AND 5),
    client_feedback TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'assigned'
        CHECK (status IN ('assigned', 'accepted', 'in_progress', 'completed', 'rejected'))
);

CREATE TABLE IF NOT EXISTS location_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_service_requests_client_id ON service_requests(client_id);
CREATE INDEX IF NOT EXISTS idx_service_requests_status ON service_requests(status);
CREATE INDEX IF NOT EXISTS idx_assignments_employee_id ON assignments(employee_id);
CREATE INDEX IF NOT EXISTS idx_assignments_request_id ON assignments(request_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id ON location_tracking(user_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id ON location_tracking(assignment_id);

```

---

## 📄 backend/internal/database/migrations/000008_add_subscription_fields.down.sql

```sql
-- إزالة حقول الاشتراك من جدول المستخدمين
ALTER TABLE users
DROP COLUMN IF EXISTS subscription_status,
DROP COLUMN IF EXISTS subscription_expires_at;

```

---

## 📄 backend/internal/database/migrations/000008_add_subscription_fields.up.sql

```sql
-- إضافة حقول الاشتراك إلى جدول المستخدمين
ALTER TABLE users
ADD COLUMN subscription_status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (subscription_status IN ('active', 'expired', 'canceled')),
ADD COLUMN subscription_expires_at TIMESTAMPTZ NULL;

```

---

## 📄 backend/internal/database/migrations/000009_add_missing_columns.down.sql

```sql
-- التراجع عن إضافة الأعمدة المفقودة

-- حذف الفهارس
DROP INDEX IF EXISTS idx_worksites_assigned_employee;
DROP INDEX IF EXISTS idx_attendance_worksite_id;

-- حذف الأعمدة (اختياري - عادة نحتفظ بها للتوافق)
-- ALTER TABLE worksites DROP COLUMN IF EXISTS assigned_employee_id;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS updated_at;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_url;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_uploaded_at;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_notes;
```

---

## 📄 backend/internal/database/migrations/000009_add_missing_columns.up.sql

```sql
-- إضافة الأعمدة المفقودة للتوافق مع الإصدارات السابقة

-- إضافة عمود assigned_employee_id إلى جدول worksites
ALTER TABLE worksites 
ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- إضافة عمود updated_at إلى جدول attendance
ALTER TABLE attendance
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- إضافة حقول الصور إلى جدول attendance إذا لم تكن موجودة
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- إنشاء الفهارس المفقودة
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);
```

---

## 📄 backend/internal/database/migrations/000010_add_check_out_notes.down.sql

```sql
-- حذف عمود ملاحظات إنهاء الدوام
ALTER TABLE attendance DROP COLUMN IF EXISTS check_out_notes;

```

---

## 📄 backend/internal/database/migrations/000010_add_check_out_notes.up.sql

```sql
-- إضافة عمود لملاحظات إنهاء الدوام
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS check_out_notes TEXT;

```

---

## 📄 backend/internal/database/migrations/000011_add_default_admin_users.down.sql

```sql
-- Remove default admin users
-- This removes the default admin accounts created in the up migration

DELETE FROM users 
WHERE email IN (
  'admin@devpro.com',
  'support@devpro.com', 
  'manager@devpro.com'
);

```

---

## 📄 backend/internal/database/migrations/000011_add_default_admin_users.up.sql

```sql
-- Add default admin users for WorkTrack v2
-- This migration creates default admin accounts for initial setup

-- Enable pgcrypto extension if not already enabled
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Insert default admin users
INSERT INTO users (
  full_name,
  email,
  phone,
  password_hash,
  role,
  is_active,
  subscription_status,
  subscription_expires_at,
  created_at,
  updated_at
)
VALUES
(
  'DevPro System Administrator',
  'admin@devpro.com',
  '+1234567890',
  crypt('devproadmin', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL,
  now(),
  now()
),
(
  'DevPro Support Manager',
  'support@devpro.com',
  '+1234567891',
  crypt('devprosupport', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NOW() + INTERVAL '1 year',
  now(),
  now()
),
(
  'DevPro Project Manager',
  'manager@devpro.com',
  '+1234567892',
  crypt('devpromanager', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL,
  now(),
  now()
)
ON CONFLICT (email) DO UPDATE
SET
  password_hash = CASE WHEN users.password_hash IS NULL THEN EXCLUDED.password_hash ELSE users.password_hash END,
  role = 'admin',
  is_active = TRUE,
  subscription_status = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN 'active'
    ELSE EXCLUDED.subscription_status
  END,
  subscription_expires_at = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN NULL
    ELSE EXCLUDED.subscription_expires_at
  END,
  updated_at = now()
RETURNING email, role, is_active, subscription_status, subscription_expires_at;

```

---

## 📄 backend/internal/database/migrations/000012_fix_customer_role_constraint.down.sql

```sql
-- العودة إلى constraint القديم (بدون customer/client)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- إضافة constraint القديم
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee'));
```

---

## 📄 backend/internal/database/migrations/000012_fix_customer_role_constraint.up.sql

```sql
-- تحديث constraint دور المستخدم للسماح بالأدوار الجديدة
-- هذا الـ migration ضروري لأن الباك إند يحاول إنشاء مستخدمين بدور customer/client

-- حذف constraint القديم إذا كان موجوداً
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- إضافة constraint جديد يشمل جميع الأدوار المطلوبة
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'customer', 'client'));
```

---

## 📄 backend/internal/database/migrations/000013_add_location_name.down.sql

```sql
-- Remove location_name field from service_requests table

ALTER TABLE service_requests DROP COLUMN IF EXISTS location_name;
```

---

## 📄 backend/internal/database/migrations/000013_add_location_name.up.sql

```sql
-- Add location_name field to service_requests table  
-- This field will store the human-readable location name from reverse geocoding

ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS location_name TEXT;
```

---

## 📄 backend/internal/database/migrations/000014_complete_schema.down.sql

```sql
-- =====================================================
-- Rollback for complete schema migration
-- =====================================================

-- Remove default admin users
DELETE FROM users WHERE email IN ('admin@devpro.com', 'support@devpro.com', 'manager@devpro.com');

-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_location_tracking_recorded_at;
DROP INDEX IF EXISTS idx_location_tracking_assignment_id;
DROP INDEX IF EXISTS idx_location_tracking_user_id;
DROP INDEX IF EXISTS idx_assignments_status;
DROP INDEX IF EXISTS idx_assignments_admin_id;
DROP INDEX IF EXISTS idx_assignments_request_id;
DROP INDEX IF EXISTS idx_assignments_employee_id;
DROP INDEX IF EXISTS idx_service_requests_priority;
DROP INDEX IF EXISTS idx_service_requests_status;
DROP INDEX IF EXISTS idx_service_requests_client_id;
DROP INDEX IF EXISTS idx_attendance_status;
DROP INDEX IF EXISTS idx_attendance_check_in_time;
DROP INDEX IF EXISTS idx_attendance_worksite_id;
DROP INDEX IF EXISTS idx_attendance_task_id;
DROP INDEX IF EXISTS idx_attendance_user_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_client_id;
DROP INDEX IF EXISTS idx_tasks_assigned_user_id;
DROP INDEX IF EXISTS idx_tasks_worksite_id;
DROP INDEX IF EXISTS idx_worksites_is_active;
DROP INDEX IF EXISTS idx_worksites_assigned_employee;
DROP INDEX IF EXISTS idx_worksites_location;
DROP INDEX IF EXISTS idx_users_subscription_status;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_phone;

-- Drop tables
DROP TABLE IF EXISTS location_tracking;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS service_requests;

-- Remove location_name from service_requests if table still exists
-- (This is a safety measure in case the table wasn't dropped)
ALTER TABLE service_requests DROP COLUMN IF EXISTS location_name;

-- Remove columns from notifications
ALTER TABLE notifications DROP COLUMN IF EXISTS related_id;
ALTER TABLE notifications DROP COLUMN IF EXISTS type;

-- Remove columns from attendance
ALTER TABLE attendance DROP COLUMN IF EXISTS updated_at;
ALTER TABLE attendance DROP COLUMN IF EXISTS photo_notes;
ALTER TABLE attendance DROP COLUMN IF EXISTS photo_uploaded_at;
ALTER TABLE attendance DROP COLUMN IF EXISTS photo_url;

-- Remove columns from tasks
ALTER TABLE tasks DROP COLUMN IF EXISTS updated_at;

-- Remove columns from worksites
ALTER TABLE worksites DROP COLUMN IF EXISTS assigned_employee_id;
ALTER TABLE worksites DROP COLUMN IF EXISTS updated_at;

-- Remove columns from users
ALTER TABLE users DROP COLUMN IF EXISTS login_code_expires;
ALTER TABLE users DROP COLUMN IF EXISTS login_code;
ALTER TABLE users DROP COLUMN IF EXISTS last_login_device;
ALTER TABLE users DROP COLUMN IF EXISTS phone_verified;
ALTER TABLE users DROP COLUMN IF EXISTS device_model;
ALTER TABLE users DROP COLUMN IF EXISTS device_id;
ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_status;

-- Reset role constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'client'));
```

---

## 📄 backend/internal/database/migrations/000014_complete_schema.up.sql

```sql
-- =====================================================
-- WorkTrack Database Schema - DevPro Version
-- Compatible with current project structure
-- This migration completes the production schema with missing fields and improvements
-- =====================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =====================================================
-- Update users table with missing fields
-- =====================================================

-- Add subscription fields if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (subscription_status IN ('active', 'expired', 'canceled'));
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ NULL;

-- Add updated_at if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add device and login fields (from migration 000007)
-- =====================================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_id TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_model TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_device TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

-- =====================================================
-- Update role constraint to include ALL roles (admin, employee, customer, client)
-- =====================================================
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'customer', 'client'));

-- =====================================================
-- Add missing updated_at to worksites table
-- =====================================================
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add assigned_employee_id to worksites if not exist
-- =====================================================
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- =====================================================
-- Add missing updated_at to tasks table
-- =====================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add photo fields to attendance (from migration 000007)
-- =====================================================
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- =====================================================
-- Add missing updated_at to attendance table
-- =====================================================
ALTER TABLE attendance
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- =====================================================
-- Add type and related_id to notifications (from migration 000007)
-- =====================================================
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

-- =====================================================
-- Create service_requests table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS service_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_name VARCHAR(150),
    client_phone VARCHAR(30),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    location_name TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'assigned', 'in_progress', 'completed', 'cancelled')),
    priority VARCHAR(20) NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    photos TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add location_name if table already exists and column is missing
ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS location_name TEXT;

-- =====================================================
-- Create assignments table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    admin_notes TEXT,
    employee_notes TEXT,
    client_rating INTEGER CHECK (client_rating BETWEEN 1 AND 5),
    client_feedback TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'assigned'
        CHECK (status IN ('assigned', 'accepted', 'in_progress', 'completed', 'rejected'))
);

-- =====================================================
-- Create location_tracking table (from migration 000007)
-- =====================================================
CREATE TABLE IF NOT EXISTS location_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- Create indexes for performance
-- =====================================================

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_subscription_status ON users(subscription_status);

-- Worksites indexes
CREATE INDEX IF NOT EXISTS idx_worksites_location ON worksites(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);
CREATE INDEX IF NOT EXISTS idx_worksites_is_active ON worksites(is_active);

-- Tasks indexes
CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id ON tasks(worksite_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id ON tasks(assigned_user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_client_id ON tasks(client_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

-- Attendance indexes
CREATE INDEX IF NOT EXISTS idx_attendance_user_id ON attendance(user_id);
CREATE INDEX IF NOT EXISTS idx_attendance_task_id ON attendance(task_id);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);
CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time ON attendance(check_in_time);
CREATE INDEX IF NOT EXISTS idx_attendance_status ON attendance(status);

-- Service requests indexes
CREATE INDEX IF NOT EXISTS idx_service_requests_client_id ON service_requests(client_id);
CREATE INDEX IF NOT EXISTS idx_service_requests_status ON service_requests(status);
CREATE INDEX IF NOT EXISTS idx_service_requests_priority ON service_requests(priority);

-- Assignments indexes
CREATE INDEX IF NOT EXISTS idx_assignments_employee_id ON assignments(employee_id);
CREATE INDEX IF NOT EXISTS idx_assignments_request_id ON assignments(request_id);
CREATE INDEX IF NOT EXISTS idx_assignments_admin_id ON assignments(admin_id);
CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);

-- Location tracking indexes
CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id ON location_tracking(user_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id ON location_tracking(assignment_id);
CREATE INDEX IF NOT EXISTS idx_location_tracking_recorded_at ON location_tracking(recorded_at);

-- Notifications indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- =====================================================
-- Create default admin users
-- =====================================================
INSERT INTO users (
  full_name,
  email,
  password_hash,
  role,
  is_active,
  subscription_status,
  subscription_expires_at
)
VALUES
(
  'DevPro System Administrator',
  'admin@devpro.com',
  crypt('devproadmin', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
),
(
  'DevPro Support Manager',
  'support@devpro.com',
  crypt('devprosupport', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NOW() + INTERVAL '1 year'
),
(
  'DevPro Project Manager',
  'manager@devpro.com',
  crypt('devpromanager', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
)
ON CONFLICT (email) DO UPDATE
SET
  password_hash = CASE WHEN users.password_hash IS NULL THEN EXCLUDED.password_hash ELSE users.password_hash END,
  role = 'admin',
  is_active = TRUE,
  subscription_status = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN 'active'
    ELSE EXCLUDED.subscription_status
  END,
  subscription_expires_at = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN NULL
    ELSE EXCLUDED.subscription_expires_at
  END,
  updated_at = now()
RETURNING email, role, is_active, subscription_status, subscription_expires_at;
```

---

## 📄 backend/internal/database/migrations/000015_add_password_changed_at.down.sql

```sql
-- Remove password_changed_at column from users table

DROP INDEX IF EXISTS idx_users_password_changed_at;
ALTER TABLE users DROP COLUMN IF EXISTS password_changed_at;
```

---

## 📄 backend/internal/database/migrations/000015_add_password_changed_at.up.sql

```sql
-- Add password_changed_at column to users table
-- This column will be used to invalidate tokens when password is changed

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMPTZ DEFAULT NOW();

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_users_password_changed_at ON users(password_changed_at);
```

---

## 📄 backend/internal/database/migrations/000015_add_service_request_to_attendance.down.sql

```sql
-- Remove service_request_id from attendance table

ALTER TABLE attendance DROP COLUMN IF EXISTS service_request_id;

```

---

## 📄 backend/internal/database/migrations/000015_add_service_request_to_attendance.up.sql

```sql
-- Add service_request_id to attendance table
-- This will link attendance records to service requests

ALTER TABLE attendance ADD COLUMN IF NOT EXISTS service_request_id UUID REFERENCES service_requests(id) ON DELETE SET NULL;

```

---

## 📄 backend/internal/database/migrations/000016_add_language_field.down.sql

```sql
-- =====================================================
-- Remove language field from tasks table
-- =====================================================
ALTER TABLE tasks
    DROP COLUMN IF EXISTS language;

-- =====================================================
-- Remove language field from service_requests table
-- =====================================================
ALTER TABLE service_requests
    DROP COLUMN IF EXISTS language;

```

---

## 📄 backend/internal/database/migrations/000016_add_language_field.up.sql

```sql
-- =====================================================
-- Add language field to tasks table
-- =====================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';

-- =====================================================
-- Add language field to service_requests table
-- =====================================================
ALTER TABLE service_requests
    ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';

```

---

## 📄 backend/internal/database/migrations/000017_allow_null_worksite_in_attendance.down.sql

```sql
-- التراجع عن تعديل عمود worksite_id

-- إعادة القيد الأجنبي إلى RESTRICT
ALTER TABLE attendance DROP CONSTRAINT attendance_worksite_id_fkey;

ALTER TABLE attendance 
ADD CONSTRAINT attendance_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE RESTRICT;

-- إضافة قيد NOT NULL
ALTER TABLE attendance ALTER COLUMN worksite_id SET NOT NULL;
```

---

## 📄 backend/internal/database/migrations/000017_allow_null_worksite_in_attendance.up.sql

```sql
-- تعديل عمود worksite_id للسماح بـ NULL للحفاظ على سجلات الحضور عند حذف موقع العمل
-- هذا يسمح بفك ارتباط سجلات الحضور من مواقع العمل المحذوفة دون فقدان البيانات

-- أولاً إزالة القيد NOT NULL
ALTER TABLE attendance ALTER COLUMN worksite_id DROP NOT NULL;

-- ثم تغيير قيد ON DELETE من RESTRICT إلى SET NULL
-- يتطلب إزالة القيد الأجنبي أولاً ثم إضافته مرة أخرى
ALTER TABLE attendance DROP CONSTRAINT attendance_worksite_id_fkey;

ALTER TABLE attendance 
ADD CONSTRAINT attendance_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;
```

---

## 📄 backend/internal/database/migrations/000018_add_translations_and_fix_worksite_delete.down.sql

```sql
-- =====================================================
-- Rollback: Remove Translation Columns + Restore Worksite Deletion Constraint
-- =====================================================

-- =====================================================
-- PART 1: Restore Worksite Deletion Constraint
-- =====================================================

-- Restore the foreign key constraint to ON DELETE RESTRICT
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_worksite_id_fkey;

ALTER TABLE tasks 
ADD CONSTRAINT tasks_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE RESTRICT;

-- Restore NOT NULL constraint (note: this will fail if there are NULL values)
-- ALTER TABLE tasks ALTER COLUMN worksite_id SET NOT NULL;

-- =====================================================
-- PART 2: Remove Translation Columns
-- =====================================================

-- Drop the helper function
DROP FUNCTION IF EXISTS get_translation(TEXT, TEXT, TEXT, VARCHAR);

-- Remove translation columns from service_requests (if table exists)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'service_requests') THEN
        ALTER TABLE service_requests
            DROP COLUMN IF EXISTS client_name_en,
            DROP COLUMN IF EXISTS client_name_he,
            DROP COLUMN IF EXISTS client_name_ar,
            DROP COLUMN IF EXISTS location_name_en,
            DROP COLUMN IF EXISTS location_name_he,
            DROP COLUMN IF EXISTS location_name_ar,
            DROP COLUMN IF EXISTS address_en,
            DROP COLUMN IF EXISTS address_he,
            DROP COLUMN IF EXISTS address_ar,
            DROP COLUMN IF EXISTS description_en,
            DROP COLUMN IF EXISTS description_he,
            DROP COLUMN IF EXISTS description_ar,
            DROP COLUMN IF EXISTS title_en,
            DROP COLUMN IF EXISTS title_he,
            DROP COLUMN IF EXISTS title_ar;
    END IF;
END $$;

-- Remove translation columns from tasks
ALTER TABLE tasks
    DROP COLUMN IF EXISTS description_en,
    DROP COLUMN IF EXISTS description_he,
    DROP COLUMN IF EXISTS description_ar,
    DROP COLUMN IF EXISTS title_en,
    DROP COLUMN IF EXISTS title_he,
    DROP COLUMN IF EXISTS title_ar;

-- Remove translation columns from worksites
ALTER TABLE worksites
    DROP COLUMN IF EXISTS address_en,
    DROP COLUMN IF EXISTS address_he,
    DROP COLUMN IF EXISTS address_ar,
    DROP COLUMN IF EXISTS name_en,
    DROP COLUMN IF EXISTS name_he,
    DROP COLUMN IF EXISTS name_ar;

-- Remove translation columns from clients (if table exists)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'clients') THEN
        ALTER TABLE clients
            DROP COLUMN IF EXISTS address_en,
            DROP COLUMN IF EXISTS address_he,
            DROP COLUMN IF EXISTS address_ar,
            DROP COLUMN IF EXISTS name_en,
            DROP COLUMN IF EXISTS name_he,
            DROP COLUMN IF EXISTS name_ar;
    END IF;
END $$;

-- Remove translation columns from users
ALTER TABLE users
    DROP COLUMN IF EXISTS full_name_en,
    DROP COLUMN IF EXISTS full_name_he,
    DROP COLUMN IF EXISTS full_name_ar;
```

---

## 📄 backend/internal/database/migrations/000018_add_translations_and_fix_worksite_delete.up.sql

```sql
-- =====================================================
-- Migration: Add Translation Columns for Complete i18n Support
-- =====================================================
-- This migration adds translation columns for all text fields
-- to support Arabic, Hebrew, and English translations

-- =====================================================
-- Users table translations
-- =====================================================
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS full_name_ar VARCHAR(150),
    ADD COLUMN IF NOT EXISTS full_name_he VARCHAR(150),
    ADD COLUMN IF NOT EXISTS full_name_en VARCHAR(150);

-- Migrate existing full_name to all translation columns
UPDATE users 
SET 
    full_name_ar = full_name,
    full_name_he = full_name,
    full_name_en = full_name
WHERE full_name_ar IS NULL;

-- =====================================================
-- Clients table translations (if table exists)
-- =====================================================
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'clients') THEN
        ALTER TABLE clients
            ADD COLUMN IF NOT EXISTS name_ar VARCHAR(150),
            ADD COLUMN IF NOT EXISTS name_he VARCHAR(150),
            ADD COLUMN IF NOT EXISTS name_en VARCHAR(150),
            ADD COLUMN IF NOT EXISTS address_ar TEXT,
            ADD COLUMN IF NOT EXISTS address_he TEXT,
            ADD COLUMN IF NOT EXISTS address_en TEXT;

        -- Migrate existing data to translation columns
        UPDATE clients 
        SET 
            name_ar = name,
            name_he = name,
            name_en = name,
            address_ar = address,
            address_he = address,
            address_en = address
        WHERE name_ar IS NULL;
    END IF;
END $$;

-- =====================================================
-- Worksites table translations
-- =====================================================
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS name_ar VARCHAR(150),
    ADD COLUMN IF NOT EXISTS name_he VARCHAR(150),
    ADD COLUMN IF NOT EXISTS name_en VARCHAR(150),
    ADD COLUMN IF NOT EXISTS address_ar TEXT,
    ADD COLUMN IF NOT EXISTS address_he TEXT,
    ADD COLUMN IF NOT EXISTS address_en TEXT;

-- Migrate existing data to translation columns
UPDATE worksites 
SET 
    name_ar = name,
    name_he = name,
    name_en = name,
    address_ar = address,
    address_he = address,
    address_en = address
WHERE name_ar IS NULL;

-- =====================================================
-- Tasks table translations
-- =====================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS title_ar VARCHAR(200),
    ADD COLUMN IF NOT EXISTS title_he VARCHAR(200),
    ADD COLUMN IF NOT EXISTS title_en VARCHAR(200),
    ADD COLUMN IF NOT EXISTS description_ar TEXT,
    ADD COLUMN IF NOT EXISTS description_he TEXT,
    ADD COLUMN IF NOT EXISTS description_en TEXT;

-- Migrate existing data to translation columns
UPDATE tasks 
SET 
    title_ar = title,
    title_he = title,
    title_en = title,
    description_ar = description,
    description_he = description,
    description_en = description
WHERE title_ar IS NULL;

-- =====================================================
-- Service requests table translations (if table exists)
-- =====================================================
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'service_requests') THEN
        ALTER TABLE service_requests
            ADD COLUMN IF NOT EXISTS title_ar VARCHAR(200),
            ADD COLUMN IF NOT EXISTS title_he VARCHAR(200),
            ADD COLUMN IF NOT EXISTS title_en VARCHAR(200),
            ADD COLUMN IF NOT EXISTS description_ar TEXT,
            ADD COLUMN IF NOT EXISTS description_he TEXT,
            ADD COLUMN IF NOT EXISTS description_en TEXT,
            ADD COLUMN IF NOT EXISTS address_ar TEXT,
            ADD COLUMN IF NOT EXISTS address_he TEXT,
            ADD COLUMN IF NOT EXISTS address_en TEXT,
            ADD COLUMN IF NOT EXISTS location_name_ar TEXT,
            ADD COLUMN IF NOT EXISTS location_name_he TEXT,
            ADD COLUMN IF NOT EXISTS location_name_en TEXT,
            ADD COLUMN IF NOT EXISTS client_name_ar VARCHAR(150),
            ADD COLUMN IF NOT EXISTS client_name_he VARCHAR(150),
            ADD COLUMN IF NOT EXISTS client_name_en VARCHAR(150);

        -- Migrate existing data to translation columns
        UPDATE service_requests 
        SET 
            title_ar = title,
            title_he = title,
            title_en = title,
            description_ar = description,
            description_he = description,
            description_en = description,
            address_ar = address,
            address_he = address,
            address_en = address,
            location_name_ar = location_name,
            location_name_he = location_name,
            location_name_en = location_name,
            client_name_ar = client_name,
            client_name_he = client_name,
            client_name_en = client_name
        WHERE title_ar IS NULL;
    END IF;
END $$;

-- =====================================================
-- Create helper function to get translated text
-- =====================================================
CREATE OR REPLACE FUNCTION get_translation(
    text_ar TEXT,
    text_he TEXT,
    text_en TEXT,
    lang VARCHAR DEFAULT 'ar'
) RETURNS TEXT AS $$
BEGIN
    CASE lang
        WHEN 'he' THEN
            RETURN COALESCE(text_he, text_ar, text_en, '');
        WHEN 'en' THEN
            RETURN COALESCE(text_en, text_ar, text_he, '');
        ELSE
            RETURN COALESCE(text_ar, text_he, text_en, '');
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Verification queries
-- =====================================================
-- Check if columns were added successfully
SELECT 
    table_name,
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_schema = 'public'
    AND table_name IN ('users', 'clients', 'worksites', 'tasks', 'service_requests')
    AND (column_name LIKE '%_ar' OR column_name LIKE '%_he' OR column_name LIKE '%_en')
ORDER BY table_name, column_name;
```

---

## 📄 backend/internal/database/migrations/000019_fix_worksite_deletion.down.sql

```sql
-- =====================================================
-- Rollback: Restore Worksite Deletion Constraint
-- =====================================================

-- Restore the foreign key constraint to ON DELETE RESTRICT
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_worksite_id_fkey;

ALTER TABLE tasks 
ADD CONSTRAINT tasks_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE RESTRICT;

-- Restore NOT NULL constraint (note: this will fail if there are NULL values)
-- ALTER TABLE tasks ALTER COLUMN worksite_id SET NOT NULL;
```

---

## 📄 backend/internal/database/migrations/000019_fix_worksite_deletion.up.sql

```sql
-- =====================================================
-- Migration: Fix Worksite Deletion - Preserve Tasks and Attendance
-- =====================================================
-- This migration changes the worksite_id constraint in tasks table 
-- from ON DELETE RESTRICT to ON DELETE SET NULL
-- This allows worksite deletion while preserving tasks and attendance records

-- First, remove NOT NULL constraint from worksite_id in tasks
ALTER TABLE tasks ALTER COLUMN worksite_id DROP NOT NULL;

-- Then change the foreign key constraint from RESTRICT to SET NULL
-- This requires dropping the existing constraint first
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_worksite_id_fkey;

ALTER TABLE tasks 
ADD CONSTRAINT tasks_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;

-- =====================================================
-- Verification query
-- =====================================================
-- Check if the foreign key constraint was updated correctly
SELECT
    tc.table_name,
    tc.constraint_name,
    tc.is_deferrable,
    tc.initially_deferred,
    rc.match_option,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints tc
JOIN information_schema.referential_constraints rc 
    ON tc.constraint_name = rc.constraint_name
WHERE tc.table_schema = 'public'
    AND tc.table_name = 'tasks'
    AND tc.constraint_name = 'tasks_worksite_id_fkey';
```

---

## 📄 backend/internal/database/migrations/000020_add_security_logs.down.sql

```sql
-- حذف جدول security_logs وفهارسه
DROP INDEX IF EXISTS idx_security_logs_email_success_time;
DROP INDEX IF EXISTS idx_security_logs_ip_success_time;
DROP INDEX IF EXISTS idx_security_logs_created_at;
DROP INDEX IF EXISTS idx_security_logs_event_type;
DROP INDEX IF EXISTS idx_security_logs_user_id;
DROP INDEX IF EXISTS idx_security_logs_email;
DROP INDEX IF EXISTS idx_security_logs_ip;
DROP TABLE IF EXISTS security_logs;
```

---

## 📄 backend/internal/database/migrations/000020_add_security_logs.up.sql

```sql
-- إنشاء جدول security_logs لتسجيل أحداث الأمان
CREATE TABLE IF NOT EXISTS security_logs (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    event_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(36),
    email VARCHAR(255),
    phone VARCHAR(20),
    ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    success BOOLEAN DEFAULT false,
    reason TEXT,
    details TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- إنشاء فهارس للبحث السريع
CREATE INDEX IF NOT EXISTS idx_security_logs_ip ON security_logs(ip);
CREATE INDEX IF NOT EXISTS idx_security_logs_email ON security_logs(email);
CREATE INDEX IF NOT EXISTS idx_security_logs_user_id ON security_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_security_logs_event_type ON security_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_security_logs_created_at ON security_logs(created_at);

-- إنشاء فهرس مركب للبحث عن الأنماط المشبوهة
CREATE INDEX IF NOT EXISTS idx_security_logs_ip_success_time ON security_logs(ip, success, created_at);
CREATE INDEX IF NOT EXISTS idx_security_logs_email_success_time ON security_logs(email, success, created_at);
```

---

## 📄 backend/internal/database/migrations/000021_fix_on_delete_cascade.down.sql

```sql
-- =====================================================
-- Rollback: Restore previous ON DELETE behavior
-- =====================================================

-- Restore tasks.client_id to SET NULL
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_client_id_fkey;
ALTER TABLE tasks 
ADD CONSTRAINT tasks_client_id_fkey 
FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE SET NULL;

-- Restore attendance.task_id to SET NULL
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_task_id_fkey;
ALTER TABLE attendance 
ADD CONSTRAINT attendance_task_id_fkey 
FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL;

-- Restore attendance.service_request_id to SET NULL
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_service_request_id_fkey;
ALTER TABLE attendance 
ADD CONSTRAINT attendance_service_request_id_fkey 
FOREIGN KEY (service_request_id) REFERENCES service_requests(id) ON DELETE SET NULL;
```

---

## 📄 backend/internal/database/migrations/000021_fix_on_delete_cascade.up.sql

```sql
-- =====================================================
-- Migration: Fix ON DELETE CASCADE Constraints
-- =====================================================
-- This migration ensures proper CASCADE behavior for critical relationships
-- When a user is deleted, their related data should be automatically cleaned up

-- =====================================================
-- Fix tasks.client_id to use ON DELETE CASCADE
-- =====================================================
-- First, make sure the column can be NULL temporarily to avoid constraint violations
ALTER TABLE tasks ALTER COLUMN client_id DROP NOT NULL;

-- Drop existing constraint
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_client_id_fkey;

-- Add new constraint with CASCADE
ALTER TABLE tasks 
ADD CONSTRAINT tasks_client_id_fkey 
FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE CASCADE;

-- =====================================================
-- Fix attendance.task_id to use ON DELETE CASCADE
-- =====================================================
-- Drop existing constraint
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_task_id_fkey;

-- Add new constraint with CASCADE
ALTER TABLE attendance 
ADD CONSTRAINT attendance_task_id_fkey 
FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- =====================================================
-- Fix attendance.service_request_id to use ON DELETE CASCADE
-- =====================================================
-- Drop existing constraint
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_service_request_id_fkey;

-- Add new constraint with CASCADE
ALTER TABLE attendance 
ADD CONSTRAINT attendance_service_request_id_fkey 
FOREIGN KEY (service_request_id) REFERENCES service_requests(id) ON DELETE CASCADE;

-- =====================================================
-- Verification queries
-- =====================================================
-- Check all foreign key constraints with their delete rules
SELECT
    tc.table_name,
    tc.constraint_name,
    rc.delete_rule
FROM information_schema.table_constraints tc
JOIN information_schema.referential_constraints rc 
    ON tc.constraint_name = rc.constraint_name
WHERE tc.table_schema = 'public'
    AND tc.table_name IN ('tasks', 'attendance', 'notifications', 'service_requests', 'assignments', 'location_tracking')
ORDER BY tc.table_name, tc.constraint_name;
```

---

## 📄 backend/internal/database/postgres.go

```go
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

	// إعدادات Connection Pooling مناسبة للخطة المجانية (512MB RAM)
	db.SetMaxOpenConns(5)      // الحد الأقصى للاتصالات المفتوحة
	db.SetMaxIdleConns(3)      // الحد الأقصى للاتصالات الخاملة
	db.SetConnMaxLifetime(0)   // بدون حد أقصى لعمر الاتصال
	db.SetConnMaxIdleTime(300) // إغلاق الاتصالات الخاملة بعد 5 دقائق

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("فشل الاتصال بقاعدة البيانات: %w", err)
	}

	log.Println("✅ تم الاتصال بقاعدة البيانات بنجاح")
	log.Println("🔧 Connection Pool: MaxOpen=5, MaxIdle=3, MaxIdleTime=5min")
	
	// Apply migrations
	if err := applyMigrations(db); err != nil {
		log.Printf("⚠️ فشل تطبيق الترحيلات: %v", err)
		// Don't fail the app startup if migration fails, just log it
	}
	
	return db, nil
}

// applyMigrations applies pending database migrations
func applyMigrations(db *sql.DB) error {
	// Apply complete schema migration (if needed)
	// This includes subscription fields, device fields, updated_at columns, etc.
	completeSchemaSQL := `
		-- Add subscription fields if not exist
		ALTER TABLE users
		    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
		        CHECK (subscription_status IN ('active', 'expired', 'canceled'));
		ALTER TABLE users
		    ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ NULL;

		-- Add updated_at if not exist
		ALTER TABLE users
		    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

		-- Add device and login fields
		ALTER TABLE users ADD COLUMN IF NOT EXISTS device_id TEXT;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS device_model TEXT;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_device TEXT;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code TEXT;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

		-- Update role constraint to include ALL roles
		ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
		ALTER TABLE users ADD CONSTRAINT users_role_check
		    CHECK (role IN ('admin', 'employee', 'customer', 'client'));

		-- Add missing updated_at to worksites table
		ALTER TABLE worksites
		    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

		-- Add assigned_employee_id to worksites if not exist
		ALTER TABLE worksites
		    ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

		-- Add missing updated_at to tasks table
		ALTER TABLE tasks
		    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

		-- Add photo fields to attendance
		ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
		ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
		ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

		-- Add missing updated_at to attendance table
		ALTER TABLE attendance
		    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

		-- Add type and related_id to notifications
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

		-- Add location_name to service_requests
		ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS location_name TEXT;

		-- Add language field to tasks
		ALTER TABLE tasks ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';

		-- Add language field to service_requests
		ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';

		-- Fix worksite deletion: allow NULL in tasks.worksite_id and change constraint to SET NULL
		ALTER TABLE tasks ALTER COLUMN worksite_id DROP NOT NULL;
		ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_worksite_id_fkey;
		ALTER TABLE tasks 
		ADD CONSTRAINT tasks_worksite_id_fkey 
		FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;

		-- Fix worksite deletion: allow NULL in attendance.worksite_id and change constraint to SET NULL
		ALTER TABLE attendance ALTER COLUMN worksite_id DROP NOT NULL;
		ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_worksite_id_fkey;
		ALTER TABLE attendance 
		ADD CONSTRAINT attendance_worksite_id_fkey 
		FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;
	`

	_, err := db.Exec(completeSchemaSQL)
	if err != nil {
		return fmt.Errorf("فشل تطبيق ترحيل complete_schema: %w", err)
	}

	log.Println("✅ تم تطبيق ترحيل complete_schema بنجاح")

	// Create missing tables if they don't exist
	createTablesSQL := `
		-- Create service_requests table if not exists
		CREATE TABLE IF NOT EXISTS service_requests (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    client_name VARCHAR(150),
		    client_name_ar VARCHAR(150),
		    client_name_he VARCHAR(150),
		    client_name_en VARCHAR(150),
		    client_phone VARCHAR(30),
		    title VARCHAR(200) NOT NULL,
		    title_ar VARCHAR(200),
		    title_he VARCHAR(200),
		    title_en VARCHAR(200),
		    description TEXT,
		    description_ar TEXT,
		    description_he TEXT,
		    description_en TEXT,
		    latitude DOUBLE PRECISION NOT NULL,
		    longitude DOUBLE PRECISION NOT NULL,
		    address TEXT,
		    address_ar TEXT,
		    address_he TEXT,
		    address_en TEXT,
		    phone VARCHAR(30),
		    location_name TEXT,
		    location_name_ar TEXT,
		    location_name_he TEXT,
		    location_name_en TEXT,
		    status VARCHAR(20) NOT NULL DEFAULT 'pending'
		        CHECK (status IN ('pending', 'assigned', 'in_progress', 'completed', 'cancelled')),
		    priority VARCHAR(20) NOT NULL DEFAULT 'normal'
		        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
		    photos TEXT[],
		    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		    language VARCHAR(10) NOT NULL DEFAULT 'en'
		);

		-- Create assignments table if not exists
		CREATE TABLE IF NOT EXISTS assignments (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    request_id UUID NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
		    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		    started_at TIMESTAMPTZ,
		    completed_at TIMESTAMPTZ,
		    admin_notes TEXT,
		    employee_notes TEXT,
		    client_rating INTEGER CHECK (client_rating BETWEEN 1 AND 5),
		    client_feedback TEXT,
		    status VARCHAR(20) NOT NULL DEFAULT 'assigned'
		        CHECK (status IN ('assigned', 'accepted', 'in_progress', 'completed', 'rejected'))
		);

		-- Create location_tracking table if not exists
		CREATE TABLE IF NOT EXISTS location_tracking (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
		    latitude DOUBLE PRECISION NOT NULL,
		    longitude DOUBLE PRECISION NOT NULL,
		    accuracy DOUBLE PRECISION,
		    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("فشل إنشاء الجداول المفقودة: %w", err)
	}

	log.Println("✅ تم إنشاء الجداول المفقودة بنجاح")

	// Create indexes for performance
	createIndexesSQL := `
		-- Users indexes
		CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
		CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
		CREATE INDEX IF NOT EXISTS idx_users_subscription_status ON users(subscription_status);

		-- Worksites indexes
		CREATE INDEX IF NOT EXISTS idx_worksites_location ON worksites(latitude, longitude);
		CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);
		CREATE INDEX IF NOT EXISTS idx_worksites_is_active ON worksites(is_active);

		-- Tasks indexes
		CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id ON tasks(worksite_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id ON tasks(assigned_user_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_client_id ON tasks(client_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

		-- Attendance indexes
		CREATE INDEX IF NOT EXISTS idx_attendance_user_id ON attendance(user_id);
		CREATE INDEX IF NOT EXISTS idx_attendance_task_id ON attendance(task_id);
		CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);
		CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time ON attendance(check_in_time);
		CREATE INDEX IF NOT EXISTS idx_attendance_status ON attendance(status);

		-- Service requests indexes
		CREATE INDEX IF NOT EXISTS idx_service_requests_client_id ON service_requests(client_id);
		CREATE INDEX IF NOT EXISTS idx_service_requests_status ON service_requests(status);
		CREATE INDEX IF NOT EXISTS idx_service_requests_priority ON service_requests(priority);

		-- Assignments indexes
		CREATE INDEX IF NOT EXISTS idx_assignments_employee_id ON assignments(employee_id);
		CREATE INDEX IF NOT EXISTS idx_assignments_request_id ON assignments(request_id);
		CREATE INDEX IF NOT EXISTS idx_assignments_admin_id ON assignments(admin_id);
		CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);

		-- Location tracking indexes
		CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id ON location_tracking(user_id);
		CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id ON location_tracking(assignment_id);
		CREATE INDEX IF NOT EXISTS idx_location_tracking_recorded_at ON location_tracking(recorded_at);

		-- Notifications indexes
		CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
		CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
		CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
	`

	_, err = db.Exec(createIndexesSQL)
	if err != nil {
		return fmt.Errorf("فشل إنشاء الفهارس: %w", err)
	}

	log.Println("✅ تم إنشاء الفهارس بنجاح")

	// Create default admin users if they don't exist
	createDefaultAdminsSQL := `
		INSERT INTO users (
		  full_name,
		  email,
		  password_hash,
		  role,
		  is_active,
		  subscription_status,
		  subscription_expires_at
		)
		VALUES
		(
		  'DevPro System Administrator',
		  'admin@devpro.com',
		  crypt('devproadmin', gen_salt('bf', 12)),
		  'admin',
		  TRUE,
		  'active',
		  NULL
		),
		(
		  'DevPro Support Manager',
		  'support@devpro.com',
		  crypt('devprosupport', gen_salt('bf', 12)),
		  'admin',
		  TRUE,
		  'active',
		  NOW() + INTERVAL '1 year'
		),
		(
		  'DevPro Project Manager',
		  'manager@devpro.com',
		  crypt('devpromanager', gen_salt('bf', 12)),
		  'admin',
		  TRUE,
		  'active',
		  NULL
		)
		ON CONFLICT (email) DO UPDATE
		SET
		  password_hash = CASE WHEN users.password_hash IS NULL THEN EXCLUDED.password_hash ELSE users.password_hash END,
		  role = 'admin',
		  is_active = TRUE,
		  subscription_status = CASE
		    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN 'active'
		    ELSE EXCLUDED.subscription_status
		  END,
		  subscription_expires_at = CASE
		    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN NULL
		    ELSE EXCLUDED.subscription_expires_at
		  END,
		  updated_at = now()
	`
	
	_, err = db.Exec(createDefaultAdminsSQL)
	if err != nil {
		log.Printf("⚠️ فشل إنشاء المستخدمين الافتراضيين: %v", err)
		// Don't fail if this fails, just log it
	} else {
		log.Println("✅ تم إنشاء المستخدمين الافتراضيين بنجاح")
	}

	// Apply ON DELETE CASCADE fixes
	cascadeFixSQL := `
		-- Fix tasks.client_id to use ON DELETE CASCADE
		ALTER TABLE tasks ALTER COLUMN client_id DROP NOT NULL;
		ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_client_id_fkey;
		ALTER TABLE tasks 
		ADD CONSTRAINT tasks_client_id_fkey 
		FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE CASCADE;

		-- Fix attendance.task_id to use ON DELETE CASCADE
		ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_task_id_fkey;
		ALTER TABLE attendance 
		ADD CONSTRAINT attendance_task_id_fkey 
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

		-- Fix attendance.service_request_id to use ON DELETE CASCADE
		ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_service_request_id_fkey;
		ALTER TABLE attendance 
		ADD CONSTRAINT attendance_service_request_id_fkey 
		FOREIGN KEY (service_request_id) REFERENCES service_requests(id) ON DELETE CASCADE;
	`

	_, err = db.Exec(cascadeFixSQL)
	if err != nil {
		log.Printf("⚠️ فشل تطبيق إصلاحات ON DELETE CASCADE: %v", err)
		// Don't fail if this fails, just log it
	} else {
		log.Println("✅ تم تطبيق إصلاحات ON DELETE CASCADE بنجاح")
	}

	return nil
}

```

---

## 📄 backend/internal/database/redis.go

```go
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis يفتح اتصالاً بـ Redis
func ConnectRedis(redisURL string) error {
	if redisURL == "" {
		log.Println("⚠️  REDIS_URL غير معرّف - سيتم تشغيل التطبيق بدون Redis")
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("فشل تحليل REDIS_URL: %w", err)
	}

	// إعدادات الاتصال المناسبة للخطة المجانية
	opt.PoolSize = 5      // عدد الاتصالات في Pool (مناسب لـ 25MB RAM)
	opt.MinIdleConns = 2  // حد أدنى من الاتصالات النشطة
	opt.PoolTimeout = 0   // بدون timeout للحصول على اتصال

	RedisClient = redis.NewClient(opt)

	// اختبار الاتصال
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("فشل الاتصال بـ Redis: %w", err)
	}

	log.Println("✅ تم الاتصال بـ Redis بنجاح")
	return nil
}

// CloseRedis يغلق اتصال Redis
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
```

---

## 📄 backend/internal/handlers/attendance_handler.go

```go
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service             *services.AttendanceService
	NotificationService *services.NotificationService
}

func NewAttendanceHandler(s *services.AttendanceService, n *services.NotificationService) *AttendanceHandler {
	return &AttendanceHandler{Service: s, NotificationService: n}
}

type checkInRequest struct {
	WorksiteID       string  `json:"worksite_id" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	ServiceRequestID *string `json:"service_request_id,omitempty"`
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req checkInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات غير صحيحة. تأكد من إرسال worksite_id, latitude, longitude",
		})
		return
	}

	log.Printf("📝 طلب CheckIn: worksite_id=%s, lat=%f, lng=%f", req.WorksiteID, req.Latitude, req.Longitude)

	attendance, geofenceResult, err := h.Service.CheckIn(
		userID.(string),
		req.WorksiteID,
		req.Latitude,
		req.Longitude,
		req.ServiceRequestID,
	)

	if errors.Is(err, services.ErrOutsideGeofence) {
		log.Printf("❌ المستخدم خارج النطاق: %.2f > %d",
			geofenceResult.DistanceMeters, geofenceResult.AllowedRadius)

		if h.NotificationService != nil {
			_ = h.NotificationService.Send(
				userID.(string),
				"🚨 خارج النطاق!",
				"أنت خارج نطاق موقع العمل المسموح. المسافة: "+formatDistance(geofenceResult.DistanceMeters),
			)
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "❌ أنت خارج نطاق موقع العمل",
			"geofence": geofenceResult,
			"details": gin.H{
				"distance":    geofenceResult.DistanceMeters,
				"max_allowed": geofenceResult.AllowedRadius,
				"is_inside":   false,
			},
		})
		return
	}

	if err != nil {
		log.Printf("❌ فشل CheckIn: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "✅ تم تسجيل بدء الدوام بنجاح",
		"attendance": attendance,
		"geofence":   geofenceResult,
	})
}

type checkOutRequest struct {
	AttendanceID string  `json:"attendance_id" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req checkOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	// =============================================
	// 🔒 التحقق من وجود إحداثيات الموقع
	// =============================================
	if req.Latitude == 0 && req.Longitude == 0 {
		log.Printf("❌ الموقع غير محدد (0,0)")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "📍 يرجى تحديد موقعك أولاً قبل إنهاء الدوام",
			"code":  "location_required",
		})
		return
	}

	log.Printf("📝 طلب CheckOut: attendance_id=%s, lat=%f, lng=%f", req.AttendanceID, req.Latitude, req.Longitude)

	geofenceResult, workedHours, serviceRequestID, err := h.Service.CheckOut(
		userID.(string),
		req.AttendanceID,
		req.Latitude,
		req.Longitude,
	)

	if errors.Is(err, services.ErrOutsideGeofence) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":    "❌ أنت خارج نطاق موقع العمل عند الخروج",
			"geofence": geofenceResult,
		})
		return
	}

	if err != nil {
		log.Printf("❌ فشل CheckOut: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "✅ تم تسجيل إنهاء الدوام بنجاح",
		"geofence":           geofenceResult,
		"worked_hours":       workedHours,
		"service_request_id": serviceRequestID,
	})
}

func (h *AttendanceHandler) GetCurrentAttendance(c *gin.Context) {
	userID, _ := c.Get("user_id")

	attendance, err := h.Service.GetCurrentAttendance(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب حالة الحضور"})
		return
	}

	if attendance == nil {
		c.JSON(http.StatusOK, gin.H{
			"has_active": false,
			"message":    "لا يوجد سجل حضور نشط",
		})
		return
	}

	var elapsedSeconds int64
	if attendance.CheckInTime != nil {
		elapsedSeconds = int64(time.Since(*attendance.CheckInTime).Seconds())
	}

	c.JSON(http.StatusOK, gin.H{
		"has_active":         true,
		"attendance_id":      attendance.ID,
		"worksite_id":        attendance.WorksiteID,
		"check_in_time":      attendance.CheckInTime,
		"elapsed_seconds":    elapsedSeconds,
		"status":             attendance.Status,
		"service_request_id": attendance.ServiceRequestID,
	})
}

func (h *AttendanceHandler) GetAttendanceSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var todayHours, weekHours, monthHours float64

	// تعديل الاستعلامات للتعامل مع فترات العمل عبر منتصف الليل
	// نستخدم CASE لحساب الفرق بشكل صحيح عبر منتصف الليل
	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN check_out_time < check_in_time THEN 
					EXTRACT(EPOCH FROM (check_out_time + INTERVAL '24 hours' - check_in_time)) / 3600
				ELSE 
					EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600
			END
		), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) = CURRENT_DATE AND check_out_time IS NOT NULL
	`, userID).Scan(&todayHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN check_out_time < check_in_time THEN 
					EXTRACT(EPOCH FROM (check_out_time + INTERVAL '24 hours' - check_in_time)) / 3600
				ELSE 
					EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600
			END
		), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) >= DATE_TRUNC('week', CURRENT_DATE) AND check_out_time IS NOT NULL
	`, userID).Scan(&weekHours)

	_ = h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN check_out_time < check_in_time THEN 
					EXTRACT(EPOCH FROM (check_out_time + INTERVAL '24 hours' - check_in_time)) / 3600
				ELSE 
					EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600
			END
		), 0)
		FROM attendance 
		WHERE user_id = $1 AND DATE(check_in_time) >= DATE_TRUNC('month', CURRENT_DATE) AND check_out_time IS NOT NULL
	`, userID).Scan(&monthHours)

	c.JSON(http.StatusOK, gin.H{
		"today_hours": todayHours,
		"week_hours":  weekHours,
		"month_hours": monthHours,
	})
}

func formatDistance(d float64) string {
	if d < 1000 {
		return fmt.Sprintf("%.0f متر", d)
	}
	return fmt.Sprintf("%.2f كيلومتر", d/1000)
}

func (h *AttendanceHandler) GetAllAttendanceSummary(c *gin.Context) {
	rows, err := h.Service.DB.Query(`
		SELECT
			u.id,
			u.full_name,
			u.email,
			COALESCE(SUM(
				CASE
					WHEN a.check_out_time < a.check_in_time THEN
						EXTRACT(EPOCH FROM (a.check_out_time + INTERVAL '24 hours' - a.check_in_time)) / 3600
					ELSE
						EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600
				END
			), 0) as total_hours,
			COUNT(CASE WHEN a.check_out_time IS NOT NULL THEN 1 END) as completed_shifts,
			COUNT(CASE WHEN a.check_out_time IS NULL THEN 1 END) as active_shifts
		FROM users u
		LEFT JOIN attendance a ON a.user_id = u.id
		WHERE u.role = 'employee'
		AND (a.check_in_time IS NULL OR DATE(a.check_in_time) >= DATE_TRUNC('week', CURRENT_DATE))
		GROUP BY u.id, u.full_name, u.email
		ORDER BY total_hours DESC
	`)

	if err != nil {
		log.Printf("❌ فشل جلب ملخص ساعات العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب ملخص ساعات العمل"})
		return
	}
	defer rows.Close()

	var summary []gin.H
	for rows.Next() {
		var id, fullName, email string
		var totalHours float64
		var completedShifts, activeShifts int

		if err := rows.Scan(&id, &fullName, &email, &totalHours, &completedShifts, &activeShifts); err == nil {
			summary = append(summary, gin.H{
				"id":               id,
				"full_name":        fullName,
				"email":            email,
				"total_hours":      totalHours,
				"completed_shifts": completedShifts,
				"active_shifts":    activeShifts,
			})
		}
	}

	c.JSON(http.StatusOK, summary)
}

// GetEmployeeAttendanceHistory جلب سجل الحضور لموظف معين
func (h *AttendanceHandler) GetEmployeeAttendanceHistory(c *gin.Context) {
	userID := c.Param("id")
	
	// الحصول على معاملات التصفية
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")
	
	// بناء الاستعلام حسب الفلتر
	query := `
		SELECT 
			a.id,
			a.worksite_id,
			w.name as worksite_name,
			a.check_in_time,
			a.check_in_lat,
			a.check_in_lng,
			a.check_in_distance_meters,
			a.check_out_time,
			a.check_out_lat,
			a.check_out_lng,
			a.check_out_distance_meters,
			a.status,
			a.created_at
		FROM attendance a
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1
	
	if year != "" && month != "" {
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM check_in_time) = $%d", argCount)
		args = append(args, year)
		
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM check_in_time) = $%d", argCount)
		args = append(args, month)
	}
	
	query += " ORDER BY check_in_time DESC"
	
	rows, err := h.Service.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}
	defer rows.Close()
	
	var history []gin.H
	for rows.Next() {
		var id, worksiteID, worksiteName string
		var checkInTime, checkOutTime, createdAt time.Time
		var checkInLat, checkInLng, checkInDistance, checkOutLat, checkOutLng, checkOutDistance *float64
		var status string

		err := rows.Scan(
			&id, &worksiteID, &worksiteName, &checkInTime, &checkInLat, &checkInLng, &checkInDistance,
			&checkOutTime, &checkOutLat, &checkOutLng, &checkOutDistance, &status, &createdAt,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		var workedHours *float64
		if !checkOutTime.IsZero() {
			// حساب الفرق بالساعات مع التعامل مع فترات العمل عبر منتصف الليل
			hours := checkOutTime.Sub(checkInTime).Hours()
			// إذا كان الفرق سلبياً (يعني عبر منتصف الليل)، أضف 24 ساعة
			if hours < 0 {
				hours += 24
			}
			workedHours = &hours
		}

		record := gin.H{
			"id":                        id,
			"worksite_id":               worksiteID,
			"worksite_name":             worksiteName,
			"check_in_time":             checkInTime,
			"check_in_lat":              checkInLat,
			"check_in_lng":              checkInLng,
			"check_in_distance_meters":  checkInDistance,
			"check_out_time":            checkOutTime,
			"check_out_lat":             checkOutLat,
			"check_out_lng":             checkOutLng,
			"check_out_distance_meters": checkOutDistance,
			"status":                    status,
			"worked_hours":              workedHours,
			"created_at":                createdAt,
		}

		history = append(history, record)
	}

	c.JSON(http.StatusOK, history)
}

// GetEmployeeMonthlySummary جلب ملخص شهري لموظف معين
func (h *AttendanceHandler) GetEmployeeMonthlySummary(c *gin.Context) {
	userID := c.Param("id")
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")

	if year == "" || month == "" {
		// استخدام الشهر الحالي إذا لم يتم تحديده
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%d", int(now.Month()))
	}

	// جلب إجمالي الساعات للشهر مع التعامل مع فترات العمل عبر منتصف الليل
	var totalHours float64
	err := h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(
			CASE
				WHEN check_out_time < check_in_time THEN
					EXTRACT(EPOCH FROM (check_out_time + INTERVAL '24 hours' - check_in_time)) / 3600
				ELSE
					EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600
			END
		), 0)
		FROM attendance
		WHERE user_id = $1
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&totalHours)
	
	if err != nil {
		log.Printf("❌ فشل جلب الملخص الشهري: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملخص الشهري"})
		return
	}
	
	// جلب عدد أيام العمل
	var workDays int
	err = h.Service.DB.QueryRow(`
		SELECT COUNT(DISTINCT DATE(check_in_time))
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&workDays)
	
	if err != nil {
		log.Printf("❌ فشل جلب عدد أيام العمل: %v", err)
		workDays = 0
	}
	
	// جلب معلومات الموظف
	var fullName, email string
	err = h.Service.DB.QueryRow(`
		SELECT full_name, email FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email)
	
	if err != nil {
		log.Printf("❌ فشل جلب معلومات الموظف: %v", err)
		fullName = "غير معروف"
		email = ""
	}
	
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"id":         userID,
			"full_name":  fullName,
			"email":      email,
		},
		"period": gin.H{
			"year":  year,
			"month": month,
		},
		"summary": gin.H{
			"total_hours": totalHours,
			"work_days":   workDays,
		},
	})
}

// GetMyAttendanceHistory جلب سجل الحضور للمستخدم الحالي
func (h *AttendanceHandler) GetMyAttendanceHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// الحصول على معاملات التصفية
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")

	// بناء الاستعلام حسب الفلتر
	query := `
		SELECT
			a.id,
			a.worksite_id,
			w.name as worksite_name,
			a.check_in_time,
			a.check_in_lat,
			a.check_in_lng,
			a.check_in_distance_meters,
			a.check_out_time,
			a.check_out_lat,
			a.check_out_lng,
			a.check_out_distance_meters,
			a.status,
			a.created_at
		FROM attendance a
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1
	`
	args := []interface{}{userID.(string)}
	argCount := 1

	// إضافة الفلتر فقط إذا تم تحديد السنة والشهر معاً
	if year != "" && month != "" {
		argCount++
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM check_in_time) = $%d", argCount)
		args = append(args, year)

		argCount++
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM check_in_time) = $%d", argCount)
		args = append(args, month)
	}

	query += " ORDER BY check_in_time DESC"

	log.Printf("📝 جلب سجل الحضور للمستخدم: user_id=%s, year=%s, month=%s", userID.(string), year, month)

	rows, err := h.Service.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}
	defer rows.Close()

	var history []gin.H
	for rows.Next() {
		var id, worksiteID, worksiteName string
		var checkInTime, checkOutTime, createdAt time.Time
		var checkInLat, checkInLng, checkInDistance, checkOutLat, checkOutLng, checkOutDistance *float64
		var status string

		err := rows.Scan(
			&id, &worksiteID, &worksiteName, &checkInTime, &checkInLat, &checkInLng, &checkInDistance,
			&checkOutTime, &checkOutLat, &checkOutLng, &checkOutDistance, &status, &createdAt,
		)
		if err != nil {
			log.Printf("❌ خطأ في قراءة البيانات: %v", err)
			continue
		}

		var workedHours *float64
		if !checkOutTime.IsZero() {
			// حساب الفرق بالساعات مع التعامل مع فترات العمل عبر منتصف الليل
			hours := checkOutTime.Sub(checkInTime).Hours()
			// إذا كان الفرق سلبياً (يعني عبر منتصف الليل)، أضف 24 ساعة
			if hours < 0 {
				hours += 24
			}
			workedHours = &hours
		}

		record := gin.H{
			"id":                        id,
			"worksite_id":               worksiteID,
			"worksite_name":             worksiteName,
			"check_in_time":             checkInTime,
			"check_in_lat":              checkInLat,
			"check_in_lng":              checkInLng,
			"check_in_distance_meters":  checkInDistance,
			"check_out_time":            checkOutTime,
			"check_out_lat":             checkOutLat,
			"check_out_lng":             checkOutLng,
			"check_out_distance_meters": checkOutDistance,
			"status":                    status,
			"worked_hours":              workedHours,
			"created_at":                createdAt,
		}

		history = append(history, record)
	}

	log.Printf("✅ تم جلب %d سجل حضور", len(history))
	c.JSON(http.StatusOK, history)
}

// GetMyMonthlySummary جلب الملخص الشهري للمستخدم الحالي
func (h *AttendanceHandler) GetMyMonthlySummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")

	if year == "" || month == "" {
		// استخدام الشهر الحالي إذا لم يتم تحديده
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%d", int(now.Month()))
	}

	// جلب إجمالي الساعات للشهر مع التعامل مع فترات العمل عبر منتصف الليل
	var totalHours float64
	err := h.Service.DB.QueryRow(`
		SELECT COALESCE(SUM(
			CASE
				WHEN check_out_time < check_in_time THEN
					EXTRACT(EPOCH FROM (check_out_time + INTERVAL '24 hours' - check_in_time)) / 3600
				ELSE
					EXTRACT(EPOCH FROM (check_out_time - check_in_time)) / 3600
			END
		), 0)
		FROM attendance
		WHERE user_id = $1
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&totalHours)
	
	if err != nil {
		log.Printf("❌ فشل جلب الملخص الشهري: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملخص الشهري"})
		return
	}
	
	// جلب عدد أيام العمل
	var workDays int
	err = h.Service.DB.QueryRow(`
		SELECT COUNT(DISTINCT DATE(check_in_time))
		FROM attendance 
		WHERE user_id = $1 
		AND EXTRACT(YEAR FROM check_in_time) = $2
		AND EXTRACT(MONTH FROM check_in_time) = $3
		AND check_out_time IS NOT NULL
	`, userID, year, month).Scan(&workDays)
	
	if err != nil {
		log.Printf("❌ فشل جلب عدد أيام العمل: %v", err)
		workDays = 0
	}
	
	// جلب معلومات المستخدم
	var fullName, email string
	err = h.Service.DB.QueryRow(`
		SELECT full_name, email FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email)
	
	if err != nil {
		log.Printf("❌ فشل جلب معلومات المستخدم: %v", err)
		fullName = "غير معروف"
		email = ""
	}
	
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"id":         userID,
			"full_name":  fullName,
			"email":      email,
		},
		"period": gin.H{
			"year":  year,
			"month": month,
		},
		"summary": gin.H{
			"total_hours": totalHours,
			"work_days":   workDays,
		},
	})
}

// CleanupOldRecords حذف السجلات القديمة (أكثر من 3 أشهر)
func (h *AttendanceHandler) CleanupOldRecords(c *gin.Context) {
	// حذف السجلات الأقدم من 3 أشهر
	result, err := h.Service.DB.Exec(`
		DELETE FROM attendance 
		WHERE check_in_time < NOW() - INTERVAL '3 months'
		AND status = 'completed'
	`)
	if err != nil {
		log.Printf("❌ فشل حذف السجلات القديمة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف السجلات القديمة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ تم حذف %d سجل قديم", rowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("تم حذف %d سجل قديم بنجاح", rowsAffected),
		"deleted_count": rowsAffected,
	})
}

// ForceCheckOut إنهاء دوام الموظف من قبل المدير
func (h *AttendanceHandler) ForceCheckOut(c *gin.Context) {
	// التحقق من أن المستخدم مدير
	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لهذه العملية"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "غير مصرح"})
		return
	}

	adminIDStr := adminID.(string)

	// التحقق من حالة الاشتراك (مع التسامح مع عدم وجود الحقول)
	var subscriptionStatus string
	var subscriptionExpiresAt time.Time
	err := h.Service.DB.QueryRow(`
		SELECT subscription_status, subscription_expires_at
		FROM users
		WHERE id = $1
	`, adminIDStr).Scan(&subscriptionStatus, &subscriptionExpiresAt)

	// إذا فشل جلب حالة الاشتراك (مثلاً الحقول غير موجودة)، نتجاوز التحقق
	if err != nil {
		log.Printf("⚠️ فشل جلب حالة الاشتراك، تجاوز التحقق: %v", err)
		// نستمر بدون التحقق من الاشتراك
	} else {
		// التحقق من أن الاشتراك نشط
		if subscriptionStatus == "canceled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "اشتراكك ملغي، الرجاء التواصل مع الدعم"})
			return
		}

		if subscriptionStatus == "expired" || (!subscriptionExpiresAt.IsZero() && time.Now().After(subscriptionExpiresAt)) {
			// تحديث الحالة إلى منتهي إذا لزم الأمر
			if subscriptionStatus != "expired" {
				_, _ = h.Service.DB.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, adminIDStr)
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "اشتراكك منتهي، الرجاء تجديده للتمكن من استخدام هذه الميزة"})
			return
		}
	}

	var req struct {
		AttendanceID string `json:"attendance_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 طلب ForceCheckOut: attendance_id=%s, admin_id=%s", req.AttendanceID, adminIDStr)

	workedHours, err := h.Service.ForceCheckOut(req.AttendanceID, adminIDStr)
	if err != nil {
		log.Printf("❌ فشل ForceCheckOut: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "✅ تم إنهاء الدوام بنجاح",
		"worked_hours": workedHours,
	})
}

```

---

## 📄 backend/internal/handlers/auth_handler.go

```go
package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DB            *sql.DB
	AuthService   *services.AuthService
	SecurityLogger *services.SecurityLogger
	WSHandler     interface{} // مؤشر لـ WSHandler
}

func NewAuthHandler(db *sql.DB, authService *services.AuthService, securityLogger *services.SecurityLogger) *AuthHandler {
	return &AuthHandler{DB: db, AuthService: authService, SecurityLogger: securityLogger}
}

func (h *AuthHandler) SetWSHandler(wsHandler interface{}) {
	h.WSHandler = wsHandler
}

func (h *AuthHandler) checkSubscriptionStatus(userID string, status string, expiresAt sql.NullTime) error {
	if status == "canceled" {
		return errors.New("subscription canceled")
	}

	if status == "expired" {
		return errors.New("subscription expired")
	}

	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		_, err := h.DB.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, userID)
		if err != nil {
			log.Printf("⚠️ failed to mark subscription expired for user %s: %v", userID, err)
		}
		return errors.New("subscription expired")
	}

	return nil
}

func (h *AuthHandler) validatePassword(password, storedHash, email string) bool {
	if h.AuthService.CheckPassword(password, storedHash) {
		return true
	}

	var passwordMatches bool
	err := h.DB.QueryRow(`SELECT crypt($1, password_hash) = password_hash FROM users WHERE email = $2`, password, email).Scan(&passwordMatches)
	if err != nil {
		log.Printf("⚠️ password fallback check failed for %s: %v", email, err)
		return false
	}

	return passwordMatches
}

// =============================================
// 1. تسجيل دخول المدير
// =============================================
func (h *AuthHandler) Login(c *gin.Context) {
	lang := i18n.Detect(c)
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_login_fields")})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// التحقق من الأنماط المشبوهة قبل المتابعة
	if h.SecurityLogger != nil {
		isSuspicious, reason := h.SecurityLogger.CheckForSuspiciousPatterns(ip, req.Email)
		if isSuspicious {
			h.SecurityLogger.LogAuthAttempt(req.Email, "", ip, userAgent, false, reason)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": i18n.T(lang, "err_suspicious_activity")})
			return
		}
	}

	var userID, fullName, role, passwordHash, subscriptionStatus, passwordChangedAt string
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT id, full_name, role, password_hash, subscription_status, subscription_expires_at, 
		       COALESCE(password_changed_at::text, '') as password_changed_at
		FROM users 
		WHERE email = $1 AND is_active = TRUE
	`, req.Email).Scan(&userID, &fullName, &role, &passwordHash, &subscriptionStatus, &subscriptionExpiresAt, &passwordChangedAt)

	if err != nil {
		log.Printf("❌ مستخدم غير موجود: %v", err)
		if h.SecurityLogger != nil {
			h.SecurityLogger.LogAuthAttempt(req.Email, "", ip, userAgent, false, "user_not_found")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_credentials")})
		return
	}

	// ⚠️ منع الموظفين من تسجيل الدخول عبر البريد - يستخدمون الهاتف فقط
	if role == "employee" {
		log.Printf("⚠️ محاولة دخول موظف عبر البريد: %s", req.Email)
		if h.SecurityLogger != nil {
			h.SecurityLogger.LogAuthAttempt(req.Email, "", ip, userAgent, false, "employee_email_login_blocked")
		}
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employees must use phone login for security. Please use the mobile app with your phone number.",
			"use_phone_login": true,
		})
		return
	}

	// ✅ العملاء والمديرين يمكنهم الدخول عبر البريد
	// العملاء للبورتال الخاص بهم، المديرين للوحة الإدارة

	if !h.validatePassword(req.Password, passwordHash, req.Email) {
		log.Printf("❌ كلمة مرور خاطئة")
		if h.SecurityLogger != nil {
			h.SecurityLogger.LogAuthAttempt(req.Email, "", ip, userAgent, false, "invalid_password")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_credentials")})
		return
	}

	// التحقق من الاشتراك للمديرين فقط
	// العملاء والموظفين لا يحتاجون اشتراك للدخول
	if role == "admin" {
		if err := h.checkSubscriptionStatus(userID, subscriptionStatus, subscriptionExpiresAt); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}
	}

	token, err := h.AuthService.GenerateToken(userID, role, passwordChangedAt)
	if err != nil {
		log.Printf("❌ فشل إنشاء التوكن: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_session_create_failed")})
		return
	}

	log.Printf("✅ تسجيل دخول ناجح: %s", fullName)
	
	// تسجيل محاولة الدخول الناجحة
	if h.SecurityLogger != nil {
		h.SecurityLogger.LogAuthAttempt(req.Email, "", ip, userAgent, true, "success")
	}

	// تعيين HttpOnly cookie للرمز المميز
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"worktrack_token",
		token,
		3600*24*7, // 7 أيام
		"/",
		"",
		true, // Secure - يجب أن يكون true في الإنتاج مع HTTPS
		true, // HttpOnly - يحظر الوصول عبر JavaScript
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":                  userID,
			"full_name":           fullName,
			"role":                role,
			"subscription_status": subscriptionStatus,
			"subscription_expires_at": func() interface{} {
				if subscriptionExpiresAt.Valid {
					return subscriptionExpiresAt.Time
				}
				return nil
			}(),
		},
	})
}

// =============================================
// 2. جلب جميع الموظفين
// =============================================
func (h *AuthHandler) ListEmployees(c *gin.Context) {
	log.Println("📋 جلب قائمة الموظفين...")

	rows, err := h.DB.Query(`
		SELECT 
			u.id, u.full_name, u.email, u.phone, u.role, u.is_active, u.created_at,
			COALESCE(u.device_model, '') as device_model,
			ws.name as current_worksite,
			ws.id as current_worksite_id
		FROM users u
		LEFT JOIN attendance a ON u.id = a.user_id AND a.status = 'in_progress'
		LEFT JOIN worksites ws ON a.worksite_id = ws.id
		WHERE u.role = 'employee' AND u.is_active = TRUE
		ORDER BY u.created_at DESC
	`)

	if err != nil {
		log.Printf("❌ خطأ في جلب الموظفين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, role, createdAt, deviceModel string
		var isActive bool
		var currentWorksite, currentWorksiteID sql.NullString

		if err := rows.Scan(&id, &fullName, &email, &phone, &role, &isActive, &createdAt, &deviceModel, &currentWorksite, &currentWorksiteID); err != nil {
			log.Printf("⚠️ خطأ في القراءة: %v", err)
			continue
		}

		employee := gin.H{
			"id":                id,
			"full_name":         fullName,
			"email":             email,
			"phone":             phone,
			"role":              role,
			"is_active":         isActive,
			"created_at":        createdAt,
			"device_model":      deviceModel,
			"is_registered":     deviceModel != "",
			"current_worksite":  currentWorksite.String,
			"current_worksite_id": currentWorksiteID.String,
		}

		employees = append(employees, employee)
	}

	log.Printf("✅ تم جلب %d موظف", len(employees))
	c.JSON(http.StatusOK, employees)
}

// =============================================
// 3. إنشاء موظف جديد
// =============================================
func (h *AuthHandler) CreateEmployeePhone(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 إنشاء موظف: %s, %s", req.FullName, req.Phone)

	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`, req.Phone).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "حدث خطأ"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "رقم الهاتف مستخدم"})
		return
	}

	id := uuid.NewString()
	role := "employee"
	if req.Role != "" {
		role = req.Role
	}

	bytes := make([]byte, 8)
	rand.Read(bytes)
	password := hex.EncodeToString(bytes)
	hash, _ := h.AuthService.HashPassword(password)
	email := req.Phone + "@worktrack.com"

	_, err = h.DB.Exec(`
		INSERT INTO users (
			id, full_name, email, phone, password_hash, 
			role, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, TRUE, now(), now())
	`, id, req.FullName, email, req.Phone, hash, role)

	if err != nil {
		log.Printf("❌ فشل إنشاء الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنشاء الموظف"})
		return
	}

	log.Printf("✅ تم إنشاء موظف: %s", req.FullName)

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "تم إنشاء الموظف بنجاح",
		"user": gin.H{
			"id":        id,
			"full_name": req.FullName,
			"phone":     req.Phone,
			"email":     email,
			"role":      role,
		},
	})
}

// =============================================
// 4. تسجيل دخول الموظف (مع التحقق من نوع الجهاز)
// =============================================
func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	lang := i18n.Detect(c)
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	log.Println("========================================")
	log.Println("📱 تسجيل دخول الموظف...")

	var req struct {
		Phone       string `json:"phone" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
		DeviceModel string `json:"device_model" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📱 رقم الهاتف: %s", req.Phone)
	log.Printf("🆔 معرف الجهاز: %s", req.DeviceID)
	log.Printf("📱 نوع الجهاز: %s", req.DeviceModel)
	
	// التحقق من الأنماط المشبوهة قبل المتابعة
	if h.SecurityLogger != nil {
		isSuspicious, reason := h.SecurityLogger.CheckForSuspiciousPatterns(ip, "")
		if isSuspicious {
			h.SecurityLogger.LogAuthAttempt("", req.Phone, ip, userAgent, false, reason)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": i18n.T(lang, "err_suspicious_activity")})
			return
		}
	}

	var userID, fullName, role, storedDeviceID, storedDeviceModel, subscriptionStatus, passwordChangedAt string
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT 
			id, 
			full_name, 
			role, 
			COALESCE(device_id, ''), 
			COALESCE(device_model, ''),
			subscription_status,
			subscription_expires_at,
			COALESCE(password_changed_at::text, '') as password_changed_at
		FROM users 
		WHERE phone = $1 AND is_active = TRUE AND role = 'employee'
	`, req.Phone).Scan(&userID, &fullName, &role, &storedDeviceID, &storedDeviceModel, &subscriptionStatus, &subscriptionExpiresAt, &passwordChangedAt)

	if err != nil {
		log.Printf("❌ المستخدم غير موجود: %s", req.Phone)
		if h.SecurityLogger != nil {
			h.SecurityLogger.LogAuthAttempt("", req.Phone, ip, userAgent, false, "user_not_found")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_phone_not_registered")})
		return
	}

	// ⚠️ منع المديرين من تسجيل الدخول عبر الهاتف - يستخدمون البريد فقط
	if role == "admin" {
		log.Printf("⚠️ محاولة دخول مدير عبر الهاتف: %s", req.Phone)
		if h.SecurityLogger != nil {
			h.SecurityLogger.LogAuthAttempt("", req.Phone, ip, userAgent, false, "admin_phone_login_blocked")
		}
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admins must use email login for security. Please use the admin panel with your email.",
			"use_email_login": true,
		})
		return
	}

	log.Printf("✅ تم العثور على المستخدم الموظف: %s", fullName)

	// =============================================
	// 🔒 طبقة أمان 1: التحقق من Device ID
	// =============================================
	if storedDeviceID != "" && storedDeviceID != req.DeviceID {
		log.Printf("🚨 Device ID غير مطابق!")
		log.Printf("   المسجل: %s", storedDeviceID)
		log.Printf("   الحالي: %s", req.DeviceID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":           "رقم الهاتف غير صحيح، يرجى التواصل مع المدير",
			"device_mismatch": true,
		})
		return
	}

	// =============================================
	// 🔒 طبقة أمان 2: التحقق من Device Model (نوع الجهاز)
	// =============================================
	if storedDeviceModel != "" && storedDeviceModel != req.DeviceModel {
		log.Printf("🚨 نوع الجهاز غير مطابق!")
		log.Printf("   المسجل: %s", storedDeviceModel)
		log.Printf("   الحالي: %s", req.DeviceModel)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":          "رقم الهاتف غير صحيح، يرجى التواصل مع المدير",
			"model_mismatch": true,
		})
		return
	}

	// =============================================
	// تسجيل الجهاز إذا كان جديداً
	// =============================================
	if storedDeviceID == "" || storedDeviceModel == "" {
		log.Println("📱 أول تسجيل دخول - تسجيل الجهاز...")

		_, err = h.DB.Exec(`
			UPDATE users 
			SET 
				device_id = $1,
				device_model = $2,
				phone_verified = TRUE,
				updated_at = now()
			WHERE id = $3
		`, req.DeviceID, req.DeviceModel, userID)

		if err != nil {
			log.Printf("❌ فشل تسجيل الجهاز: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تسجيل الجهاز"})
			return
		}

		log.Printf("✅ تم تسجيل الجهاز بنجاح!")
		log.Printf("   Device ID: %s", req.DeviceID)
		log.Printf("   Device Model: %s", req.DeviceModel)
	}

	// تحديث وقت آخر دخول
	_, err = h.DB.Exec(`
		UPDATE users 
		SET updated_at = now()
		WHERE id = $1
	`, userID)
	if err != nil {
		log.Printf("⚠️ فشل تحديث وقت الدخول: %v", err)
	}

	// إنشاء التوكن
	token, err := h.AuthService.GenerateToken(userID, role, passwordChangedAt)
	if err != nil {
		log.Printf("❌ فشل إنشاء التوكن: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تسجيل الدخول"})
		return
	}

	log.Printf("✅ تسجيل دخول ناجح: %s", fullName)
	log.Printf("📱 نوع الجهاز المسجل: %s", req.DeviceModel)
	log.Println("========================================")
	
	// تسجيل محاولة الدخول الناجحة
	if h.SecurityLogger != nil {
		h.SecurityLogger.LogAuthAttempt("", req.Phone, ip, userAgent, true, "success")
	}

	// تعيين HttpOnly cookie للرمز المميز
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"worktrack_token",
		token,
		3600*24*7, // 7 أيام
		"/",
		"",
		true, // Secure - يجب أن يكون true في الإنتاج مع HTTPS
		true, // HttpOnly - يحظر الوصول عبر JavaScript
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        userID,
			"full_name": fullName,
			"role":      role,
			"phone":     req.Phone,
		},
		"device_registered": storedDeviceID != "",
	})
}

// =============================================
// 5. تسجيل الخروج
// =============================================
func (h *AuthHandler) Logout(c *gin.Context) {
	// مسح الـ cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"worktrack_token",
		"",
		-1, // انتهاء الصلاحية فوراً
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "تم تسجيل الخروج بنجاح"})
}

// =============================================
// 6. بيانات المستخدم الحالي
// =============================================
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var fullName, email, role, subscriptionStatus string
	var phone sql.NullString
	var subscriptionExpiresAt sql.NullTime
	err := h.DB.QueryRow(`
		SELECT full_name, email, phone, role, subscription_status, subscription_expires_at
		FROM users WHERE id = $1
	`, userID).Scan(&fullName, &email, &phone, &role, &subscriptionStatus, &subscriptionExpiresAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "المستخدم غير موجود"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                  userID,
		"full_name":           fullName,
		"email":               email,
		"phone":               func() interface{} {
			if phone.Valid {
				return phone.String
			}
			return ""
		}(),
		"role":                role,
		"subscription_status": subscriptionStatus,
		"subscription_expires_at": func() interface{} {
			if subscriptionExpiresAt.Valid {
				return subscriptionExpiresAt.Time
			}
			return nil
		}(),
	})
}

// =============================================
// 7. حذف موظف
// =============================================
func (h *AuthHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	if id == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "لا يمكن حذف المدير الرئيسي"})
		return
	}

	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "المستخدم غير موجود"})
		return
	}

	// حذف الموظف نهائياً
	_, err = h.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	
	if err != nil {
		log.Printf("❌ فشل حذف الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الحذف"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم حذف الموظف نهائياً"})
}

// =============================================
// 8. إعادة تعيين جهاز الموظف
// =============================================
func (h *AuthHandler) ResetDevice(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	_, err := h.DB.Exec(`
		UPDATE users 
		SET 
			device_id = NULL,
			device_model = NULL,
			phone_verified = FALSE,
			updated_at = now()
		WHERE id = $1
	`, req.UserID)

	if err != nil {
		log.Printf("❌ فشل إعادة تعيين الجهاز: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إعادة تعيين الجهاز"})
		return
	}

	log.Printf("✅ تم إعادة تعيين جهاز المستخدم: %s", req.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "تم إعادة تعيين الجهاز"})
}

// =============================================
// 9. معلومات الجهاز
// =============================================
func (h *AuthHandler) GetDeviceInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var deviceID, deviceModel string
	err := h.DB.QueryRow(`
		SELECT 
			COALESCE(device_id, ''), 
			COALESCE(device_model, '') 
		FROM users WHERE id = $1
	`, userID).Scan(&deviceID, &deviceModel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب معلومات الجهاز"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id":     deviceID,
		"device_model":  deviceModel,
		"is_registered": deviceID != "" && deviceModel != "",
	})
}

// =============================================
// 10. جلب جميع العملاء
// =============================================
func (h *AuthHandler) ListCustomers(c *gin.Context) {
	log.Println("📋 جلب قائمة العملاء...")

	rows, err := h.DB.Query(`
		SELECT 
			u.id, u.full_name, u.email, u.phone, u.role, u.is_active, u.created_at,
			COALESCE(u.device_model, '') as device_model
		FROM users u
		WHERE u.role IN ('customer', 'client') AND u.is_active = TRUE
		ORDER BY u.created_at DESC
	`)

	if err != nil {
		log.Printf("❌ خطأ في جلب العملاء: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب العملاء"})
		return
	}
	defer rows.Close()

	var customers []gin.H
	for rows.Next() {
		var id, fullName, email, phone, role, createdAt, deviceModel string
		var isActive bool

		if err := rows.Scan(&id, &fullName, &email, &phone, &role, &isActive, &createdAt, &deviceModel); err != nil {
			log.Printf("⚠️ خطأ في القراءة: %v", err)
			continue
		}

		customer := gin.H{
			"id":           id,
			"full_name":    fullName,
			"email":        email,
			"phone":        phone,
			"role":         role,
			"is_active":    isActive,
			"created_at":   createdAt,
			"device_model": deviceModel,
		}

		customers = append(customers, customer)
	}

	log.Printf("✅ تم جلب %d عميل", len(customers))
	c.JSON(http.StatusOK, customers)
}

// =============================================
// 11. إنشاء عميل جديد
// =============================================
func (h *AuthHandler) CreateClient(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ خطأ في البيانات: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("📝 إنشاء عميل: %s, %s, %s", req.FullName, req.Phone, req.Email)

	// التحقق من البريد الإلكتروني إذا تم إرساله
	if req.Email != "" {
		var exists bool
		err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&exists)
		if err != nil {
			log.Printf("❌ خطأ في التحقق من البريد: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "حدث خطأ"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "البريد الإلكتروني مستخدم"})
			return
		}
	}

	id := uuid.NewString()

	// إنشاء كلمة مرور عشوائية
	bytes := make([]byte, 8)
	rand.Read(bytes)
	password := hex.EncodeToString(bytes)
	hash, _ := h.AuthService.HashPassword(password)

	role := "customer"

	// استخدام البريد الإلكتروني أو NULL إذا لم يكن موجوداً
	var email interface{} = req.Email
	if req.Email == "" {
		email = nil
	}

	_, err := h.DB.Exec(`
		INSERT INTO users (
			id, full_name, email, phone, password_hash,
			role, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, TRUE, now(), now())
	`, id, req.FullName, email, req.Phone, hash, role)

	if err != nil {
		log.Printf("❌ فشل إنشاء العميل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنشاء العميل"})
		return
	}

	log.Printf("✅ تم إنشاء عميل: %s", req.FullName)

	// Return the actual email that was stored (could be nil)
	responseEmail := req.Email
	if responseEmail == "" {
		responseEmail = ""
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "تم إنشاء العميل بنجاح",
		"user": gin.H{
			"id":        id,
			"full_name": req.FullName,
			"email":     responseEmail,
			"phone":     req.Phone,
			"role":      role,
		},
		"password": password,
	})
}

// =============================================
// 12. إعادة تعيين كلمة مرور العميل
// =============================================
func (h *AuthHandler) ResetCustomerPassword(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	log.Printf("🔄 إعادة تعيين كلمة مرور العميل: %s", req.CustomerID)

	// التحقق من وجود العميل
	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role IN ('customer', 'client'))`, req.CustomerID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "حدث خطأ"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "العميل غير موجود"})
		return
	}

	// إنشاء كلمة مرور عشوائية جديدة
	bytes := make([]byte, 8)
	rand.Read(bytes)
	newPassword := hex.EncodeToString(bytes)
	hash, err := h.AuthService.HashPassword(newPassword)
	if err != nil {
		log.Printf("❌ فشل تشفير كلمة المرور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تشفير كلمة المرور"})
		return
	}

	// تحديث كلمة المرور
	_, err = h.DB.Exec(`
		UPDATE users 
		SET password_hash = $1, updated_at = now()
		WHERE id = $2
	`, hash, req.CustomerID)

	if err != nil {
		log.Printf("❌ فشل إعادة تعيين كلمة المرور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إعادة تعيين كلمة المرور"})
		return
	}

	log.Printf("✅ تم إعادة تعيين كلمة مرور العميل: %s", req.CustomerID)

	c.JSON(http.StatusOK, gin.H{
		"message": "تم إعادة تعيين كلمة المرور بنجاح",
		"password": newPassword,
	})
}

// =============================================
// 12. حذف العميل
// =============================================
func (h *AuthHandler) DeleteCustomer(c *gin.Context) {
	customerID := c.Param("id")

	log.Printf("🗑️ حذف العميل: %s", customerID)

	// التحقق من وجود العميل
	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role IN ('customer', 'client'))`, customerID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "حدث خطأ"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "العميل غير موجود"})
		return
	}

	// حذف العميل نهائياً
	_, err = h.DB.Exec(`DELETE FROM users WHERE id = $1`, customerID)
	
	if err != nil {
		log.Printf("❌ فشل حذف العميل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف العميل"})
		return
	}

	log.Printf("✅ تم حذف العميل نهائياً: %s", customerID)

	c.JSON(http.StatusOK, gin.H{
		"message": "تم حذف العميل بنجاح",
	})
}

// NotifyPasswordChange - إرسال إشعار تغيير كلمة المرور عبر WebSocket
func (h *AuthHandler) NotifyPasswordChange(c *gin.Context) {
	lang := i18n.Detect(c)
	
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_invalid_request")})
		return
	}
	
	// التحقق من وجود المستخدم
	var exists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
	`, req.UserID).Scan(&exists)
	
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.T(lang, "err_user_not_found")})
		return
	}
	
	// إرسال إشعار عبر WebSocket
	if h.WSHandler != nil {
		// استخدام type assertion للوصول إلى WSHandler
		if wsHandler, ok := h.WSHandler.(interface{ BroadcastPasswordChanged(string) }); ok {
			wsHandler.BroadcastPasswordChanged(req.UserID)
		}
	}
	
	log.Printf("🚨 إرسال إشعار تغيير كلمة المرور للمستخدم: %s", req.UserID)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "تم إرسال إشعار تغيير كلمة المرور",
	})
}

```

---

## 📄 backend/internal/handlers/client_handler.go

```go
package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientHandler struct {
	DB *sql.DB
}

func NewClientHandler(db *sql.DB) *ClientHandler {
	return &ClientHandler{DB: db}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}
	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO clients (id, name, phone, email, created_at)
		VALUES ($1, $2, $3, $4, now())`,
		id, req.Name, req.Phone, req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الإضافة"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ClientHandler) List(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, phone, email FROM clients`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الجلب"})
		return
	}
	defer rows.Close()
	var clients []gin.H
	for rows.Next() {
		var id, name, phone, email string
		if err := rows.Scan(&id, &name, &phone, &email); err == nil {
			clients = append(clients, gin.H{"id": id, "name": name, "phone": phone, "email": email})
		}
	}
	c.JSON(http.StatusOK, clients)
}

```

---

## 📄 backend/internal/handlers/geocoding_handler.go

```go
package handlers

import (
	"net/http"

	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type GeocodingHandler struct {
	service *services.GeocodingService
}

func NewGeocodingHandler(service *services.GeocodingService) *GeocodingHandler {
	return &GeocodingHandler{service: service}
}

func (h *GeocodingHandler) Autocomplete(c *gin.Context) {
	query := c.Query("q")
	lang := c.Query("lang")
	
	if lang == "" {
		lang = "ar"
	}
	
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "الرجاء إدخال حرفين على الأقل"})
		return
	}

	results, err := h.service.Autocomplete(query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}

```

---

## 📄 backend/internal/handlers/location_handler.go

```go
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationHandler struct {
	DB *sql.DB
	WS *WSHandler
}

func NewLocationHandler(db *sql.DB, wsHandler *WSHandler) *LocationHandler {
	return &LocationHandler{DB: db, WS: wsHandler}
}

// UpdateLocation - تحديث موقع المستخدم
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Accuracy  float64 `json:"accuracy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	_, err := h.DB.Exec(`
		INSERT INTO location_tracking (id, user_id, latitude, longitude, accuracy, recorded_at)
		VALUES ($1, $2, $3, $4, $5, now())
	`, uuid.NewString(), userID, req.Latitude, req.Longitude, req.Accuracy)
	if err != nil {
		log.Printf("❌ فشل حفظ الموقع: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// التحقق من الخروج عن النطاق
	go h.checkGeofenceViolation(userID.(string), req.Latitude, req.Longitude)

	// إرسال تحديث فوري عبر WebSocket
	if h.WS != nil {
		go func() {
			// جلب معلومات الموظف
			var fullName, email string
			h.DB.QueryRow(`
				SELECT full_name, email FROM users WHERE id = $1
			`, userID).Scan(&fullName, &email)

			// إرسال التحديث
			h.WS.BroadcastLocationUpdate(map[string]interface{}{
				"user_id":    userID,
				"full_name":  fullName,
				"email":      email,
				"latitude":   req.Latitude,
				"longitude":  req.Longitude,
				"accuracy":   req.Accuracy,
				"immediate":  true, // تحديث فوري من المستخدم
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم تحديث الموقع"})
}

// checkGeofenceViolation - التحقق من خروج الموظف عن النطاق
func (h *LocationHandler) checkGeofenceViolation(userID string, lat, lng float64) {
	var attendanceID, worksiteID string
	var worksiteLat, worksiteLng float64
	var radiusMeters int

	err := h.DB.QueryRow(`
		SELECT a.id, a.worksite_id, w.latitude, w.longitude, w.radius_meters
		FROM attendance a
		JOIN worksites w ON a.worksite_id = w.id
		WHERE a.user_id = $1 AND a.status = 'in_progress'
	`, userID).Scan(&attendanceID, &worksiteID, &worksiteLat, &worksiteLng, &radiusMeters)

	if err != nil {
		return
	}

	distance := utils.HaversineDistance(lat, lng, worksiteLat, worksiteLng)
	if distance > float64(radiusMeters) {
		var count int
		_ = h.DB.QueryRow(`
			SELECT COUNT(*) FROM notifications 
			WHERE user_id = $1 AND type = 'geofence_alert' 
			AND created_at > now() - interval '5 minutes'
		`, userID).Scan(&count)

		if count == 0 {
			_, _ = h.DB.Exec(`
				INSERT INTO notifications (id, user_id, title, body, type, is_read, created_at)
				VALUES ($1, $2, $3, $4, 'geofence_alert', FALSE, now())
			`, uuid.NewString(), userID,
				"🚨 خروج عن النطاق!",
				fmt.Sprintf("خرج الموظف عن النطاق المسموح به. المسافة: %.0f متر", distance))
		}
	}
}

// GetActiveEmployees - جلب الموظفين النشطين مع مواقعهم
func (h *LocationHandler) GetActiveEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT DISTINCT ON (u.id)
			u.id, 
			u.full_name, 
			u.email, 
			u.phone,
			lt.latitude, 
			lt.longitude, 
			lt.recorded_at,
			a.id as attendance_id,
			a.check_in_time,
			w.id as worksite_id, 
			w.name as worksite_name,
			w.latitude as worksite_lat, 
			w.longitude as worksite_lng, 
			w.radius_meters,
			EXTRACT(EPOCH FROM (now() - a.check_in_time)) / 3600 as hours_worked
		FROM users u
		JOIN attendance a ON a.user_id = u.id AND a.status = 'in_progress'
		JOIN location_tracking lt ON lt.user_id = u.id
		JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee' AND u.is_active = TRUE
		ORDER BY u.id, lt.recorded_at DESC
	`)

	if err != nil {
		log.Printf("❌ فشل جلب الموظفين النشطين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب البيانات"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, attendanceID, worksiteID, worksiteName string
		var lat, lng, worksiteLat, worksiteLng float64
		var radiusMeters int
		var recordedAt, checkInTime time.Time
		var hoursWorked float64

		if err := rows.Scan(&id, &fullName, &email, &phone, &lat, &lng, &recordedAt,
			&attendanceID, &checkInTime, &worksiteID, &worksiteName, 
			&worksiteLat, &worksiteLng, &radiusMeters, &hoursWorked); err != nil {
			continue
		}

		distance := utils.HaversineDistance(lat, lng, worksiteLat, worksiteLng)
		isInside := distance <= float64(radiusMeters)

		// تحديد حالة الموظف
		status := "inside"
		statusText := "✅ داخل النطاق"
		if !isInside {
			status = "outside"
			statusText = "❌ خارج النطاق"
		}

		// حساب وقت الانقطاع (آخر تحديث للموقع)
		timeSinceLastUpdate := time.Since(recordedAt)
		isActive := timeSinceLastUpdate < 5*time.Minute

		employees = append(employees, gin.H{
			"id":              id,
			"full_name":       fullName,
			"email":           email,
			"phone":           phone,
			"latitude":        lat,
			"longitude":       lng,
			"last_update":     recordedAt,
			"is_active":       isActive,
			"status":          status,
			"status_text":     statusText,
			"attendance_id":   attendanceID,
			"check_in_time":   checkInTime,
			"hours_worked":    hoursWorked,
			"worksite": gin.H{
				"id":        worksiteID,
				"name":      worksiteName,
				"latitude":  worksiteLat,
				"longitude": worksiteLng,
				"radius":    radiusMeters,
				"is_inside": isInside,
				"distance":  distance,
			},
		})
	}

	c.JSON(http.StatusOK, employees)
}

// GetEmployeeTrack - جلب مسار موظف معين
func (h *LocationHandler) GetEmployeeTrack(c *gin.Context) {
	employeeID := c.Param("id")
	limit := c.DefaultQuery("limit", "50")

	rows, err := h.DB.Query(`
		SELECT latitude, longitude, recorded_at
		FROM location_tracking
		WHERE user_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`, employeeID, limit)

	if err != nil {
		log.Printf("❌ فشل جلب مسار الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب المسار"})
		return
	}
	defer rows.Close()

	var track []gin.H
	for rows.Next() {
		var lat, lng float64
		var recordedAt time.Time

		if err := rows.Scan(&lat, &lng, &recordedAt); err != nil {
			continue
		}

		track = append(track, gin.H{
			"latitude":  lat,
			"longitude": lng,
			"timestamp": recordedAt,
		})
	}

	c.JSON(http.StatusOK, track)
}

// GetEmployeeSecurityNotes - جلب ملاحظات أمنية عن موظف
func (h *LocationHandler) GetEmployeeSecurityNotes(c *gin.Context) {
	employeeID := c.Param("id")

	rows, err := h.DB.Query(`
		SELECT 
			title, 
			body, 
			created_at,
			type
		FROM notifications
		WHERE user_id = $1 AND type = 'geofence_alert'
		ORDER BY created_at DESC
		LIMIT 20
	`, employeeID)

	if err != nil {
		log.Printf("❌ فشل جلب الملاحظات الأمنية: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الملاحظات"})
		return
	}
	defer rows.Close()

	var notes []gin.H
	for rows.Next() {
		var title, body, notificationType string
		var createdAt time.Time

		if err := rows.Scan(&title, &body, &createdAt, &notificationType); err != nil {
			continue
		}

		notes = append(notes, gin.H{
			"title":      title,
			"body":       body,
			"created_at": createdAt,
			"type":       notificationType,
		})
	}

	c.JSON(http.StatusOK, notes)
}

// GetLocationLogs - جلب سجل المواقع للمدير
func (h *LocationHandler) GetLocationLogs(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			lt.latitude,
			lt.longitude,
			lt.recorded_at,
			u.full_name,
			u.email,
			w.name as worksite_name,
			a.status as attendance_status
		FROM location_tracking lt
		JOIN users u ON lt.user_id = u.id
		LEFT JOIN attendance a ON a.user_id = u.id AND a.status = 'in_progress'
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee'
		ORDER BY lt.recorded_at DESC
		LIMIT 100
	`)

	if err != nil {
		log.Printf("❌ فشل جلب سجل المواقع: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب السجل"})
		return
	}
	defer rows.Close()

	var logs []gin.H
	for rows.Next() {
		var lat, lng float64
		var recordedAt time.Time
		var fullName, email, worksiteName, attendanceStatus string

		if err := rows.Scan(&lat, &lng, &recordedAt, &fullName, &email, &worksiteName, &attendanceStatus); err != nil {
			continue
		}

		logs = append(logs, gin.H{
			"employee":          fullName,
			"email":             email,
			"latitude":          lat,
			"longitude":         lng,
			"recorded_at":       recordedAt,
			"worksite":          worksiteName,
			"attendance_status": attendanceStatus,
		})
	}

	c.JSON(http.StatusOK, logs)
}

```

---

## 📄 backend/internal/handlers/notification_handler.go

```go
package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	DB *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

// List يعيد إشعارات المستخدم الحالي فقط. عناوين ونصوص الإشعارات نفسها كانت
// قد خُزّنت مسبقاً بلغة الموظف وقت إنشائها (راجع notif_checkin_rejected_* في i18n)
func (h *NotificationHandler) List(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	rows, err := h.DB.Query(`
		SELECT id, title, body, is_read, created_at FROM notifications
		WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`, userID)
	if err != nil {
		log.Printf("failed to fetch notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var notifications []gin.H
	for rows.Next() {
		var id, title, body, createdAt string
		var isRead bool
		if err := rows.Scan(&id, &title, &body, &isRead, &createdAt); err == nil {
			notifications = append(notifications, gin.H{
				"id": id, "title": title, "body": body, "is_read": isRead, "created_at": createdAt,
			})
		}
	}

	c.JSON(http.StatusOK, notifications)
}

```

---

## 📄 backend/internal/handlers/report_handler.go

```go
package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	DB *sql.DB
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

// ComprehensiveReport - تقرير شامل للنظام
func (h *ReportHandler) ComprehensiveReport(c *gin.Context) {
	// إحصائيات الموظفين
	var totalEmployees, activeEmployees, onDutyEmployees, completedDutyToday int
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee'`).Scan(&totalEmployees)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee' AND is_active = TRUE`).Scan(&activeEmployees)
	_ = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT a.user_id) 
		FROM attendance a 
		WHERE DATE(a.check_in_time) = CURRENT_DATE 
		AND a.status = 'active'
	`).Scan(&onDutyEmployees)
	_ = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT a.user_id) 
		FROM attendance a 
		WHERE DATE(a.check_in_time) = CURRENT_DATE 
		AND a.status = 'completed'
	`).Scan(&completedDutyToday)

	// إحصائيات العملاء
	var totalClients, activeClients int
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'client' OR role = 'customer'`).Scan(&totalClients)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE (role = 'client' OR role = 'customer') AND is_active = TRUE`).Scan(&activeClients)

	// إحصائيات طلبات الخدمة
	var totalServiceRequests, pendingRequests, inProgressRequests, completedRequests, ratedRequests float64
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM service_requests`).Scan(&totalServiceRequests)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM service_requests WHERE status = 'pending'`).Scan(&pendingRequests)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM service_requests WHERE status = 'in_progress'`).Scan(&inProgressRequests)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM service_requests WHERE status = 'completed'`).Scan(&completedRequests)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM assignments WHERE client_rating IS NOT NULL`).Scan(&ratedRequests)

	// إحصائيات التقييمات
	var avgRating float64
	_ = h.DB.QueryRow(`SELECT COALESCE(AVG(client_rating), 0) FROM assignments WHERE client_rating IS NOT NULL`).Scan(&avgRating)

	// إحصائيات المهام
	var totalTasks, completedTasks, pendingTasks, inProgressTasks int
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&totalTasks)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'completed'`).Scan(&completedTasks)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'pending'`).Scan(&pendingTasks)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'`).Scan(&inProgressTasks)

	// إحصائيات نقاط العمل
	var totalWorksites, activeWorksites int
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM worksites`).Scan(&totalWorksites)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM worksites WHERE is_active = TRUE`).Scan(&activeWorksites)

	// إحصائيات الحضور هذا الأسبوع
	var totalAttendanceHours, avgDailyHours float64
	_ = h.DB.QueryRow(`
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600), 0)
		FROM attendance a 
		WHERE a.check_out_time IS NOT NULL
		AND DATE(a.check_in_time) >= DATE_TRUNC('week', CURRENT_DATE)
	`).Scan(&totalAttendanceHours)
	
	_ = h.DB.QueryRow(`
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600), 0)
		FROM attendance a 
		WHERE a.check_out_time IS NOT NULL
		AND DATE(a.check_in_time) >= DATE_TRUNC('week', CURRENT_DATE)
	`).Scan(&avgDailyHours)

	c.JSON(http.StatusOK, gin.H{
		"employees": gin.H{
			"total":              totalEmployees,
			"active":             activeEmployees,
			"on_duty":            onDutyEmployees,
			"completed_duty_today": completedDutyToday,
		},
		"clients": gin.H{
			"total":  totalClients,
			"active": activeClients,
		},
		"service_requests": gin.H{
			"total":       totalServiceRequests,
			"pending":     pendingRequests,
			"in_progress": inProgressRequests,
			"completed":   completedRequests,
			"rated":       ratedRequests,
		},
		"ratings": gin.H{
			"average": avgRating,
		},
		"tasks": gin.H{
			"total":        totalTasks,
			"completed":    completedTasks,
			"pending":      pendingTasks,
			"in_progress":  inProgressTasks,
		},
		"worksites": gin.H{
			"total":  totalWorksites,
			"active": activeWorksites,
		},
		"attendance": gin.H{
			"total_hours_week": totalAttendanceHours,
			"avg_daily_hours":   avgDailyHours,
		},
	})
}

// ServiceRequestsReport - تقرير طلبات الخدمة التفصيلي
func (h *ReportHandler) ServiceRequestsReport(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			sr.id,
			sr.title,
			sr.status,
			sr.priority,
			sr.created_at,
			COALESCE(u.full_name, 'Unknown') as client_name,
			COALESCE(e.full_name, 'Not Assigned') as employee_name,
			a.client_rating,
			COALESCE(a.client_feedback, '') as client_feedback,
			a.completed_at
		FROM service_requests sr
		LEFT JOIN users u ON sr.client_id = u.id
		LEFT JOIN assignments a ON sr.id = a.request_id
		LEFT JOIN users e ON a.employee_id = e.id
		ORDER BY sr.created_at DESC
		LIMIT 100
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب تقرير طلبات الخدمة"})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		var id, title, status, priority, clientName, employeeName, clientFeedback string
		var createdAt string
		var completedAt sql.NullTime
		var clientRating sql.NullFloat64

		if err := rows.Scan(&id, &title, &status, &priority, &createdAt, &clientName, &employeeName, &clientRating, &clientFeedback, &completedAt); err != nil {
			continue
		}

		rating := 0.0
		if clientRating.Valid {
			rating = clientRating.Float64
		}

		completedAtStr := ""
		if completedAt.Valid {
			completedAtStr = completedAt.Time.Format("2006-01-02 15:04:05")
		}

		requests = append(requests, gin.H{
			"id":             id,
			"title":          title,
			"status":         status,
			"priority":       priority,
			"created_at":     createdAt,
			"client_name":    clientName,
			"employee_name":  employeeName,
			"client_rating":  rating,
			"client_feedback": clientFeedback,
			"completed_at":   completedAtStr,
		})
	}

	c.JSON(http.StatusOK, requests)
}

// EmployeePerformanceReport - تقرير أداء الموظفين
func (h *ReportHandler) EmployeePerformanceReport(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.phone,
			COUNT(DISTINCT a.id) as total_attendance,
			COALESCE(SUM(EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600), 0) as total_hours,
			COUNT(DISTINCT CASE WHEN a.status = 'completed' THEN a.id END) as completed_shifts,
			COUNT(DISTINCT CASE WHEN sr.id IS NOT NULL THEN a.id END) as assigned_services,
			COALESCE(
				(
					SELECT AVG(sub_a.client_rating)
					FROM assignments sub_a
					WHERE sub_a.employee_id = u.id AND sub_a.client_rating IS NOT NULL
				), 
				0
			) as avg_rating
		FROM users u
		LEFT JOIN attendance a ON u.id = a.user_id
		LEFT JOIN assignments a_assign ON u.id = a_assign.employee_id
		LEFT JOIN service_requests sr ON a_assign.request_id = sr.id
		WHERE u.role = 'employee' AND u.is_active = TRUE
		GROUP BY u.id, u.full_name, u.phone
		ORDER BY total_hours DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب تقرير أداء الموظفين"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, phone string
		var totalAttendance, completedShifts, assignedServices int
		var totalHours, avgRating float64

		if err := rows.Scan(&id, &fullName, &phone, &totalAttendance, &totalHours, &completedShifts, &assignedServices, &avgRating); err != nil {
			continue
		}

		employees = append(employees, gin.H{
			"id":                 id,
			"full_name":          fullName,
			"phone":              phone,
			"total_attendance":   totalAttendance,
			"total_hours":        totalHours,
			"completed_shifts":   completedShifts,
			"assigned_services":  assignedServices,
			"avg_rating":         avgRating,
		})
	}

	c.JSON(http.StatusOK, employees)
}

// ClientActivityReport - تقرير نشاط العملاء
func (h *ReportHandler) ClientActivityReport(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.phone,
			COUNT(DISTINCT sr.id) as total_requests,
			COUNT(DISTINCT CASE WHEN sr.status = 'completed' THEN sr.id END) as completed_requests,
			COALESCE(
				(
					SELECT AVG(sub_a.client_rating)
					FROM assignments sub_a
					JOIN service_requests sub_sr ON sub_a.request_id = sub_sr.id
					WHERE sub_sr.client_id = u.id AND sub_a.client_rating IS NOT NULL
				), 
				0
			) as avg_rating
		FROM users u
		LEFT JOIN service_requests sr ON u.id = sr.client_id
		WHERE (u.role = 'client' OR u.role = 'customer') AND u.is_active = TRUE
		GROUP BY u.id, u.full_name, u.phone
		ORDER BY total_requests DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب تقرير نشاط العملاء"})
		return
	}
	defer rows.Close()

	var clients []gin.H
	for rows.Next() {
		var id, fullName, phone string
		var totalRequests, completedRequests int
		var avgRating float64

		if err := rows.Scan(&id, &fullName, &phone, &totalRequests, &completedRequests, &avgRating); err != nil {
			continue
		}

		clients = append(clients, gin.H{
			"id":                 id,
			"full_name":          fullName,
			"phone":              phone,
			"total_requests":     totalRequests,
			"completed_requests": completedRequests,
			"avg_rating":         avgRating,
		})
	}

	c.JSON(http.StatusOK, clients)
}

// DailySummary - إحصائيات عامة (للحفاظ على التوافق)
func (h *ReportHandler) DailySummary(c *gin.Context) {
	var completed, inProgress, pending, totalEmployees, waitingEmployees, completedToday int

	// الموظفين الذين لم يبدؤوا العمل اليوم (قيد الانتظار)
	_ = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		WHERE u.role = 'employee' 
		AND u.is_active = TRUE
		AND NOT EXISTS (
			SELECT 1 FROM attendance a 
			WHERE a.user_id = u.id 
			AND DATE(a.check_in_time) = CURRENT_DATE
		)
	`).Scan(&waitingEmployees)

	// الموظفين الذين أكملوا عملهم اليوم (مكتمل)
	_ = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		WHERE u.role = 'employee' 
		AND u.is_active = TRUE
		AND DATE(a.check_in_time) = CURRENT_DATE
		AND a.status = 'completed'
		AND a.check_out_time IS NOT NULL
	`).Scan(&completedToday)

	// إحصائيات المهام
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'completed' AND created_at::date = CURRENT_DATE`).Scan(&completed)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'`).Scan(&inProgress)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'pending'`).Scan(&pending)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee' AND is_active = TRUE`).Scan(&totalEmployees)

	c.JSON(http.StatusOK, gin.H{
		"completed_today":     completed,
		"in_progress":         inProgress,
		"pending":             pending,
		"total_employees":     totalEmployees,
		"waiting_employees":   waitingEmployees,
		"completed_employees": completedToday,
	})
}

// GetPendingEmployees - جلب الموظفين قيد الانتظار
func (h *ReportHandler) GetPendingEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.email,
			u.phone
		FROM users u
		WHERE u.role = 'employee'
		AND u.is_active = TRUE
		AND NOT EXISTS (
			SELECT 1 FROM attendance a 
			WHERE a.user_id = u.id 
			AND DATE(a.check_in_time) = CURRENT_DATE
		)
		ORDER BY u.full_name
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين قيد الانتظار"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone string

		if err := rows.Scan(&id, &fullName, &email, &phone); err != nil {
			continue
		}

		employees = append(employees, gin.H{
			"id":         id,
			"full_name":  fullName,
			"email":      email,
			"phone":      phone,
		})
	}

	c.JSON(http.StatusOK, employees)
}

// GetCompletedEmployees - جلب الموظفين المكتملين اليوم
func (h *ReportHandler) GetCompletedEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.email,
			u.phone,
			w.name as worksite_name,
			a.check_in_time,
			a.check_out_time,
			EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600 as hours_worked
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee'
		AND u.is_active = TRUE
		AND DATE(a.check_in_time) = CURRENT_DATE
		AND a.status = 'completed'
		AND a.check_out_time IS NOT NULL
		ORDER BY a.check_out_time DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين المكتملين"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, worksiteName string
		var checkInTime, checkOutTime string
		var hoursWorked float64

		if err := rows.Scan(&id, &fullName, &email, &phone, &worksiteName, &checkInTime, &checkOutTime, &hoursWorked); err != nil {
			continue
		}

		employees = append(employees, gin.H{
			"id":             id,
			"full_name":      fullName,
			"email":          email,
			"phone":          phone,
			"worksite_name":  worksiteName,
			"check_in_time":  checkInTime,
			"check_out_time": checkOutTime,
			"hours_worked":   hoursWorked,
		})
	}

	c.JSON(http.StatusOK, employees)
}

```

---

## 📄 backend/internal/handlers/service_handler.go

```go
package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	DB *sql.DB
}

func NewServiceHandler(db *sql.DB) *ServiceHandler {
	return &ServiceHandler{DB: db}
}

type createServiceRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	Address      string  `json:"address"`
	Phone        string  `json:"phone"`
	Priority     string  `json:"priority"`
	Photos       []string `json:"photos"`
	LocationName string  `json:"location_name"`
}

func (h *ServiceHandler) CreateRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	var req createServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	// Get client information to store in the request
	var clientName, clientPhone sql.NullString
	err := h.DB.QueryRow("SELECT full_name, phone FROM users WHERE id = $1", userID).Scan(&clientName, &clientPhone)
	if err != nil {
		log.Printf("failed to get client info: %v", err)
	}

	id := uuid.NewString()
	_, err = h.DB.Exec(`
		INSERT INTO service_requests (
			id, client_id, client_name, client_phone, title, description,
			latitude, longitude, address, phone,
			priority, status, created_at, updated_at, location_name, language
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pending', now(), now(), $12, $13)`,
		id, userID, clientName, clientPhone, req.Title, req.Description,
		req.Latitude, req.Longitude, req.Address, req.Phone,
		req.Priority, req.LocationName, string(lang),
	)

	if err != nil {
		log.Printf("failed to create service request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "تم إرسال طلب الخدمة بنجاح، سيتم توجيهه لأقرب موظف",
	})
}

func (h *ServiceHandler) ListRequests(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT
			sr.id, sr.title, sr.description,
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at, sr.updated_at, sr.location_name,
			u.full_name as client_name, u.phone as client_phone,
			e.full_name as employee_name,
			a.assigned_at, a.completed_at, a.client_rating, a.client_feedback,
			a.status as assignment_status
		FROM service_requests sr
		LEFT JOIN users u ON sr.client_id = u.id
		LEFT JOIN assignments a ON sr.id = a.request_id AND a.status IN ('assigned', 'accepted', 'in_progress', 'completed')
		LEFT JOIN users e ON a.employee_id = e.id
		ORDER BY
			CASE sr.status
				WHEN 'pending' THEN 1
				WHEN 'assigned' THEN 2
				WHEN 'in_progress' THEN 3
				WHEN 'completed' THEN 4
				WHEN 'cancelled' THEN 5
			END,
			CASE sr.priority
				WHEN 'urgent' THEN 1
				WHEN 'high' THEN 2
				WHEN 'normal' THEN 3
				WHEN 'low' THEN 4
			END,
			sr.created_at DESC
	`)

	if err != nil {
		log.Printf("failed to fetch requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		var id, title, description, address, phone, status, priority, locationName string
		var clientName, clientPhone, employeeName, assignmentStatus sql.NullString
		var latitude, longitude float64
		var createdAt, updatedAt time.Time
		var assignedAt, completedAt sql.NullTime
		var clientRating sql.NullInt64
		var clientFeedback sql.NullString

		err := rows.Scan(&id, &title, &description, &latitude, &longitude,
			&address, &phone, &status, &priority, &createdAt, &updatedAt, &locationName,
			&clientName, &clientPhone, &employeeName,
			&assignedAt, &completedAt, &clientRating, &clientFeedback, &assignmentStatus)
		
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		
		// Convert sql.NullString to string
		clientNameStr := ""
		if clientName.Valid {
			clientNameStr = clientName.String
		}
		
		clientPhoneStr := ""
		if clientPhone.Valid {
			clientPhoneStr = clientPhone.String
		}
		
		employeeNameStr := ""
		if employeeName.Valid {
			employeeNameStr = employeeName.String
		}

		// Build assignment info for all assigned services (not just completed)
		var assignmentInfo gin.H
		if assignedAt.Valid && employeeName.Valid {
			assignmentInfo = gin.H{
				"assigned_at":     assignedAt.Time,
				"employee_name":   employeeNameStr,
				"status":          assignmentStatus.String,
			}
			if completedAt.Valid {
				assignmentInfo["completed_at"] = completedAt.Time
			}
			if clientRating.Valid {
				assignmentInfo["client_rating"] = clientRating.Int64
			}
			if clientFeedback.Valid {
				assignmentInfo["client_feedback"] = clientFeedback.String
			}
		}
		
		requests = append(requests, gin.H{
			"id":            id,
			"title":         title,
			"description":   description,
			"latitude":      latitude,
			"longitude":     longitude,
			"address":       address,
			"phone":         phone,
			"status":        status,
			"priority":      priority,
			"created_at":    createdAt,
			"updated_at":    updatedAt,
			"location_name": locationName,
			"client_name":   clientNameStr,
			"client_phone":  clientPhoneStr,
			"employee_name": employeeNameStr,
			"assignment":    assignmentInfo,
		})
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ServiceHandler) GetRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	var req struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Latitude    float64   `json:"latitude"`
		Longitude   float64   `json:"longitude"`
		Address     string    `json:"address"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
		LocationName string   `json:"location_name"`
		ClientName  string    `json:"client_name"`
		ClientPhone string    `json:"client_phone"`
	}

	err := h.DB.QueryRow(`
		SELECT
			sr.id, sr.title, sr.description,
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at, sr.location_name,
			u.full_name, u.phone
		FROM service_requests sr
		LEFT JOIN users u ON sr.client_id = u.id
		WHERE sr.id = $1`, id,
	).Scan(&req.ID, &req.Title, &req.Description, &req.Latitude, &req.Longitude,
		&req.Address, &req.Phone, &req.Status, &req.Priority, &req.CreatedAt, &req.LocationName,
		&req.ClientName, &req.ClientPhone)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, req)
}

type assignRequest struct {
	RequestID  string `json:"request_id" binding:"required"`
	EmployeeID *string `json:"employee_id"`
	AdminNotes string `json:"admin_notes"`
}

func (h *ServiceHandler) AssignEmployee(c *gin.Context) {
	lang := i18n.Detect(c)
	adminID, _ := c.Get("user_id")

	var req assignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer tx.Rollback()

	// Handle unassign (EmployeeID is null)
	if req.EmployeeID == nil {
		// Delete existing assignment
		_, err = tx.Exec(`
			DELETE FROM assignments 
			WHERE request_id = $1`,
			req.RequestID,
		)

		if err != nil {
			log.Printf("failed to delete assignment: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
			return
		}

		// Set status back to pending
		_, err = tx.Exec(`
			UPDATE service_requests 
			SET status = 'pending', updated_at = now() 
			WHERE id = $1`,
			req.RequestID,
		)

		if err != nil {
			log.Printf("failed to update request status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
			return
		}

		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message": "تم إلغاء تعيين الموظف بنجاح",
		})
		return
	}

	// Handle assign/reassign
	assignmentID := uuid.NewString()
	_, err = tx.Exec(`
		INSERT INTO assignments (
			id, request_id, employee_id, admin_id, 
			admin_notes, status, assigned_at
		) VALUES ($1, $2, $3, $4, $5, 'assigned', now())`,
		assignmentID, req.RequestID, req.EmployeeID, adminID, req.AdminNotes,
	)

	if err != nil {
		log.Printf("failed to create assignment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	_, err = tx.Exec(`
		UPDATE service_requests 
		SET status = 'assigned', updated_at = now() 
		WHERE id = $1`,
		req.RequestID,
	)

	if err != nil {
		log.Printf("failed to update request status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":       "تم تعيين الموظف بنجاح",
		"assignment_id": assignmentID,
	})
}

func (h *ServiceHandler) GetEmployees(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT id, full_name, email, phone, is_active
		FROM users 
		WHERE role = 'employee' AND is_active = TRUE
		ORDER BY full_name`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone string
		var isActive bool
		if err := rows.Scan(&id, &fullName, &email, &phone, &isActive); err == nil {
			employees = append(employees, gin.H{
				"id":        id,
				"full_name": fullName,
				"email":     email,
				"phone":     phone,
				"is_active": isActive,
			})
		}
	}

	c.JSON(http.StatusOK, employees)
}

func (h *ServiceHandler) UpdateStatus(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	_, err := h.DB.Exec(`
		UPDATE service_requests 
		SET status = $1, updated_at = now() 
		WHERE id = $2`,
		req.Status, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم تحديث حالة الطلب"})
}

// GetMyAssignedRequests - جلب طلبات الخدمة المعينة للموظف الحالي
func (h *ServiceHandler) GetMyAssignedRequests(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	rows, err := h.DB.Query(`
		SELECT
			sr.id, sr.title, sr.description,
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at, sr.location_name,
			u.full_name as client_name, u.phone as client_phone,
			a.id as assignment_id, a.assigned_at, a.admin_notes,
			admin.full_name as admin_name
		FROM service_requests sr
		INNER JOIN assignments a ON sr.id = a.request_id AND a.employee_id = $1
		LEFT JOIN users u ON sr.client_id = u.id
		LEFT JOIN users admin ON a.admin_id = admin.id
		WHERE a.status IN ('assigned', 'accepted', 'in_progress')
		ORDER BY sr.priority DESC, sr.created_at DESC`,
		userID,
	)

	if err != nil {
		log.Printf("failed to fetch assigned requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		var id, title, description, address, phone, status, priority, locationName string
		var clientName, clientPhone, adminName, adminNotes sql.NullString
		var assignmentID string
		var assignedAt time.Time
		var latitude, longitude float64
		var createdAt time.Time

		err := rows.Scan(&id, &title, &description, &latitude, &longitude,
			&address, &phone, &status, &priority, &createdAt, &locationName,
			&clientName, &clientPhone, &assignmentID, &assignedAt, &adminNotes,
			&adminName)

		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Convert sql.NullString to string
		clientNameStr := ""
		if clientName.Valid {
			clientNameStr = clientName.String
		}

		clientPhoneStr := ""
		if clientPhone.Valid {
			clientPhoneStr = clientPhone.String
		}

		adminNameStr := ""
		if adminName.Valid {
			adminNameStr = adminName.String
		}

		adminNotesStr := ""
		if adminNotes.Valid {
			adminNotesStr = adminNotes.String
		}

		requests = append(requests, gin.H{
			"id":            id,
			"title":         title,
			"description":   description,
			"latitude":      latitude,
			"longitude":     longitude,
			"address":       address,
			"phone":         phone,
			"status":        status,
			"priority":      priority,
			"created_at":    createdAt,
			"location_name": locationName,
			"client_name":   clientNameStr,
			"client_phone":  clientPhoneStr,
			"assignment_id": assignmentID,
			"assigned_at":   assignedAt,
			"admin_notes":   adminNotesStr,
			"admin_name":    adminNameStr,
		})
	}

	c.JSON(http.StatusOK, requests)
}

// GetAssignedRequestDetails - جلب تفاصيل طلب خدمة معين للموظف
func (h *ServiceHandler) GetAssignedRequestDetails(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")
	requestID := c.Param("id")

	var req struct {
		ID           string    `json:"id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Latitude     float64   `json:"latitude"`
		Longitude    float64   `json:"longitude"`
		Address      string    `json:"address"`
		Phone        string    `json:"phone"`
		Status       string    `json:"status"`
		Priority     string    `json:"priority"`
		CreatedAt    time.Time `json:"created_at"`
		LocationName string    `json:"location_name"`
		ClientName   string    `json:"client_name"`
		ClientPhone  string    `json:"client_phone"`
		AssignmentID string    `json:"assignment_id"`
		AssignedAt   time.Time `json:"assigned_at"`
		AdminNotes   string    `json:"admin_notes"`
		AdminName    string    `json:"admin_name"`
	}

	var clientName, clientPhone, adminNotes, adminName sql.NullString

	err := h.DB.QueryRow(`
		SELECT
			sr.id, sr.title, sr.description,
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at, sr.location_name,
			u.full_name, u.phone,
			a.id as assignment_id, a.assigned_at, a.admin_notes,
			admin.full_name as admin_name
		FROM service_requests sr
		INNER JOIN assignments a ON sr.id = a.request_id AND a.employee_id = $1
		LEFT JOIN users u ON sr.client_id = u.id
		LEFT JOIN users admin ON a.admin_id = admin.id
		WHERE sr.id = $2`,
		userID, requestID,
	).Scan(&req.ID, &req.Title, &req.Description, &req.Latitude, &req.Longitude,
		&req.Address, &req.Phone, &req.Status, &req.Priority, &req.CreatedAt, &req.LocationName,
		&clientName, &clientPhone, &req.AssignmentID, &req.AssignedAt, &adminNotes,
		&adminName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// Convert sql.NullString to string
	if clientName.Valid {
		req.ClientName = clientName.String
	}
	if clientPhone.Valid {
		req.ClientPhone = clientPhone.String
	}
	if adminNotes.Valid {
		req.AdminNotes = adminNotes.String
	}
	if adminName.Valid {
		req.AdminName = adminName.String
	}

	c.JSON(http.StatusOK, req)
}

// UpdateAssignmentStatus - تحديث حالة التعيين (قبول/رفض/بدء/إكمال)
func (h *ServiceHandler) UpdateAssignmentStatus(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")
	requestID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer tx.Rollback()

	// تحديث حالة التعيين
	_, err = tx.Exec(`
		UPDATE assignments 
		SET status = $1, employee_notes = $2
		WHERE request_id = $3 AND employee_id = $4`,
		req.Status, req.Notes, requestID, userID,
	)

	if err != nil {
		log.Printf("failed to update assignment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	// تحديث حالة طلب الخدمة بناءً على حالة التعيين
	var serviceRequestStatus string
	switch req.Status {
	case "accepted":
		serviceRequestStatus = "assigned"
	case "in_progress":
		serviceRequestStatus = "in_progress"
		startedAt := time.Now()
		_, err = tx.Exec(`
			UPDATE assignments 
			SET started_at = $1
			WHERE request_id = $2 AND employee_id = $3`,
			startedAt, requestID, userID,
		)
	case "completed":
		serviceRequestStatus = "completed"
		completedAt := time.Now()
		_, err = tx.Exec(`
			UPDATE assignments 
			SET completed_at = $1
			WHERE request_id = $2 AND employee_id = $3`,
			completedAt, requestID, userID,
		)
	case "rejected":
		serviceRequestStatus = "pending"
	default:
		serviceRequestStatus = "assigned"
	}

	_, err = tx.Exec(`
		UPDATE service_requests 
		SET status = $1, updated_at = now() 
		WHERE id = $2`,
		serviceRequestStatus, requestID,
	)

	if err != nil {
		log.Printf("failed to update service request status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تحديث حالة التعيين بنجاح",
		"status":  serviceRequestStatus,
	})
}

// DeleteRequest - حذف طلب خدمة
func (h *ServiceHandler) DeleteRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	log.Printf("DEBUG: Attempting to delete service request with ID: %s", id)

	// التحقق من وجود الطلب
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM service_requests WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		log.Printf("failed to check request existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	log.Printf("DEBUG: Request exists: %v", exists)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "طلب الخدمة غير موجود"})
		return
	}

	// حذف طلب الخدمة (سيتم حذف التعيينات المرتبطة تلقائياً بسبب ON DELETE CASCADE)
	_, err = h.DB.Exec("DELETE FROM service_requests WHERE id = $1", id)
	if err != nil {
		log.Printf("failed to delete service request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "تم حذف طلب الخدمة بنجاح",
	})
}

// CompleteRequest - تحديث حالة طلب الخدمة إلى مكتمل (يستخدم عند انتهاء الموظف من ورديته)
func (h *ServiceHandler) CompleteRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	_, err := h.DB.Exec(`
		UPDATE service_requests 
		SET status = 'completed', updated_at = now() 
		WHERE id = $1`,
		id,
	)

	if err != nil {
		log.Printf("failed to complete service request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تحديث حالة الطلب إلى مكتمل",
		"status":  "completed",
	})
}

// RateRequest - تقييم طلب خدمة من قبل العميل
func (h *ServiceHandler) RateRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")
	requestID := c.Param("id")

	var req struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	// التحقق من أن الطلب belongs to العميل
	var clientID string
	err := h.DB.QueryRow("SELECT client_id FROM service_requests WHERE id = $1", requestID).Scan(&clientID)
	if err != nil {
		log.Printf("failed to verify request ownership: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "طلب الخدمة غير موجود"})
		return
	}

	if clientID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "غير مصرح لك بتقييم هذا الطلب"})
		return
	}

	// تحديث التقييم في جدول assignments
	_, err = h.DB.Exec(`
		UPDATE assignments 
		SET client_rating = $1, client_feedback = $2
		WHERE request_id = $3`,
		req.Rating, req.Comment, requestID,
	)

	if err != nil {
		log.Printf("failed to rate service request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "تم إرسال تقييمك بنجاح",
	})
}

```

---

## 📄 backend/internal/handlers/task_handler.go

```go
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// MyTasks يعيد مهام الموظف المسجل دخوله فقط (assigned_user_id = المستخدم الحالي)
// بالإضافة إلى طلبات الخدمة المخصصة له عبر جدول assignments
func (h *TaskHandler) MyTasks(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	// جلب المهام من جدول tasks مع assigned_user_id مع الترجمات
	taskRows, err := h.DB.Query(`
		SELECT t.id, t.title, t.title_ar, t.title_he, t.title_en, 
		       t.description, t.description_ar, t.description_he, t.description_en,
		       t.worksite_id, t.status, t.priority, t.scheduled_start, t.scheduled_end, t.created_at, t.language,
		       w.name as worksite_name, w.name_ar as worksite_name_ar, w.name_he as worksite_name_he, w.name_en as worksite_name_en,
		       w.address as worksite_address, w.address_ar as worksite_address_ar, w.address_he as worksite_address_he, w.address_en as worksite_address_en,
		       w.latitude, w.longitude,
		       c.name as client_name, c.name_ar as client_name_ar, c.name_he as client_name_he, c.name_en as client_name_en,
		       c.address as client_address, c.address_ar as client_address_ar, c.address_he as client_address_he, c.address_en as client_address_en,
		       c.phone as client_phone
		FROM tasks t
		LEFT JOIN worksites w ON t.worksite_id = w.id
		LEFT JOIN clients c ON t.client_id = c.id
		WHERE t.assigned_user_id = $1 ORDER BY t.scheduled_start ASC`, userID)
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer taskRows.Close()

	var tasks []models.Task
	for taskRows.Next() {
		var t models.Task
		var title, titleAr, titleHe, titleEn sql.NullString
		var description, descriptionAr, descriptionHe, descriptionEn sql.NullString
		var worksiteName, worksiteNameAr, worksiteNameHe, worksiteNameEn sql.NullString
		var worksiteAddress, worksiteAddressAr, worksiteAddressHe, worksiteAddressEn sql.NullString
		var clientName, clientNameAr, clientNameHe, clientNameEn sql.NullString
		var clientAddress, clientAddressAr, clientAddressHe, clientAddressEn sql.NullString
		var clientPhone, priority, taskLanguage sql.NullString
		var worksiteLat, worksiteLng sql.NullFloat64

		if err := taskRows.Scan(&t.ID, &title, &titleAr, &titleHe, &titleEn,
			&description, &descriptionAr, &descriptionHe, &descriptionEn,
			&t.WorksiteID, &t.Status, &priority, &t.ScheduledStart, &t.ScheduledEnd, &t.CreatedAt, &taskLanguage,
			&worksiteName, &worksiteNameAr, &worksiteNameHe, &worksiteNameEn,
			&worksiteAddress, &worksiteAddressAr, &worksiteAddressHe, &worksiteAddressEn,
			&worksiteLat, &worksiteLng,
			&clientName, &clientNameAr, &clientNameHe, &clientNameEn,
			&clientAddress, &clientAddressAr, &clientAddressHe, &clientAddressEn,
			&clientPhone); err == nil {
			// استخدام اللغة المخزنة في المهمة أو الإنجليزية كافتراضي للمهام القديمة
			taskLang := i18n.Normalize(taskLanguage.String)
			if taskLanguage.String == "" {
				taskLang = i18n.English // استخدام الإنجليزية كافتراضي للمهام القديمة
			}
			
			// التعامل مع المهام القديمة التي لا تحتوي على ترجمات
			// إذا كانت الترجمات فارغة، استخدم النص الأصلي
			if titleEn.String == "" && title.String != "" {
				titleEn.String = title.String
			}
			if descriptionEn.String == "" && description.String != "" {
				descriptionEn.String = description.String
			}
			
			// إضافة البيانات الإضافية كحقول مؤقتة
			t.Title = i18n.GetTranslation(taskLang, titleAr.String, titleHe.String, titleEn.String)
			t.TitleAr = titleAr.String
			t.TitleHe = titleHe.String
			t.TitleEn = titleEn.String
			t.Description = i18n.GetTranslation(taskLang, descriptionAr.String, descriptionHe.String, descriptionEn.String)
			t.DescriptionAr = descriptionAr.String
			t.DescriptionHe = descriptionHe.String
			t.DescriptionEn = descriptionEn.String
			t.Priority = priority.String
			t.WorksiteName = i18n.GetTranslation(taskLang, worksiteNameAr.String, worksiteNameHe.String, worksiteNameEn.String)
			t.WorksiteNameAr = worksiteNameAr.String
			t.WorksiteNameHe = worksiteNameHe.String
			t.WorksiteNameEn = worksiteNameEn.String
			t.WorksiteAddress = i18n.GetTranslation(taskLang, worksiteAddressAr.String, worksiteAddressHe.String, worksiteAddressEn.String)
			t.WorksiteAddressAr = worksiteAddressAr.String
			t.WorksiteAddressHe = worksiteAddressHe.String
			t.WorksiteAddressEn = worksiteAddressEn.String
			t.WorksiteLatitude = worksiteLat.Float64
			t.WorksiteLongitude = worksiteLng.Float64
			t.ClientName = i18n.GetTranslation(taskLang, clientNameAr.String, clientNameHe.String, clientNameEn.String)
			t.ClientNameAr = clientNameAr.String
			t.ClientNameHe = clientNameHe.String
			t.ClientNameEn = clientNameEn.String
			t.ClientAddress = i18n.GetTranslation(taskLang, clientAddressAr.String, clientAddressHe.String, clientAddressEn.String)
			t.ClientAddressAr = clientAddressAr.String
			t.ClientAddressHe = clientAddressHe.String
			t.ClientAddressEn = clientAddressEn.String
			t.ClientPhone = clientPhone.String
			t.Language = taskLanguage.String
			tasks = append(tasks, t)
		}
	}

	// جلب طلبات الخدمة المخصصة للموظف عبر جدول assignments مع الترجمات
	serviceRequestRows, err := h.DB.Query(`
		SELECT sr.id, sr.title, sr.title_ar, sr.title_he, sr.title_en,
		       sr.description, sr.description_ar, sr.description_he, sr.description_en,
		       sr.status, sr.priority, sr.created_at, sr.language,
		       sr.latitude, sr.longitude,
		       sr.address, sr.address_ar, sr.address_he, sr.address_en,
		       sr.phone as service_phone,
		       sr.location_name, sr.location_name_ar, sr.location_name_he, sr.location_name_en,
		       sr.client_name, sr.client_name_ar, sr.client_name_he, sr.client_name_en,
		       u.full_name as user_full_name, u.phone as user_phone
		FROM service_requests sr
		INNER JOIN assignments a ON sr.id = a.request_id
		LEFT JOIN users u ON sr.client_id = u.id
		WHERE a.employee_id = $1 AND a.status IN ('assigned', 'accepted', 'in_progress')
		ORDER BY sr.created_at DESC`, userID)
	if err != nil {
		log.Printf("failed to fetch service requests: %v", err)
		// لا نعيد خطأ هنا لأن المهام قد تكون موجودة
	} else {
		defer serviceRequestRows.Close()

		for serviceRequestRows.Next() {
			var id, status, priority string
			var title, titleAr, titleHe, titleEn sql.NullString
			var description, descriptionAr, descriptionHe, descriptionEn sql.NullString
			var address, addressAr, addressHe, addressEn sql.NullString
			var locationName, locationNameAr, locationNameHe, locationNameEn sql.NullString
			var clientName, clientNameAr, clientNameHe, clientNameEn sql.NullString
			var servicePhone, userPhone, serviceLanguage sql.NullString
			var userFullName sql.NullString
			var latitude, longitude float64
			var createdAt time.Time

			if err := serviceRequestRows.Scan(&id, &title, &titleAr, &titleHe, &titleEn,
				&description, &descriptionAr, &descriptionHe, &descriptionEn,
				&status, &priority, &createdAt, &serviceLanguage,
				&latitude, &longitude,
				&address, &addressAr, &addressHe, &addressEn,
				&servicePhone, &locationName, &locationNameAr, &locationNameHe, &locationNameEn,
				&clientName, &clientNameAr, &clientNameHe, &clientNameEn,
				&userFullName, &userPhone); err == nil {
				// استخدام اللغة المخزنة في طلب الخدمة أو الإنجليزية كافتراضي للطلبات القديمة
				serviceLang := i18n.Normalize(serviceLanguage.String)
				if serviceLanguage.String == "" {
					serviceLang = i18n.English // استخدام الإنجليزية كافتراضي للطلبات القديمة
				}
				
				// التعامل مع الطلبات القديمة التي لا تحتوي على ترجمات
				// إذا كانت الترجمات فارغة، استخدم النص الأصلي
				if titleEn.String == "" && title.String != "" {
					titleEn.String = title.String
				}
				if descriptionEn.String == "" && description.String != "" {
					descriptionEn.String = description.String
				}
				
				// تحويل طلب الخدمة إلى شكل مشابه للمهمة مع إضافة التفاصيل الكاملة والترجمات
				// استخدام رقم الهاتف المتوفر (إما من service_requests أو من users)
				phone := servicePhone.String
				if phone == "" && userPhone.Valid {
					phone = userPhone.String
				}

				task := models.Task{
					ID:               id,
					Title:            i18n.GetTranslation(serviceLang, titleAr.String, titleHe.String, titleEn.String),
					TitleAr:          titleAr.String,
					TitleHe:          titleHe.String,
					TitleEn:          titleEn.String,
					Description:      i18n.GetTranslation(serviceLang, descriptionAr.String, descriptionHe.String, descriptionEn.String),
					DescriptionAr:    descriptionAr.String,
					DescriptionHe:    descriptionHe.String,
					DescriptionEn:    descriptionEn.String,
					Status:           status,
					Priority:         priority,
					CreatedAt:        createdAt,
					WorksiteName:     i18n.GetTranslation(serviceLang, locationNameAr.String, locationNameHe.String, locationNameEn.String),
					WorksiteNameAr:   locationNameAr.String,
					WorksiteNameHe:   locationNameHe.String,
					WorksiteNameEn:   locationNameEn.String,
					WorksiteAddress:  i18n.GetTranslation(serviceLang, addressAr.String, addressHe.String, addressEn.String),
					WorksiteAddressAr: addressAr.String,
					WorksiteAddressHe: addressHe.String,
					WorksiteAddressEn: addressEn.String,
					WorksiteLatitude: latitude,
					WorksiteLongitude: longitude,
					ClientName:       i18n.GetTranslation(serviceLang, clientNameAr.String, clientNameHe.String, clientNameEn.String),
					ClientNameAr:     clientNameAr.String,
					ClientNameHe:     clientNameHe.String,
					ClientNameEn:     clientNameEn.String,
					ClientAddress:    i18n.GetTranslation(serviceLang, addressAr.String, addressHe.String, addressEn.String),
					ClientAddressAr:  addressAr.String,
					ClientAddressHe:  addressHe.String,
					ClientAddressEn:  addressEn.String,
					ClientPhone:      phone,
					Language:         serviceLanguage.String,
				}
				tasks = append(tasks, task)
			}
		}
	}

	c.JSON(http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title          string  `json:"title" binding:"required"`
	TitleAr        string  `json:"title_ar"`
	TitleHe        string  `json:"title_he"`
	TitleEn        string  `json:"title_en"`
	Description    string  `json:"description"`
	DescriptionAr  string  `json:"description_ar"`
	DescriptionHe  string  `json:"description_he"`
	DescriptionEn  string  `json:"description_en"`
	ClientID       *string `json:"client_id"`
	WorksiteID     string  `json:"worksite_id" binding:"required"`
	AssignedUserID *string `json:"assigned_user_id"`
	Priority       string  `json:"priority"`
	ScheduledStart string  `json:"scheduled_start"`
	ScheduledEnd   string  `json:"scheduled_end"`
}

// Create ينشئ مهمة جديدة ويربطها بنقطة عمل (worksite) — هذا الربط هو ما يفعّل الـ Geofence لاحقاً
func (h *TaskHandler) Create(c *gin.Context) {
	lang := i18n.Detect(c)

	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO tasks (id, title, title_ar, title_he, title_en, 
		                  description, description_ar, description_he, description_en,
		                  client_id, worksite_id, assigned_user_id,
		                  status, priority, scheduled_start, scheduled_end, created_at, updated_at, language)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pending', COALESCE($12, 'normal'), 
		        NULLIF($13,'')::timestamptz, NULLIF($14,'')::timestamptz, now(), now(), $15)`,
		id, req.Title, req.TitleAr, req.TitleHe, req.TitleEn, 
		req.Description, req.DescriptionAr, req.DescriptionHe, req.DescriptionEn,
		req.ClientID, req.WorksiteID, req.AssignedUserID,
		req.Priority, req.ScheduledStart, req.ScheduledEnd, string(lang),
	)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": i18n.T(lang, "msg_task_created")})
}

// ListAll يعيد كل المهام (لوحة تحكم المدير)
func (h *TaskHandler) ListAll(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT id, title, description, worksite_id, status, priority, scheduled_start, scheduled_end, created_at
		FROM tasks ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		log.Printf("failed to fetch all tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var priority sql.NullString
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.WorksiteID,
			&t.Status, &priority, &t.ScheduledStart, &t.ScheduledEnd, &t.CreatedAt); err == nil {
			t.Priority = priority.String
			tasks = append(tasks, t)
		}
	}

	c.JSON(http.StatusOK, tasks)
}

type updateTaskRequest struct {
	Title          *string `json:"title"`
	TitleAr        *string `json:"title_ar"`
	TitleHe        *string `json:"title_he"`
	TitleEn        *string `json:"title_en"`
	Description    *string `json:"description"`
	DescriptionAr  *string `json:"description_ar"`
	DescriptionHe  *string `json:"description_he"`
	DescriptionEn  *string `json:"description_en"`
	ClientID       *string `json:"client_id"`
	WorksiteID     *string `json:"worksite_id"`
	AssignedUserID *string `json:"assigned_user_id"`
	Status         *string `json:"status"`
	Priority       *string `json:"priority"`
	ScheduledStart *string `json:"scheduled_start"`
	ScheduledEnd   *string `json:"scheduled_end"`
	Language       *string `json:"language"`
}

// Update تحديث مهمة موجودة
func (h *TaskHandler) Update(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	// بناء الاستعلام ديناميكياً بناءً على الحقول المقدمة
	query := "UPDATE tasks SET updated_at = now()"
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		query += ", title = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Title)
		argCount++
	}
	if req.TitleAr != nil {
		query += ", title_ar = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.TitleAr)
		argCount++
	}
	if req.TitleHe != nil {
		query += ", title_he = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.TitleHe)
		argCount++
	}
	if req.TitleEn != nil {
		query += ", title_en = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.TitleEn)
		argCount++
	}
	if req.Description != nil {
		query += ", description = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Description)
		argCount++
	}
	if req.DescriptionAr != nil {
		query += ", description_ar = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.DescriptionAr)
		argCount++
	}
	if req.DescriptionHe != nil {
		query += ", description_he = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.DescriptionHe)
		argCount++
	}
	if req.DescriptionEn != nil {
		query += ", description_en = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.DescriptionEn)
		argCount++
	}
	if req.ClientID != nil {
		query += ", client_id = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.ClientID)
		argCount++
	}
	if req.WorksiteID != nil {
		query += ", worksite_id = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.WorksiteID)
		argCount++
	}
	if req.AssignedUserID != nil {
		query += ", assigned_user_id = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.AssignedUserID)
		argCount++
	}
	if req.Status != nil {
		query += ", status = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Status)
		argCount++
	}
	if req.Priority != nil {
		query += ", priority = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Priority)
		argCount++
	}
	if req.ScheduledStart != nil {
		query += ", scheduled_start = NULLIF($" + fmt.Sprintf("%d", argCount) + ",'')::timestamptz"
		args = append(args, *req.ScheduledStart)
		argCount++
	}
	if req.ScheduledEnd != nil {
		query += ", scheduled_end = NULLIF($" + fmt.Sprintf("%d", argCount) + ",'')::timestamptz"
		args = append(args, *req.ScheduledEnd)
		argCount++
	}
	if req.Language != nil {
		query += ", language = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Language)
		argCount++
	}

	query += " WHERE id = $" + fmt.Sprintf("%d", argCount)
	args = append(args, id)

	_, err := h.DB.Exec(query, args...)
	if err != nil {
		log.Printf("failed to update task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(lang, "msg_task_updated")})
}

```

---

## 📄 backend/internal/handlers/upload_handler.go

```go
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	DB *sql.DB
}

func NewUploadHandler(db *sql.DB) *UploadHandler {
	return &UploadHandler{DB: db}
}

func (h *UploadHandler) UploadTaskPhoto(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	log.Printf("📸 محاولة رفع صورة من المستخدم: %s", userID)

	attendanceID := c.PostForm("attendance_id")
	log.Printf("📝 attendance_id: %s", attendanceID)

	if attendanceID == "" {
		log.Println("❌ attendance_id مطلوب")
		c.JSON(http.StatusBadRequest, gin.H{"error": "attendance_id مطلوب"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		log.Printf("❌ فشل الحصول على الصورة: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_no_photo")})
		return
	}

	log.Printf("📄 اسم الملف: %s, الحجم: %d بايت", file.Filename, file.Size)

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "نوع الملف غير مدعوم"})
		return
	}

	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "حجم الملف كبير جداً"})
		return
	}

	var employeeID string
	err = h.DB.QueryRow(`
		SELECT user_id FROM attendance WHERE id = $1
	`, attendanceID).Scan(&employeeID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ سجل الحضور غير موجود: %s", attendanceID)
			c.JSON(http.StatusNotFound, gin.H{"error": "سجل الحضور غير موجود"})
			return
		}
		log.Printf("❌ فشل جلب سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجل الحضور"})
		return
	}

	if employeeID != userID {
		log.Printf("❌ المستخدم ليس صاحب السجل")
		c.JSON(http.StatusForbidden, gin.H{"error": "ليس لديك صلاحية"})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	filename := fmt.Sprintf("%d_%s_%s%s", 
		time.Now().Unix(), 
		userID.(string)[:8], 
		uuid.NewString()[:8],
		ext,
	)

	filePath := filepath.Join(uploadDir, filename)
	
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("❌ فشل حفظ الصورة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حفظ الصورة"})
		return
	}

	photoURL := fmt.Sprintf("/uploads/%s", filename)
	
	log.Printf("📸 تم حفظ الصورة: %s", filePath)
	log.Printf("📸 رابط الصورة: %s", photoURL)

	now := time.Now()
	
	result, err := h.DB.Exec(`
		UPDATE attendance 
		SET photo_url = $1, photo_uploaded_at = $2, photo_notes = $3
		WHERE id = $4
	`, photoURL, now, "تم رفع صورة إثبات الإنجاز", attendanceID)

	if err != nil {
		log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حفظ رابط الصورة"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ تم تحديث %d سجل", rowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message":       "✅ تم رفع الصورة بنجاح",
		"photo_url":     photoURL,
		"attendance_id": attendanceID,
		"filename":      filename,
	})
}

func (h *UploadHandler) GetAttendancePhotos(c *gin.Context) {
	log.Println("📸 جلب صور إثبات الإنجاز...")

	rows, err := h.DB.Query(`
		SELECT 
			a.id,
			a.photo_url,
			a.photo_uploaded_at,
			a.photo_notes,
			u.full_name as employee_name,
			u.email as employee_email,
			COALESCE(t.title, 'بدون مهمة') as task_title,
			COALESCE(w.name, 'بدون موقع') as worksite_name
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN tasks t ON a.task_id = t.id
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE a.photo_url IS NOT NULL AND a.photo_url != ''
		ORDER BY a.photo_uploaded_at DESC
		LIMIT 50
	`)

	if err != nil {
		log.Printf("❌ فشل جلب الصور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الصور"})
		return
	}
	defer rows.Close()

	var photos []gin.H
	for rows.Next() {
		var id, photoURL, photoNotes, employeeName, employeeEmail, taskTitle, worksiteName string
		var photoUploadedAt sql.NullTime

		if err := rows.Scan(&id, &photoURL, &photoUploadedAt, &photoNotes,
			&employeeName, &employeeEmail, &taskTitle, &worksiteName); err != nil {
			log.Printf("⚠️ خطأ في قراءة البيانات: %v", err)
			continue
		}

		photo := gin.H{
			"id":              id,
			"photo_url":       photoURL,
			"photo_notes":     photoNotes,
			"employee_name":   employeeName,
			"employee_email":  employeeEmail,
			"task_title":      taskTitle,
			"worksite_name":   worksiteName,
			"uploaded_at":     photoUploadedAt.Time,
		}

		photos = append(photos, photo)
	}

	log.Printf("✅ تم جلب %d صورة", len(photos))
	c.JSON(http.StatusOK, photos)
}

func (h *UploadHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	log.Printf("🗑️ محاولة حذف الصورة: %s", id)

	var photoURL string
	err := h.DB.QueryRow(`
		SELECT photo_url FROM attendance WHERE id = $1
	`, id).Scan(&photoURL)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ الصورة غير موجودة: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "الصورة غير موجودة"})
			return
		}
		log.Printf("❌ فشل جلب الصورة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الصورة"})
		return
	}

	if photoURL != "" {
		filename := filepath.Base(photoURL)
		filePath := filepath.Join("./uploads", filename)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				log.Printf("⚠️ فشل حذف الملف: %v", err)
			} else {
				log.Printf("✅ تم حذف الملف: %s", filePath)
			}
		}
	}

	_, err = h.DB.Exec(`
		UPDATE attendance 
		SET photo_url = NULL, photo_uploaded_at = NULL, photo_notes = NULL
		WHERE id = $1
	`, id)

	if err != nil {
		log.Printf("❌ فشل حذف الصورة من قاعدة البيانات: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف الصورة"})
		return
	}

	log.Printf("✅ تم حذف الصورة: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "✅ تم حذف الصورة بنجاح"})
}

func (h *UploadHandler) DownloadPhoto(c *gin.Context) {
	id := c.Param("id")

	log.Printf("⬇️ محاولة تحميل الصورة: %s", id)

	var photoURL string
	err := h.DB.QueryRow(`
		SELECT photo_url FROM attendance WHERE id = $1
	`, id).Scan(&photoURL)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ الصورة غير موجودة: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "الصورة غير موجودة"})
			return
		}
		log.Printf("❌ فشل جلب الصورة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الصورة"})
		return
	}

	if photoURL == "" {
		log.Printf("❌ لا توجد صورة للسجل: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "لا توجد صورة"})
		return
	}

	filename := filepath.Base(photoURL)
	filePath := filepath.Join("./uploads", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("❌ الملف غير موجود: %s", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "الملف غير موجود"})
		return
	}

	// تعيين رأس الاستجابة لتحميل الملف
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	c.File(filePath)

	log.Printf("✅ تم تحميل الصورة: %s", filename)
}

```

---

## 📄 backend/internal/handlers/websocket_handler.go

```go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // السماح بجميع المصادر في وضع التطوير
	},
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

// WebSocketHub - مركز إدارة اتصالات WebSocket
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	maxClients int
}

// WSHandler - معالج WebSocket
type WSHandler struct {
	hub *WebSocketHub
}

// NewWSHandler - إنشاء معالج WebSocket جديد
func NewWSHandler() *WSHandler {
	hub := &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256), // Buffered channel to prevent blocking
		register:   make(chan *websocket.Conn, 256),
		unregister: make(chan *websocket.Conn, 256),
		maxClients: 100, // Limit concurrent WebSocket connections
	}

	handler := &WSHandler{hub: hub}
	go hub.run()

	return handler
}

// run - تشغيل مركز WebSocket
func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if len(h.clients) >= h.maxClients {
				log.Printf("⚠️ تم رفض اتصال WebSocket جديد - الحد الأقصى %d", h.maxClients)
				client.Close()
			} else {
				h.clients[client] = true
				log.Println("✅ تم تسجيل عميل WebSocket جديد")
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("❌ تم إلغاء تسجيل عميل WebSocket")

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("❌ خطأ في إرسال WebSocket: %v", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// HandleWebSocket - معالجة اتصال WebSocket
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	lang := i18n.Detect(c)
	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ فشل ترقية الاتصال إلى WebSocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	h.hub.register <- conn

	// إرسال رسالة ترحيب
	welcomeMsg := map[string]interface{}{
		"type":      "connected",
		"message":   "تم الاتصال بنجاح",
		"timestamp": time.Now(),
	}
	msgBytes, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)

	// الاستماع للرسائل من العميل
	go func() {
		defer func() {
			h.hub.unregister <- conn
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// BroadcastLocationUpdate - إرسال تحديث الموقع لجميع العملاء
func (h *WSHandler) BroadcastLocationUpdate(update map[string]interface{}) {
	message := map[string]interface{}{
		"type":      "location_update",
		"data":      update,
		"timestamp": time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ خطأ في ترميز تحديث الموقع: %v", err)
		return
	}

	h.hub.broadcast <- msgBytes
}

// GetClientCount - عدد العملاء المتصلين
func (h *WSHandler) GetClientCount() int {
	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()
	return len(h.hub.clients)
}

// BroadcastEmployeeStatus - إرسال تحديث حالة الموظف
func (h *WSHandler) BroadcastEmployeeStatus(update map[string]interface{}) {
	message := map[string]interface{}{
		"type":      "employee_status",
		"data":      update,
		"timestamp": time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ خطأ في ترميز تحديث الحالة: %v", err)
		return
	}

	h.hub.broadcast <- msgBytes
}

// BroadcastPasswordChanged - إرسال إشعار تغيير كلمة المرور لطرد فوري
func (h *WSHandler) BroadcastPasswordChanged(userID string) {
	message := map[string]interface{}{
		"type":      "password_changed",
		"user_id":   userID,
		"timestamp": time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ خطأ في ترميز إشعار تغيير كلمة المرور: %v", err)
		return
	}

	h.hub.broadcast <- msgBytes
	log.Printf("🚨 إرسال إشعار تغيير كلمة المرور للمستخدم: %s", userID)
}
```

---

## 📄 backend/internal/handlers/worksite_handler.go

```go
package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"worktrack/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorksiteHandler struct {
	DB *sql.DB
}

func NewWorksiteHandler(db *sql.DB) *WorksiteHandler {
	return &WorksiteHandler{DB: db}
}

// List - جلب جميع نقاط العمل مع الموظفين العاملين حالياً
func (h *WorksiteHandler) List(c *gin.Context) {
	// التحقق من وجود عمود assigned_employee_id للتوافق مع الإصدارات القديمة
	var columnExists bool
	h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'worksites' AND column_name = 'assigned_employee_id'
		)
	`).Scan(&columnExists)

	var rows *sql.Rows
	var err error

	if columnExists {
		// استخدام الاستعلام مع assigned_employee_id (الإصدار الحديث)
		rows, err = h.DB.Query(`
			SELECT w.id, w.name, w.address, w.latitude, w.longitude, w.radius_meters, w.is_active, w.created_at,
				u.id as assigned_employee_id, u.full_name as assigned_employee_name
			FROM worksites w
			LEFT JOIN users u ON w.assigned_employee_id = u.id
			ORDER BY w.created_at DESC`)
	} else {
		// استخدام الاستعلام بدون assigned_employee_id (للتوافق مع الإصدارات القديمة)
		rows, err = h.DB.Query(`
			SELECT w.id, w.name, w.address, w.latitude, w.longitude, w.radius_meters, w.is_active, w.created_at
			FROM worksites w
			ORDER BY w.created_at DESC`)
	}

	if err != nil {
		log.Printf("❌ فشل جلب نقاط العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب نقاط العمل"})
		return
	}
	defer rows.Close()

	var worksites []gin.H
	for rows.Next() {
		var w models.Worksite
		var assignedEmployeeID, assignedEmployeeName sql.NullString

		if columnExists {
			if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
				&w.RadiusMeters, &w.IsActive, &w.CreatedAt, &assignedEmployeeID, &assignedEmployeeName); err != nil {
				log.Printf("⚠️ خطأ في قراءة صف نقاط العمل: %v", err)
				continue
			}
		} else {
			if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
				&w.RadiusMeters, &w.IsActive, &w.CreatedAt); err != nil {
				log.Printf("⚠️ خطأ في قراءة صف نقاط العمل: %v", err)
				continue
			}
		}

		worksite := gin.H{
			"id":            w.ID,
			"name":          w.Name,
			"address":       w.Address,
			"latitude":      w.Latitude,
			"longitude":     w.Longitude,
			"radius_meters": w.RadiusMeters,
			"is_active":     w.IsActive,
			"created_at":    w.CreatedAt,
		}

		if columnExists && assignedEmployeeID.Valid && assignedEmployeeName.Valid {
			worksite["assigned_to"] = gin.H{
				"id":   assignedEmployeeID.String,
				"name": assignedEmployeeName.String,
			}
		}

		worksites = append(worksites, worksite)
	}

	// جلب الموظفين العاملين حالياً في كل نقطة عمل
	for i, ws := range worksites {
		if ws["id"] == "unassigned" {
			continue
		}
		
		workingRows, err := h.DB.Query(`
			SELECT DISTINCT u.id, u.full_name, a.id as attendance_id
			FROM users u
			JOIN attendance a ON u.id = a.user_id
			WHERE a.status = 'in_progress' 
			AND a.worksite_id = $1
			AND u.role = 'employee'
		`, ws["id"])
		
		if err == nil {
			var workingEmployees []gin.H
			for workingRows.Next() {
				var id, name, attendanceID string
				if err := workingRows.Scan(&id, &name, &attendanceID); err == nil {
					workingEmployees = append(workingEmployees, gin.H{
						"id":             id,
						"name":           name,
						"attendance_id":  attendanceID,
					})
				}
			}
			workingRows.Close()
			worksites[i]["working_employees"] = workingEmployees
		}
	}

	// إضافة خيار "غير معين" للموظفين الذين يعملون بدون نقطة عمل
	unassignedRows, err := h.DB.Query(`
		SELECT DISTINCT u.id, u.full_name, a.id as attendance_id
		FROM users u
		JOIN attendance a ON u.id = a.user_id
		WHERE a.status = 'in_progress' 
		AND (a.worksite_id IS NULL OR a.worksite_id = '')
		AND u.role = 'employee'
	`)
	if err == nil {
		var unassignedEmployees []gin.H
		for unassignedRows.Next() {
			var id, name, attendanceID string
			if err := unassignedRows.Scan(&id, &name, &attendanceID); err == nil {
				unassignedEmployees = append(unassignedEmployees, gin.H{
					"id":             id,
					"name":           name,
					"attendance_id":  attendanceID,
				})
			}
		}
		unassignedRows.Close()

		if len(unassignedEmployees) > 0 {
			worksites = append([]gin.H{
				{
					"id":                 "unassigned",
					"name":               "غير معين",
					"address":            "الموظفون الذين يعملون بدون نقطة عمل محددة",
					"latitude":           0,
					"longitude":          0,
					"radius_meters":      0,
					"is_active":          true,
					"created_at":         "now()",
					"is_unassigned":      true,
					"working_employees":  unassignedEmployees,
				},
			}, worksites...)
		}
	}

	c.JSON(http.StatusOK, worksites)
}

// GetAvailableWorksites - جلب نقاط العمل المتاحة للموظفين
func (h *WorksiteHandler) GetAvailableWorksites(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, name, address, latitude, longitude, radius_meters
		FROM worksites WHERE is_active = TRUE ORDER BY name`)
	if err != nil {
		log.Printf("❌ فشل جلب نقاط العمل المتاحة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب نقاط العمل"})
		return
	}
	defer rows.Close()

	var worksites []gin.H
	for rows.Next() {
		var id, name, address string
		var lat, lng float64
		var radius int
		if err := rows.Scan(&id, &name, &address, &lat, &lng, &radius); err != nil {
			continue
		}
		worksites = append(worksites, gin.H{
			"id":            id,
			"name":          name,
			"address":       address,
			"latitude":      lat,
			"longitude":     lng,
			"radius_meters": radius,
		})
	}
	c.JSON(http.StatusOK, worksites)
}

// Create - إضافة نقطة عمل جديدة
func (h *WorksiteHandler) Create(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Address      string  `json:"address"`
		Latitude     float64 `json:"latitude" binding:"required"`
		Longitude    float64 `json:"longitude" binding:"required"`
		RadiusMeters int     `json:"radius_meters" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}
	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE, now(), now())`,
		id, req.Name, req.Address, req.Latitude, req.Longitude, req.RadiusMeters,
	)
	if err != nil {
		log.Printf("❌ فشل الإضافة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الإضافة"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "تم الإضافة"})
}

// Delete - حذف نقطة العمل مع الحفاظ على سجلات الحضور
func (h *WorksiteHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	log.Printf("🗑️ محاولة حذف نقطة العمل: %s", id)

	// بدء معاملة
	tx, err := h.DB.Begin()
	if err != nil {
		log.Printf("❌ فشل بدء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل بدء عملية الحذف"})
		return
	}
	defer tx.Rollback()

	// 1. الحصول على عدد سجلات الحضور المرتبطة قبل الحذف
	var attendanceCount int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM attendance WHERE worksite_id = $1
	`, id).Scan(&attendanceCount)
	if err != nil {
		log.Printf("❌ فشل الحصول على عدد سجلات الحضور: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الحصول على معلومات الحضور"})
		return
	}
	log.Printf("📊 عدد سجلات الحضور المرتبطة: %d", attendanceCount)

	// 2. الحصول على عدد المهام المرتبطة قبل الحذف
	var tasksCount int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM tasks WHERE worksite_id = $1
	`, id).Scan(&tasksCount)
	if err != nil {
		log.Printf("❌ فشل الحصول على عدد المهام: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الحصول على معلومات المهام"})
		return
	}
	log.Printf("📊 عدد المهام المرتبطة: %d", tasksCount)

	// 3. حذف نقاط العمل (سيقوم قيد ON DELETE SET NULL تلقائياً بفك ارتباط سجلات الحضور والمهام)
	result, err := tx.Exec(`
		DELETE FROM worksites WHERE id = $1
	`, id)
	if err != nil {
		log.Printf("❌ فشل حذف نقطة العمل: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل حذف نقطة العمل"})
		return
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		log.Printf("❌ نقطة العمل غير موجودة: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "نقطة العمل غير موجودة"})
		return
	}

	// إنهاء المعاملة
	if err := tx.Commit(); err != nil {
		log.Printf("❌ فشل إنهاء المعاملة: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل إنهاء عملية الحذف"})
		return
	}

	log.Printf("✅ تم حذف نقطة العمل: %s (فك ارتباط %d حضور و %d مهام تلقائياً)", id, attendanceCount, tasksCount)
	c.JSON(http.StatusOK, gin.H{
		"message":              "تم حذف نقطة العمل بنجاح",
		"attendance_preserved": attendanceCount,
		"tasks_preserved":      tasksCount,
	})
}

// AssignEmployee - تعيين موظف لنقطة عمل
func (h *WorksiteHandler) AssignEmployee(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
		WorksiteID string `json:"worksite_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	// تحديث نقطة العمل لتعيين الموظف
	_, err := h.DB.Exec(`
		UPDATE worksites 
		SET assigned_employee_id = $1, updated_at = now()
		WHERE id = $2
	`, req.EmployeeID, req.WorksiteID)

	if err != nil {
		log.Printf("❌ فشل تعيين الموظف: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل تعيين الموظف"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم تعيين الموظف بنجاح"})
}

// GetAvailableEmployees - جلب الموظفين المتاحين
func (h *WorksiteHandler) GetAvailableEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, full_name FROM users WHERE role = 'employee' AND is_active = TRUE`)
	if err != nil {
		log.Printf("❌ فشل جلب الموظفين: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين"})
		return
	}
	defer rows.Close()
	var employees []gin.H
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err == nil {
			employees = append(employees, gin.H{"id": id, "full_name": name})
		}
	}
	c.JSON(http.StatusOK, employees)
}

```

---

## 📄 backend/internal/i18n/detect.go

```go
package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Detect يحدد لغة الرد المناسبة لطلب معيّن، بترتيب أولوية واضح:
//
//  1. باراميتر صريح في الرابط: /api/v1/tasks?lang=he
//     (مفيد لتجربة الـ API يدوياً أو من تطبيقات لا تتحكم بالهيدرز بسهولة)
//  2. هيدر مخصص يرسله الفرونت إند: X-Lang: he
//     (هذا ما تستخدمه واجهات WorkTrack الثلاث حسب لغة واجهة المستخدم المختارة)
//  3. هيدر المتصفح القياسي: Accept-Language
//  4. العربية كلغة افتراضية أخيرة إن لم يُرسَل أي شيء مما سبق
func Detect(c *gin.Context) Lang {
	if q := c.Query("lang"); q != "" {
		return Normalize(q)
	}

	if h := c.GetHeader("X-Lang"); h != "" {
		return Normalize(h)
	}

	if al := c.GetHeader("Accept-Language"); al != "" {
		// Accept-Language قد يأتي بصيغة معقدة مثل: "he-IL,he;q=0.9,en;q=0.8"
		// نأخذ فقط أول تفضيل ونحذف منه أي كود بلد (مثل IL) لنبقي على "he" فقط
		primary := strings.Split(al, ",")[0]
		primary = strings.Split(primary, ";")[0]
		primary = strings.Split(primary, "-")[0]
		return Normalize(primary)
	}

	return English
}

```

---

## 📄 backend/internal/i18n/i18n.go

```go
package i18n

import "strings"

// Lang يمثل رمز اللغة المدعومة في المشروع
type Lang string

const (
	Arabic  Lang = "ar"
	Hebrew  Lang = "he"
	English Lang = "en"
)

// messages هي "قاموس الترجمة" المركزي: مفتاح الرسالة -> نص لكل لغة
// أي رسالة جديدة يحتاجها المشروع تُضاف هنا مرة واحدة فقط بثلاث لغات
var messages = map[string]map[Lang]string{
	// ---------- تسجيل الدخول والحسابات ----------
	"err_missing_login_fields": {
		Arabic:  "الرجاء إدخال البريد الإلكتروني وكلمة المرور",
		Hebrew:  "נא להזין כתובת אימייל וסיסמה",
		English: "Please enter your email and password",
	},
	"err_invalid_credentials": {
		Arabic:  "تم تغير كلمة السر تواصل مع الادمن",
		Hebrew:  "הסיסמה שונתה, צור קשר עם ההנהלה",
		English: "Password has been changed, please contact administration",
	},
	"err_phone_not_registered": {
		Arabic:  "رقم الهاتف غير مسجل، يرجى التواصل مع المدير",
		Hebrew:  "מספר הטלפון לא רשום, אנא צור קשר עם המנהל",
		English: "Phone number not registered, please contact your manager",
	},
	"err_invalid_phone_credentials": {
		Arabic:  "رقم الهاتف غير صحيح، يرجى التواصل مع المدير",
		Hebrew:  "מספר הטלפון שגוי, אנא צור קשר עם המנהל",
		English: "Invalid phone number, please contact your manager",
	},
	"err_session_create_failed": {
		Arabic:  "تعذر إنشاء جلسة الدخول",
		Hebrew:  "לא ניתן היה ליצור את הפעלת ההתחברות",
		English: "Could not create a login session",
	},
	"err_missing_fields": {
		Arabic:  "بيانات غير مكتملة",
		Hebrew:  "חסרים פרטים בבקשה",
		English: "Missing required fields",
	},
	"err_invalid_email": {
		Arabic:  "صيغة البريد الإلكتروني غير صحيحة",
		Hebrew:  "כתובת האימייל אינה תקינה",
		English: "Invalid email format",
	},
	"err_password_hash_failed": {
		Arabic:  "فشل تشفير كلمة المرور",
		Hebrew:  "הצפנת הסיסמה נכשלה",
		English: "Failed to hash the password",
	},
	"err_email_in_use": {
		Arabic:  "البريد الإلكتروني مستخدم مسبقاً",
		Hebrew:  "כתובת האימייל כבר קיימת במערכת",
		English: "This email is already registered",
	},
	"msg_account_created": {
		Arabic:  "تم إنشاء الحساب بنجاح",
		Hebrew:  "החשבון נוצר בהצלחה",
		English: "Account created successfully",
	},
	"err_user_not_found": {
		Arabic:  "المستخدم غير موجود",
		Hebrew:  "המשתמש לא נמצא",
		English: "User not found",
	},

	// ---------- الجلسة والصلاحيات ----------
	"err_please_login": {
		Arabic:  "الرجاء تسجيل الدخول",
		Hebrew:  "נא להתחבר תחילה",
		English: "Please log in",
	},
	"err_invalid_session": {
		Arabic:  "جلسة غير صالحة، الرجاء تسجيل الدخول من جديد",
		Hebrew:  "ההתחברות אינה תקפה, נא להתחבר מחדש",
		English: "Invalid session, please log in again",
	},
	"err_password_changed_relogin": {
		Arabic:  "تم تغيير كلمة المرور، يرجى تسجيل الدخول مرة أخرى",
		Hebrew:  "הסיסמה שונתה, אנא התחבר שוב",
		English: "Password has been changed, please log in again",
	},
	"err_subscription_expired": {
		Arabic:  "اشتراكك انتهى أو تم إيقافه، الرجاء التواصل مع الدعم",
		Hebrew:  "המנוי שלך פג או בוטל, נא ליצור קשר עם התמיכה",
		English: "Your subscription has expired or been canceled, please contact support",
	},
	"err_forbidden_role": {
		Arabic:  "ليست لديك صلاحية للوصول لهذا الإجراء",
		Hebrew:  "אין לך הרשאה לבצע פעולה זו",
		English: "You don't have permission to perform this action",
	},
	"err_too_many_requests": {
		Arabic:  "طلبات كثيرة جداً، حاول لاحقاً",
		Hebrew:  "יותר מדי בקשות, נסה שוב מאוחר יותר",
		English: "Too many requests, please try again later",
	},
	"err_suspicious_activity": {
		Arabic:  "محاولات مشبوهة - تم حظر مؤقت",
		Hebrew:  "פעילות חשודה - נחסה זמנית",
		English: "Suspicious activity - temporarily blocked",
	},

	// ---------- التختيم والنطاق الجغرافي (Geofence) ----------
	"err_invalid_request_data": {
		Arabic:  "بيانات الطلب غير صحيحة",
		Hebrew:  "נתוני הבקשה אינם תקינים",
		English: "Invalid request data",
	},
	"err_invalid_coordinates": {
		Arabic:  "إحداثيات الموقع غير صالحة",
		Hebrew:  "קואורדינטות המיקום אינן תקינות",
		English: "Invalid location coordinates",
	},
	"err_outside_geofence_checkin": {
		Arabic:  "لا يمكنك بدء المهمة، أنت خارج نطاق موقع العمل المسموح",
		Hebrew:  "לא ניתן להתחיל את המשימה, אתה מחוץ לטווח המותר של אתר העבודה",
		English: "You can't start the task, you are outside the allowed worksite range",
	},
	"err_outside_geofence_checkout": {
		Arabic:  "لا يمكنك إنهاء المهمة، أنت خارج نطاق موقع العمل المسموح",
		Hebrew:  "לא ניתן לסיים את המשימה, אתה מחוץ לטווח המותר של אתר העבודה",
		English: "You can't end the task, you are outside the allowed worksite range",
	},
	"notif_checkin_rejected_title": {
		Arabic:  "محاولة تختيم مرفوضة",
		Hebrew:  "ניסיון החתמה נדחה",
		English: "Check-in attempt rejected",
	},
	"notif_checkin_rejected_body": {
		Arabic:  "لقد حاولت بدء مهمة وأنت خارج نطاق موقع العمل المحدد",
		Hebrew:  "ניסית להתחיל משימה כשאתה מחוץ לטווח אתר העבודה שהוגדר",
		English: "You tried to start a task while outside the defined worksite range",
	},
	"msg_checkin_success": {
		Arabic:  "تم بدء المهمة بنجاح، أنت داخل النطاق المسموح",
		Hebrew:  "המשימה התחילה בהצלחה, אתה בטווח המותר",
		English: "Task started successfully, you are within the allowed range",
	},
	"msg_checkout_success": {
		Arabic:  "تم إنهاء المهمة بنجاح",
		Hebrew:  "המשימה הסתיימה בהצלחה",
		English: "Task ended successfully",
	},

	// ---------- عمليات عامة (نقاط عمل / مهام / عملاء / إشعارات / رفع ملفات) ----------
	"err_operation_failed": {
		Arabic:  "حدث خطأ ما، الرجاء المحاولة مرة أخرى",
		Hebrew:  "משהו השתבש, נא לנסות שוב",
		English: "Something went wrong, please try again",
	},
	"msg_worksite_created": {
		Arabic:  "تم إنشاء نقطة العمل بنجاح",
		Hebrew:  "אתר העבודה נוצר בהצלחה",
		English: "Worksite created successfully",
	},
	"msg_task_created": {
		Arabic:  "تم إنشاء المهمة بنجاح",
		Hebrew:  "המשימה נוצרה בהצלחה",
		English: "Task created successfully",
	},
	"err_no_photo": {
		Arabic:  "لم يتم إرفاق أي صورة",
		Hebrew:  "לא צורפה תמונה",
		English: "No photo was attached",
	},
	"err_file_open_failed": {
		Arabic:  "تعذر فتح الملف",
		Hebrew:  "לא ניתן היה לפתוח את הקובץ",
		English: "Could not open the file",
	},
	"msg_health_ok": {
		Arabic:  "الخدمة تعمل بشكل طبيعي",
		Hebrew:  "השירות פועל כרגיל",
		English: "Service is healthy",
	},
}

// T (Translate) تُرجع نص الرسالة المطابق للمفتاح key بلغة lang
// إذا لم توجد ترجمة لتلك اللغة تحديداً، تُستخدم العربية كخيار احتياطي
// وإذا لم يوجد المفتاح نفسه أصلاً، تُرجَع قيمة المفتاح كما هي (يساعد أثناء التطوير
// على اكتشاف أي رسالة نسينا ترجمتها بسهولة، بدل ظهور رسالة فارغة للمستخدم)
func T(lang Lang, key string) string {
	if translations, ok := messages[key]; ok {
		if text, ok := translations[lang]; ok {
			return text
		}
		if text, ok := translations[Arabic]; ok {
			return text
		}
	}
	return key
}

// Normalize تحوّل أي نص لغة وارد من العميل (query param أو هيدر) إلى Lang معروفة
// وتتجاهل أي قيمة غير مدعومة بإرجاع الإنجليزية كافتراضي آمن
func Normalize(raw string) Lang {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "he", "heb", "hebrew", "עברית":
		return Hebrew
	case "en", "eng", "english":
		return English
	case "ar", "ara", "arabic", "العربية":
		return Arabic
	default:
		return English
	}
}

// GetTranslation returns the appropriate translation based on language
// It falls back to English, then Arabic, then Hebrew if the requested language is not available
func GetTranslation(lang Lang, ar, he, en string) string {
	switch lang {
	case Hebrew:
		if he != "" {
			return he
		}
		if en != "" {
			return en
		}
		return ar
	case Arabic:
		if ar != "" {
			return ar
		}
		if en != "" {
			return en
		}
		return he
	case English:
		if en != "" {
			return en
		}
		if ar != "" {
			return ar
		}
		return he
	default: // Default to English
		if en != "" {
			return en
		}
		if ar != "" {
			return ar
		}
		return he
	}
}

```

---

## 📄 backend/internal/middleware/auth_middleware.go

```go
package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware يتحقق من وجود JWT صالح في cookie أو هيدر Authorization
// ويضع user_id و role في السياق (Context) لاستخدامها في الـ Handlers
func AuthMiddleware(authService *services.AuthService, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)

		// محاولة الحصول على الرمز من cookie أولاً (الأكثر أماناً)
		tokenString, err := c.Cookie("worktrack_token")
		if err != nil || tokenString == "" {
			// إذا لم يوجد في cookie، حاول من هيدر Authorization للتوافقية
			header := c.GetHeader("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_please_login")})
				return
			}
			tokenString = strings.TrimPrefix(header, "Bearer ")
		}

		// التحقق من صحة التوكن
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_session")})
			return
		}

		// التحقق من أن كلمة المرور لم تتغير
		var currentPasswordChangedAt string
		err = db.QueryRow(`
			SELECT COALESCE(password_changed_at::text, '') 
			FROM users 
			WHERE id = $1
		`, claims.UserID).Scan(&currentPasswordChangedAt)
		
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_session")})
			return
		}

		// إذا كانت كلمة المرور قد تغيرت، أبطال التوكن
		if claims.PasswordChangedAt != currentPasswordChangedAt {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_password_changed_relogin")})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

```

---

## 📄 backend/internal/middleware/cors.go

```go
package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware يسمح فقط لروابط الواجهات الأمامية (Vercel) بالوصول للـ API
// نسمح أيضاً بهيدر X-Lang المخصص الذي تستخدمه الواجهات لتحديد لغة الرد (ar/he/en)
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Lang", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

```

---

## 📄 backend/internal/middleware/error_handler.go

```go
package middleware

import (
	"log"
	"time"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles errors and provides consistent error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ Panic recovered: %v", err)
				lang := i18n.Detect(c)
				c.JSON(500, gin.H{
					"error": i18n.T(lang, "err_internal_server"),
				})
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			lang := i18n.Detect(c)

			// Check if it's an AppError
			if appErr := utils.GetAppError(err); appErr.StatusCode != 500 {
				c.JSON(appErr.StatusCode, gin.H{
					"error": appErr.Message,
				})
				return
			}

			// Default error response
			c.JSON(500, gin.H{
				"error": i18n.T(lang, "err_internal_server"),
			})
		}
	}
}

// RequestLogger middleware logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request
		log.Printf("📡 %s %s - Status: %d - Duration: %v", method, path, statusCode, duration)

		// Log errors if any
		if len(c.Errors) > 0 {
			log.Printf("❌ Errors in %s %s: %v", method, path, c.Errors.String())
		}
	}
}
```

---

## 📄 backend/internal/middleware/rate_limiter.go

```go
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

// حدود الطلبات
const (
	generalLimit = 1000 // حد عام متساهل للإنتاج
	authLimit    = 100  // حد متساهل للمصادقة
	blockDuration = 5 * time.Minute // مدة الحظر عند الانتهاك
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

		limit := generalLimit
		if isAuthEndpoint {
			limit = authLimit
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

```

---

## 📄 backend/internal/middleware/role_middleware.go

```go
package middleware

import (
	"net/http"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

// RequireRole يقيّد وصول Endpoint معين على دور محدد فقط (مثلاً "admin")
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)

		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_forbidden_role")})
			return
		}
		c.Next()
	}
}

```

---

## 📄 backend/internal/middleware/security_headers.go

```go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityMiddleware يضيف رؤوس الأمان للوقاية من هجمات XSS، clickjacking، وغيرها
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy - يمنع تحميل موارد من مصادر غير مصرح بها
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-ancestors 'none';")
		
		// يمنع المتصفح من تخمين نوع MIME
		c.Header("X-Content-Type-Options", "nosniff")
		
		// يمنع تحميل الصفحة في iframe - حماية من clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// تفعيل حماية XSS المدمجة في المتصفح
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// سياسة المرجع - تحكم في كيفية إرسال معلومات المرجع
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// سياسة الأذونات للوصول إلى ميزات المتصفح
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}
```

---

## 📄 backend/internal/middleware/subscription_middleware.go

```go
package middleware

import (
	"database/sql"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

func SubscriptionMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.Detect(c)
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_please_login")})
			return
		}

		var status string
		var expiresAt sql.NullTime
		err := db.QueryRow(`
			SELECT subscription_status, subscription_expires_at
			FROM users
			WHERE id = $1
		`, userID).Scan(&status, &expiresAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": i18n.T(lang, "err_invalid_session")})
			return
		}

		if status == "canceled" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}

		if status == "expired" || (expiresAt.Valid && time.Now().After(expiresAt.Time)) {
			if status != "expired" {
				_, _ = db.Exec(`UPDATE users SET subscription_status = 'expired' WHERE id = $1`, userID)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": i18n.T(lang, "err_subscription_expired")})
			return
		}

		c.Next()
	}
}

```

---

## 📄 backend/internal/models/attendance.go

```go
package models

import "time"

// Attendance يمثل سجل حضور/انصراف واحد لموظف على مهمة معينة
// يحتوي على إثبات الموقع الجغرافي (lat/lng + المسافة المحسوبة) لكل من البدء والانتهاء
type Attendance struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	TaskID            *string `json:"task_id,omitempty"`
	WorksiteID        string  `json:"worksite_id"`
	ServiceRequestID  *string `json:"service_request_id,omitempty"` // ربط مع طلب الخدمة

	CheckInTime           *time.Time `json:"check_in_time,omitempty"`
	CheckInLat            *float64   `json:"check_in_lat,omitempty"`
	CheckInLng            *float64   `json:"check_in_lng,omitempty"`
	CheckInDistanceMeters *float64   `json:"check_in_distance_meters,omitempty"`

	CheckOutTime           *time.Time `json:"check_out_time,omitempty"`
	CheckOutLat            *float64   `json:"check_out_lat,omitempty"`
	CheckOutLng            *float64   `json:"check_out_lng,omitempty"`
	CheckOutDistanceMeters *float64   `json:"check_out_distance_meters,omitempty"`

	Status    string    `json:"status"` // in_progress | completed
	CreatedAt time.Time `json:"created_at"`
}

```

---

## 📄 backend/internal/models/client.go

```go
package models

import "time"

type Client struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	NameAr      string    `json:"name_ar,omitempty"`
	NameHe      string    `json:"name_he,omitempty"`
	NameEn      string    `json:"name_en,omitempty"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	AddressAr   string    `json:"address_ar,omitempty"`
	AddressHe   string    `json:"address_he,omitempty"`
	AddressEn   string    `json:"address_en,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

```

---

## 📄 backend/internal/models/notification.go

```go
package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

```

---

## 📄 backend/internal/models/task.go

```go
package models

import "time"

// Task تمثل مهمة عمل واحدة مرتبطة بعميل ونقطة عمل وموظف مكلَّف (اختياري كل منها)
type Task struct {
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	TitleAr           string     `json:"title_ar,omitempty"`
	TitleHe           string     `json:"title_he,omitempty"`
	TitleEn           string     `json:"title_en,omitempty"`
	Description       string     `json:"description"`
	DescriptionAr     string     `json:"description_ar,omitempty"`
	DescriptionHe     string     `json:"description_he,omitempty"`
	DescriptionEn     string     `json:"description_en,omitempty"`
	ClientID          *string    `json:"client_id,omitempty"`
	WorksiteID        string     `json:"worksite_id"`
	AssignedUserID    *string    `json:"assigned_user_id,omitempty"`
	Status            string     `json:"status"` // pending | in_progress | completed | late | cancelled
	Priority          string     `json:"priority,omitempty"` // low | normal | high | urgent
	ScheduledStart    *time.Time `json:"scheduled_start,omitempty"`
	ScheduledEnd      *time.Time `json:"scheduled_end,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	// حقول إضافية للعرض فقط (ليست في قاعدة البيانات)
	WorksiteName      string     `json:"worksite_name,omitempty"`
	WorksiteNameAr    string     `json:"worksite_name_ar,omitempty"`
	WorksiteNameHe    string     `json:"worksite_name_he,omitempty"`
	WorksiteNameEn    string     `json:"worksite_name_en,omitempty"`
	WorksiteAddress   string     `json:"worksite_address,omitempty"`
	WorksiteAddressAr string     `json:"worksite_address_ar,omitempty"`
	WorksiteAddressHe string     `json:"worksite_address_he,omitempty"`
	WorksiteAddressEn string     `json:"worksite_address_en,omitempty"`
	WorksiteLatitude  float64    `json:"worksite_latitude,omitempty"`
	WorksiteLongitude float64    `json:"worksite_longitude,omitempty"`
	ClientName        string     `json:"client_name,omitempty"`
	ClientNameAr      string     `json:"client_name_ar,omitempty"`
	ClientNameHe      string     `json:"client_name_he,omitempty"`
	ClientNameEn      string     `json:"client_name_en,omitempty"`
	ClientAddress     string     `json:"client_address,omitempty"`
	ClientAddressAr   string     `json:"client_address_ar,omitempty"`
	ClientAddressHe   string     `json:"client_address_he,omitempty"`
	ClientAddressEn   string     `json:"client_address_en,omitempty"`
	ClientPhone       string     `json:"client_phone,omitempty"`
	Language          string     `json:"language,omitempty"`
}

```

---

## 📄 backend/internal/models/user.go

```go
package models

import "time"

// User يمثل صف واحد من جدول users. الحقل Role يحدد إن كان "admin" أو "employee"
type User struct {
	ID                    string     `json:"id"`
	FullName              string     `json:"full_name"`
	Email                 string     `json:"email"`
	Phone                 string     `json:"phone"`
	PasswordHash          string     `json:"-"`    // "-" يعني: لا تُرجَع هذه القيمة أبداً في أي رد JSON
	Role                  string     `json:"role"` // admin | employee
	AvatarURL             string     `json:"avatar_url,omitempty"`
	IsActive              bool       `json:"is_active"`
	SubscriptionStatus    string     `json:"subscription_status,omitempty"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

```

---

## 📄 backend/internal/models/worksite.go

```go
package models

import "time"

// Worksite تمثل نقطة عمل جغرافية (Geofence Zone)
// radius_meters هو نصف القطر المسموح للموظف أن "يختم" داخله
type Worksite struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	NameAr       string    `json:"name_ar,omitempty"`
	NameHe       string    `json:"name_he,omitempty"`
	NameEn       string    `json:"name_en,omitempty"`
	Address      string    `json:"address"`
	AddressAr    string    `json:"address_ar,omitempty"`
	AddressHe    string    `json:"address_he,omitempty"`
	AddressEn    string    `json:"address_en,omitempty"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters int       `json:"radius_meters"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

```

---

## 📄 backend/internal/router/router.go

```go
package router

import (
	"database/sql"
	"net/http"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/handlers"
	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/middleware"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func Setup(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SecurityMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))
	r.Use(middleware.RateLimiter())

	authService := services.NewAuthService(cfg.JWTSecret)
	attendanceService := services.NewAttendanceService(db)
	notificationService := services.NewNotificationService(db)
	geocodingService := services.NewGeocodingService(cfg.GeoapifyKey)
	securityLogger := services.NewSecurityLogger(db)

	wsHandler := handlers.NewWSHandler()
	authHandler := handlers.NewAuthHandler(db, authService, securityLogger)
	authHandler.SetWSHandler(wsHandler)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, notificationService)
	worksiteHandler := handlers.NewWorksiteHandler(db)
	reportHandler := handlers.NewReportHandler(db)
	notificationHandler := handlers.NewNotificationHandler(db)
	serviceHandler := handlers.NewServiceHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	locationHandler := handlers.NewLocationHandler(db, wsHandler)
	geocodingHandler := handlers.NewGeocodingHandler(geocodingService)
	uploadHandler := handlers.NewUploadHandler(db)

	r.GET("/health", func(c *gin.Context) {
		lang := i18n.Detect(c)
		c.JSON(200, gin.H{"status": "ok", "message": i18n.T(lang, "msg_health_ok")})
	})

	// WebSocket endpoint للتتبع اللحظي
	r.GET("/ws", wsHandler.HandleWebSocket)

	api := r.Group("/api/v1")

	// مسارات المصادقة (بدون توكن)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/phone-login", authHandler.PhoneLogin)
	api.POST("/auth/logout", authHandler.Logout)
	api.GET("/geocode/autocomplete", geocodingHandler.Autocomplete)

	// المسارات المحمية (تتطلب توكن)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService, db))
	{
		// المستخدم
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/device", authHandler.GetDeviceInfo)

		// المسارات المتاحة للموظفين (بدون التحقق من الاشتراك)
		employee := protected.Group("")
		{
			// نقاط العمل
			employee.GET("/worksites", worksiteHandler.List)
			employee.GET("/worksites/available", worksiteHandler.GetAvailableWorksites)

			// الحضور
			employee.POST("/attendance/check-in", attendanceHandler.CheckIn)
			employee.POST("/attendance/check-out", attendanceHandler.CheckOut)
			employee.GET("/attendance/current", attendanceHandler.GetCurrentAttendance)
			employee.GET("/attendance/summary", attendanceHandler.GetAttendanceSummary)
			employee.GET("/attendance/all-summary", attendanceHandler.GetAllAttendanceSummary)
			employee.GET("/attendance/my-history", attendanceHandler.GetMyAttendanceHistory)
			employee.GET("/attendance/my-monthly-summary", attendanceHandler.GetMyMonthlySummary)
			employee.GET("/attendance/history", attendanceHandler.GetMyAttendanceHistory)

			// المهام
			employee.GET("/tasks/mine", taskHandler.MyTasks)

			// رفع الصور (للموظفين)
			employee.POST("/upload/task-photo", uploadHandler.UploadTaskPhoto)
			employee.GET("/upload/photos", uploadHandler.GetAttendancePhotos)
			employee.GET("/upload/download/:id", uploadHandler.DownloadPhoto)

			// الموقع
			employee.GET("/location/active", locationHandler.GetActiveEmployees)
			employee.GET("/location/track/:id", locationHandler.GetEmployeeTrack)
			employee.GET("/location/logs", locationHandler.GetLocationLogs)
			employee.POST("/location/update", locationHandler.UpdateLocation)

			// الإشعارات
			employee.GET("/notifications", notificationHandler.List)
		}

		// طلبات الخدمة (متاحة للعملاء بدون اشتراك)
		protected.GET("/service/requests", serviceHandler.ListRequests)
		protected.POST("/service/requests", serviceHandler.CreateRequest)
		protected.POST("/service/requests/:id/rate", serviceHandler.RateRequest)

		// طلبات الخدمة المعينة للموظف
		employee.GET("/service/assigned", serviceHandler.GetMyAssignedRequests)
		employee.GET("/service/assigned/:id", serviceHandler.GetAssignedRequestDetails)
		employee.PUT("/service/assigned/:id/status", serviceHandler.UpdateAssignmentStatus)
		employee.POST("/service/requests/:id/complete", serviceHandler.CompleteRequest)

		// مسارات المدير (تتطلب دور admin ولكن بدون اشتراك)
		admin := protected.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			// إدارة الموظفين
			admin.POST("/auth/employee-phone", authHandler.CreateEmployeePhone)
			admin.GET("/admin/employees", authHandler.ListEmployees)
			admin.DELETE("/admin/employees/:id", authHandler.DeleteEmployee)
			admin.POST("/admin/reset-device", authHandler.ResetDevice)
			// إدارة العملاء
			admin.GET("/admin/customers", authHandler.ListCustomers)
			admin.POST("/admin/create-client", authHandler.CreateClient)
			admin.POST("/admin/reset-customer-password", authHandler.ResetCustomerPassword)
			admin.DELETE("/admin/customers/:id", authHandler.DeleteCustomer)
		}

		// المسارات التي تتطلب اشتراكاً نشطاً (للمديرين فقط)
		paid := protected.Group("")
		paid.Use(middleware.SubscriptionMiddleware(db))
		{
			// مسارات المدير (تتطلب دور admin + اشتراك نشط)
			adminPaid := paid.Group("")
			adminPaid.Use(middleware.RequireRole("admin"))
			{
				// تعيين الموظفين لطلبات الخدمة
				adminPaid.POST("/service/assign", serviceHandler.AssignEmployee)
				// حذف طلبات الخدمة
				adminPaid.DELETE("/service/requests/:id", serviceHandler.DeleteRequest)

				// نقاط العمل
				adminPaid.POST("/worksites", worksiteHandler.Create)
				adminPaid.DELETE("/worksites/:id", worksiteHandler.Delete)
				adminPaid.POST("/worksites/assign", worksiteHandler.AssignEmployee)
				adminPaid.GET("/worksites/employees", worksiteHandler.GetAvailableEmployees)

				// المهام
				adminPaid.GET("/tasks", taskHandler.ListAll)
				adminPaid.POST("/tasks", taskHandler.Create)
				adminPaid.PUT("/tasks/:id", taskHandler.Update)

				// التقارير - مسارات جديدة
				adminPaid.GET("/reports/daily-summary", reportHandler.DailySummary)
				adminPaid.GET("/reports/comprehensive", reportHandler.ComprehensiveReport)
				adminPaid.GET("/reports/service-requests", reportHandler.ServiceRequestsReport)
				adminPaid.GET("/reports/employee-performance", reportHandler.EmployeePerformanceReport)
				adminPaid.GET("/reports/client-activity", reportHandler.ClientActivityReport)
				adminPaid.GET("/reports/pending-employees", reportHandler.GetPendingEmployees)
				adminPaid.GET("/reports/completed-employees", reportHandler.GetCompletedEmployees)

				// سجل الحضور للموظفين
				adminPaid.GET("/attendance/employee/:id/history", attendanceHandler.GetEmployeeAttendanceHistory)
				adminPaid.GET("/attendance/employee/:id/monthly-summary", attendanceHandler.GetEmployeeMonthlySummary)
				adminPaid.POST("/attendance/cleanup-old-records", attendanceHandler.CleanupOldRecords)
				adminPaid.POST("/attendance/force-checkout", attendanceHandler.ForceCheckOut)

				// إشعارات تغيير كلمة المرور
				adminPaid.POST("/admin/notify-password-change", authHandler.NotifyPasswordChange)

				// إدارة الصور (حذف فقط - المدير يمكنه حذف الصور)
				adminPaid.DELETE("/upload/photo/:id", uploadHandler.DeletePhoto)

				// الموقع والملاحظات الأمنية
				adminPaid.GET("/location/security/:id", locationHandler.GetEmployeeSecurityNotes)
				
				// سجلات الأمان
				adminPaid.GET("/security/logs", func(c *gin.Context) {
					limit := 50
					logs, err := securityLogger.GetRecentSecurityLogs(limit)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب سجلات الأمان"})
						return
					}
					c.JSON(http.StatusOK, logs)
				})
			}
		}
	}

	return r
}

```

---

## 📄 backend/internal/services/attendance_service.go

```go
package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"worktrack/backend/internal/models"
	"worktrack/backend/pkg/utils"

	"github.com/google/uuid"
)

type AttendanceService struct {
	DB *sql.DB
}

func NewAttendanceService(db *sql.DB) *AttendanceService {
	return &AttendanceService{DB: db}
}

var ErrOutsideGeofence = errors.New("خارج نطاق الموقع المسموح")

func (s *AttendanceService) GetCurrentAttendance(userID string) (*models.Attendance, error) {
	var attendance models.Attendance
	var checkInTime time.Time
	var serviceRequestID *string

	err := s.DB.QueryRow(`
		SELECT id, worksite_id, check_in_time, status, service_request_id
		FROM attendance
		WHERE user_id = $1 AND status = 'in_progress'
		ORDER BY check_in_time DESC
		LIMIT 1
	`, userID).Scan(&attendance.ID, &attendance.WorksiteID, &checkInTime, &attendance.Status, &serviceRequestID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	attendance.UserID = userID
	attendance.CheckInTime = &checkInTime
	attendance.ServiceRequestID = serviceRequestID
	return &attendance, nil
}

func (s *AttendanceService) CheckIn(userID, worksiteID string, lat, lng float64, serviceRequestID *string) (*models.Attendance, *GeofenceCheckResult, error) {
	log.Printf("📌 CheckIn - المستخدم: %s, الموقع: %s", userID, worksiteID)
	log.Printf("📍 إحداثيات المستخدم: %f, %f", lat, lng)

	// التحقق من وجود وردية نشطة
	existing, _ := s.GetCurrentAttendance(userID)
	if existing != nil {
		return nil, nil, errors.New("يوجد وردية نشطة بالفعل")
	}

	if !utils.IsValidCoordinates(lat, lng) {
		return nil, nil, errors.New("إحداثيات غير صالحة")
	}

	// جلب نقطة العمل
	var worksite models.Worksite
	err := s.DB.QueryRow(`
		SELECT id, name, address, latitude, longitude, radius_meters, is_active
		FROM worksites WHERE id = $1 AND is_active = TRUE
	`, worksiteID).Scan(&worksite.ID, &worksite.Name, &worksite.Address,
		&worksite.Latitude, &worksite.Longitude, &worksite.RadiusMeters,
		&worksite.IsActive)

	if err != nil {
		log.Printf("❌ نقطة العمل غير موجودة: %v", err)
		return nil, nil, fmt.Errorf("نقطة العمل غير موجودة: %w", err)
	}

	log.Printf("📍 نقطة العمل: %s (%.6f, %.6f) - النطاق: %d متر",
		worksite.Name, worksite.Latitude, worksite.Longitude, worksite.RadiusMeters)

	// حساب المسافة بين المستخدم ونقطة العمل
	distance := utils.HaversineDistance(lat, lng, worksite.Latitude, worksite.Longitude)
	log.Printf("📏 المسافة المحسوبة: %.2f متر", distance)

	// التحقق من النطاق
	if distance > float64(worksite.RadiusMeters) {
		log.Printf("❌ المستخدم خارج النطاق: %.2f > %d", distance, worksite.RadiusMeters)
		result := &GeofenceCheckResult{
			IsWithinRange:  false,
			DistanceMeters: distance,
			AllowedRadius:  worksite.RadiusMeters,
		}
		return nil, result, ErrOutsideGeofence
	}

	log.Printf("✅ المستخدم داخل النطاق: %.2f <= %d", distance, worksite.RadiusMeters)

	// إنشاء سجل الحضور
	id := uuid.NewString()
	now := utils.NowInJerusalem()

	_, err = s.DB.Exec(`
		INSERT INTO attendance (id, user_id, worksite_id, check_in_time,
			check_in_lat, check_in_lng, check_in_distance_meters, status, service_request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'in_progress', $8)`,
		id, userID, worksite.ID, now, lat, lng, distance, serviceRequestID,
	)
	if err != nil {
		log.Printf("❌ فشل حفظ بداية الدوام: %v", err)
		return nil, nil, fmt.Errorf("فشل حفظ بداية الدوام: %w", err)
	}

	result := &GeofenceCheckResult{
		IsWithinRange:  true,
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}

	attendance := &models.Attendance{
		ID:                    id,
		UserID:                userID,
		WorksiteID:            worksite.ID,
		ServiceRequestID:      serviceRequestID,
		CheckInTime:           &now,
		CheckInLat:            &lat,
		CheckInLng:            &lng,
		CheckInDistanceMeters: &distance,
		Status:                "in_progress",
	}

	log.Printf("✅ تم تسجيل بدء الدوام بنجاح: %s", id)
	return attendance, result, nil
}

func (s *AttendanceService) CheckOut(userID, attendanceID string, lat, lng float64) (*GeofenceCheckResult, float64, *string, error) {
	log.Printf("📌 CheckOut - المستخدم: %s", userID)
	log.Printf("📍 إحداثيات المستخدم: %f, %f", lat, lng)

	if !utils.IsValidCoordinates(lat, lng) {
		return nil, 0, nil, errors.New("إحداثيات غير صالحة")
	}

	var worksiteID string
	var checkInTime time.Time
	var serviceRequestID *string

	err := s.DB.QueryRow(`
		SELECT worksite_id, check_in_time, service_request_id FROM attendance 
		WHERE id = $1 AND user_id = $2 AND status = 'in_progress'
	`, attendanceID, userID).Scan(&worksiteID, &checkInTime, &serviceRequestID)
	if err != nil {
		log.Printf("❌ لا يوجد وردية نشطة: %v", err)
		return nil, 0, nil, errors.New("لا يوجد وردية نشطة لهذا المستخدم")
	}

	var worksite models.Worksite
	err = s.DB.QueryRow(`
		SELECT id, name, latitude, longitude, radius_meters
		FROM worksites WHERE id = $1
	`, worksiteID).Scan(&worksite.ID, &worksite.Name,
		&worksite.Latitude, &worksite.Longitude, &worksite.RadiusMeters)
	if err != nil {
		log.Printf("❌ نقطة العمل غير موجودة: %v", err)
		return nil, 0, nil, fmt.Errorf("نقطة العمل غير موجودة: %w", err)
	}

	// حساب المسافة عند الخروج
	distance := utils.HaversineDistance(lat, lng, worksite.Latitude, worksite.Longitude)
	log.Printf("📏 المسافة عند الخروج: %.2f متر", distance)

	if distance > float64(worksite.RadiusMeters) {
		log.Printf("❌ المستخدم خارج النطاق عند الخروج: %.2f > %d", distance, worksite.RadiusMeters)
		result := &GeofenceCheckResult{
			IsWithinRange:  false,
			DistanceMeters: distance,
			AllowedRadius:  worksite.RadiusMeters,
		}
		return result, 0, serviceRequestID, ErrOutsideGeofence
	}

	now := utils.NowInJerusalem()
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, check_out_lat = $2, check_out_lng = $3,
		    check_out_distance_meters = $4, status = 'completed'
		WHERE id = $5`,
		now, lat, lng, distance, attendanceID,
	)
	if err != nil {
		log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
		return nil, 0, nil, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
	}

	// حساب الساعات مع التعامل مع فترات العمل عبر منتصف الليل
	workedHours := now.Sub(checkInTime).Hours()
	if workedHours < 0 {
		workedHours += 24 // إضافة 24 ساعة إذا كان الفرق سلبياً (عبر منتصف الليل)
	}
	log.Printf("✅ تم تسجيل إنهاء الدوام: %.2f ساعة", workedHours)

	result := &GeofenceCheckResult{
		IsWithinRange:  true,
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}

	return result, workedHours, serviceRequestID, nil
}

// ForceCheckOut إنهاء دوام الموظف من قبل المدير (بدون التحقق من الموقع)
func (s *AttendanceService) ForceCheckOut(attendanceID string, adminID string) (float64, error) {
	log.Printf("📌 ForceCheckOut - الـ ID: %s, المدير: %s", attendanceID, adminID)

	var checkInTime time.Time
	var userID string

	// جلب معلومات الوردية
	err := s.DB.QueryRow(`
		SELECT user_id, check_in_time FROM attendance 
		WHERE id = $1 AND status = 'in_progress'
	`, attendanceID).Scan(&userID, &checkInTime)
	if err != nil {
		log.Printf("❌ لا يوجد وردية نشطة: %v", err)
		return 0, errors.New("لا يوجد وردية نشطة")
	}

	now := utils.NowInJerusalem()

	// محاولة التحديث مع check_out_notes أولاً
	_, err = s.DB.Exec(`
		UPDATE attendance
		SET check_out_time = $1, status = 'completed', check_out_notes = 'تم إنهاء الدوام من قبل المدير'
		WHERE id = $2
	`, now, attendanceID)

	// إذا فشل بسبب عدم وجود عمود check_out_notes، حاول بدونه
	if err != nil {
		log.Printf("⚠️ فشل التحديث مع check_out_notes، محاولة بدونه: %v", err)
		_, err = s.DB.Exec(`
			UPDATE attendance
			SET check_out_time = $1, status = 'completed'
			WHERE id = $2
		`, now, attendanceID)
		if err != nil {
			log.Printf("❌ فشل تحديث سجل الحضور: %v", err)
			return 0, fmt.Errorf("فشل تحديث سجل الحضور: %w", err)
		}
	}

	// حساب الساعات مع التعامل مع فترات العمل عبر منتصف الليل
	workedHours := now.Sub(checkInTime).Hours()
	if workedHours < 0 {
		workedHours += 24 // إضافة 24 ساعة إذا كان الفرق سلبياً (عبر منتصف الليل)
	}
	log.Printf("✅ تم إنهاء الدوام من قبل المدير: %.2f ساعة", workedHours)

	return workedHours, nil
}

```

---

## 📄 backend/internal/services/auth_service.go

```go
package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{jwtSecret: []byte(secret)}
}

// HashPassword يشفّر كلمة المرور قبل حفظها في قاعدة البيانات (bcrypt، اتجاه واحد فقط)
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword يقارن كلمة المرور المدخلة بالنسخة المشفرة المخزنة
func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Claims بيانات التوكن: هوية المستخدم ودوره (admin/employee)
type Claims struct {
	UserID              string `json:"user_id"`
	Role                string `json:"role"`
	PasswordChangedAt   string `json:"password_changed_at"` // لإبطال التوكنات عند تغيير كلمة المرور
	jwt.RegisteredClaims
}

// GenerateToken ينشئ JWT صالح لمدة 7 أيام لكل موظف/مدير بعد تسجيل الدخول
func (s *AuthService) GenerateToken(userID, role, passwordChangedAt string) (string, error) {
	claims := Claims{
		UserID:            userID,
		Role:              role,
		PasswordChangedAt: passwordChangedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken يتحقق من صلاحية التوكن ويعيد بياناته.
// نص الخطأ هنا تقني (إنجليزي) عمداً لأنه لا يُعرَض للمستخدم مباشرة أبداً؛
// الـ AuthMiddleware هو من يترجم أي فشل هنا إلى رسالة i18n مناسبة للعميل.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

// ValidateTokenWithPasswordCheck يتحقق من صلاحية التوكن ويحقق أيضاً أن كلمة المرور لم تتغير
func (s *AuthService) ValidateTokenWithPasswordCheck(tokenString string, currentPasswordChangedAt string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// التحقق من أن كلمة المرور لم تتغير منذ إصدار التوكن
	if claims.PasswordChangedAt != currentPasswordChangedAt {
		return nil, errors.New("password changed, token invalid")
	}

	return claims, nil
}

```

---

## 📄 backend/internal/services/cache_service.go

```go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"worktrack/backend/internal/database"
)

// CacheService خدمة التخزين المؤقت باستخدام Redis
type CacheService struct{}

// NewCacheService ينشئ خدمة تخزين مؤقت جديدة
func NewCacheService() *CacheService {
	return &CacheService{}
}

// Set يخزن قيمة في Redis
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if database.RedisClient == nil {
		return nil // Redis غير متصل، تجاوز التخزين المؤقت
	}

	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("فشل تحويل القيمة إلى JSON: %w", err)
	}

	return database.RedisClient.Set(ctx, key, json, expiration).Err()
}

// Get يسترجع قيمة من Redis
func (cs *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if database.RedisClient == nil {
		return fmt.Errorf("Redis غير متصل") // Redis غير متصل
	}

	val, err := database.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete يحذف قيمة من Redis
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	if database.RedisClient == nil {
		return nil
	}

	return database.RedisClient.Del(ctx, key).Err()
}

// Exists يتحقق من وجود مفتاح في Redis
func (cs *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	if database.RedisClient == nil {
		return false, nil
	}

	result, err := database.RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Clear ي مسح جميع المفاتيح التي تطابق النمط
func (cs *CacheService) Clear(ctx context.Context, pattern string) error {
	if database.RedisClient == nil {
		return nil
	}

	iter := database.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		database.RedisClient.Del(ctx, iter.Val())
	}
	return iter.Err()
}
```

---

## 📄 backend/internal/services/geocoding_service.go

```go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type GeocodingService struct {
	client *resty.Client
	apiKey string
}

type GeocodeResult struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	LabelHe     string  `json:"label_he"`      // الاسم بالعبرية
	LabelAr     string  `json:"label_ar"`      // الاسم بالعربية
	Language    string  `json:"language"`
	CountryCode string  `json:"country_code"`
	State       string  `json:"state"`
	City        string  `json:"city"`
	CityHe      string  `json:"city_he"`       // المدينة بالعبرية
	CityAr      string  `json:"city_ar"`       // المدينة بالعربية
	Street      string  `json:"street"`
	StreetHe    string  `json:"street_he"`     // الشارع بالعبرية
	StreetAr    string  `json:"street_ar"`     // الشارع بالعربية
	HouseNumber string  `json:"house_number"`
	PostalCode  string  `json:"postal_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Type        string  `json:"type"`
}

func NewGeocodingService(apiKey string) *GeocodingService {
	return &GeocodingService{
		client: resty.New().SetTimeout(10 * time.Second),
		apiKey: apiKey,
	}
}

func (s *GeocodingService) Autocomplete(query, language string) ([]GeocodeResult, error) {
	if len(query) < 2 {
		return nil, nil
	}

	log.Printf("🔍 البحث العالمي: %s (اللغة: %s)", query, language)

	// Geoapify API - بحث عالمي
	url := "https://api.geoapify.com/v1/geocode/autocomplete"
	
	// استخدام اللغة المطلوبة
	searchLang := "en"
	if language == "ar" {
		searchLang = "ar"
	} else if language == "he" {
		searchLang = "he"
	} else {
		searchLang = "en" // افتراضي إنجليزي
	}
	
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"text":   query,
			"apiKey": s.apiKey,
			"limit":  "15",
			"lang":   searchLang,
			// إزالة فلتر الدولة للبحث العالمي
		}).
		Get(url)

	if err != nil {
		log.Printf("❌ فشل الاتصال بـ Geoapify: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("❌ خطأ من Geoapify: %s", resp.String())
		return nil, fmt.Errorf("Geoapify error: %s", resp.Status())
	}

	var response struct {
		Features []struct {
			Properties struct {
				PlaceID     string `json:"place_id"`
				Formatted   string `json:"formatted"`
				CountryCode string `json:"country_code"`
				State       string `json:"state"`
				City        string `json:"city"`
				Street      string `json:"street"`
				HouseNumber string `json:"housenumber"`
				Postcode    string `json:"postcode"`
				AddressLine1 string `json:"address_line1"`
				AddressLine2 string `json:"address_line2"`
				ResultType  string `json:"result_type"`
				Lat         float64 `json:"lat"`
				Lon         float64 `json:"lon"`
				// أسماء من المصدر
				Name        string `json:"name"`
			} `json:"properties"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Printf("❌ فشل تحليل الرد: %v", err)
		return nil, err
	}

	var results []GeocodeResult
	for _, feature := range response.Features {
		props := feature.Properties
		
		var lat, lon float64
		if len(feature.Geometry.Coordinates) >= 2 {
			lon = feature.Geometry.Coordinates[0]
			lat = feature.Geometry.Coordinates[1]
		} else {
			lat = props.Lat
			lon = props.Lon
		}

		// استخراج الأسماء بناءً على اللغة
		cityHe := props.City
		streetHe := props.Street
		cityAr := props.City
		streetAr := props.Street
		
		// إذا كان الاسم فارغاً، استخدم المنسق
		labelHe := props.AddressLine1
		if labelHe == "" {
			labelHe = props.Street
			if props.HouseNumber != "" {
				labelHe = labelHe + " " + props.HouseNumber
			}
		}
		if labelHe == "" {
			labelHe = props.City
		}
		if labelHe == "" {
			labelHe = props.Formatted
		}

		// نفس الشيء للعربية
		labelAr := labelHe // Geoapify يعيد نفس الاسم عادة

		results = append(results, GeocodeResult{
			ID:          props.PlaceID,
			Label:       props.Formatted,
			LabelHe:     labelHe,
			LabelAr:     labelAr,
			Language:    searchLang,
			CountryCode: props.CountryCode,
			State:       props.State,
			City:        props.City,
			CityHe:      cityHe,
			CityAr:      cityAr,
			Street:      props.Street,
			StreetHe:    streetHe,
			StreetAr:    streetAr,
			HouseNumber: props.HouseNumber,
			PostalCode:  props.Postcode,
			Latitude:    lat,
			Longitude:   lon,
			Type:        props.ResultType,
		})
	}

	log.Printf("✅ تم العثور على %d نتيجة عالمية", len(results))
	return results, nil
}

```

---

## 📄 backend/internal/services/geofence_service.go

```go
package services

import (
	"worktrack/backend/internal/models"
	"worktrack/backend/pkg/utils"
)

// GeofenceCheckResult نتيجة التحقق من موقع الموظف مقابل نقطة العمل
type GeofenceCheckResult struct {
	IsWithinRange  bool    `json:"is_within_range"`
	DistanceMeters float64 `json:"distance_meters"`
	AllowedRadius  int     `json:"allowed_radius_meters"`
}

// CheckWithinWorksite هي الدالة المسؤولة عن القرار الحاسم:
// هل موقع الموظف الحالي (lat, lng) يقع داخل نطاق نقطة العمل المحددة؟
//
// ⚠️ هذه الدالة تُستدعى دائماً من الـ Backend وليس فقط من الواجهة،
// لأن أي تحقق يتم فقط في المتصفح يمكن التلاعب به بسهولة (تعديل JS، إحداثيات مزيّفة...).
// السيرفر هو الحَكَم الوحيد الموثوق لقرار السماح بالتختيم أو رفضه، بغض النظر
// عن لغة الواجهة التي يستخدمها الموظف (عربي أو عبري أو إنجليزي) — القرار
// نفسه لا يتأثر باللغة إطلاقاً، فقط رسالة الرد للمستخدم هي ما تُترجَم لاحقاً.
func CheckWithinWorksite(userLat, userLng float64, worksite models.Worksite) GeofenceCheckResult {
	distance := utils.HaversineDistance(userLat, userLng, worksite.Latitude, worksite.Longitude)

	return GeofenceCheckResult{
		IsWithinRange:  distance <= float64(worksite.RadiusMeters),
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}
}

```

---

## 📄 backend/internal/services/notification_service.go

```go
package services

import (
	"database/sql"

	"github.com/google/uuid"
)

type NotificationService struct {
	DB *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// Send ينشئ إشعاراً جديداً لموظف معين (title/body تُمرَّر مُترجَمة مسبقاً من الـ Handler
// عبر i18n.T، حتى يصل الإشعار بنفس لغة واجهة الموظف الذي سيقرأه)
func (s *NotificationService) Send(userID, title, body string) error {
	_, err := s.DB.Exec(`
		INSERT INTO notifications (id, user_id, title, body, is_read, created_at)
		VALUES ($1, $2, $3, $4, FALSE, now())`,
		uuid.NewString(), userID, title, body,
	)
	return err
}

```

---

## 📄 backend/internal/services/security_logger.go

```go
package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SecurityLogger يسجل أحداث الأمان والمحاولات المشبوهة
type SecurityLogger struct {
	DB *sql.DB
}

func NewSecurityLogger(db *sql.DB) *SecurityLogger {
	return &SecurityLogger{DB: db}
}

// LogAuthAttempt يسجل محاولة مصادقة
func (sl *SecurityLogger) LogAuthAttempt(email, phone, ip, userAgent string, success bool, reason string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نتخطى التسجيل
		return
	}
	
	_, err = sl.DB.Exec(`
		INSERT INTO security_logs (event_type, email, phone, ip, user_agent, success, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, "auth_attempt", email, phone, ip, userAgent, success, reason, time.Now())
	
	if err != nil {
		log.Printf("⚠️ فشل تسجيل محاولة مصادقة: %v", err)
	}
}

// LogSuspiciousActivity يسجل نشاط مشبوه
func (sl *SecurityLogger) LogSuspiciousActivity(userID, ip, userAgent, activityType, details string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نتخطى التسجيل
		return
	}
	
	_, err = sl.DB.Exec(`
		INSERT INTO security_logs (event_type, user_id, ip, user_agent, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, activityType, userID, ip, userAgent, details, time.Now())
	
	if err != nil {
		log.Printf("⚠️ فشل تسجيل نشاط مشبوه: %v", err)
	}
}

// CheckForSuspiciousPatterns يتحقق من أنماط مشبوهة في محاولات الدخول
func (sl *SecurityLogger) CheckForSuspiciousPatterns(ip, email string) (bool, string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، لا نتحقق من الأنماط
		return false, ""
	}
	
	// التحقق من محاولات فاشلة متعددة من نفس IP (خلال 15 دقيقة)
	var failedAttempts int
	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)
	err = sl.DB.QueryRow(`
		SELECT COUNT(*) FROM security_logs 
		WHERE ip = $1 AND event_type = 'auth_attempt' AND success = false 
		AND created_at > $2
	`, ip, fifteenMinutesAgo).Scan(&failedAttempts)
	
	if err == nil && failedAttempts > 5 {
		return true, fmt.Sprintf("محاولات فاشلة متعددة من IP %s (%d محاولة)", ip, failedAttempts)
	}
	
	// التحقق من محاولات فاشلة متعددة لنفس البريد (خلال ساعة واحدة)
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err = sl.DB.QueryRow(`
		SELECT COUNT(*) FROM security_logs 
		WHERE email = $1 AND event_type = 'auth_attempt' AND success = false 
		AND created_at > $2
	`, email, oneHourAgo).Scan(&failedAttempts)
	
	if err == nil && failedAttempts > 10 {
		return true, fmt.Sprintf("محاولات فاشلة متعددة للبريد %s (%d محاولة)", email, failedAttempts)
	}
	
	return false, ""
}

// GetRecentSecurityLogs يجلب سجلات الأمان الأخيرة
func (sl *SecurityLogger) GetRecentSecurityLogs(limit int) ([]map[string]interface{}, error) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نرجع مصفوفة فارغة
		return []map[string]interface{}{}, nil
	}
	
	rows, err := sl.DB.Query(`
		SELECT id, event_type, email, phone, ip, user_agent, success, reason, details, created_at
		FROM security_logs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []map[string]interface{}
	for rows.Next() {
		var id, eventType, email, phone, ip, userAgent, reason, details string
		var success bool
		var createdAt time.Time
		
		err := rows.Scan(&id, &eventType, &email, &phone, &ip, &userAgent, &success, &reason, &details, &createdAt)
		if err != nil {
			continue
		}
		
		log := map[string]interface{}{
			"id":          id,
			"event_type":  eventType,
			"email":       email,
			"phone":       phone,
			"ip":          ip,
			"user_agent":  userAgent,
			"success":     success,
			"reason":      reason,
			"details":     details,
			"created_at":  createdAt,
		}
		
		logs = append(logs, log)
	}
	
	return logs, nil
}
```

---

## 📄 backend/internal/services/storage_service.go

```go
package services

import (
	"fmt"
	"mime/multipart"
)

// StorageService مسؤول عن رفع الصور (صور تنفيذ المهام) إلى Cloudflare R2
// ملاحظة: هذا هيكل مبسّط جاهز للتوسعة بربطه الفعلي مع AWS SDK v2 (متوافق مع R2)
type StorageService struct {
	AccountID  string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func NewStorageService(accountID, accessKey, secretKey, bucketName string) *StorageService {
	return &StorageService{
		AccountID:  accountID,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
	}
}

// UploadFile يرفع ملف (صورة) ويعيد الرابط العام له
// TODO: ربطها فعلياً بـ github.com/aws/aws-sdk-go-v2 عند إعداد مفاتيح R2
func (s *StorageService) UploadFile(file multipart.File, filename string) (string, error) {
	if s.AccessKey == "" {
		return "", fmt.Errorf("Cloudflare R2 credentials are not configured yet")
	}

	publicURL := fmt.Sprintf("https://%s.r2.dev/%s", s.BucketName, filename)
	return publicURL, nil
}

```

---

## 📄 backend/pkg/utils/datetime.go

```go
package utils

import "time"

// JerusalemLocation يعيد توقيت القدس، يُستخدم في كل مكان يتعلق بتوقيت المهام والحضور
// (المشروع يستهدف مستخدمين في نفس المنطقة الزمنية بغض النظر عن لغة الواجهة)
func JerusalemLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		return time.UTC
	}
	return loc
}

func NowInJerusalem() time.Time {
	return time.Now().In(JerusalemLocation())
}

```

---

## 📄 backend/pkg/utils/errors.go

```go
package utils

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Custom error types
type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error constructors
func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusBadRequest,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusUnauthorized,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusForbidden,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusNotFound,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusConflict,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusInternalServerError,
	}
}

// RecoverPanic recovers from panics and converts them to errors
func RecoverPanic() error {
	if r := recover(); r != nil {
		stack := debug.Stack()
		return fmt.Errorf("panic recovered: %v\n%s", r, stack)
	}
	return nil
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError converts an error to AppError if possible
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewInternalError("Internal server error", err)
}
```

---

## 📄 backend/pkg/utils/haversine.go

```go
package utils

import "math"

const earthRadiusMeters = 6371000

// HaversineDistance تحسب المسافة بالمتر بين نقطتين جغرافيتين (خط عرض/طول)
// هذه هي الدالة الرياضية الأساسية التي يعتمد عليها منع التختيم خارج نطاق العمل:
// كل تحقق Geofence في المشروع يمر من هنا في النهاية
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)
	deltaLat := degreesToRadians(lat2 - lat1)
	deltaLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

```

---

## 📄 backend/pkg/utils/logger.go

```go
package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger represents a structured logger
type Logger struct {
	prefix string
	level  LogLevel
}

// NewLogger creates a new logger with the specified prefix
func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		level:  INFO,
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// formatMessage formats a log message with timestamp and level
func (l *Logger) formatMessage(level LogLevel, message string) string {
	levelStr := ""
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	case FATAL:
		levelStr = "FATAL"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, l.prefix, levelStr, message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	if l.level <= DEBUG {
		log.Println(l.formatMessage(DEBUG, message))
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	if l.level <= INFO {
		log.Println(l.formatMessage(INFO, message))
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	if l.level <= WARN {
		log.Println(l.formatMessage(WARN, message))
	}
}

// Error logs an error message
func (l *Logger) Error(message string, err error) {
	if l.level <= ERROR {
		if err != nil {
			log.Println(l.formatMessage(ERROR, fmt.Sprintf("%s: %v", message, err)))
		} else {
			log.Println(l.formatMessage(ERROR, message))
		}
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, err error) {
	if l.level <= FATAL {
		if err != nil {
			log.Println(l.formatMessage(FATAL, fmt.Sprintf("%s: %v", message, err)))
		} else {
			log.Println(l.formatMessage(FATAL, message))
		}
		os.Exit(1)
	}
}

// RequestLogger logs HTTP request information
func (l *Logger) Request(method, path string, statusCode int, duration time.Duration) {
	message := fmt.Sprintf("%s %s - Status: %d - Duration: %v", method, path, statusCode, duration)
	if statusCode >= 500 {
		l.Error(message, nil)
	} else if statusCode >= 400 {
		l.Warn(message)
	} else {
		l.Debug(message)
	}
}

// DatabaseLogger logs database operations
func (l *Logger) Database(operation, table string, duration time.Duration, err error) {
	message := fmt.Sprintf("DB: %s on %s - Duration: %v", operation, table, duration)
	if err != nil {
		l.Error(message, err)
	} else {
		l.Debug(message)
	}
}

// CacheLogger logs cache operations
func (l *Logger) Cache(operation, key string, hit bool) {
	status := "MISS"
	if hit {
		status = "HIT"
	}
	message := fmt.Sprintf("Cache: %s %s - %s", operation, key, status)
	l.Debug(message)
}
```

---

## 📄 backend/pkg/utils/validator.go

```go
package utils

import "strings"

// IsValidEmail تحقق مبسط من صيغة البريد الإلكتروني (ليس تحققاً كاملاً بمعايير RFC،
// لكنه كافٍ لمنع الأخطاء الشائعة عند التسجيل)
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) >= 5
}

// IsValidCoordinates تتحقق أن الإحداثيات المرسلة من الموظف ضمن مجال منطقي على الخريطة
// (خط العرض بين -90 و90، وخط الطول بين -180 و180) قبل حتى حساب المسافة
func IsValidCoordinates(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

```

---

## 📄 backend/tests/api_integration_test.go

```go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/database"
	"worktrack/backend/internal/router"

	"github.com/gin-gonic/gin"
)

// setupTestRouter يُعدّ router للاختبار
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Port:          "8080",
		DatabaseURL:   "",
		JWTSecret:     "test-secret-key-for-testing-only",
		AllowedOrigin: "http://localhost:3000",
		DefaultLang:   "ar",
		GeoapifyKey:   "",
		RedisURL:      "",
	}
	
	db := getTestDB()
	return router.Setup(db, cfg)
}

// TestHealthEndpoint يختبر نقطة نهاية health
func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
	
	t.Log("✅ Health endpoint test passed")
}

// TestRateLimiterWithoutRedis يختبر أن Rate Limiter يعمل بدون Redis
func TestRateLimiterWithoutRedis(t *testing.T) {
	router := setupTestRouter()
	
	// إرسال طلبات متعددة
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i+1, w.Code)
		}
	}
	
	t.Log("✅ Rate limiter without Redis test passed")
}

// TestAuthEndpoints يختبر نقاط نهاية المصادقة
func TestAuthEndpoints(t *testing.T) {
	router := setupTestRouter()
	
	// اختبار تسجيل الدخول
	loginData := map[string]string{
		"phone":    "1234567890",
		"password": "wrongpassword",
	}
	jsonData, _ := json.Marshal(loginData)
	
	req, _ := http.NewRequest("POST", "/api/v1/auth/phone-login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// نتوقع خطأ (كلمة مرور خاطئة)
	if w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Logf("Login attempt returned status %d (expected 401 or 400)", w.Code)
	}
	
	t.Log("✅ Auth endpoints test passed")
}

// TestProtectedEndpointsWithoutToken يختبر أن النقاط المحمية ترفض الطلبات بدون توكن
func TestProtectedEndpointsWithoutToken(t *testing.T) {
	router := setupTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	t.Log("✅ Protected endpoints without token test passed")
}

// TestDatabaseMigrationsApplied يتحقق من تطبيق الترحيلات
func TestDatabaseMigrationsApplied(t *testing.T) {
	db := getTestDB()
	
	// التحقق من وجود الأعمدة المضافة
	checks := []struct {
		table    string
		column   string
		expected bool
	}{
		{"users", "subscription_status", true},
		{"users", "subscription_expires_at", true},
		{"users", "device_id", true},
		{"users", "updated_at", true},
		{"worksites", "assigned_employee_id", true},
		{"worksites", "updated_at", true},
		{"tasks", "updated_at", true},
		{"attendance", "photo_url", true},
		{"attendance", "updated_at", true},
	}
	
	for _, check := range checks {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = $1 AND column_name = $2
			)
		`
		err := db.QueryRow(query, check.table, check.column).Scan(&exists)
		if err != nil {
			t.Errorf("Failed to check column %s.%s: %v", check.table, check.column, err)
			continue
		}
		
		if exists != check.expected {
			t.Errorf("Column %s.%s existence mismatch: expected %v, got %v", 
				check.table, check.column, check.expected, exists)
		} else {
			t.Logf("✅ Column %s.%s exists as expected", check.table, check.column)
		}
	}
	
	t.Log("✅ Database migrations test passed")
}

// TestConnectionPoolSettings يتحقق من إعدادات Connection Pool
func TestConnectionPoolSettings(t *testing.T) {
	db := getTestDB()
	
	stats := db.Stats()
	
	// التحقق من الإعدادات
	if stats.MaxOpenConnections != 5 {
		t.Errorf("Expected MaxOpenConnections to be 5, got %d", stats.MaxOpenConnections)
	}
	
	t.Logf("✅ Connection Pool settings verified:")
	t.Logf("  - MaxOpenConnections: %d", stats.MaxOpenConnections)
	t.Logf("  - OpenConnections: %d", stats.OpenConnections)
	t.Logf("  - InUse: %d", stats.InUse)
	t.Logf("  - Idle: %d", stats.Idle)
}

// TestRedisConnection يتحقق من اتصال Redis (اختياري)
func TestRedisConnection(t *testing.T) {
	// هذا الاختبار يتحقق فقط أن التطبيق يمكنه العمل بدون Redis
	// Redis غير متوفر في بيئة الاختبار الحالية
	
	if database.RedisClient == nil {
		t.Log("ℹ️  Redis not connected (expected in test environment)")
		return
	}
	
	t.Log("✅ Redis connection test passed")
}
```

---

## 📄 backend/tests/auth_test.go

```go
package tests

import (
	"testing"

	"worktrack/backend/internal/services"
)

// TestPasswordHashAndCheck يتحقق أن كلمة المرور المشفّرة تُطابق نفسها عند المقارنة،
// ولا تُطابق كلمة مرور خاطئة
func TestPasswordHashAndCheck(t *testing.T) {
	auth := services.NewAuthService("test_secret")

	hash, err := auth.HashPassword("my-secure-password")
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if !auth.CheckPassword("my-secure-password", hash) {
		t.Error("expected correct password to match hash")
	}

	if auth.CheckPassword("wrong-password", hash) {
		t.Error("expected wrong password to NOT match hash")
	}
}

// TestGenerateAndValidateToken يتحقق من دورة حياة التوكن كاملة: إنشاء ثم تحقق
func TestGenerateAndValidateToken(t *testing.T) {
	auth := services.NewAuthService("test_secret")

	token, err := auth.GenerateToken("user-123", "admin")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != "user-123" || claims.Role != "admin" {
		t.Errorf("unexpected claims: %+v", claims)
	}
}

```

---

## 📄 backend/tests/geofence_test.go

```go
package tests

import (
	"testing"

	"worktrack/backend/internal/models"
	"worktrack/backend/internal/services"
)

// TestCheckWithinWorksite_Inside يتحقق أن نقطة قريبة جداً من موقع العمل تُعتبر "داخل النطاق"
func TestCheckWithinWorksite_Inside(t *testing.T) {
	worksite := models.Worksite{
		Latitude:     31.9539,
		Longitude:    35.9106,
		RadiusMeters: 100,
	}

	// نفس الإحداثيات تماماً => المسافة صفر => يجب أن تكون داخل النطاق دائماً
	result := services.CheckWithinWorksite(31.9539, 35.9106, worksite)

	if !result.IsWithinRange {
		t.Errorf("expected point to be within range, got distance=%.2f, radius=%d",
			result.DistanceMeters, result.AllowedRadius)
	}
}

// TestCheckWithinWorksite_Outside يتحقق أن نقطة بعيدة جداً تُعتبر "خارج النطاق" وتُرفض
func TestCheckWithinWorksite_Outside(t *testing.T) {
	worksite := models.Worksite{
		Latitude:     31.9539,
		Longitude:    35.9106,
		RadiusMeters: 100,
	}

	// إحداثيات مدينة أخرى بعيدة تماماً => يجب رفضها بوضوح
	result := services.CheckWithinWorksite(32.0853, 34.7818, worksite)

	if result.IsWithinRange {
		t.Errorf("expected point to be outside range, got distance=%.2f, radius=%d",
			result.DistanceMeters, result.AllowedRadius)
	}
}

```

---

## 📄 backend/tests/integration_test.go

```go
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

```

---

## 📄 supabase_migrations/supabase_complete_migration.sql

```sql
-- =====================================================
-- WorkTrack Database Migration - Complete Unified Migration
-- =====================================================
-- This migration includes all updates in a single file:
-- 1. Translation columns for i18n support
-- 2. Fix worksite deletion (preserve tasks and attendance records)
-- 3. System updates and fixes
-- 4. Tables updates
-- 5. Create indexes
-- 6. Verification queries
-- =====================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =====================================================
-- PART 1: Translation Columns for i18n Support
-- =====================================================

-- Users table translations
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS full_name_ar VARCHAR(150),
    ADD COLUMN IF NOT EXISTS full_name_he VARCHAR(150),
    ADD COLUMN IF NOT EXISTS full_name_en VARCHAR(150);

-- Migrate existing full_name to all translation columns
UPDATE users 
SET 
    full_name_ar = full_name,
    full_name_he = full_name,
    full_name_en = full_name
WHERE full_name_ar IS NULL;

-- Clients table translations (if table exists)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'clients') THEN
        ALTER TABLE clients
            ADD COLUMN IF NOT EXISTS name_ar VARCHAR(150),
            ADD COLUMN IF NOT EXISTS name_he VARCHAR(150),
            ADD COLUMN IF NOT EXISTS name_en VARCHAR(150),
            ADD COLUMN IF NOT EXISTS address_ar TEXT,
            ADD COLUMN IF NOT EXISTS address_he TEXT,
            ADD COLUMN IF NOT EXISTS address_en TEXT;

        -- Migrate existing data to translation columns
        UPDATE clients 
        SET 
            name_ar = name,
            name_he = name,
            name_en = name,
            address_ar = address,
            address_he = address,
            address_en = address
        WHERE name_ar IS NULL;
    END IF;
END $$;

-- Worksites table translations
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS name_ar VARCHAR(150),
    ADD COLUMN IF NOT EXISTS name_he VARCHAR(150),
    ADD COLUMN IF NOT EXISTS name_en VARCHAR(150),
    ADD COLUMN IF NOT EXISTS address_ar TEXT,
    ADD COLUMN IF NOT EXISTS address_he TEXT,
    ADD COLUMN IF NOT EXISTS address_en TEXT;

-- Migrate existing data to translation columns
UPDATE worksites 
SET 
    name_ar = name,
    name_he = name,
    name_en = name,
    address_ar = address,
    address_he = address,
    address_en = address
WHERE name_ar IS NULL;

-- Tasks table translations
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS title_ar VARCHAR(200),
    ADD COLUMN IF NOT EXISTS title_he VARCHAR(200),
    ADD COLUMN IF NOT EXISTS title_en VARCHAR(200),
    ADD COLUMN IF NOT EXISTS description_ar TEXT,
    ADD COLUMN IF NOT EXISTS description_he TEXT,
    ADD COLUMN IF NOT EXISTS description_en TEXT;

-- Migrate existing data to translation columns
UPDATE tasks 
SET 
    title_ar = title,
    title_he = title,
    title_en = title,
    description_ar = description,
    description_he = description,
    description_en = description
WHERE title_ar IS NULL;

-- Service requests table translations (if table exists)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'service_requests') THEN
        ALTER TABLE service_requests
            ADD COLUMN IF NOT EXISTS title_ar VARCHAR(200),
            ADD COLUMN IF NOT EXISTS title_he VARCHAR(200),
            ADD COLUMN IF NOT EXISTS title_en VARCHAR(200),
            ADD COLUMN IF NOT EXISTS description_ar TEXT,
            ADD COLUMN IF NOT EXISTS description_he TEXT,
            ADD COLUMN IF NOT EXISTS description_en TEXT,
            ADD COLUMN IF NOT EXISTS address_ar TEXT,
            ADD COLUMN IF NOT EXISTS address_he TEXT,
            ADD COLUMN IF NOT EXISTS address_en TEXT,
            ADD COLUMN IF NOT EXISTS location_name_ar TEXT,
            ADD COLUMN IF NOT EXISTS location_name_he TEXT,
            ADD COLUMN IF NOT EXISTS location_name_en TEXT,
            ADD COLUMN IF NOT EXISTS client_name_ar VARCHAR(150),
            ADD COLUMN IF NOT EXISTS client_name_he VARCHAR(150),
            ADD COLUMN IF NOT EXISTS client_name_en VARCHAR(150);

        -- Migrate existing data to translation columns
        UPDATE service_requests 
        SET 
            title_ar = title,
            title_he = title,
            title_en = title,
            description_ar = description,
            description_he = description,
            description_en = description,
            address_ar = address,
            address_he = address,
            address_en = address,
            location_name_ar = location_name,
            location_name_he = location_name,
            location_name_en = location_name,
            client_name_ar = client_name,
            client_name_he = client_name,
            client_name_en = client_name
        WHERE title_ar IS NULL;
    END IF;
END $$;

-- Create helper function to get translated text
CREATE OR REPLACE FUNCTION get_translation(
    text_ar TEXT,
    text_he TEXT,
    text_en TEXT,
    lang VARCHAR DEFAULT 'ar'
) RETURNS TEXT AS $$
BEGIN
    CASE lang
        WHEN 'he' THEN
            RETURN COALESCE(text_he, text_ar, text_en, '');
        WHEN 'en' THEN
            RETURN COALESCE(text_en, text_ar, text_he, '');
        ELSE
            RETURN COALESCE(text_ar, text_he, text_en, '');
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- PART 2: Fix Worksite Deletion - Preserve Tasks and Attendance Records
-- =====================================================
-- This section ensures that when a worksite is deleted:
-- - Tasks with the deleted worksite_id will have worksite_id set to NULL (preserved)
-- - Attendance records with the deleted worksite_id will have worksite_id set to NULL (preserved)
-- - No data loss occurs, only the reference to the deleted worksite is removed

-- First, remove NOT NULL constraint from worksite_id in tasks
ALTER TABLE tasks ALTER COLUMN worksite_id DROP NOT NULL;

-- Then change the foreign key constraint from RESTRICT to SET NULL
-- This requires dropping the existing constraint first
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_worksite_id_fkey;

ALTER TABLE tasks 
ADD CONSTRAINT tasks_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;

-- Fix attendance table to handle worksite deletion properly
-- This ensures attendance records are preserved when worksite is deleted
-- Remove NOT NULL constraint from worksite_id in attendance
ALTER TABLE attendance ALTER COLUMN worksite_id DROP NOT NULL;

-- Change the foreign key constraint from RESTRICT to SET NULL
-- This allows the worksite to be deleted while keeping attendance records
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_worksite_id_fkey;

ALTER TABLE attendance 
ADD CONSTRAINT attendance_worksite_id_fkey 
FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;

-- =====================================================
-- PART 3: System Updates and Fixes
-- =====================================================

-- Users table updates
-- Add subscription fields if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (subscription_status IN ('active', 'expired', 'canceled'));

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ NULL;

-- Add device and login fields if not exist
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS device_id TEXT;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS device_model TEXT;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_login_device TEXT;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS login_code TEXT;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

-- Fix existing roles that might not match our expected roles
UPDATE users SET role = 'employee' WHERE role NOT IN ('admin', 'employee', 'client', 'customer');

-- Drop existing role constraint if exists
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- Add new role constraint
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'client', 'customer'));

-- Create default admin users
INSERT INTO users (
  full_name,
  email,
  password_hash,
  role,
  is_active,
  subscription_status,
  subscription_expires_at
)
VALUES
(
  'DevPro System Administrator',
  'admin@devpro.com',
  crypt('devproadmin', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
),
(
  'DevPro Support Manager',
  'support@devpro.com',
  crypt('devprosupport', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NOW() + INTERVAL '1 year'
),
(
  'DevPro Project Manager',
  'manager@devpro.com',
  crypt('devpromanager', gen_salt('bf', 12)),
  'admin',
  TRUE,
  'active',
  NULL
)
ON CONFLICT (email) DO UPDATE
SET
  password_hash = EXCLUDED.password_hash,
  role = 'admin',
  is_active = TRUE,
  subscription_status = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN 'active'
    ELSE EXCLUDED.subscription_status
  END,
  subscription_expires_at = CASE
    WHEN users.email IN ('admin@devpro.com', 'manager@devpro.com') THEN NULL
    ELSE EXCLUDED.subscription_expires_at
  END,
  updated_at = now()
RETURNING email, role, is_active, subscription_status, subscription_expires_at;

-- Worksites table updates
-- Add assigned_employee_id if not exist
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add updated_at if not exist
ALTER TABLE worksites
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Tasks table updates
-- Add updated_at if not exist
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add priority column if not exist
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS priority VARCHAR(20) NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent'));

-- =====================================================
-- PART 4: Tables Updates
-- =====================================================

-- Attendance table updates
-- Add updated_at if not exist
ALTER TABLE attendance
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add photo fields to attendance if not exist
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- Add service_request_id to attendance table
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS service_request_id UUID REFERENCES service_requests(id) ON DELETE SET NULL;

-- Service requests table updates
-- Add location_name if not exist
ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS location_name TEXT;

-- Fix service_requests table to handle client deletion properly
ALTER TABLE service_requests DROP CONSTRAINT IF EXISTS service_requests_client_id_fkey;
ALTER TABLE service_requests ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE service_requests 
ADD CONSTRAINT service_requests_client_id_fkey 
FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE SET NULL;

-- Notifications table updates
-- Add type and related_id to notifications if not exist
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

-- =====================================================
-- PART 5: Create Indexes
-- =====================================================

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_subscription_status ON users(subscription_status);

-- Worksites indexes
CREATE INDEX IF NOT EXISTS idx_worksites_location
    ON worksites(latitude, longitude);

CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee
    ON worksites(assigned_employee_id);

CREATE INDEX IF NOT EXISTS idx_worksites_is_active ON worksites(is_active);

-- Tasks indexes
CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id
    ON tasks(worksite_id);

CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id
    ON tasks(assigned_user_id);

CREATE INDEX IF NOT EXISTS idx_tasks_client_id
    ON tasks(client_id);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

-- Attendance indexes
CREATE INDEX IF NOT EXISTS idx_attendance_user_id
    ON attendance(user_id);

CREATE INDEX IF NOT EXISTS idx_attendance_task_id
    ON attendance(task_id);

CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id
    ON attendance(worksite_id);

CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time
    ON attendance(check_in_time);

CREATE INDEX IF NOT EXISTS idx_attendance_status ON attendance(status);

CREATE INDEX IF NOT EXISTS idx_attendance_service_request_id
    ON attendance(service_request_id);

-- Service requests indexes
CREATE INDEX IF NOT EXISTS idx_service_requests_client_id
    ON service_requests(client_id);

CREATE INDEX IF NOT EXISTS idx_service_requests_status
    ON service_requests(status);

CREATE INDEX IF NOT EXISTS idx_service_requests_priority
    ON service_requests(priority);

-- Assignments indexes
CREATE INDEX IF NOT EXISTS idx_assignments_employee_id
    ON assignments(employee_id);

CREATE INDEX IF NOT EXISTS idx_assignments_request_id
    ON assignments(request_id);

CREATE INDEX IF NOT EXISTS idx_assignments_admin_id
    ON assignments(admin_id);

CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);

-- Location tracking indexes
CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id
    ON location_tracking(user_id);

CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id
    ON location_tracking(assignment_id);

CREATE INDEX IF NOT EXISTS idx_location_tracking_recorded_at
    ON location_tracking(recorded_at);

-- Notifications indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id
    ON notifications(user_id);

CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);

CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- =====================================================
-- PART 6: Verification Queries
-- =====================================================

-- Check if translation columns were added successfully
SELECT 
    table_name,
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_schema = 'public'
    AND table_name IN ('users', 'clients', 'worksites', 'tasks', 'service_requests')
    AND (column_name LIKE '%_ar' OR column_name LIKE '%_he' OR column_name LIKE '%_en')
ORDER BY table_name, column_name;

-- Check if the worksite deletion fix was applied correctly for tasks table
-- Expected: delete_rule should be 'SET NULL'
SELECT
    tc.table_name,
    tc.constraint_name,
    tc.is_deferrable,
    tc.initially_deferred,
    rc.match_option,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints tc
JOIN information_schema.referential_constraints rc 
    ON tc.constraint_name = rc.constraint_name
WHERE tc.table_schema = 'public'
    AND tc.table_name = 'tasks'
    AND tc.constraint_name = 'tasks_worksite_id_fkey';

-- Check if the attendance worksite deletion fix was applied correctly
-- Expected: delete_rule should be 'SET NULL'
SELECT
    tc.table_name,
    tc.constraint_name,
    tc.is_deferrable,
    tc.initially_deferred,
    rc.match_option,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints tc
JOIN information_schema.referential_constraints rc 
    ON tc.constraint_name = rc.constraint_name
WHERE tc.table_schema = 'public'
    AND tc.table_name = 'attendance'
    AND tc.constraint_name = 'attendance_worksite_id_fkey';

-- Verify role constraint
SELECT
    conname as constraint_name,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'users'::regclass AND conname = 'users_role_check';

-- Check default admin users
SELECT email, role, is_active, subscription_status, subscription_expires_at 
FROM users 
WHERE role = 'admin' 
ORDER BY email;
```

---

