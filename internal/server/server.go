package server

import (
	"blog-app/internal/config"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
)

func NewServer(
	ctx context.Context,
	dbConn *pgx.Conn,
	conf *config.Config,
	// stores

	// services
) (*fiber.App, error) {
	srv := fiber.New(fiber.Config{
		// https://docs.gofiber.io/guide/faster-fiber
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// =============== CORS ================
	srv.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	addRoutes(
		ctx,
		srv,
		conf,
		dbConn,
	)

	err := srv.Listen(fmt.Sprintf(":%s", conf.Port))
	return srv, err
}
