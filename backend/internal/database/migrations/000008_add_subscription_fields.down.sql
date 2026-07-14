-- إزالة حقول الاشتراك من جدول المستخدمين
ALTER TABLE users
DROP COLUMN IF EXISTS subscription_status,
DROP COLUMN IF EXISTS subscription_expires_at;
