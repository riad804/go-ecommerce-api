package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/models"
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
	count, status, err := h.adminService.GetUserCount()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", count)
}

func (h *AdminhHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	status, err := h.adminService.DeleteUser(id)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "user deleted successfully", nil)
}

func (h *AdminhHandler) AddCategory(c *fiber.Ctx) error {
	var req models.CategoryCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	imageUrl, ok := c.Locals("image_url").(string)
	if !ok {
		return models.Error(c, fiber.StatusBadRequest, "Not getting image url from middleware")
	}
	cat := models.Category{
		Name:  req.Name,
		Color: req.Color,
		Image: imageUrl,
	}
	category, status, err := h.adminService.AddCategory(cat)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Category added successfully", category)
}

func (h *AdminhHandler) EditCategory(c *fiber.Ctx) error {
	var req models.CategoryCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	imageUrl, ok := c.Locals("image_url").(string)
	if !ok {
		return models.Error(c, fiber.StatusBadRequest, "Not getting image url from middleware")
	}
	cat := models.Category{
		Name:  req.Name,
		Color: req.Color,
		Image: imageUrl,
	}
	result, status, err := h.adminService.EditCategory(cat)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "updated successfully", result)
}

func (h *AdminhHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	status, err := h.adminService.DeleteCategory(id)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "deleted successfully", nil)
}

// func (h *AdminhHandler) GetProductsCount(c *fiber.Ctx) error {
// }

// func (h *AdminhHandler) AddProduct(c *fiber.Ctx) error {
// }

// func (h *AdminhHandler) EditProduct(c *fiber.Ctx) error {
// }

// func (h *AdminhHandler) DeleteProductImages(c *fiber.Ctx) error {
// }

// func (h *AdminhHandler) DeleteProduct(c *fiber.Ctx) error {
// }

func (h *AdminhHandler) GetOrders(c *fiber.Ctx) error {
	orders, status, err := h.adminService.GetOrders()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", orders)
}

// func (h *AdminhHandler) GetOrderCount(c *fiber.Ctx) error {
// }

// func (h *AdminhHandler) ChangeOrderStatus(c *fiber.Ctx) error {
// }
