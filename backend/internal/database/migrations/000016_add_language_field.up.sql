-- =====================================================
-- Add language field to tasks table
-- =====================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';

-- =====================================================
-- Add language field to service_requests table
-- =====================================================
ALTER TABLE service_requests
    ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT 'en';
