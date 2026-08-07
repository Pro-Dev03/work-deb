-- إضافة مهام تجريبية سريعة لحل مشكلة توزيع المهام
-- تنفذ هذا الملف في Supabase SQL Editor

-- أولاً، حذف المهام القديمة إن وجدت
DELETE FROM tasks WHERE id IN (
    '660e8400-e29b-41d4-a716-446655440001',
    '660e8400-e29b-41d4-a716-446655440002',
    '660e8400-e29b-41d4-a716-446655440003',
    '660e8400-e29b-41d4-a716-446655440004',
    '660e8400-e29b-41d4-a716-446655440005',
    '660e8400-e29b-41d4-a716-446655440006',
    '660e8400-e29b-41d4-a716-446655440007',
    '660e8400-e29b-41d4-a716-446655440008',
    '660e8400-e29b-41d4-a716-446655440009',
    '660e8400-e29b-41d4-a716-446655440010',
    '660e8400-e29b-41d4-a716-446655440011',
    '660e8400-e29b-41d4-a716-446655440012',
    '660e8400-e29b-41d4-a716-446655440013',
    '660e8400-e29b-41d4-a716-446655440014'
);

-- 1. التأكد من وجود نقاط عمل (باستخدام UUID صحيح)
INSERT INTO worksites (id, name, address, latitude, longitude, radius_meters, is_active, is_deleted, created_at, updated_at)
VALUES 
  ('550e8400-e29b-41d4-a716-446655440001', 'موقع رئيسي', 'عمان، الأردن', 31.9539, 35.9106, 100, true, false, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'موقع فرعي', 'إربد، الأردن', 32.5556, 35.8519, 100, true, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. إضافة مهام تجريبية بجميع الحالات مع تواريخ مختلفة لفحص الفلاتر
INSERT INTO tasks (id, title, description, worksite_id, status, scheduled_start, scheduled_end, created_at, updated_at)
VALUES 
  -- مهام مكتملة (اليوم)
  ('660e8400-e29b-41d4-a716-446655440001', 'توصيل طلب مكتمل اليوم', 'تم توصيل الطلب بنجاح اليوم', '550e8400-e29b-41d4-a716-446655440001', 'completed', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '1 hour', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440002', 'صيانة مكتملة اليوم', 'تم إصلاح العطل اليوم', '550e8400-e29b-41d4-a716-446655440002', 'completed', NOW() - INTERVAL '8 hours', NOW() - INTERVAL '2 hours', NOW(), NOW()),
  
  -- مهام مكتملة (الأسبوع الماضي)
  ('660e8400-e29b-41d4-a716-446655440003', 'توصيل مكتمل الأسبوع الماضي', 'تم التوصيل الأسبوع الماضي', '550e8400-e29b-41d4-a716-446655440001', 'completed', NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
  ('660e8400-e29b-41d4-a716-446655440004', 'صيانة مكتملة الأسبوع الماضي', 'تم الإصلاح الأسبوع الماضي', '550e8400-e29b-41d4-a716-446655440002', 'completed', NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
  
  -- مهام مكتملة (الشهر الماضي)
  ('660e8400-e29b-41d4-a716-446655440005', 'توصيل مكتمل الشهر الماضي', 'تم التوصيل الشهر الماضي', '550e8400-e29b-41d4-a716-446655440001', 'completed', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
  
  -- مهام جارية (اليوم)
  ('660e8400-e29b-41d4-a716-446655440006', 'توصيل طلب جاري', 'جاري التوصيل حالياً', '550e8400-e29b-41d4-a716-446655440001', 'in_progress', NOW() - INTERVAL '1 hour', NOW() + INTERVAL '2 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440007', 'صيانة جارية', 'جاري الإصلاح', '550e8400-e29b-41d4-a716-446655440002', 'in_progress', NOW(), NOW() + INTERVAL '3 hours', NOW(), NOW()),
  
  -- مهام جارية (الأسبوع الماضي)
  ('660e8400-e29b-41d4-a716-446655440008', 'توصيل جاري الأسبوع الماضي', 'كان جاري التوصيل', '550e8400-e29b-41d4-a716-446655440001', 'in_progress', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
  
  -- مهام قيد الانتظار (اليوم والمستقبل)
  ('660e8400-e29b-41d4-a716-446655440009', 'طلب جديد', 'طلب توصيل جديد', '550e8400-e29b-41d4-a716-446655440001', 'pending', NOW() + INTERVAL '1 hour', NOW() + INTERVAL '4 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440010', 'صيانة مجدولة', 'صيانة مجدولة لغداً', '550e8400-e29b-41d4-a716-446655440002', 'pending', NOW() + INTERVAL '5 hours', NOW() + INTERVAL '8 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440011', 'طلب للغد', 'طلب مجدول للغد', '550e8400-e29b-41d4-a716-446655440001', 'pending', NOW() + INTERVAL '1 day', NOW() + INTERVAL '2 days', NOW(), NOW()),
  
  -- مهام متأخرة (الأيام الماضية)
  ('660e8400-e29b-41d4-a716-446655440012', 'طلب متأخر', 'طلب تأخر عن الموعد', '550e8400-e29b-41d4-a716-446655440001', 'late', NOW() - INTERVAL '1 day', NOW() - INTERVAL '5 hours', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440013', 'صيانة متأخرة', 'صيانة تأخرت عن الموعد', '550e8400-e29b-41d4-a716-446655440002', 'late', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440014', 'طلب متأخر جداً', 'طلب تأخر بأيام', '550e8400-e29b-41d4-a716-446655440001', 'late', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days')
ON CONFLICT (id) DO NOTHING;

-- 3. عرض تقرير شامل للفلاتر
SELECT 
    'التقرير الشامل' as report,
    -- فترة اليوم
    (SELECT COUNT(*) FROM tasks WHERE status = 'completed' AND created_at::date = CURRENT_DATE) as completed_today,
    (SELECT COUNT(*) FROM tasks WHERE status = 'in_progress' AND created_at::date = CURRENT_DATE) as in_progress_today,
    (SELECT COUNT(*) FROM tasks WHERE status = 'pending' AND created_at::date = CURRENT_DATE) as pending_today,
    (SELECT COUNT(*) FROM tasks WHERE status = 'late' AND created_at::date = CURRENT_DATE) as late_today,
    -- فترة الأسبوع
    (SELECT COUNT(*) FROM tasks WHERE status = 'completed' AND created_at >= CURRENT_DATE - INTERVAL '7 days') as completed_week,
    (SELECT COUNT(*) FROM tasks WHERE status = 'in_progress' AND created_at >= CURRENT_DATE - INTERVAL '7 days') as in_progress_week,
    (SELECT COUNT(*) FROM tasks WHERE status = 'pending' AND created_at >= CURRENT_DATE - INTERVAL '7 days') as pending_week,
    (SELECT COUNT(*) FROM tasks WHERE status = 'late' AND created_at >= CURRENT_DATE - INTERVAL '7 days') as late_week,
    -- فترة الشهر
    (SELECT COUNT(*) FROM tasks WHERE status = 'completed' AND created_at >= CURRENT_DATE - INTERVAL '30 days') as completed_month,
    (SELECT COUNT(*) FROM tasks WHERE status = 'in_progress' AND created_at >= CURRENT_DATE - INTERVAL '30 days') as in_progress_month,
    (SELECT COUNT(*) FROM tasks WHERE status = 'pending' AND created_at >= CURRENT_DATE - INTERVAL '30 days') as pending_month,
    (SELECT COUNT(*) FROM tasks WHERE status = 'late' AND created_at >= CURRENT_DATE - INTERVAL '30 days') as late_month,
    -- الإجمالي
    (SELECT COUNT(*) FROM tasks) as total_tasks;