package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/handlers"
)

func SetupRouter() *fiber.App {
	app := fiber.New()

	app.Get("/users", handlers.ListUsers)
	app.Get("/merda/:id", handlers.GetUser)
	app.Post("/users", handlers.CreateUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)
	app.Get("/users/:userId/posts/:postId", handlers.GetUserPost)
	app.Get("/search", handlers.SearchUsers)

	return app
}
