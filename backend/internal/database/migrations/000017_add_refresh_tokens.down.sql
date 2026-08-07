-- حذف indexes
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_is_revoked;

-- حذف العمود من جدول users
ALTER TABLE users DROP COLUMN IF EXISTS current_refresh_token_id;

-- حذف جدول refresh_tokens
DROP TABLE IF EXISTS refresh_tokens;
