package main

import (
	"blog-app/internal/config"
	"blog-app/internal/database"
	db "blog-app/internal/database/queries"
	"blog-app/internal/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func main() {
	// Viper
	conf, err := config.SetupConfig()
	if err != nil {
		log.Fatal().Msg("> error to load the configuration.")
		return
	}
	// Zerolog
	logger.SetupLogger()

	// DB
	ctx, conn, err := database.SetupDB(conf)
	if err != nil {
		log.Fatal().Msg("> Can' connect to database")
		return
	}

	queries := db.New(conn)

	// Convert string to UUID
	userID, err := uuid.Parse("12345678-1234-1234-1234-123456789abc")
	if err != nil {
		log.Error().Err(err).Msg("> Invalid UUID format")
		return
	}

	// Create pgtype.UUID
	pgUUID := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}

	user, err := queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		log.Error().Err(err).Msg("> Failed to get user")
		return
	}

	log.Info().Str("email", user.Email).Msg("> User retrieved successfully")
}
