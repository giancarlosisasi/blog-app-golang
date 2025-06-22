package handlers

import (
	"blog-app/internal/config"
	"blog-app/internal/models"
	"blog-app/internal/security"
	"blog-app/internal/services"
	"blog-app/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *services.UserService
	appConfig   *config.Config
}

func NewUserHandler(userService *services.UserService, config *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		appConfig:   config,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// parse the request to get the email and password data
	rawUserData := new(models.CreateUserRequest)

	if err := c.BodyParser(rawUserData); err != nil {
		return err
	}

	u, err := h.userService.RegisterUser(rawUserData.Email, rawUserData.Password)
	if err != nil {
		log.Error().Err(err).Msg("userService.RegisterUser: error to register user")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return nil
	}

	jwtConfig := security.JWTDefaultConfig(h.appConfig)

	tokenString, err := security.GenerateJWT(
		jwtConfig,
		u.ID,
		u.Email,
		u.Username,
	)
	if err != nil {
		log.Error().Err(err).Msg("GenerateJWT: error to generate the jwt token")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewCustomError(
			utils.INTERNAL_SERVER_ERROR,
			"internal server error",
		))
	}

	security.SetJWTCookie(c, jwtConfig, tokenString)

	return c.Status(fiber.StatusCreated).JSON(&models.CreateUserResponse{
		Success: true,
		Data:    u,
	})

}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	// parse the request to get the email and password data
	return nil
}
