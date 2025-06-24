package middleware

import (
	"blog-app/internal/security"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware to validate JWT from cookies
func AuthMiddleware(config *security.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString, err := security.GetJWTFromCookie(c, config)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// validate token
		claims, err := security.ValidateJWT(config, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token session",
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
