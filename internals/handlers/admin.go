package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

type AdminhHandler struct {
	adminService *service.AdminService
	validate     *validator.Validate
}

func NewAdminHandler(adminService *service.AdminService, validate *validator.Validate) *AdminhHandler {
	return &AdminhHandler{
		adminService: adminService,
		validate:     validate,
	}
}

func (h *AdminhHandler) GetUserCount(c *fiber.Ctx) error {
}

func (h *AdminhHandler) DeleteUser(c *fiber.Ctx) error {
}

func (h *AdminhHandler) AddCategory(c *fiber.Ctx) error {
}

func (h *AdminhHandler) EditCategory(c *fiber.Ctx) error {
}

func (h *AdminhHandler) DeleteCategory(c *fiber.Ctx) error {
}

func (h *AdminhHandler) GetProductsCount(c *fiber.Ctx) error {
}

func (h *AdminhHandler) AddProduct(c *fiber.Ctx) error {
}

func (h *AdminhHandler) EditProduct(c *fiber.Ctx) error {
}

func (h *AdminhHandler) DeleteProductImages(c *fiber.Ctx) error {
}

func (h *AdminhHandler) DeleteProduct(c *fiber.Ctx) error {
}

func (h *AdminhHandler) GetOrders(c *fiber.Ctx) error {
}

func (h *AdminhHandler) GetOrderCount(c *fiber.Ctx) error {
}

func (h *AdminhHandler) ChangeOrderStatus(c *fiber.Ctx) error {
}
