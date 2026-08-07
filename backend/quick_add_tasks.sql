-- إضافة مهام تجريبية سريعة لحل مشكلة توزيع المهام
-- تنفذ هذا الملف في Supabase SQL Editor

-- 1. أولاً، التأكد من وجود نقاط عمل
INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, is_deleted, created_at, updated_at)
VALUES 
  ('sample-worksite-1', 'موقع رئيسي', 'عمان، الأردن', 31.9539, 35.9106, 100, true, false, NOW(), NOW()),
  ('sample-worksite-2', 'موقع فرعي', 'إربد، الأردن', 32.5556, 35.8519, 100, true, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. إضافة مهام تجريبية بجميع الحالات
INSERT INTO tasks (id, title, description, worksite_id, status, scheduled_start, scheduled_end, created_at, updated_at)
VALUES 
  -- مهام مكتملة (اليوم)
  ('task-completed-1', 'توصيل طلب مكتمل', 'تم توصيل الطلب بنجاح', 'sample-worksite-1', 'completed', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '1 hour', NOW(), NOW()),
  ('task-completed-2', 'صيانة مكتملة', 'تم إصلاح العطل', 'sample-worksite-2', 'completed', NOW() - INTERVAL '8 hours', NOW() - INTERVAL '2 hours', NOW(), NOW()),
  
  -- مهام جارية
  ('task-inprogress-1', 'توصيل طلب جاري', 'جاري التوصيل حالياً', 'sample-worksite-1', 'in_progress', NOW() - INTERVAL '1 hour', NOW() + INTERVAL '2 hours', NOW(), NOW()),
  ('task-inprogress-2', 'صيانة جارية', 'جاري الإصلاح', 'sample-worksite-2', 'in_progress', NOW(), NOW() + INTERVAL '3 hours', NOW(), NOW()),
  
  -- مهام قيد الانتظار
  ('task-pending-1', 'طلب جديد', 'طلب توصيل جديد', 'sample-worksite-1', 'pending', NOW() + INTERVAL '1 hour', NOW() + INTERVAL '4 hours', NOW(), NOW()),
  ('task-pending-2', 'صيانة مجدولة', 'صيانة مجدولة لغداً', 'sample-worksite-2', 'pending', NOW() + INTERVAL '5 hours', NOW() + INTERVAL '8 hours', NOW(), NOW()),
  
  -- مهام متأخرة
  ('task-late-1', 'طلب متأخر', 'طلب تأخر عن الموعد', 'sample-worksite-1', 'late', NOW() - INTERVAL '1 day', NOW() - INTERVAL '5 hours', NOW(), NOW()),
  ('task-late-2', 'صيانة متأخرة', 'صيانة تأخرت عن الموعد', 'sample-worksite-2', 'late', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. عرض التقرير النهائي
SELECT 
    'التقرير النهائي' as report,
    (SELECT COUNT(*) FROM tasks WHERE status = 'completed') as completed,
    (SELECT COUNT(*) FROM tasks WHERE status = 'in_progress') as in_progress,
    (SELECT COUNT(*) FROM tasks WHERE status = 'pending') as pending,
    (SELECT COUNT(*) FROM tasks WHERE status = 'late') as late,
    (SELECT COUNT(*) FROM tasks) as total_tasks;