-- حذف المهام التجريبية ونقاط العمل التجريبية
-- تنفذ هذا الملف في Supabase SQL Editor

-- حذف المهام التجريبية
DELETE FROM tasks WHERE id LIKE '660e8400%';

-- حذف نقاط العمل التجريبية
DELETE FROM worksites WHERE id LIKE '550e8400%';

-- عرض التقرير النهائي
SELECT 
    'تم الحذف' as status,
    (SELECT COUNT(*) FROM tasks) as remaining_tasks,
    (SELECT COUNT(*) FROM worksites) as remaining_worksites;