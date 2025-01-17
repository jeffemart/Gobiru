package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/gobiru/examples/gorilla/handlers"
)

func SetupPublicRoutes(r *mux.Router) {
	// API v1 pública
	api := r.PathPrefix("/api/v1").Subrouter()

	// Autenticação
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.HandleFunc("/caralho", handlers.Register).Methods("POST")
	auth.HandleFunc("/forgot-password", handlers.ForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", handlers.ResetPassword).Methods("POST")

	// Documentação e informações públicas
	docs := api.PathPrefix("/docs").Subrouter()
	docs.HandleFunc("", handlers.GetAPIDocumentation).Methods("GET")
	docs.HandleFunc("/swagger", handlers.GetSwaggerUI).Methods("GET")

	// Status da API
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	api.HandleFunc("/status", handlers.APIStatus).Methods("GET")
}
