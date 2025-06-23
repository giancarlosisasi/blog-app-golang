package main

import (
	"blog-app/internal/config"
	"blog-app/internal/database"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Load configuration
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup database connection
	ctx, conn, err := database.SetupDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Check if we should clear existing data
	if len(os.Args) > 1 && os.Args[1] == "--clear" {
		log.Println("Clearing existing roles and permissions data...")
		if err := clearExistingData(ctx, conn); err != nil {
			log.Fatalf("Failed to clear existing data: %v", err)
		}
	}

	// Run seeding
	if err := database.SeedData(ctx, conn); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

// clearExistingData removes existing roles, permissions, and related data
func clearExistingData(ctx context.Context, conn *pgx.Conn) error {
	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Clear data in the correct order due to foreign key constraints
	queries := []string{
		"DELETE FROM user_roles",
		"DELETE FROM role_permissions",
		"DELETE FROM permissions",
		"DELETE FROM roles",
	}

	for _, query := range queries {
		if _, err := tx.Exec(ctx, query); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
