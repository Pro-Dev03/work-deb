-- =====================================================
-- سكريبت قاعدة بيانات WorkTrack
-- يتوافق مع هيكل المشروع الحالي
-- =====================================================

-- تفعيل الامتدادات المطلوبة
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- جدول المستخدمين
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    phone VARCHAR(30),
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee'
        CHECK (role IN ('admin', 'employee', 'client')),
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (subscription_status IN ('active', 'expired', 'canceled')),
    subscription_expires_at TIMESTAMPTZ NULL,

    device_id TEXT,
    device_model TEXT,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_device TEXT,
    login_code TEXT,
    login_code_expires TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- إضافة حقول الاشتراك إذا لم تكن موجودة (للتوافق مع قواعد البيانات الموجودة)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (subscription_status IN ('active', 'expired', 'canceled'));
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ NULL;

-- إضافة حقول الجهاز وتسجيل الدخول إذا لم تكن موجودة
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

-- تحديث قيد الأدوار لإضافة دور client
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'employee', 'client'));

-- =====================================================
-- إنشاء مستخدمي المديرين الافتراضيين
-- =====================================================
-- ملاحظة: morad@worktrack.com لديه اشتراك مدى الحياة (NULL)
--         admin@worktrack.com لديه اشتراك لمدة سنة
--         morad-admin@worktrack.com لديه اشتراك مدى الحياة (NULL)
--
-- كلمات المرور:
-- - morad@worktrack.com → adminmorad
-- - admin@worktrack.com → adminadmin
-- - morad-admin@worktrack.com → admin123
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
  'مدير النظام - مراد',
  'morad@worktrack.com',
  crypt('adminmorad', gen_salt('bf', 12)),  -- كلمة المرور: adminmorad
  'admin',
  TRUE,
  'active',
  NULL  -- اشتراك مدى الحياة
),
(
  'مدير النظام - الرئيسي',
  'admin@worktrack.com',
  crypt('adminadmin', gen_salt('bf', 12)),  -- كلمة المرور: adminadmin
  'admin',
  TRUE,
  'active',
  NOW() + INTERVAL '1 year'  -- اشتراك لمدة سنة
),
(
  'مدير النظام - مراد (إضافي)',
  'morad-admin@worktrack.com',
  crypt('admin123', gen_salt('bf', 12)),  -- كلمة المرور: admin123
  'admin',
  TRUE,
  'active',
  NULL  -- اشتراك مدى الحياة
)
ON CONFLICT (email) DO UPDATE
SET
  password_hash = EXCLUDED.password_hash,
  role = 'admin',
  is_active = TRUE,
  subscription_status = CASE
    WHEN users.email IN ('morad@worktrack.com', 'morad-admin@worktrack.com') THEN 'active'
    ELSE EXCLUDED.subscription_status
  END,
  subscription_expires_at = CASE
    WHEN users.email IN ('morad@worktrack.com', 'morad-admin@worktrack.com') THEN NULL  -- اشتراك مدى الحياة
    ELSE EXCLUDED.subscription_expires_at
  END,
  updated_at = now()
RETURNING email, role, is_active, subscription_status, subscription_expires_at;

-- التحقق من المدراء والاشتراكات
SELECT 
  id,
  full_name,
  email,
  role,
  is_active,
  subscription_status,
  subscription_expires_at,
  CASE
    WHEN subscription_expires_at IS NULL THEN 'دائم'
    WHEN subscription_expires_at > NOW() THEN 'نشط'
    ELSE 'منتهي'
  END as current_status
FROM users
WHERE role = 'admin';

-- التحقق من مستخدم محدد (مثال: morad@worktrack.com)
-- يمكنك تغيير البريد الإلكتروني للتحقق من أي مستخدم آخر
SELECT
  '=== التحقق من حالة اشتراك morad@worktrack.com ===' as info;

SELECT
  id,
  full_name,
  email,
  role,
  is_active,
  subscription_status,
  subscription_expires_at,
  CASE
    WHEN subscription_expires_at IS NULL THEN 'دائم (Lifetime)'
    WHEN subscription_expires_at > NOW() THEN 'نشط - المتبقي: ' ||
      EXTRACT(DAY FROM subscription_expires_at - NOW()) || ' يوم'
    ELSE 'منتهي'
  END as current_status,
  CASE
    WHEN subscription_expires_at IS NULL THEN TRUE
    WHEN subscription_expires_at > NOW() THEN TRUE
    ELSE FALSE
  END as can_login
