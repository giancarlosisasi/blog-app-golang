package repositories

import (
	database "blog-app/internal/database/queries"
)

type UserRepository interface {
	RegisterUser(
		email string,
		password string,
	) (*database.User, error)
	LoginUser(
		email string,
		password string,
	) (user *database.User, token string, err error)
}
