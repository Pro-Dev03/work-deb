-- حذف الفهارس المركبة
DROP INDEX IF EXISTS idx_security_logs_email_success_time;
DROP INDEX IF EXISTS idx_security_logs_ip_success_time;

-- حذف الفهارس البسيطة
DROP INDEX IF EXISTS idx_security_logs_created_at;
DROP INDEX IF EXISTS idx_security_logs_event_type;
DROP INDEX IF EXISTS idx_security_logs_user_id;
DROP INDEX IF EXISTS idx_security_logs_email;
DROP INDEX IF EXISTS idx_security_logs_ip;

-- حذف جدول security_logs
DROP TABLE IF EXISTS security_logs;
