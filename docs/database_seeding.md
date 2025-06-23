# Database Seeding

This document explains how to use the database seeding functionality to populate your database with initial roles, permissions, and sample data.

## Overview

The seeding system creates a comprehensive role-based access control (RBAC) system with:

- **7 Roles**: super_admin, admin, moderator, author, editor, user, guest
- **30+ Permissions**: covering users, posts, comments, categories, tags, and system management
- **6 Sample Users**: one for each role with realistic data

## Roles and Permissions

### Roles Hierarchy

1. **super_admin** - Full system access
2. **admin** - Management capabilities (no system-level access)
3. **moderator** - Content moderation powers
4. **editor** - Can create, edit, and publish content
5. **author** - Can create and edit their own content
6. **user** - Basic access (read posts, create comments)
7. **guest** - Read-only access

### Permission Categories

- **User Management**: create, read, update, delete, manage_roles
- **Post Management**: create, read, update, delete, publish, approve
- **Comment Management**: create, read, update, delete, moderate
- **Category Management**: create, read, update, delete
- **Tag Management**: create, read, update, delete
- **System Management**: settings, logs, backup

## Usage

### Basic Seeding

To seed the database with initial data:

```bash
task seed
```

Or run directly:

```bash
go run cmd/seed/main.go
```

### Clear and Reseed

To clear existing roles and permissions data before seeding:

```bash
task seed:clear
```

Or run directly:

```bash
go run cmd/seed/main.go --clear
```

## Sample Data Created

### Sample Users

| Email | Username | Role | Description |
|-------|----------|------|-------------|
| admin@blogapp.com | admin | super_admin | John Admin |
| moderator@blogapp.com | moderator | moderator | Sarah Moderator |
| author@blogapp.com | author | author | Mike Author |
| editor@blogapp.com | editor | editor | Lisa Editor |
| user@blogapp.com | user | user | Tom User |
| guest@blogapp.com | guest | guest | Guest User |

### Sample Roles

- **super_admin**: All permissions
- **admin**: All permissions except system-level
- **moderator**: User read, post approval, comment moderation
- **author**: Post CRUD, comment CRUD, read categories/tags
- **editor**: Post CRUD + publish, comment CRUD, read categories/tags
- **user**: Read posts, comment CRUD, read categories/tags
- **guest**: Read-only access to posts, comments, categories, tags

## Database Schema

The seeding creates data in these tables:

- `roles` - Role definitions
- `permissions` - Permission definitions
- `role_permissions` - Role-permission assignments
- `user_roles` - User-role assignments

## Integration with Your Application

### Getting User Permissions

```go
import "blog-app/internal/database"

// Get all permissions for a user
permissions, err := database.GetUserPermissions(ctx, conn, userID)
if err != nil {
    // handle error
}

// Check if user has specific permission
hasPermission := false
for _, perm := range permissions {
    if perm.Name == "posts.create" {
        hasPermission = true
        break
    }
}
```

### Permission Checking Helper

You can create a helper function to check permissions:

```go
func HasPermission(permissions []database.Permission, resource, action string) bool {
    for _, perm := range permissions {
        if perm.Resource == resource && perm.Action == action {
            return true
        }
    }
    return false
}

// Usage
if HasPermission(userPermissions, "posts", "create") {
    // User can create posts
}
```

## Security Notes

- Sample users have dummy password hashes - change them in production
- The `--clear` flag removes ALL roles and permissions data
- Always backup your database before running seeding in production
- Consider creating a separate seeding script for production with real passwords

## Customization

To customize the seeding data:

1. Modify the `createRoles()` function in `internal/database/seed.go`
2. Modify the `createPermissions()` function to add/remove permissions
3. Modify the `assignPermissionsToRoles()` function to change role assignments
4. Modify the `createSampleUsersAndRoles()` function to change sample users

## Troubleshooting

### Common Issues

1. **Foreign Key Violations**: Ensure migrations are run before seeding
2. **Duplicate Key Errors**: Use `--clear` flag to remove existing data
3. **Connection Issues**: Check your database configuration in `.env`

### Logs

The seeding process provides detailed logs showing:
- Each role created with its ID
- Each permission created with its ID
- Each user created with their role assignment
- Success/failure messages for each step 