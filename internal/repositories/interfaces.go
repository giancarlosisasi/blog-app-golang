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

type PostRepository interface {
	CreatePost(
		post models.Post,
	) (*models.Post, error)
	GetPostBySlug(
		slug string,
	) (*models.Post, error)
	GetPostByID(
		id string,
	) (*models.Post, error)
	GetPostsByAuthorID(
		authorId string,
	) ([]*models.Post, error)
	UpdatePostByID(
		id string,
	) (*models.Post, error)
	DeletePostByID(
		id string,
	) error
}
