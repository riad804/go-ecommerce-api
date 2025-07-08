package routes

import (
	"github.com/riad804/go_ecommerce_api/internals/handlers"
	"github.com/riad804/go_ecommerce_api/internals/middlewares"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

func (routes *Routes) NewUserRoutes() {
	userRepo := repositories.NewUserRepository(routes.Mongo.Database)
	userService := service.NewUserService(*routes.tokenMaker, userRepo, routes.cfg, routes.distributor)
	userHandler := handlers.NewUserHandler(userService, routes.Validator)

	authMiddleware := middlewares.AuthMiddleware(*routes.tokenMaker)
	api := routes.api.Group("/users", authMiddleware)
	api.Get("/", userHandler.GetUsers)
	api.Get("/:id", userHandler.GetUserById)
	api.Put("/:id", userHandler.UpdateUser)
}
