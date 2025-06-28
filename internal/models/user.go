package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	AvatarUrl *string   `json:"avatar_url"`
	Bio       *string   `json:"bio"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// only public info
type UserProfile struct {
	ID        string  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type UserProfileRequestBody struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BaseUser struct {
	ID       pgtype.UUID `json:"id"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
}

type UserCreated struct {
	BaseUser
}

type CreateUserResponse struct {
	Success bool         `json:"success"`
	Data    *UserCreated `json:"data"`
}

type UserLogged struct {
	BaseUser
}

type UserWithPasswordHash struct {
	BaseUser
	PasswordHash string `json:"-"`
}

type LoginUserResponse struct {
	Success bool `json:"success"`
}
