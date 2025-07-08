package service

import (
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/riad804/go_ecommerce_api/workers"
)

type AdminService struct {
	tokenMaker  token.Maker
	userRepo    repositories.UserRepository
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewAdminService(tokenMaker token.Maker, userRepo repositories.UserRepository, cfg *config.Config, distributor workers.TaskDistributor) *AdminService {
	return &AdminService{
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		cfg:         cfg,
		distributor: distributor,
	}
}