FROM users
WHERE email = 'morad@worktrack.com';

-- التحقق من المستخدم الجديد morad-admin
SELECT
  '=== التحقق من حالة اشتراك morad-admin@worktrack.com ===' as info;

SELECT
  id,
  full_name,
  email,
  role,
  is_active,
  subscription_status,
  subscription_expires_at,
  CASE
    WHEN subscription_expires_at IS NULL THEN 'دائم (Lifetime)'
    WHEN subscription_expires_at > NOW() THEN 'نشط - المتبقي: ' ||
      EXTRACT(DAY FROM subscription_expires_at - NOW()) || ' يوم'
    ELSE 'منتهي'
  END as current_status,
  CASE
    WHEN subscription_expires_at IS NULL THEN TRUE
    WHEN subscription_expires_at > NOW() THEN TRUE
    ELSE FALSE
  END as can_login
FROM users
WHERE email = 'morad-admin@worktrack.com';

-- استعلام عام للتحقق من أي مستخدم (غير التعليق واستبدل البريد)
-- SELECT * FROM users WHERE email = 'your_email@example.com';

-- جدول العملاء
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(150),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- جدول نقاط العمل
CREATE TABLE IF NOT EXISTS worksites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius_meters INTEGER NOT NULL DEFAULT 100
        CHECK (radius_meters > 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- جدول المهام
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
    worksite_id UUID NOT NULL REFERENCES worksites(id) ON DELETE RESTRICT,
    assigned_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'in_progress', 'completed', 'late', 'cancelled')),
    scheduled_start TIMESTAMPTZ,
    scheduled_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- جدول الحضور والانصراف
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- إضافة حقول الصور إلى جدول الحضور إذا لم تكن موجودة
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- جدول طلبات الخدمة
CREATE TABLE IF NOT EXISTS service_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- جدول تعيينات الموظفين
CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- جدول تتبع الموقع
CREATE TABLE IF NOT EXISTS location_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- جدول الإشعارات
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'info',
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    related_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- إضافة حقول type و related_id إلى جدول الإشعارات إذا لم تكن موجودة
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info';
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS related_id UUID;

-- إنشاء الفهارس
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE INDEX IF NOT EXISTS idx_worksites_location
    ON worksites(latitude, longitude);

CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id
    ON tasks(worksite_id);

CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id
    ON tasks(assigned_user_id);

CREATE INDEX IF NOT EXISTS idx_attendance_user_id
    ON attendance(user_id);

CREATE INDEX IF NOT EXISTS idx_attendance_task_id
    ON attendance(task_id);

CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time
    ON attendance(check_in_time);

CREATE INDEX IF NOT EXISTS idx_service_requests_client_id
    ON service_requests(client_id);

CREATE INDEX IF NOT EXISTS idx_service_requests_status
    ON service_requests(status);

CREATE INDEX IF NOT EXISTS idx_assignments_employee_id
    ON assignments(employee_id);

CREATE INDEX IF NOT EXISTS idx_assignments_request_id
    ON assignments(request_id);

CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id
    ON location_tracking(user_id);

CREATE INDEX IF NOT EXISTS idx_location_tracking_assignment_id
    ON location_tracking(assignment_id);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id
    ON notifications(user_id);

-- =====================================================
-- تعليمات الاستخدام
-- =====================================================
-- 1. تشغيل السكريبت على قاعدة بيانات جديدة:
--    psql -U username -d database_name -f schema_updated.sql
--
-- 2. تشغيل السكريبت على قاعدة بيانات موجودة (آمن):
--    سيتم تجاهل الجداول الموجودة وإضافة الحقول المفقودة فقط
--
-- 3. بيانات تسجيل الدخول الافتراضية:
--    - morad@worktrack.com / adminmorad (اشتراك مدى الحياة)
--    - admin@worktrack.com / adminadmin (اشتراك سنة واحدة)
--    - morad-admin@worktrack.com / admin123 (اشتراك مدى الحياة) - جديد
--
-- 4. السكريبت آمن للتشغيل المتعدد - يستخدم IF NOT EXISTS و ON CONFLICT
-- =====================================================
