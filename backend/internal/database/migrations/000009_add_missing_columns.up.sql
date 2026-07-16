-- إضافة الأعمدة المفقودة للتوافق مع الإصدارات السابقة

-- إضافة عمود assigned_employee_id إلى جدول worksites
ALTER TABLE worksites 
ADD COLUMN IF NOT EXISTS assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- إضافة عمود updated_at إلى جدول attendance
ALTER TABLE attendance
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- إضافة حقول الصور إلى جدول attendance إذا لم تكن موجودة
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_url TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_uploaded_at TIMESTAMPTZ;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS photo_notes TEXT;

-- إنشاء الفهارس المفقودة
CREATE INDEX IF NOT EXISTS idx_worksites_assigned_employee ON worksites(assigned_employee_id);
CREATE INDEX IF NOT EXISTS idx_attendance_worksite_id ON attendance(worksite_id);