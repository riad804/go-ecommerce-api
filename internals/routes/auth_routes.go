package routes

import (
	"github.com/riad804/go_ecommerce_api/internals/handlers"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

func (routes *Routes) NewAuthRoutes() {
	userRepo := repositories.NewUserRepository(routes.Mongo.Database)
	authService := service.NewAuthService(*routes.tokenMaker, userRepo)
	authHandler := handlers.NewAuthHandler(authService, routes.Validator)

	user := routes.api.Group("/")

	user.Post("/login", authHandler.Login)
	user.Post("/register", authHandler.Register)
	user.Post("/forgot-password", authHandler.ForgotPassword)
	user.Post("/verify-otp", authHandler.VerifyOtp)
	user.Post("/reset-password", authHandler.ResetPassword)
}
