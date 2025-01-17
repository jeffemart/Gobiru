package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/handlers"
)

// SetupProductRoutes configura as rotas de produtos
func SetupProductRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	products := api.Group("/products")

	products.Get("", handlers.ListProducts)
	products.Get("/:id", handlers.GetProduct)
	products.Post("", handlers.CreateProduct)
}
