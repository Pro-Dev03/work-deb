-- تحديث اشتراك مدير معين بمدة محددة

-- المثال 1: تفعيل اشتراك لمدة شهر
UPDATE users 
SET 
  subscription_status = 'active',
  subscription_expires_at = NOW() + INTERVAL '1 month'
WHERE email = 'admin@worktrack.com';

-- المثال 2: تفعيل اشتراك لمدة سنة
UPDATE users 
SET 
  subscription_status = 'active',
  subscription_expires_at = NOW() + INTERVAL '1 year'
WHERE email = 'morad@worktrack.com';

-- المثال 3: تفعيل اشتراك لمدة محددة (تاريخ معين)
UPDATE users 
SET 
  subscription_status = 'active',
  subscription_expires_at = '2026-12-31 23:59:59'
WHERE email = 'admin@worktrack.com';

-- المثال 4: إلغاء الاشتراك
UPDATE users 
SET 
  subscription_status = 'canceled',
  subscription_expires_at = NULL
WHERE email = 'admin@worktrack.com';

-- المثال 5: تفعيل اشتراك دائم (بدون تاريخ انتهاء)
UPDATE users 
SET 
  subscription_status = 'active',
  subscription_expires_at = NULL
WHERE email = 'admin@worktrack.com';

-- التحقق من التحديث
SELECT 
  id, 
  full_name, 
  email, 
  role, 
  subscription_status, 
  subscription_expires_at,
  CASE 
    WHEN subscription_expires_at IS NULL THEN 'دائم'
    WHEN subscription_expires_at > NOW() THEN 'نشط'
    ELSE 'منتهي'
  END as current_status
FROM users 
WHERE role = 'admin';
