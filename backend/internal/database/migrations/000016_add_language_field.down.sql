-- =====================================================
-- Remove language field from tasks table
-- =====================================================
ALTER TABLE tasks
    DROP COLUMN IF EXISTS language;

-- =====================================================
-- Remove language field from service_requests table
-- =====================================================
ALTER TABLE service_requests
    DROP COLUMN IF EXISTS language;
