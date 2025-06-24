package routes

import (
	"blog-app/internal/config"
	"blog-app/internal/handlers"
	"blog-app/internal/middleware"
	"blog-app/internal/repositories"
	"blog-app/internal/security"
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
	jwtConfig *security.JWTConfig,
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
	postRepository := repositories.NewPostPostgresRepository(ctx, dbConn)

	// services
	userService := services.NewUserService(userRepository)
	postService := services.NewPostService(postRepository)

	// handlers
	userHandler := handlers.NewUserHandler(userService, config)
	postHandler := handlers.NewPostHandler(postService, config)

	// user endpoints
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// posts endpoints
	// -- protected, Only the author can create, or update/delete his own posts
	v1.Post("/posts", middleware.AuthMiddleware(jwtConfig), postHandler.CreatePost)
	v1.Put("/posts/:id", middleware.AuthMiddleware(jwtConfig), postHandler.UpdatePostByID)
	v1.Delete("/posts/:id", middleware.AuthMiddleware(jwtConfig), postHandler.DeletePostByID)

	v1.Get("/posts/s/:slug", postHandler.GetPostBySlug)
	v1.Get("/posts/:slug", postHandler.GetPostByID)
	v1.Get("/posts/:authorId", postHandler.GetPostsByAuthorID)

}
