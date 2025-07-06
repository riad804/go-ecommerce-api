package service

import (
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
)

type AuthService struct {
	tokenMaker token.Maker
	userRepo   repositories.UserRepository
}

func NewAuthService(tokenMaker token.Maker, userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		tokenMaker: tokenMaker,
	}
}
