package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jeffemart/gobiru/examples/fiber/routes"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
