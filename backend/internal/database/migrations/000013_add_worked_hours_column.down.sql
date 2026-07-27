-- التراجع عن إضافة حقل worked_hours

-- حذف الفهرس
DROP INDEX IF EXISTS idx_attendance_worked_hours;

-- حذف العمود
ALTER TABLE attendance DROP COLUMN IF EXISTS worked_hours;