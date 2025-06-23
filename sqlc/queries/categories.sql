-- name: ListCategories :many
SELECT * FROM categories ORDER BY updated_at DESC;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (name, slug, description, color, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, slug = $3, description = $4, color = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: DeleteCategoryBySlug :exec
DELETE FROM categories WHERE slug = $1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories WHERE slug = $1;