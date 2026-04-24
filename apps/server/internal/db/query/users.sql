-- name: InsertUser :one
INSERT INTO users(name, email, email_verified, phone, phone_verified, role, password_hash, google_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: InsertUserOnConflict :one
INSERT INTO users(name, email, email_verified, phone, phone_verified, role, password_hash, google_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (email)
    DO NOTHING
RETURNING
    *;

-- name: GetUserByEmail :one
SELECT
    id,
    name,
    email,
    email_verified,
    phone,
    phone_verified,
    ROLE,
    password_hash,
    google_id,
    refresh_token_hash,
    refresh_token_expires_at,
    created_at,
    updated_at
FROM
    users
WHERE
    email = $1;

-- name: GetUserByID :one
SELECT
    id,
    name,
    email,
    email_verified,
    phone,
    phone_verified,
    ROLE,
    password_hash,
    google_id,
    refresh_token_hash,
    refresh_token_expires_at,
    created_at,
    updated_at
FROM
    users
WHERE
    id = $1;

-- name: GetUserByIDForUpdate :one
SELECT
    id,
    name,
    email,
    email_verified,
    phone,
    phone_verified,
    ROLE,
    password_hash,
    google_id,
    refresh_token_hash,
    refresh_token_expires_at,
    created_at,
    updated_at
FROM
    users
WHERE
    id = $1
FOR UPDATE;

-- name: GetUserByGoogleID :one
SELECT
    id,
    name,
    email,
    email_verified,
    phone,
    phone_verified,
    ROLE,
    password_hash,
    google_id,
    refresh_token_hash,
    refresh_token_expires_at,
    created_at,
    updated_at
FROM
    users
WHERE
    google_id = $1;

-- name: UpdateUserVerification :one
UPDATE
    users
SET
    email_verified = $2,
    phone_verified = $3
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateRefreshToken :exec
UPDATE
    users
SET
    refresh_token_hash = $2,
    refresh_token_expires_at = $3
WHERE
    id = $1;

-- name: GetUserByRefreshTokenHash :one
SELECT
    id,
    name,
    email,
    email_verified,
    phone,
    phone_verified,
    ROLE,
    password_hash,
    google_id,
    refresh_token_hash,
    refresh_token_expires_at,
    created_at,
    updated_at
FROM
    users
WHERE
    refresh_token_hash = $1;

-- name: ClearRefreshToken :exec
UPDATE
    users
SET
    refresh_token_hash = NULL,
    refresh_token_expires_at = NULL
WHERE
    id = $1;

-- name: UpdateUserRole :one
UPDATE
    users
SET
    ROLE = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserPassword :exec
UPDATE
    users
SET
    password_hash = $2
WHERE
    id = $1;

-- AttachGoogleToUser :one
UPDATE
    users
SET
    google_id = $2 email_verified = TRUE
WHERE
    id = $1
    AND google_id IS NULL
RETURNING
    *;

