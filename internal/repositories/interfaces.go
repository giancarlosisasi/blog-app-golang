package repositories

import (
	database "blog-app/internal/database/queries"
	"blog-app/internal/models"
)

type UserRepository interface {
	RegisterUser(
		email string,
		hashedPassword string,
		username string,
	) (*models.UserCreated, error)
	LoginUser(
		email string,
		rawPassword string,
	) (user *database.User, token string, err error)
	CheckIfEmailExists(
		email string,
	) (bool, error)
}
