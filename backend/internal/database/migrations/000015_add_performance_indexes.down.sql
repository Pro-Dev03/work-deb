-- حذف indexes لتحسين الأداء

-- حذف indexes من جدول users
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_phone;

-- حذف indexes من جدول worksites
DROP INDEX IF EXISTS idx_worksites_is_active;
DROP INDEX IF EXISTS idx_worksites_assigned_employee_id;

-- حذف indexes من جدول attendance
DROP INDEX IF EXISTS idx_attendance_user_id;
DROP INDEX IF EXISTS idx_attendance_worksite_id;
DROP INDEX IF EXISTS idx_attendance_status;
DROP INDEX IF EXISTS idx_attendance_check_in_time;
DROP INDEX IF EXISTS idx_attendance_user_status;
DROP INDEX IF EXISTS idx_attendance_worksite_status;

-- حذف indexes من جدول tasks
DROP INDEX IF EXISTS idx_tasks_worksite_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_created_at;

-- حذف indexes من جدول notifications
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_created_at;

-- حذف composite index
DROP INDEX IF EXISTS idx_attendance_user_checkin;
