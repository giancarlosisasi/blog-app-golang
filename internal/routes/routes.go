package routes

import (
	"blog-app/internal/config"
	"blog-app/internal/handlers"
	"blog-app/internal/repositories"
	"blog-app/internal/services"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(
	ctx context.Context,
	srv *fiber.App,
	config *config.Config,
	dbConn *pgx.Conn,
	// stores

	// services
) {

	srv.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	api := srv.Group("/api")
	v1 := api.Group("/v1")

	// DB health check
	v1.Get("/db-health-check", func(c *fiber.Ctx) error {
		if err := dbConn.Ping(ctx); err != nil {
			log.Fatal().Msg("> Can' connect to database")
		}
		return c.SendString("db ok")
	})

	// V1 endpoint health check
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("v1 api")
	})

	// storages (a.k.a repositories)
	userRepository := repositories.NewUserPostgresRepository(ctx, dbConn)

	// services
	userService := services.NewUserService(userRepository)

	// handlers
	userHandler := handlers.NewUserHandler(userService)

	// user endpoints
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)
}
