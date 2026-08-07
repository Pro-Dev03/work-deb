-- إضافة مهام تجريبية لفحص توزيع المهام
-- هذا السكريبت يقوم بإضافة مهام بجميع الحالات المختلفة

-- أولاً، التأكد من وجود بعض نقاط العمل
INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, is_deleted)
VALUES 
  ('550e8400-e29b-41d4-a716-446655440001', 'موقع تجريبي 1', 'عمان، الأردن', 31.9539, 35.9106, 100, true, false),
  ('550e8400-e29b-41d4-a716-446655440002', 'موقع تجريبي 2', 'إربد، الأردن', 32.5556, 35.8519, 100, true, false)
ON CONFLICT (id) DO NOTHING;

-- التأكد من وجود بعض المستخدمين (موظفين)
-- نفترض أن هناك مستخدمين موجودين، سنستخدم معرفات موجودة

-- إضافة مهام تجريبية بجميع الحالات
INSERT INTO tasks (id, title, description, worksite_id, assigned_user_id, status, scheduled_start, scheduled_end)
VALUES 
  -- مهام مكتملة
  ('660e8400-e29b-41d4-a716-446655440001', 'مهمة مكتملة 1', 'هذه مهمة مكتملة', '550e8400-e29b-41d4-a716-446655440001', NULL, 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
  ('660e8400-e29b-41d4-a716-446655440002', 'مهمة مكتملة 2', 'هذه مهمة مكتملة أخرى', '550e8400-e29b-41d4-a716-446655440002', NULL, 'completed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
  
  -- مهام جارية
  ('660e8400-e29b-41d4-a716-446655440003', 'مهمة جارية 1', 'هذه مهمة قيد التنفيذ', '550e8400-e29b-41d4-a716-446655440001', NULL, 'in_progress', NOW(), NOW() + INTERVAL '1 day'),
  ('660e8400-e29b-41d4-a716-446655440004', 'مهمة جارية 2', 'هذه مهمة جارية أخرى', '550e8400-e29b-41d4-a716-446655440002', NULL, 'in_progress', NOW() - INTERVAL '1 hour', NOW() + INTERVAL '5 hours'),
  
  -- مهام قيد الانتظار
  ('660e8400-e29b-41d4-a716-446655440005', 'مهمة معلقة 1', 'هذه مهمة قيد الانتظار', '550e8400-e29b-41d4-a716-446655440001', NULL, 'pending', NOW() + INTERVAL '1 day', NOW() + INTERVAL '2 days'),
  ('660e8400-e29b-41d4-a716-446655440006', 'مهمة معلقة 2', 'هذه مهمة معلقة أخرى', '550e8400-e29b-41d4-a716-446655440002', NULL, 'pending', NOW() + INTERVAL '2 days', NOW() + INTERVAL '3 days'),
  
  -- مهام متأخرة
  ('660e8400-e29b-41d4-a716-446655440007', 'مهمة متأخرة 1', 'هذه مهمة متأخرة', '550e8400-e29b-41d4-a716-446655440001', NULL, 'late', NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day'),
  ('660e8400-e29b-41d4-a716-446655440008', 'مهمة متأخرة 2', 'هذه مهمة متأخرة أخرى', '550e8400-e29b-41d4-a716-446655440002', NULL, 'late', NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days')
ON CONFLICT (id) DO NOTHING;

-- عرض ملخص المهام المضافة
SELECT 
    status,
    COUNT(*) as count,
    COUNT(CASE WHEN created_at::date = CURRENT_DATE THEN 1 END) as created_today
FROM tasks 
GROUP BY status
ORDER BY status;