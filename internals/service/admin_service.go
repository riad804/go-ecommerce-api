package service

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/helpers"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/riad804/go_ecommerce_api/workers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	tokenMaker  token.Maker
	userRepo    repositories.UserRepository
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewAdminService(
	tokenMaker token.Maker,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	cfg *config.Config,
	distributor workers.TaskDistributor,
) *AdminService {
	return &AdminService{
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cfg:         cfg,
		distributor: distributor,
	}
}

func (s *AdminService) GetUserCount() (*int64, int, error) {
	count, err := s.userRepo.CountAll()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return &count, fiber.StatusOK, nil
}

func (s *AdminService) DeleteUser(id string) (int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	user, err := s.userRepo.FindByID(obId)
	if err != nil {
		return fiber.StatusNotFound, fmt.Errorf("user not found")
	}
	orders, _ := s.orderRepo.FindOrderByUserId(user.ID)
	var orderItemIds []primitive.ObjectID
	for _, order := range orders {
		orderItemIds = append(orderItemIds, order.OrderItems...)
	}
	err = s.orderRepo.DeleteOrderByUserId(user.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}
	if len(orderItemIds) > 0 {
		err := s.orderRepo.DeleteOrderItems(orderItemIds)
		if err != nil {
			return fiber.StatusInternalServerError, err
		}
	}
	err = s.orderRepo.DeleteCartByUserId(user.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	err = s.userRepo.DeleteByID(user.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusNoContent, nil
}

func (s *AdminService) AddCategory(cat models.Category) (*models.Category, int, error) {
	result, err := s.productRepo.CategorySave(cat)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("InsertedID is not an ObjectID")
	}
	category, err := s.productRepo.CategoryFindOne(oid)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("Category not found after created")
	}
	return category, fiber.StatusCreated, nil
}

func (s *AdminService) EditCategory(category models.Category) (*models.Category, int, error) {
	cat, err := s.productRepo.CategoryUpdate(category)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return cat, fiber.StatusAccepted, nil
}

func (s *AdminService) DeleteCategory(id string) (int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	err = s.productRepo.CategoryDeleteById(obId)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}
	return fiber.StatusNoContent, nil
}

func (s *AdminService) GetProductsCount() (*int64, int, error) {
	count, err := s.productRepo.CountAllProducts()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return &count, fiber.StatusOK, nil
}

func (s *AdminService) GetProductDetails() error {
	return nil //todo implementation
}

func (s *AdminService) AddProduct(req models.ProductCreateRequest, image string, gallery []string) (*models.Product, int, error) {
	obId, err := primitive.ObjectIDFromHex(req.CategoryId)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("invalid category id")
	}
	category, err := s.productRepo.CategoryFindOne(obId)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("invalid category")
	}
	if category.MarkedForDelete {
		return nil, fiber.StatusNotFound, fmt.Errorf("Category marked for deletion, you cannot add products to this category.")
	}
	product := models.Product{
		Name:              req.Name,
		Description:       req.Description,
		Price:             req.Price,
		Colors:            req.Colors,
		Image:             image,
		Images:            gallery,
		Sizes:             req.Sizes,
		GenderAgeCategory: req.GenderAgeCategory,
		CountInStock:      req.CountInStock,
		DateAdded:         time.Now(),
	}
	product.Category = append(product.Category, obId)

	result, err := s.productRepo.ProductSave(product)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("InsertedID is not an ObjectID")
	}
	productResult, err := s.productRepo.ProductFindOne(oid)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("Product not found after created")
	}
	return productResult, fiber.StatusCreated, nil
}

func (s *AdminService) EditProduct(id string, req models.ProductUpdateRequest) (*models.Product, int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("invalid product id")
	}
	if req.CategoryId != nil {
		catId, err := primitive.ObjectIDFromHex(*req.CategoryId)
		if err != nil {
			return nil, fiber.StatusBadRequest, fmt.Errorf("invalid category id")
		}
		category, err := s.productRepo.CategoryFindOne(catId)
		if err != nil {
			return nil, fiber.StatusNotFound, fmt.Errorf("invalid category")
		}
		if category.MarkedForDelete {
			return nil, fiber.StatusNotFound, fmt.Errorf("Category marked for deletion, you cannot add products to this category.")
		}
	}
	product, err := s.productRepo.ProductFindOne(obId)
	if err != nil {
		return nil, fiber.StatusNotFound, err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Colors != nil {
		product.Colors = *req.Colors
	}
	if req.Sizes != nil {
		product.Sizes = *req.Sizes
	}
	if req.GenderAgeCategory != nil {
		product.GenderAgeCategory = *req.GenderAgeCategory
	}
	if req.CountInStock != nil {
		product.CountInStock = *req.CountInStock
	}
	if req.ImageUrl != nil {
		product.Image = *req.ImageUrl
	}
	if req.Gallery != nil {
		product.Images = *req.Gallery
	}
	updatedProduct, err := s.productRepo.ProductUpdate(*product)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return updatedProduct, fiber.StatusAccepted, nil
}

func (s *AdminService) DeleteProductImages() error {
	return nil //todo implementation
}

func (s *AdminService) DeleteProduct() error {
	return nil //todo implementation
}

func (s *AdminService) GetOrders() ([]models.OrderResponse, int, error) {
	orders, err := s.orderRepo.FindAllOrders()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return orders, fiber.StatusOK, nil
}

func (s *AdminService) GetOrderCount() (*int64, int, error) {
	count, err := s.orderRepo.FindOrdersCount()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return count, fiber.StatusOK, nil
}

func (s *AdminService) ChangeOrderStatus(id string, status models.OrderStatus) (int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	order, err := s.orderRepo.FindOrderById(obId)
	if err != nil {
		return fiber.StatusNotFound, fmt.Errorf("Order not found")
	}
	if !helpers.Contains(order.StatusHistory, order.Status) {
		order.StatusHistory = append(order.StatusHistory, order.Status)
	}
	order.Status = status
	_, err = s.orderRepo.UpdateOrder(*order)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}
	return fiber.StatusAccepted, nil
}

func (s *AdminService) DeleteOrder(id string) (int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	result := s.orderRepo.DeleteOrderById(obId)
	if err := result.Err(); err != nil {
		return fiber.StatusNotFound, err
	}
	var order models.Order
	err = result.Decode(order)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("order decoding failed")
	}
	err = s.orderRepo.DeleteOrderItems(order.OrderItems)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}
	return fiber.StatusNoContent, nil
}
