package database

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// SeedData populates the database with initial roles, permissions, and sample data
func SeedData(ctx context.Context, conn *pgx.Conn) error {
	log.Println("Starting database seeding...")

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create roles
	roles, err := createRoles(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to create roles: %w", err)
	}

	// Create permissions
	permissions, err := createPermissions(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to create permissions: %w", err)
	}

	// Assign permissions to roles
	err = assignPermissionsToRoles(ctx, tx, roles, permissions)
	if err != nil {
		return fmt.Errorf("failed to assign permissions to roles: %w", err)
	}

	// Create sample users and assign roles
	err = createSampleUsersAndRoles(ctx, tx, roles)
	if err != nil {
		return fmt.Errorf("failed to create sample users and roles: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// createRoles creates the basic roles for the application
func createRoles(ctx context.Context, tx pgx.Tx) (map[string]uuid.UUID, error) {
	roles := map[string]uuid.UUID{}

	roleData := []struct {
		name        string
		description string
	}{
		{"super_admin", "Super Administrator with full system access"},
		{"admin", "Administrator with management capabilities"},
		{"moderator", "Content moderator with limited administrative powers"},
		{"author", "Content author who can create and edit their own posts"},
		{"editor", "Editor who can create, edit, and publish content"},
		{"user", "Regular user with basic access"},
		{"guest", "Guest user with read-only access"},
	}

	for _, role := range roleData {
		var roleID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO roles (name, description) 
			VALUES ($1, $2) 
			RETURNING id
		`, role.name, role.description).Scan(&roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to create role %s: %w", role.name, err)
		}
		roles[role.name] = roleID
		log.Printf("Created role: %s (ID: %s)", role.name, roleID)
	}

	return roles, nil
}

// createPermissions creates the permissions for different resources
func createPermissions(ctx context.Context, tx pgx.Tx) (map[string]uuid.UUID, error) {
	permissions := map[string]uuid.UUID{}

	permissionData := []struct {
		name        string
		resource    string
		action      string
		description string
	}{
		// User management permissions
		{"users.create", "users", "create", "Create new users"},
		{"users.read", "users", "read", "View user information"},
		{"users.update", "users", "update", "Update user information"},
		{"users.delete", "users", "delete", "Delete users"},
		{"users.manage_roles", "users", "manage_roles", "Assign roles to users"},

		// Post management permissions
		{"posts.create", "posts", "create", "Create new posts"},
		{"posts.read", "posts", "read", "View posts"},
		{"posts.update", "posts", "update", "Update posts"},
		{"posts.delete", "posts", "delete", "Delete posts"},
		{"posts.publish", "posts", "publish", "Publish posts"},
		{"posts.approve", "posts", "approve", "Approve posts for publication"},

		// Comment management permissions
		{"comments.create", "comments", "create", "Create comments"},
		{"comments.read", "comments", "read", "View comments"},
		{"comments.update", "comments", "update", "Update comments"},
		{"comments.delete", "comments", "delete", "Delete comments"},
		{"comments.moderate", "comments", "moderate", "Moderate comments"},

		// Category management permissions
		{"categories.create", "categories", "create", "Create categories"},
		{"categories.read", "categories", "read", "View categories"},
		{"categories.update", "categories", "update", "Update categories"},
		{"categories.delete", "categories", "delete", "Delete categories"},

		// Tag management permissions
		{"tags.create", "tags", "create", "Create tags"},
		{"tags.read", "tags", "read", "View tags"},
		{"tags.update", "tags", "update", "Update tags"},
		{"tags.delete", "tags", "delete", "Delete tags"},

		// System permissions
		{"system.settings", "system", "settings", "Manage system settings"},
		{"system.logs", "system", "logs", "View system logs"},
		{"system.backup", "system", "backup", "Create system backups"},
	}

	for _, perm := range permissionData {
		var permID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO permissions (name, resource, action, description) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id
		`, perm.name, perm.resource, perm.action, perm.description).Scan(&permID)
		if err != nil {
			return nil, fmt.Errorf("failed to create permission %s: %w", perm.name, err)
		}
		permissions[perm.name] = permID
		log.Printf("Created permission: %s (ID: %s)", perm.name, permID)
	}

	return permissions, nil
}

// assignPermissionsToRoles assigns permissions to different roles
func assignPermissionsToRoles(ctx context.Context, tx pgx.Tx, roles map[string]uuid.UUID, permissions map[string]uuid.UUID) error {
	// Super Admin gets all permissions
	for _, permID := range permissions {
		_, err := tx.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_id) 
			VALUES ($1, $2)
		`, roles["super_admin"], permID)
		if err != nil {
			return fmt.Errorf("failed to assign permission to super_admin: %w", err)
		}
	}

	// Admin gets most permissions except system-level ones
	adminPermissions := []string{
		"users.create", "users.read", "users.update", "users.delete", "users.manage_roles",
		"posts.create", "posts.read", "posts.update", "posts.delete", "posts.publish", "posts.approve",
		"comments.create", "comments.read", "comments.update", "comments.delete", "comments.moderate",
		"categories.create", "categories.read", "categories.update", "categories.delete",
		"tags.create", "tags.read", "tags.update", "tags.delete",
	}

	for _, permName := range adminPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["admin"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to admin: %w", permName, err)
			}
		}
	}

	// Moderator permissions
	moderatorPermissions := []string{
		"users.read",
		"posts.read", "posts.approve",
		"comments.create", "comments.read", "comments.update", "comments.delete", "comments.moderate",
		"categories.read",
		"tags.read",
	}

	for _, permName := range moderatorPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["moderator"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to moderator: %w", permName, err)
			}
		}
	}

	// Author permissions
	authorPermissions := []string{
		"posts.create", "posts.read", "posts.update", "posts.delete",
		"comments.create", "comments.read", "comments.update", "comments.delete",
		"categories.read",
		"tags.read",
	}

	for _, permName := range authorPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["author"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to author: %w", permName, err)
			}
		}
	}

	// Editor permissions
	editorPermissions := []string{
		"posts.create", "posts.read", "posts.update", "posts.delete", "posts.publish",
		"comments.create", "comments.read", "comments.update", "comments.delete",
		"categories.read",
		"tags.read",
	}

	for _, permName := range editorPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["editor"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to editor: %w", permName, err)
			}
		}
	}

	// User permissions
	userPermissions := []string{
		"posts.read",
		"comments.create", "comments.read", "comments.update", "comments.delete",
		"categories.read",
		"tags.read",
	}

	for _, permName := range userPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["user"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to user: %w", permName, err)
			}
		}
	}

	// Guest permissions
	guestPermissions := []string{
		"posts.read",
		"comments.read",
		"categories.read",
		"tags.read",
	}

	for _, permName := range guestPermissions {
		if permID, exists := permissions[permName]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id) 
				VALUES ($1, $2)
			`, roles["guest"], permID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to guest: %w", permName, err)
			}
		}
	}

	log.Println("Permissions assigned to roles successfully")
	return nil
}

