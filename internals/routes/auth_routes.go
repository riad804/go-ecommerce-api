package routes

import (
	"github.com/riad804/go_ecommerce_api/internals/handlers"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

func (routes *Routes) NewAuthRoutes() {
	userRepo := repositories.NewUserRepository(routes.Mongo.Database)
	authService := service.NewAuthService(*routes.tokenMaker, userRepo, routes.cfg, routes.distributor)
	authHandler := handlers.NewAuthHandler(authService, routes.Validator)

	public := routes.api.Group("/")
	public.Post("/login", authHandler.Login)
	public.Post("/register", authHandler.Register)
	public.Post("/forgot-password", authHandler.ForgotPassword)
	public.Post("/verify-otp", authHandler.VerifyOtp)
	public.Post("/reset-password", authHandler.ResetPassword)
}
