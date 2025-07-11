package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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

// func (s *AdminService) GetProductsCount() error {
// }

// func (s *AdminService) AddProduct() error {
// }

// func (s *AdminService) EditProduct() error {
// }

// func (s *AdminService) DeleteProductImages() error {
// }

// func (s *AdminService) DeleteProduct() error {
// }

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

func (s *AdminService) ChangeOrderStatus(id string) (int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
}
