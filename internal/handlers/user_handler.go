package handlers

import (
	"blog-app/internal/models"
	"blog-app/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// parse the request to get the email and password data
	user := new(models.CreateUserRequest)

	if err := c.BodyParser(user); err != nil {
		return err
	}

	log.Info().Msg(fmt.Sprintf("%v", user))
	ok, err := h.userService.RegisterUser(user.Email, user.Password)
	if err != nil {
		log.Error().Err(err).Msg("userService.RegisterUser: error to register user")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return nil
	}

	if !ok {
		log.Error().Msg("userService.RegisterUser ok false: error to register user")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return nil
	}

	return c.JSON(&models.CreateUserResponse{
		Success: ok,
	})

}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	// parse the request to get the email and password data
	return nil
}
