-- تفعيل توليد UUID تلقائياً (مطلوب لكل الجداول التي تعتمد uuid_generate_v4)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- جدول المستخدمين: يضم كلاً من المدير (admin) والموظف (employee) في جدول واحد
-- ويُفرَّق بينهما عبر عمود role فقط، لتبسيط منطق تسجيل الدخول (Endpoint واحد للجميع)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(30),
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee' CHECK (role IN ('admin', 'employee')),
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
