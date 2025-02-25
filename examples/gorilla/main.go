package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffemart/gobiru/examples/gorilla/routes"
)

func main() {
	r := mux.NewRouter()

	// Setup das rotas
	routes.SetupRoutes(r)

	// Iniciando o servidor
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
