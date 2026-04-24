-- name: InsertUser :one
INSERT INTO users (
    name,
    email,
    email_verified,
    role,
    password_hash,
    google_id,
    avatar_url
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: InsertUserOnConflict :one
INSERT INTO users (
    name,
    email,
    email_verified,
    role,
    password_hash,
    google_id,
    avatar_url
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (email) DO NOTHING
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByIDForUpdate :one
SELECT *
FROM users
WHERE id = $1
FOR UPDATE;

-- name: GetUserByGoogleID :one
SELECT *
FROM users
WHERE google_id = $1;

-- name: UpdateUserVerification :one
UPDATE users
SET email_verified = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: AttachGoogleToUser :one
UPDATE users
SET
    google_id = $2,
    email_verified = TRUE
WHERE
    id = $1
    AND google_id IS NULL
RETURNING *;
