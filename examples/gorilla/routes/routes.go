package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/gobiru/examples/gorilla/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// API v1
	api := r.PathPrefix("/api/v1").Subrouter()

	// Stores and Categories
	stores := api.PathPrefix("/stores/{storeId}").Subrouter()
	categories := stores.PathPrefix("/categories/{categoryId}").Subrouter()

	// Products
	products := categories.PathPrefix("/products").Subrouter()
	products.HandleFunc("", handlers.CreateProduct).Methods("POST")

	return r
}
