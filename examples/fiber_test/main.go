package main

import (
	"log"
)

func main() {
	app := SetupRouter()

	// Serve static files from docs directory
	app.Static("/docs", "./docs")

	log.Printf("Server running at http://localhost:8080")
	log.Printf("Documentation available at:")
	log.Printf("- Swagger UI: http://localhost:8080/docs/index.html")
	log.Printf("- OpenAPI JSON: http://localhost:8080/docs/openapi.json")
	log.Printf("- Routes JSON: http://localhost:8080/docs/routes.json")

	log.Fatal(app.Listen(":8080"))
}
