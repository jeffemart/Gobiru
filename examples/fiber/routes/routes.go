package routes

import (
	"github.com/jeffemart/gobiru/examples/fiber/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// API v1
	v1 := app.Group("/api/v1")

	// Organizações e Times
	orgs := v1.Group("/organizations/:orgId")
	teams := orgs.Group("/teams/:teamId")

	// Usuários
	users := teams.Group("/users")
	users.Post("/", handlers.CreateUser) // Criar usuário
}
