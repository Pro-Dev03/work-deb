-- نقاط العمل (Geofence Zones) — هذا الجدول هو أساس ميزة منع التختيم خارج النطاق.
-- radius_meters هو نصف القطر المسموح به حول (latitude, longitude) لتلك النقطة
CREATE TABLE IF NOT EXISTS worksites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius_meters INTEGER NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_employee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
