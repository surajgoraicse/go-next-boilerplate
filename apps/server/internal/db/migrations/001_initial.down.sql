-- order matters for dependencies

SET search_path TO boilerplate, public;

-- =====================================================
-- Drop Triggers (must be dropped before tables and functions)
-- =====================================================
DROP TRIGGER IF EXISTS trg_verification_sessions_token_updated_at ON users;
DROP TRIGGER IF EXISTS trg_refresh_tokens_updated_at ON refresh_tokens;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

-- =====================================================
-- Drop Indexes
-- =====================================================
DROP INDEX IF EXISTS idx_verification_sessions_token_user_id;
DROP INDEX IF EXISTS idx_verification_sessions_token_expires_at;
DROP INDEX IF EXISTS idx_verification_sessions_token_hash;
DROP INDEX IF EXISTS idx_verification_sessions_token_active_user_purpose;
DROP INDEX IF EXISTS idx_verification_sessions_token_otp_id;
DROP INDEX IF EXISTS idx_verification_otps_user_id;
DROP INDEX IF EXISTS idx_verification_otps_expires_at;
DROP INDEX IF EXISTS idx_verification_otps_rate_limit;
DROP INDEX IF EXISTS idx_verification_otps_active_user_purpose;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

-- =====================================================
-- Drop Tables (order matters for foreign keys)
-- =====================================================
DROP TABLE IF EXISTS verification_sessions_token;
DROP TABLE IF EXISTS verification_otps;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

-- =====================================================
-- Drop Functions
-- =====================================================
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- =====================================================
-- Drop Types
-- =====================================================
DROP TYPE IF EXISTS verification_purpose;
DROP TYPE IF EXISTS user_role;
