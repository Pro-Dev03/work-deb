-- =============================================
-- إضافة أعمدة الأمان للتحقق من الجهاز
-- =============================================

-- إضافة أعمدة جديدة
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_fingerprint TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_user_agent TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_ip_address TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_first_seen TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_last_seen TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_ip TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_user_agent TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_count INTEGER DEFAULT 0;

-- عرض الأعمدة الجديدة
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name LIKE 'device_%' OR column_name LIKE 'last_login_%' OR column_name = 'login_count'
ORDER BY ordinal_position;

-- عرض المستخدمين مع معلومات الأمان
SELECT 
    full_name,
    phone,
    device_model,
    CASE 
        WHEN device_fingerprint IS NOT NULL THEN '✅ مسجل'
        ELSE '❌ غير مسجل'
    END as fingerprint_status,
    login_count,
    device_last_seen
FROM users 
WHERE role = 'employee'
ORDER BY created_at DESC;
