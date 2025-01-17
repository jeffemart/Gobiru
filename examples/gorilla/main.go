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
	routes.SetupPublicRoutes(r)  // Rotas públicas (auth, docs, status)
	routes.SetupPrivateRoutes(r) // Rotas privadas (produtos, pedidos)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
