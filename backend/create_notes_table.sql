-- جدول الملاحظات
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- إنشاء الفهارس
CREATE INDEX IF NOT EXISTS idx_notes_admin_id ON notes(admin_id);
CREATE INDEX IF NOT EXISTS idx_notes_employee_id ON notes(employee_id);
CREATE INDEX IF NOT EXISTS idx_notes_is_read ON notes(is_read);
CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
