package main

import (
	"blog-app/internal/config"
	"blog-app/internal/database"
	"blog-app/internal/logger"
	"blog-app/internal/security"
	"blog-app/internal/server"
	"context"
	"github.com/rs/zerolog/log"
)

func main() {
	// Viper
	conf, err := config.SetupConfig()
	if err != nil {
		log.Fatal().Msg("> error to load the configuration.")
		return
	}

	jwtConfig := security.JWTDefaultConfig(
		conf,
	)
	// ============== Zerolog =================
	logger.SetupLogger()

	// ========== DB ===============
	ctx, dbConn, err := database.SetupDB(conf)
	if err != nil {
		log.Fatal().Msg("> Can' connect to database")
		return
	}
	if err := dbConn.Ping(ctx); err != nil {
		log.Fatal().Msg("> Can' connect to database")
		return
	}
	defer dbConn.Close(context.Background())
	log.Info().Msg("> ✅ Success to connect to database")

	// Setup stores

	// Server
	app, err := server.NewServer(ctx, dbConn, conf, jwtConfig)

	if err != nil {
		log.Fatal().Err(err)
		log.Fatal().Msg("the server has crashed!")

		app.Shutdown()

	}
}
