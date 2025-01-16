package fiber

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestFiberAnalyzer(t *testing.T) {
	app := fiber.New()

	// Setup test routes
	app.Get("/users", func(c *fiber.Ctx) error { return nil })
	app.Get("/users/:id", func(c *fiber.Ctx) error { return nil })
	app.Post("/users", func(c *fiber.Ctx) error { return nil })

	analyzer := NewAnalyzer()
	routes := analyzer.GetRoutes()

	// Verify number of routes
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Verify path parameters
	for _, route := range routes {
		if route.Path == "/users/:id" {
			if len(route.Parameters) != 1 {
				t.Errorf("Expected 1 parameter for /users/:id, got %d", len(route.Parameters))
			}
			if route.Parameters[0].Name != "id" {
				t.Errorf("Expected parameter name 'id', got '%s'", route.Parameters[0].Name)
			}
		}
	}
}
