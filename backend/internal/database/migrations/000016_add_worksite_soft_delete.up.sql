-- إضافة عمود is_deleted لجدول worksites (Soft Delete)
ALTER TABLE worksites ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- إضافة index على is_deleted
CREATE INDEX IF NOT EXISTS idx_worksites_is_deleted ON worksites(is_deleted);
