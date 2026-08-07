-- إضافة مهام تجريبية سريعة لحل مشكلة توزيع المهام
-- تنفذ هذا الملف في Supabase SQL Editor

-- 1. أولاً، التأكد من وجود نقاط عمل (باستخدام UUID صحيح)
INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, is_deleted, created_at, updated_at)
VALUES 
  ('550e8400-e29b-41d4-a716-446655440001', 'موقع رئيسي', 'عمان، الأردن', 31.9539, 35.9106, 100, true, false, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'موقع فرعي', 'إربد، الأردن', 32.5556, 35.8519, 100, true, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. إضافة مهام تجريبية بجميع الحالات (باستخدام UUID صحيح)
INSERT INTO tasks (id, title, description, worksite_id, status, scheduled_start, scheduled_end, created_at, updated_at)
VALUES 
  -- مهام مكتملة (اليوم)
  ('660e8400-e29b-41d4-a716-446655440001', 'توصيل طلب مكتمل', 'تم توصيل الطلب بنجاح', '550e8400-e29b-41d4-a716-446655440001', 'completed', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '1 hour', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440002', 'صيانة مكتملة', 'تم إصلاح العطل', '550e8400-e29b-41d4-a716-446655440002', 'completed', NOW() - INTERVAL '8 hours', NOW() - INTERVAL '2 hours', NOW(), NOW()),
  
  -- مهام جارية
  ('660e8400-e29b-41d4-a716-446655440003', 'توصيل طلب جاري', 'جاري التوصيل حالياً', '550e8400-e29b-41d4-a716-446655440001', 'in_progress', NOW() - INTERVAL '1 hour', NOW() + INTERVAL '2 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440004', 'صيانة جارية', 'جاري الإصلاح', '550e8400-e29b-41d4-a716-446655440002', 'in_progress', NOW(), NOW() + INTERVAL '3 hours', NOW(), NOW()),
  
  -- مهام قيد الانتظار
  ('660e8400-e29b-41d4-a716-446655440005', 'طلب جديد', 'طلب توصيل جديد', '550e8400-e29b-41d4-a716-446655440001', 'pending', NOW() + INTERVAL '1 hour', NOW() + INTERVAL '4 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440006', 'صيانة مجدولة', 'صيانة مجدولة لغداً', '550e8400-e29b-41d4-a716-446655440002', 'pending', NOW() + INTERVAL '5 hours', NOW() + INTERVAL '8 hours', NOW(), NOW()),
  
  -- مهام متأخرة
  ('660e8400-e29b-41d4-a716-446655440007', 'طلب متأخر', 'طلب تأخر عن الموعد', '550e8400-e29b-41d4-a716-446655440001', 'late', NOW() - INTERVAL '1 day', NOW() - INTERVAL '5 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440008', 'صيانة متأخرة', 'صيانة تأخرت عن الموعد', '550e8400-e29b-41d4-a716-446655440002', 'late', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. عرض التقرير النهائي
SELECT 
    'التقرير النهائي' as report,
    (SELECT COUNT(*) FROM tasks WHERE status = 'completed') as completed,
    (SELECT COUNT(*) FROM tasks WHERE status = 'in_progress') as in_progress,
    (SELECT COUNT(*) FROM tasks WHERE status = 'pending') as pending,
    (SELECT COUNT(*) FROM tasks WHERE status = 'late') as late,
    (SELECT COUNT(*) FROM tasks) as total_tasks;