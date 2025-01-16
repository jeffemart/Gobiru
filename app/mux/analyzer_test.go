package mux

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestMuxAnalyzer(t *testing.T) {
	router := mux.NewRouter()

	// Setup test routes
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")

	analyzer := NewAnalyzer()
	routes := analyzer.GetRoutes()

	// Verify number of routes
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Verify path parameters
	for _, route := range routes {
		if route.Path == "/users/{id}" {
			if len(route.Parameters) != 1 {
				t.Errorf("Expected 1 parameter for /users/{id}, got %d", len(route.Parameters))
			}
			if route.Parameters[0].Name != "id" {
				t.Errorf("Expected parameter name 'id', got '%s'", route.Parameters[0].Name)
			}
		}
	}
}
