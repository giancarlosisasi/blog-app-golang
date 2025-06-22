package database

import (
	"blog-app/internal/config"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func SetupDB(config *config.Config) (context.Context, *pgx.Conn, error) {
	ctx := context.Background()
	log.Debug().Msg(config.DBUrl)
	conn, err := pgx.Connect(ctx, config.DBUrl)

	return ctx, conn, err
}
