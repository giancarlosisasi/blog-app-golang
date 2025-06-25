-- name: ListPosts :many
SELECT * FROM posts ORDER BY updated_at DESC;

-- name: ListPublishedPosts :many
SELECT * FROM posts WHERE status = 'published' ORDER BY updated_at DESC;

-- name: ListDraftPosts :many
SELECT * FROM posts WHERE status = 'draft' ORDER BY updated_at DESC;

-- name: ListArchivedPosts :many
SELECT * FROM posts WHERE status = 'archived' ORDER BY updated_at DESC;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = $1;

-- name: GetPostsByAuthorID :many
SELECT * FROM posts WHERE author_id = $1 ORDER BY updated_at DESC;

-- -- name: GetPostsByStatus :many
-- SELECT * FROM posts WHERE status = $1 ORDER BY updated_at DESC;

-- name: CreatePost :one
INSERT INTO posts (title, slug, content, excerpt, featured_image_url, status, author_id, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
SET title = $2, slug = $3, content = $4, excerpt = $5, featured_image_url = $6, status = $7, author_id = $8, published_at = $9, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: DeletePostBySlug :exec
DELETE FROM posts WHERE slug = $1;

-- name: GetPostsByCategorySlug :many
SELECT p.* FROM posts p
JOIN post_categories pc ON p.id = pc.post_id
JOIN categories c ON pc.category_id = c.id
WHERE c.slug = $1
ORDER BY p.updated_at DESC;

-- name: GetPostsByTagSlug :many
SELECT p.* FROM posts p
JOIN post_tags pt ON p.id = pt.post_id
JOIN tags t ON pt.tag_id = t.id
WHERE t.slug = $1
ORDER BY p.updated_at DESC;

-- name: AssignCategoryToPost :exec
INSERT INTO post_categories (post_id, category_id)
VALUES ($1, $2);

-- name: AssignTagToPost :exec
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2);

-- name: UnassignCategoryFromPost :exec
DELETE FROM post_categories WHERE post_id = $1 AND category_id = $2;

-- name: UnassignTagFromPost :exec
DELETE FROM post_tags WHERE post_id = $1 AND tag_id = $2;

-- name: UpdatePostStatus :exec
UPDATE posts SET status = $2 WHERE id = $1;

-- name: GetPostByIDAndAuthorID :one
SELECT * FROM posts WHERE id = $1 AND author_id = $2;