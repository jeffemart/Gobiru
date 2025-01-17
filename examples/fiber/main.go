package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/routes"
)

func main() {
	app := fiber.New()

	// Setup all routes
	routes.SetupAuthRoutes(app)
	routes.SetupProductRoutes(app)
	routes.SetupOrderRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
