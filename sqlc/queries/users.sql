-- name: CreateUser :one
INSERT INTO users (
  email,
  password_hash,
  username
) VALUES ($1, $2, $3) RETURNING id, email, username;

-- name: GetUserByEmail :one
SELECT id, email, username FROM users where email = $1 LIMIT 1;

-- name: GetUserByEmailWithPassword :one
SELECT id, email, username, password_hash FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, username, first_name, last_name, avatar_url, bio, is_active, created_at, updated_at FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  first_name = $1,
  last_name = $2,
  avatar_url = $3,
  bio = $4,
  updated_at = NOW()
WHERE id = $5
  AND is_active = true
RETURNING id, email, username, first_name, last_name, avatar_url, bio, updated_at;

-- name: UpdateIsActive :one
UPDATE users
SET
  is_active = $1
WHERE id = $2
RETURNING id, email, username, first_name, last_name, avatar_url, bio, updated_at;

-- name: UpdatePassword :one
UPDATE users
SET
  password_hash = $1
WHERE id = $2 AND is_active = true
RETURNING id, email, username, first_name, last_name, avatar_url, bio, updated_at;

-- name: DeleteUser :one
UPDATE users
SET
  is_active = false
WHERE id = $1 AND is_active = true
RETURNING id;