package gin

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinAnalyzer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup test routes
	router.GET("/users", func(c *gin.Context) {})
	router.GET("/users/:id", func(c *gin.Context) {})
	router.POST("/users", func(c *gin.Context) {})

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

func TestGinAnalyzerWithRouterFile(t *testing.T) {
	analyzer := NewAnalyzer()
	routes, err := analyzer.AnalyzeFile("../../examples/gin_test/routes.go")
	if err != nil {
		t.Fatalf("Failed to analyze routes: %v", err)
	}

	// Verify routes were parsed
	if len(routes) == 0 {
		t.Error("Expected routes to be parsed, got none")
	}

	// Verify specific routes
	foundUserRoute := false
	for _, route := range routes {
		if route.Path == "/users/:id" {
			foundUserRoute = true
			if len(route.Parameters) != 1 {
				t.Errorf("Expected 1 parameter for /users/:id, got %d", len(route.Parameters))
			}
			if route.Parameters[0].Name != "id" {
				t.Errorf("Expected parameter name 'id', got '%s'", route.Parameters[0].Name)
			}
		}
	}

	if !foundUserRoute {
		t.Error("Expected to find /users/:id route")
	}
}
