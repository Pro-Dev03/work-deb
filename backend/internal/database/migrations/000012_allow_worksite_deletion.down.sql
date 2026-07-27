-- التراجع عن التغييرات - إعادة القيد الأصلي
-- إزالة حقل worksite_name_for_history
-- جعل worksite_id NOT NULL مرة أخرى وتغيير القيد إلى ON DELETE RESTRICT

-- أولاً: تغيير قيد الحذف إلى RESTRICT
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS attendance_worksite_id_fkey;
ALTER TABLE attendance ADD CONSTRAINT attendance_worksite_id_fkey
    FOREIGN KEY (worksite_id) REFERENCES worksites(id) ON DELETE RESTRICT;

-- ثانياً: جعل العمود NOT NULL
ALTER TABLE attendance ALTER COLUMN worksite_id SET NOT NULL;

-- ثالثاً: إزالة حقل اسم نقطة العمل للسجل التاريخي
ALTER TABLE attendance DROP COLUMN IF EXISTS worksite_name_for_history;
