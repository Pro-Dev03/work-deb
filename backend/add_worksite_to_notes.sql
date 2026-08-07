-- إضافة عمود worksite_id إلى جدول الملاحظات
ALTER TABLE notes ADD COLUMN IF NOT EXISTS worksite_id UUID REFERENCES worksites(id) ON DELETE SET NULL;

-- إنشاء فهرس للعمود الجديد
CREATE INDEX IF NOT EXISTS idx_notes_worksite_id ON notes(worksite_id);
