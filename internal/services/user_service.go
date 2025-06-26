package services

import (
	"blog-app/internal/models"
	"blog-app/internal/repositories"
	"blog-app/internal/security"
	"blog-app/internal/utils"
	"blog-app/internal/validators"
	"strings"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterUser(email string, password string) (*models.UserCreated, error) {
	// business logic
	// email format validation
	email = validators.SanitizeInput(email)
	ok := validators.ValidateEmail(email)
	if !ok {
		return nil, utils.NewCustomError(
			utils.CREATE_USER_INVALID_EMAIL_ERROR,
			"invalid email",
		)
	}

	password = validators.SanitizeInput(password)
	// password validation (length, strength)
	ok = validators.HasAtLeastOneCharacter(password)
	if !ok {
		return nil, utils.NewCustomError(
			utils.CREATE_USER_INVALID_PASSWORD_ERROR,
			"the password must have at least on character",
		)
	}
	ok = validators.ValidateLength(password, 6, 20)
	if !ok {
		return nil, utils.NewCustomError(
			utils.CREATE_USER_INVALID_PASSWORD_ERROR,
			"password must be between 6 and 20 characters",
		)
	}
	ok = validators.HasAtLeastOneNumber(password)
	if !ok {
		return nil, utils.NewCustomError(
			utils.CREATE_USER_INVALID_PASSWORD_ERROR,
			"password must include at least one number",
		)
	}

	// check if email is already in db using the repository checkEmailExists()
	exists, err := s.userRepository.CheckIfEmailExists(email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, utils.NewCustomError(
			utils.USER_ALREADY_EXISTS_IN_DB_ERROR,
			"the user is already register",
		)
	}

	// hash password
	hashed_password, err := security.HashPassword(password)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, utils.NewCustomError(
			utils.CANNOT_GET_USER_WITH_EMAIL_ERROR,
			"error to process password",
		)
	}

	// create username
	// the username by default will be the username plus the domain
	username := strings.ReplaceAll(email, "@", "-")
	username = strings.ReplaceAll(username, ".", "-")

	user, err := s.userRepository.RegisterUser(email, hashed_password, username)
	if err != nil {
		log.Error().Err(err).Msg("userRepository: failed to register user")
		return nil, utils.NewCustomError(
			utils.CREATE_USER_FAILED_TO_CREATE_ERROR,
			"failed to register user",
		)
	}

	return user, nil

}

func (s *UserService) LoginUser(email string, password string) (*models.UserLogged, error) {
	// check if user exists in the db
	user, err := s.userRepository.GetUserByEmailWithPassword(email)
	if err != nil {
		return nil, utils.NewCustomError(
			utils.USER_NOT_FOUND_ERROR,
			"user not found",
		)
	}
	// check if the password is valid
	valid, err := security.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to verifyPassword")
		return nil, utils.NewCustomError(
			utils.USER_INVALID_PASSWORD_ERROR,
			"invalid password",
		)
	}

	if !valid {
		return nil, utils.NewCustomError(
			utils.USER_INVALID_PASSWORD_ERROR,
			"invalid password",
		)
	}

	// return the user logged
	return &models.UserLogged{
		BaseUser: models.BaseUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
