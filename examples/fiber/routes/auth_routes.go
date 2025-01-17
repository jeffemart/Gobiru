package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/handlers"
)

// SetupAuthRoutes configura as rotas de autenticação
func SetupAuthRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	users := api.Group("/users")

	users.Get("/:id", handlers.GetUser)
	users.Post("", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
}
