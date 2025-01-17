package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/handlers"
)

// SetupOrderRoutes configura as rotas de pedidos
func SetupOrderRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	orders := api.Group("/orders")

	orders.Post("", handlers.CreateOrder)
	orders.Get("/:id", handlers.GetOrder)
	orders.Patch("/:id/status", handlers.UpdateOrderStatus)
}
