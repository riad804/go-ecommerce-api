package routes

import "github.com/riad804/go_ecommerce_api/internals/handlers"

func (routes *Routes) NewProductsRoute() {
	productsHandler := handlers.NewProductsHandler()

	products := routes.api.Group("/products")

	products.Get("/count", productsHandler.GetCount)
	products.Get("/:id", productsHandler.GetDetails)
	products.Delete("/:id", productsHandler.Delete)
	products.Put("/:id", productsHandler.Update)
}
