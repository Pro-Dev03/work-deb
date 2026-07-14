-- إضافة حقول الاشتراك إلى جدول المستخدمين
ALTER TABLE users
ADD COLUMN subscription_status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (subscription_status IN ('active', 'expired', 'canceled')),
ADD COLUMN subscription_expires_at TIMESTAMPTZ NULL;
