package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	result, status, err := h.authService.FindUser(req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}

	return models.Success(c, status, "login success", result)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, status, err := h.authService.CreateUser(req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "register success", user)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req models.ForgotPassRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}

	_, status, err := h.authService.ForgotPassword(req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "otp sent successfully", nil)
}

func (h *AuthHandler) VerifyOtp(c *fiber.Ctx) error {
	var req models.VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}

	_, status, err := h.authService.OtpVerify(req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "otp confirmed successfully", nil)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}

	_, status, err := h.authService.ResetPassword(req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Password reset successfully", nil)
}
