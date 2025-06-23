-- name: ListTags :many
SELECT * FROM tags ORDER BY updated_at DESC;

-- name: GetTagByID :one
SELECT * FROM tags WHERE id = $1;

-- name: CreateTag :one
INSERT INTO tags (name, slug, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: UpdateTag :one
UPDATE tags
SET name = $2, slug = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: DeleteTagBySlug :exec
DELETE FROM tags WHERE slug = $1;

-- name: GetTagBySlug :one
SELECT * FROM tags WHERE slug = $1;