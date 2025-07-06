package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	tokenMaker token.Maker
	userRepo   repositories.UserRepository
}

func NewAuthService(tokenMaker token.Maker, userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		tokenMaker: tokenMaker,
		userRepo:   userRepo,
	}
}

func (s *AuthService) CreateUser(req models.RegisterRequest) (*models.User, int, error) {
	hashPassword, err := config.HashPassword(req.Password)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPassword,
		Phone:    req.Phone,
		IsAdmin:  false,
	}
	r, err := s.userRepo.Create(user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	oid, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("InsertedID is not an ObjectID")
	}
	result, err := s.userRepo.FindByID(oid)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return result, fiber.StatusAccepted, nil
}