// createSampleUsersAndRoles creates sample users and assigns them roles
func createSampleUsersAndRoles(ctx context.Context, tx pgx.Tx, roles map[string]uuid.UUID) error {
	// Create sample users
	users := []struct {
		email     string
		username  string
		firstName string
		lastName  string
		role      string
	}{
		{"admin@blogapp.com", "admin", "John", "Admin", "super_admin"},
		{"moderator@blogapp.com", "moderator", "Sarah", "Moderator", "moderator"},
		{"author@blogapp.com", "author", "Mike", "Author", "author"},
		{"editor@blogapp.com", "editor", "Lisa", "Editor", "editor"},
		{"user@blogapp.com", "user", "Tom", "User", "user"},
		{"guest@blogapp.com", "guest", "Guest", "User", "guest"},
	}

	for _, user := range users {
		// Create user
		var userID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO users (email, username, password_hash, first_name, last_name, is_active) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id
		`, user.email, user.username, "$2a$10$dummy.hash.for.seeding", user.firstName, user.lastName, true).Scan(&userID)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.email, err)
		}

		// Assign role to user
		if roleID, exists := roles[user.role]; exists {
			_, err := tx.Exec(ctx, `
				INSERT INTO user_roles (user_id, role_id, assigned_by) 
				VALUES ($1, $2, $3)
			`, userID, roleID, userID) // Self-assigned for seeding
			if err != nil {
				return fmt.Errorf("failed to assign role %s to user %s: %w", user.role, user.email, err)
			}
		}

		log.Printf("Created user: %s (%s) with role: %s", user.email, user.username, user.role)
	}

	return nil
}

// GetUserPermissions retrieves all permissions for a given user
func GetUserPermissions(ctx context.Context, conn *pgx.Conn, userID uuid.UUID) ([]Permission, error) {
	rows, err := conn.Query(ctx, `
		SELECT DISTINCT p.id, p.name, p.resource, p.action, p.description, p.created_at
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1
		ORDER BY p.resource, p.action
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// Permission represents a permission in the system
type Permission struct {
	ID          pgtype.UUID        `json:"id"`
	Name        string             `json:"name"`
	Resource    string             `json:"resource"`
	Action      string             `json:"action"`
	Description string             `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
}
