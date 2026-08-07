-- حذف index على is_deleted
DROP INDEX IF EXISTS idx_worksites_is_deleted;

-- حذف العمود is_deleted من جدول worksites
ALTER TABLE worksites DROP COLUMN IF EXISTS is_deleted;
