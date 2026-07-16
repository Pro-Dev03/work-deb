-- التراجع عن إضافة الأعمدة المفقودة

-- حذف الفهارس
DROP INDEX IF EXISTS idx_worksites_assigned_employee;
DROP INDEX IF EXISTS idx_attendance_worksite_id;

-- حذف الأعمدة (اختياري - عادة نحتفظ بها للتوافق)
-- ALTER TABLE worksites DROP COLUMN IF EXISTS assigned_employee_id;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS updated_at;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_url;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_uploaded_at;
-- ALTER TABLE attendance DROP COLUMN IF EXISTS photo_notes;