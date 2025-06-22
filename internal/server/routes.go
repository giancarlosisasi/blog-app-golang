package server

import (
	"blog-app/internal/config"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func addRoutes(
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
}
