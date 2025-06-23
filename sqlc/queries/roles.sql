-- name: CreateRole :one
INSERT INTO roles (
  name,
  description
) VALUES ($1, $2) RETURNING id, name, description, created_at;

-- name: GetRoleByID :one
SELECT id, name, description, created_at FROM roles WHERE id = $1 LIMIT 1;

-- name: GetRoleByName :one
SELECT id, name, description, created_at FROM roles WHERE name = $1 LIMIT 1;

-- name: ListRoles :many
SELECT id, name, description, created_at FROM roles ORDER BY name;

-- name: CreatePermission :one
INSERT INTO permissions (
  name,
  resource,
  action,
  description
) VALUES ($1, $2, $3, $4) RETURNING id, name, resource, action, description, created_at;

-- name: GetPermissionByID :one
SELECT id, name, resource, action, description, created_at FROM permissions WHERE id = $1 LIMIT 1;

-- name: GetPermissionByName :one
SELECT id, name, resource, action, description, created_at FROM permissions WHERE name = $1 LIMIT 1;

-- name: ListPermissions :many
SELECT id, name, resource, action, description, created_at FROM permissions ORDER BY resource, action;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (
  user_id,
  role_id,
  assigned_by
) VALUES ($1, $2, $3);

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (
  role_id,
  permission_id
) VALUES ($1, $2);

-- name: GetUserRoles :many
SELECT r.id, r.name, r.description, ur.assigned_at, ur.assigned_by
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1;

-- name: GetRolePermissions :many
SELECT p.id, p.name, p.resource, p.action, p.description
FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: GetUserPermissions :many
SELECT DISTINCT p.id, p.name, p.resource, p.action, p.description
FROM user_roles ur
JOIN role_permissions rp ON ur.role_id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE ur.user_id = $1; 