package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/riad804/go_ecommerce_api/workers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	tokenMaker  token.Maker
	userRepo    repositories.UserRepository
	orderRepo   repositories.OrderRepository
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewAdminService(
	tokenMaker token.Maker,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	cfg *config.Config,
	distributor workers.TaskDistributor,
) *AdminService {
	return &AdminService{
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		orderRepo:   orderRepo,
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
	orders, _ := s.orderRepo.FindByUserId(user.ID)
	orders[0].OrderItems
}

func (s *AdminService) AddCategory() error {
}

func (s *AdminService) EditCategory() error {
}

func (s *AdminService) DeleteCategory() error {
}

func (s *AdminService) GetProductsCount() error {
}

func (s *AdminService) AddProduct() error {
}

func (s *AdminService) EditProduct() error {
}

func (s *AdminService) DeleteProductImages() error {
}

func (s *AdminService) DeleteProduct() error {
}

func (s *AdminService) GetOrders() error {
}

func (s *AdminService) GetOrderCount() error {
}

func (s *AdminService) ChangeOrderStatus() error {
}
