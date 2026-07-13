-- =============================================
-- إضافة عمود device_model لتخزين نوع الجهاز
-- =============================================

-- إضافة العمود
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_model TEXT;

-- عرض الأعمدة الجديدة
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('device_id', 'device_model', 'phone_verified')
ORDER BY ordinal_position;

-- عرض المستخدمين مع معلومات أجهزتهم
SELECT 
    full_name,
    phone,
    device_id,
    device_model,
    phone_verified
FROM users 
WHERE role = 'employee'
ORDER BY created_at DESC;
