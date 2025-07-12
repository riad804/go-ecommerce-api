package routes

import (
	"github.com/riad804/go_ecommerce_api/internals/handlers"
	"github.com/riad804/go_ecommerce_api/internals/middlewares"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

func (routes *Routes) NewAdminRoutes() {
	userRepo := repositories.NewUserRepository(routes.Mongo.Database)
	orderRepo := repositories.NewOrderRepository(routes.Mongo.Database)
	productRepo := repositories.NewProductRepository(routes.Mongo.Database)
	adminService := service.NewAdminService(*routes.tokenMaker, userRepo, orderRepo, productRepo, routes.cfg, routes.distributor)
	adminHandler := handlers.NewAdminHandler(adminService, routes.Validator)

	authMiddleware := middlewares.AuthMiddleware(*routes.tokenMaker)
	imageMiddleware := middlewares.ImageMiddleware()

	api := routes.api.Group("/admin", authMiddleware)

	//users
	api.Get("/users/count", adminHandler.GetUserCount)
	api.Delete("/users/:id", adminHandler.DeleteUser)

	// categories
	api.Post("/categories", imageMiddleware, adminHandler.AddCategory)
	api.Put("/categories/:id", adminHandler.EditCategory)
	api.Delete("/categories/:id", adminHandler.DeleteCategory)

	//products
	api.Get("/products/count", adminHandler.GetProductsCount)
	api.Post("/products/count", adminHandler.AddProduct)
	api.Put("/products/:id", adminHandler.EditProduct)
	api.Delete("/products/:id/images", adminHandler.DeleteProductImages)
	api.Delete("/products/:id", adminHandler.DeleteProduct)

	//orders
	api.Get("/orders", adminHandler.GetOrders)
	api.Get("/orders/count", adminHandler.GetOrderCount)
	api.Put("/orders/:id", adminHandler.ChangeOrderStatus)
	api.Delete("/orders/:id", adminHandler.DeleteOrder)
}
