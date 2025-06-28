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

	GetUserByID(
		id string,
	) (user *models.User, err error)

	CheckIfEmailExists(
		email string,
	) (bool, error)

	GetUserByEmailWithPassword(
		email string,
	) (*models.UserWithPasswordHash, error)
	UpdateUserProfileById(id string, updateProfile models.UserProfile) (*models.UserProfile, error)
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
	) ([]models.Post, error)
	UpdatePostByID(
		authorId string,
		id string,
		data *models.UpdatePostData,
	) (*models.Post, error)
	DeletePostByID(
		id string,
	) error
}
