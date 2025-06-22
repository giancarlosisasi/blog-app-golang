package handlers

import (
	"blog-app/internal/models"
	"blog-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	// parse the request to get the email and password data
	return ctx.JSON(models.CreateUserResponse{
		Success: true,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	// parse the request to get the email and password data
	return nil
}
