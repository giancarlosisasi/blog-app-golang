package repositories

import (
	"blog-app/internal/models"
)

type UserRepository interface {
	RegisterUser(
		email string,
		hashedPassword string,
		username string,
	) (*models.UserCreated, error)

	GetUserByEmail(
		email string,
	) (user *models.BaseUser, err error)

	CheckIfEmailExists(
		email string,
	) (bool, error)

	GetUserByEmailWithPassword(
		email string,
	) (*models.UserWithPasswordHash, error)
}
