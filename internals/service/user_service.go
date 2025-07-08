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

type UserService struct {
	tokenMaker  token.Maker
	userRepo    repositories.UserRepository
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewUserService(tokenMaker token.Maker, userRepo repositories.UserRepository, cfg *config.Config, distributor workers.TaskDistributor) *UserService {
	return &UserService{
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		cfg:         cfg,
		distributor: distributor,
	}
}

func (s *UserService) GetAllUser() ([]models.User, int, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return users, fiber.StatusOK, nil
}

func (s *UserService) GetUserById(id string) (*models.User, int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	user, err := s.userRepo.FindByID(obId)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}
	user.Password = ""
	user.ResetPasswordOtp = nil
	user.ResetPasswordOtpExpires = nil
	return user, fiber.StatusOK, nil
}

func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, int, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("invalid user id")
	}
	user, err := s.userRepo.FindByID(obId)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone

	updatedUser, err := s.userRepo.Update(*user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	updatedUser.Password = ""
	return updatedUser, fiber.StatusAccepted, nil
}
