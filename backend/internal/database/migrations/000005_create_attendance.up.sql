-- سجل الحضور/الانصراف — كل سجل يحفظ موقع الموظف الفعلي والمسافة عن نقطة العمل
-- لحظة كل تختيم (بدء/إنهاء)، كإثبات أنه تم داخل النطاق المسموح
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    worksite_id UUID NOT NULL REFERENCES worksites(id) ON DELETE RESTRICT,

    check_in_time TIMESTAMPTZ,
    check_in_lat DOUBLE PRECISION,
    check_in_lng DOUBLE PRECISION,
    check_in_distance_meters DOUBLE PRECISION,

    check_out_time TIMESTAMPTZ,
    check_out_lat DOUBLE PRECISION,
    check_out_lng DOUBLE PRECISION,
    check_out_distance_meters DOUBLE PRECISION,

    status VARCHAR(20) NOT NULL DEFAULT 'in_progress'
        CHECK (status IN ('in_progress', 'completed')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
