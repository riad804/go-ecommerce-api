package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

type UserhHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService, validate *validator.Validate) *UserhHandler {
	return &UserhHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UserhHandler) GetUsers(c *fiber.Ctx) error {
	result, status, err := h.userService.GetAllUser()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", result)
}

func (h *UserhHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	result, status, err := h.userService.GetUserById(id)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", result)
}

func (h *UserhHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	result, status, err := h.userService.UpdateUser(id, req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "user updated successfully", result)
}
