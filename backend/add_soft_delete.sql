-- =====================================================
-- سكريبت إضافة Soft Delete إلى جدول worksites
-- =====================================================

-- إضافة حقل is_deleted إلى جدول worksites
ALTER TABLE worksites 
ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- تحديث القيم الحالية للتأكد من عدم وجود NULL
UPDATE worksites 
SET is_deleted = FALSE 
WHERE is_deleted IS NULL;

-- إضافة فهرس لتحسين الأداء على الاستعلامات
CREATE INDEX IF NOT EXISTS idx_worksites_is_deleted 
ON worksites(is_deleted);

-- =====================================================
-- التحقق من الإضافة
-- =====================================================

-- عرض هيكل الجدول الجديد
SELECT 
    column_name, 
    data_type, 
    is_nullable, 
    column_default
FROM information_schema.columns 
WHERE table_name = 'worksites' 
AND column_name = 'is_deleted';

-- عرض الفهارس
SELECT 
    indexname, 
    indexdef 
FROM pg_indexes 
WHERE tablename = 'worksites' 
AND indexname = 'idx_worksites_is_deleted';

-- =====================================================
-- الاستعلامات المفيدة
-- =====================================================

-- عرض جميع نقاط العمل النشطة
SELECT * FROM worksites WHERE is_deleted = FALSE;

-- عرض جميع نقاط العمل المحذوفة
SELECT * FROM worksites WHERE is_deleted = TRUE;

-- استعادة نقطة عمل محذوفة
-- UPDATE worksites SET is_deleted = FALSE WHERE id = 'your-worksite-id';

-- =====================================================
-- ملاحظات:
-- =====================================================
-- 1. الحقل is_deleted = FALSE يعني أن نقطة العمل نشطة
-- 2. الحقل is_deleted = TRUE يعني أن نقطة العمل محذوفة (Soft Delete)
-- 3. يتم استخدام Soft Delete للحفاظ على بيانات الحضور التاريخية
-- 4. يمكن استعادة نقاط العمل المحذوفة عن طريق تعيين is_deleted = FALSE
-- 5. الفهرس idx_worksites_is_deleted يحسن أداء الاستعلامات
-- =====================================================