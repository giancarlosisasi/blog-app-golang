package server

import (
	"blog-app/internal/config"
	"blog-app/internal/routes"
	"blog-app/internal/security"
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
)

func NewServer(
	ctx context.Context,
	dbConn *pgx.Conn,
	conf *config.Config,
	jwtConfig *security.JWTConfig,
) (*fiber.App, error) {
	srv := fiber.New(fiber.Config{
		// https://docs.gofiber.io/guide/faster-fiber
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// =============== CORS ================
	srv.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.SetupRoutes(
		ctx,
		srv,
		conf,
		dbConn,
		jwtConfig,
	)

	err := srv.Listen(fmt.Sprintf(":%s", conf.Port))
	return srv, err
}
