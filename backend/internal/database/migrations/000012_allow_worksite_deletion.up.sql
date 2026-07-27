-- السماح بحذف نقاط العمل مع الحفاظ على سجلات الحضور
-- تغيير القيد من ON DELETE RESTRICT إلى ON DELETE SET NULL
-- وجعل worksite_id nullable
-- إضافة حقل لحفظ اسم نقطة العمل للسجل التاريخي

-- أولاً: إضافة حقل لحفظ اسم نقطة العمل للسجل التاريخي
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS worksite_name_for_history TEXT;

-- ثانياً: جعل العمود nullable
ALTER TABLE attendance ALTER COLUMN worksite_id DROP NOT NULL;

-- ثالثاً: تغيير قيد الحذف
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_worksite_id_fkey;
ALTER TABLE attendance ADD CONSTRAINT attendance_worksite_id_fkey
    FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE SET NULL;
