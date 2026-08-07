DROP TRIGGER IF EXISTS update_notes_updated_at ON notes;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_notes_created_at;
DROP INDEX IF EXISTS idx_notes_employee_id;
DROP INDEX IF EXISTS idx_notes_admin_id;
DROP TABLE IF EXISTS notes;
