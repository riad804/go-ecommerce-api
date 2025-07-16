package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/middlewares"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	imageUrl, _, err := middlewares.ImageMiddleware(c)
	if err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if imageUrl == nil {
		return models.Error(c, fiber.StatusInternalServerError, "image not uploaded")
	}

	cat := models.Category{
		Name:  req.Name,
		Color: req.Color,
		Image: *imageUrl,
	}
	category, status, err := h.adminService.AddCategory(cat)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Category added successfully", category)
}

func (h *AdminhHandler) EditCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Error(c, fiber.StatusBadRequest, "invalid category id")
	}

	var req models.CategoryCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	imageUrl, _, err := middlewares.ImageMiddleware(c)
	if err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}

	cat := models.Category{
		ID:    obId,
		Name:  req.Name,
		Color: req.Color,
	}
	if imageUrl != nil {
		cat.Image = *imageUrl
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

func (h *AdminhHandler) GetProductsCount(c *fiber.Ctx) error {
	count, status, err := h.adminService.GetProductsCount()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", count)
}

func (h *AdminhHandler) GetProducts(c *fiber.Ctx) error {
	return nil //todo implementation
}

func (h *AdminhHandler) AddProduct(c *fiber.Ctx) error {
	var req models.ProductCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	imageUrl, gallery, err := middlewares.ImageMiddleware(c)
	if err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if imageUrl == nil {
		return models.Error(c, fiber.StatusInternalServerError, "image not uploaded")
	}
	if gallery == nil {
		return models.Error(c, fiber.StatusInternalServerError, "images not uploaded")
	}
	product, status, err := h.adminService.AddProduct(req, *imageUrl, *gallery)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Product added successfully", product)
}

func (h *AdminhHandler) EditProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.ProductUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.validate.Struct(&req); err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	imageUrl, gallery, err := middlewares.ImageMiddleware(c)
	if err != nil {
		return models.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if imageUrl != nil {
		req.ImageUrl = imageUrl
	}
	if gallery != nil {
		req.Gallery = gallery
	}
	product, status, err := h.adminService.EditProduct(id, req)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Product updated successfully", product)
}

func (h *AdminhHandler) DeleteProductImages(c *fiber.Ctx) error {
	return nil //todo implementation
}

func (h *AdminhHandler) DeleteProduct(c *fiber.Ctx) error {
	return nil //todo implementation
}

func (h *AdminhHandler) GetOrders(c *fiber.Ctx) error {
	orders, status, err := h.adminService.GetOrders()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", orders)
}

func (h *AdminhHandler) GetOrderCount(c *fiber.Ctx) error {
	count, status, err := h.adminService.GetOrderCount()
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "success", count)
}

func (h *AdminhHandler) ChangeOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	status := c.FormValue("status")
	result, err := h.adminService.ChangeOrderStatus(id, models.OrderStatus(status))
	if err != nil {
		return models.Error(c, result, err.Error())
	}
	return models.Success(c, result, "success", nil)
}

func (h *AdminhHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	status, err := h.adminService.DeleteOrder(id)
	if err != nil {
		return models.Error(c, status, err.Error())
	}
	return models.Success(c, status, "Deleted successfully", nil)
}
