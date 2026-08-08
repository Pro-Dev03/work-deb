-- إنشاء جدول security_logs لتسجيل أحداث الأمان
CREATE TABLE IF NOT EXISTS security_logs (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    event_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(36),
    email VARCHAR(255),
    phone VARCHAR(20),
    ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    success BOOLEAN DEFAULT false,
    reason TEXT,
    details TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- إنشاء فهارس للبحث السريع
CREATE INDEX IF NOT EXISTS idx_security_logs_ip ON security_logs(ip);
CREATE INDEX IF NOT EXISTS idx_security_logs_email ON security_logs(email);
CREATE INDEX IF NOT EXISTS idx_security_logs_user_id ON security_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_security_logs_event_type ON security_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_security_logs_created_at ON security_logs(created_at);

-- إنشاء فهرس مركب للبحث عن الأنماط المشبوهة
CREATE INDEX IF NOT EXISTS idx_security_logs_ip_success_time ON security_logs(ip, success, created_at);
CREATE INDEX IF NOT EXISTS idx_security_logs_email_success_time ON security_logs(email, success, created_at);
