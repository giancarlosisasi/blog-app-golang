package services

import (
	"blog-app/internal/models"
	"blog-app/internal/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterUser(email string, password string) (*models.CreateUserResponse, error) {
	// business logic
	// email format validation

	// password validation (length, strength)

	// check if email is already in db using the repository checkEmailExists()

	// hash password

	_, err := s.userRepository.RegisterUser(email, password)
	if err != nil {
		return nil, err
	}

	resp := models.CreateUserResponse{
		Success: true,
	}

	return &resp, nil

}
