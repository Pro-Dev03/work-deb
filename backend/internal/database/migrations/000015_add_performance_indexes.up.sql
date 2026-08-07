-- إضافة indexes لتحسين الأداء على الجداول الرئيسية

-- indexes على جدول users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);

-- indexes على جدول worksites
CREATE INDEX IF NOT EXISTS idx_worksites_is_active ON worksites(is_active);
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee_id ON worksites(assigned_employee_id);

-- indexes على جدول attendance
CREATE INDEX IF NOT EXISTS idx_attendance_user_id ON attendance(user_id);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);
CREATE INDEX IF NOT EXISTS idx_attendance_status ON attendance(status);
CREATE INDEX IF NOT EXISTS idx_attendance_check_in_time ON attendance(check_in_time);
CREATE INDEX IF NOT EXISTS idx_attendance_user_status ON attendance(user_id, status);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_status ON attendance(worksite_id, status);

-- indexes على جدول tasks
CREATE INDEX IF NOT EXISTS idx_tasks_worksite_id ON tasks(worksite_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- indexes على جدول notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- composite index للعديد من الاستعلامات المشتركة
CREATE INDEX IF NOT EXISTS idx_attendance_user_checkin ON attendance(user_id, check_in_time DESC);
