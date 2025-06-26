package middleware

import (
	"blog-app/internal/security"
	"blog-app/internal/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware to validate JWT from cookies
func AuthMiddleware(config *security.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString, err := security.GetJWTFromCookie(c, config)
		if err != nil {
			if errors.Is(err, utils.ErrAuthTokenMissing) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Authentication token is required",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// validate token
		claims, err := security.ValidateJWT(config, tokenString)
		if err != nil {
			statusCode, message := handleJWTError(err)
			return c.Status(statusCode).JSON(fiber.Map{
				"error": message,
			})
		}

		user := UserContext{
			ID:       claims.UserID,
			Email:    claims.Email,
			Username: claims.Email,
		}
		SetUserContext(c, user)
		return c.Next()
	}
}

func handleJWTError(err error) (statusCode int, message string) {
	switch {
	case errors.Is(err, utils.ErrAuthTokenMissing):
		return fiber.StatusUnauthorized, "Authentication token is required"
	case errors.Is(err, utils.ErrAuthTokenExpired):
		return fiber.StatusUnauthorized, "Authentication token has expired"
	case errors.Is(err, utils.ErrAuthTokenInvalid):
		return fiber.StatusUnauthorized, "Invalid authentication token"
	default:
		return fiber.StatusInternalServerError, "Internal server error"
	}
}
