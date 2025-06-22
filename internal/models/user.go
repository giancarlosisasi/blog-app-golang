package models

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Bio       *string   `json:"bio"`
	AvatarURL *string   `json:"avatar_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// no password hash here!
}

// only public info
type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Bio       *string   `json:"bio"`
	AvatarURL *string   `json:"avatar_url"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Success bool `json:"success"`
}
