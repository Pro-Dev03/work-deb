-- إضافة عمود assigned_employee_id إلى جدول worksites
-- هذا العمود لتعيين موظف محدد لنقطة عمل

ALTER TABLE worksites 
ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- إنشاء فهرس للأداء
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);