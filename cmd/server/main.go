package main

import (
	"blog-app/internal/config"
	"blog-app/internal/database"
	"blog-app/internal/logger"
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
	// ============== Zerolog =================
	logger.SetupLogger()

	// ========== DB ===============
	ctx, conn, err := database.SetupDB(conf)
	if err != nil {
		log.Fatal().Msg("> Can' connect to database")
		return
	}
	if err := conn.Ping(ctx); err != nil {
		log.Fatal().Msg("> Can' connect to database")
		return
	}
	defer conn.Close(context.Background())
	log.Info().Msg("> ✅ Success to connect to database")
}
