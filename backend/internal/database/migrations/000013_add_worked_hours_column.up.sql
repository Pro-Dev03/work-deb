-- إضافة حقل worked_hours لحساب إجمالي الساعات (night_hours + day_hours)
-- هذا الحقل يُحسب ويُخزن لضمان الاتساق وتسهيل الحسابات المحاسبية

ALTER TABLE attendance ADD COLUMN IF NOT EXISTS worked_hours DOUBLE PRECISION;

-- تحديث السجلات الموجودة لحساب worked_hours
UPDATE attendance 
SET worked_hours = COALESCE(night_hours, 0) + COALESCE(day_hours, 0)
WHERE worked_hours IS NULL 
AND night_hours IS NOT NULL 
AND day_hours IS NOT NULL;

-- إنشاء فهرس لتحسين استعلامات الملخصات والحسابات المحاسبية
CREATE INDEX IF NOT EXISTS idx_attendance_worked_hours 
ON attendance(worked_hours) 
WHERE worked_hours IS NOT NULL;