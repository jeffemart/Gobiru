package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffemart/gobiru/examples/gorilla/routes"
)

func main() {
	r := mux.NewRouter()

	// Setup das rotas
	routes.SetupRoutes(r)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
