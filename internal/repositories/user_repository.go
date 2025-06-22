package repositories

import (
	database "blog-app/internal/database/queries"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *UserPostgresRepository) RegisterUser(email string, hashedPassword string, username string) (pgtype.UUID, error) {
	u, err := database.New(r.dbConn).CreateUser(r.ctx, database.CreateUserParams{
		Email:        email,
		PasswordHash: hashedPassword,
		Username:     username,
	})

	if err != nil {
		var id pgtype.UUID
		return id, err
	}

	return u, nil
}

func (r *UserPostgresRepository) LoginUser(email, rawPassword string) (user *database.User, token string, err error) {
	return nil, "", nil
}

func (r *UserPostgresRepository) CheckIfEmailExists(email string) (bool, error) {
	userId, err := database.New(r.dbConn).GetUserByEmail(r.ctx, email)
	if err != nil {
		return false, nil
	}

	return userId.Valid, nil
}
