package gin

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinAnalyzer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup test routes
	router.GET("/users", func(c *gin.Context) {})
	router.GET("/users/:id", func(c *gin.Context) {})
	router.POST("/users", func(c *gin.Context) {})

	analyzer := NewGinAnalyzer()
	err := analyzer.AnalyzeRoutes(router)
	if err != nil {
		t.Fatalf("Failed to analyze routes: %v", err)
	}

	routes := analyzer.GetRoutes()
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Verify path parameters are extracted correctly
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