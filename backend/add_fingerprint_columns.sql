-- =============================================
-- إضافة أعمدة بصمة الجهاز
-- =============================================

-- إضافة العمود إذا لم يكن موجوداً
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_fingerprint TEXT;

-- عرض الأعمدة
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('device_id', 'device_model', 'device_fingerprint')
ORDER BY ordinal_position;

-- عرض المستخدمين مع بصماتهم
SELECT 
    full_name,
    phone,
    device_model,
    CASE 
        WHEN device_fingerprint IS NOT NULL THEN '✅ مسجل'
        ELSE '❌ غير مسجل'
    END as fingerprint_status
FROM users 
WHERE role = 'employee'
ORDER BY created_at DESC;
