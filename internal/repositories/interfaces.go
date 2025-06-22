package repositories

import (
	database "blog-app/internal/database/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	RegisterUser(
		email string,
		hashedPassword string,
		username string,
	) (pgtype.UUID, error)
	LoginUser(
		email string,
		rawPassword string,
	) (user *database.User, token string, err error)
	CheckIfEmailExists(
		email string,
	) (bool, error)
}
