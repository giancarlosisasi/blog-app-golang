package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type UserContext struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

const UserContextKey = "user_context"

func SetUserContext(c *fiber.Ctx, user UserContext) {
	c.Locals(UserContextKey, user)
}

func GetUserContext(c *fiber.Ctx) (*UserContext, bool) {
	user, ok := c.Locals(UserContextKey).(UserContext)
	if !ok {
		return nil, false
	}

	return &user, true
}
