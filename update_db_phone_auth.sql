-- =============================================
-- إضافة أعمدة لتسجيل الدخول عبر الهاتف
-- =============================================

-- إضافة عمود device_id لتخزين معرف الجهاز
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_id TEXT;

-- إضافة عمود phone_verified لتأكيد رقم الهاتف
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN DEFAULT FALSE;

-- إضافة عمود last_login_device لتسجيل آخر جهاز استخدم
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_device TEXT;

-- إضافة عمود login_code لتخزين رمز التحقق المؤقت
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_code_expires TIMESTAMPTZ;

-- عرض الأعمدة الجديدة
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
ORDER BY ordinal_position;
