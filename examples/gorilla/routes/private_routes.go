package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/gobiru/examples/gorilla/handlers"
)

func SetupPrivateRoutes(r *mux.Router) {
	// API v1 privada
	api := r.PathPrefix("/api/v1").Subrouter()

	// Stores and Categories
	stores := api.PathPrefix("/outrocaralho/{storeId}").Subrouter()
	categories := stores.PathPrefix("/categories/{categoryId}").Subrouter()

	// Products
	products := categories.PathPrefix("/products").Subrouter()
	products.HandleFunc("", handlers.CreateProduct).Methods("POST")
	products.HandleFunc("/{productId}", handlers.GetProduct).Methods("GET")
	products.HandleFunc("/{productId}", handlers.UpdateProduct).Methods("PUT")
	products.HandleFunc("/{productId}", handlers.DeleteProduct).Methods("DELETE")

	// Orders
	orders := stores.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("", handlers.CreateOrder).Methods("POST")
	orders.HandleFunc("/{orderId}", handlers.GetOrder).Methods("GET")
	orders.HandleFunc("/{orderId}/status", handlers.UpdateOrderStatus).Methods("PATCH")

	// Employees
	employees := stores.PathPrefix("/employees").Subrouter()
	employees.HandleFunc("", handlers.CreateEmployee).Methods("POST")
	employees.HandleFunc("/{employeeId}", handlers.GetEmployee).Methods("GET")
	employees.HandleFunc("/{employeeId}", handlers.UpdateEmployee).Methods("PUT")
	employees.HandleFunc("/{employeeId}/status", handlers.UpdateEmployeeStatus).Methods("PATCH")
}
