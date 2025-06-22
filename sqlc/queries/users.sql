-- name: CreateUser :one
INSERT INTO users (
  email,
  password_hash,
  username
) VALUES ($1, $2, $3) RETURNING id, email, username;

-- name: GetUserByEmail :one
SELECT id FROM users where email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  username = $1,
  first_name = $2,
  last_name = $3,
  avatar_url = $4,
  bio = $5,
  updated_at = NOW()
WHERE id = $6
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