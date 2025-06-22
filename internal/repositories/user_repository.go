package repositories

import (
	database "blog-app/internal/database/queries"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserPostgresRepository struct {
	ctx    context.Context
	dbConn *pgx.Conn
}

func NewUserPostgresRepository(ctx context.Context, dbConn *pgx.Conn) *UserPostgresRepository {
	return &UserPostgresRepository{
		ctx:    ctx,
		dbConn: dbConn,
	}
}

func (r *UserPostgresRepository) RegisterUser(email string, password string) (*database.User, error) {
	return nil, nil
}

func (r *UserPostgresRepository) LoginUser(email, password string) (user *database.User, token string, err error) {
	return nil, "", nil
}
