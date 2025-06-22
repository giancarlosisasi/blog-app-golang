-- name: CreateUser :one
INSERT INTO users (
  email,
  username,
  password_hash
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
