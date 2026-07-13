-- =============================================
-- التحقق من أعمدة الجهاز
-- =============================================

-- عرض الأعمدة الموجودة
SELECT column_name, data_type, is_nullable
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('device_id', 'device_model', 'phone_verified', 'device_registered_at')
ORDER BY ordinal_position;

-- إذا كان العمود device_model غير موجود، أضفه
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'device_model'
    ) THEN
        ALTER TABLE users ADD COLUMN device_model TEXT;
        RAISE NOTICE '✅ تم إضافة عمود device_model';
    END IF;
END $$;

-- عرض المستخدمين مع معلومات أجهزتهم
SELECT 
    full_name,
    phone,
    device_id,
    device_model,
    phone_verified,
    CASE 
        WHEN device_id IS NOT NULL AND device_model IS NOT NULL THEN '✅ مسجل بالكامل'
        WHEN device_id IS NOT NULL THEN '⚠️ مسجل جزئياً (بدون نوع)'
        ELSE '❌ غير مسجل'
    END as device_status
FROM users 
WHERE role = 'employee'
ORDER BY created_at DESC;
