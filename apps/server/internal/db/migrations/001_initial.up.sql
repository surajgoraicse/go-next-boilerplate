SET search_path TO boilerplate, public;

-- =====================================================
-- ENUM Types
-- =====================================================
CREATE TYPE user_role AS ENUM(
    'admin',
    'user',
    'super_admin'
);

CREATE TYPE verification_purpose AS ENUM(
    'password_reset',
    'email_verification'
);

-- =====================================================
-- Helper Functions
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- ============================================================
-- 1. AUTH MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name varchar(100) NOT NULL,
    email varchar(255) UNIQUE,
    email_verified boolean NOT NULL DEFAULT FALSE,
    password_hash text,
    role user_role NOT NULL DEFAULT 'user',
    google_id varchar(255) UNIQUE,
    avatar_url text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- At least one auth method must be present
    CONSTRAINT chk_auth_method CHECK (password_hash IS NOT NULL OR google_id IS NOT NULL)
);

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- Entity: Refresh Session
-- Description: Refresh tokens for sessions
-- =====================================================
CREATE TABLE IF NOT EXISTS refresh_sessions (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id uuid NOT NULL,

    refresh_token_hash text NOT NULL UNIQUE,

    issued_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz,
    expires_at timestamptz NOT NULL,
    last_used_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT refresh_sessions_user_device_uniq UNIQUE (user_id, device_id),
    CONSTRAINT refresh_sessions_expires_after_issued CHECK (expires_at > issued_at)
);

CREATE INDEX IF NOT EXISTS refresh_sessions_user_status_idx
    ON refresh_sessions (user_id, revoked_at, expires_at);

CREATE INDEX IF NOT EXISTS refresh_sessions_user_device_idx
    ON refresh_sessions (user_id, device_id);


-- =====================================================
-- Entity: Verification OTP
-- Description: OTP-based verification flows (password reset, email verification, MFA)
-- =====================================================
CREATE TABLE IF NOT EXISTS verification_otps(
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email varchar(255) NOT NULL,
    purpose verification_purpose NOT NULL,
    otp_hash varchar(255) NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz DEFAULT NULL,
    attempt_count integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT verification_otps_attempt_count_positive CHECK (attempt_count >= 0),
    CONSTRAINT verification_otps_expires_after_created CHECK (expires_at > created_at),
    CONSTRAINT verification_otps_used_after_created CHECK (used_at IS NULL OR used_at >= created_at)
);
CREATE TRIGGER trg_verification_otps_updated_at
    BEFORE UPDATE ON verification_otps
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Unique index: Only one active OTP per user per purpose
CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_otps_active_user_purpose ON verification_otps(user_id, purpose)
WHERE
    used_at IS NULL;

-- Index for rate limit queries (user_id, email, purpose, created_at DESC)
CREATE INDEX IF NOT EXISTS idx_verification_otps_rate_limit ON verification_otps(user_id, email, purpose, created_at DESC);

-- Index for cleanup queries (expires_at)
CREATE INDEX IF NOT EXISTS idx_verification_otps_expires_at ON verification_otps(expires_at);

-- Index for user_id lookups
CREATE INDEX IF NOT EXISTS idx_verification_otps_user_id ON verification_otps(user_id);

-- =====================================================
-- Entity: Verification Session Token
-- Description: Short-lived verification sessions after OTP verification
-- =====================================================
CREATE TABLE IF NOT EXISTS verification_sessions_token(
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    otp_id uuid NOT NULL REFERENCES verification_otps(id) ON DELETE CASCADE,
    purpose verification_purpose NOT NULL,
    session_token_hash varchar(255) NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz DEFAULT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT verification_sessions_token_expires_after_created CHECK (expires_at > created_at),
    CONSTRAINT verification_sessions_token_used_after_created CHECK (used_at IS NULL OR used_at >= created_at)
);
CREATE TRIGGER trg_verification_sessions_token_updated_at
    BEFORE UPDATE ON verification_sessions_token
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Unique constraint: One session per OTP
CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_sessions_token_otp_id ON verification_sessions_token(otp_id);

-- Unique index: Only one active session per user per purpose
CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_sessions_token_active_user_purpose ON verification_sessions_token(user_id, purpose)
WHERE
    used_at IS NULL;

-- Unique index: Fast lookups by session token hash
CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_sessions_token_hash ON verification_sessions_token(session_token_hash);

-- Index for cleanup queries (expires_at)
CREATE INDEX IF NOT EXISTS idx_verification_sessions_token_expires_at ON verification_sessions_token(expires_at);

-- Index for user_id lookups
CREATE INDEX IF NOT EXISTS idx_verification_sessions_token_user_id ON verification_sessions_token(user_id);
