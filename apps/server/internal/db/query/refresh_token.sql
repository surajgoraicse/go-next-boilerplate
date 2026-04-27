-- name: UpsertRefreshSession :one
INSERT INTO refresh_sessions(id, user_id, device_id, refresh_token_hash, issued_at, revoked_at, expires_at, last_used_at, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), NULL, $5, now(), now(), now())
ON CONFLICT (user_id, device_id)
    DO UPDATE SET
        refresh_token_hash = EXCLUDED.refresh_token_hash,
        issued_at = now(),
        revoked_at = NULL,
        expires_at = EXCLUDED.expires_at,
        last_used_at = now(),
        updated_at = now()
    RETURNING
        *;

-- name: GetActiveSessionByTokenHash :one
SELECT
    *
FROM
    refresh_sessions
WHERE
    refresh_token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > now()
FOR UPDATE;

-- name: CountActiveSessionsByUser :one
SELECT
    count(*)::bigint
FROM
    refresh_sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
    AND expires_at > now();

-- name: GetOldestActiveSessionByUser :one
SELECT
    *
FROM
    refresh_sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
    AND expires_at > now()
ORDER BY
    last_used_at ASC NULLS FIRST,
    created_at ASC
LIMIT 1;

-- name: RotateRefreshToken :one
UPDATE
    refresh_sessions
SET
    refresh_token_hash = $2,
    issued_at = now(),
    expires_at = $3,
    last_used_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND revoked_at IS NULL
    AND expires_at > now()
RETURNING
    *;

-- name: TouchRefreshSession :one
UPDATE
    refresh_sessions
SET
    last_used_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND revoked_at IS NULL
    AND expires_at > now()
RETURNING
    *;

-- name: RevokeSessionByID :one
UPDATE
    refresh_sessions
SET
    revoked_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND revoked_at IS NULL
RETURNING
    *;

-- name: RevokeAllSessionsByUser :many
UPDATE
    refresh_sessions
SET
    revoked_at = now(),
    updated_at = now()
WHERE
    user_id = $1
    AND revoked_at IS NULL
RETURNING
    id;

-- name: DeleteExpiredSessions :exec
DELETE FROM refresh_sessions
WHERE expires_at <= now();

-- name: ListActiveSessionsByUser :many
SELECT
    *
FROM
    refresh_sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
    AND expires_at > now()
ORDER BY
    last_used_at DESC NULLS LAST,
    created_at DESC;

