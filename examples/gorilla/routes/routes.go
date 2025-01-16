package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gobiru/examples/gorilla/complete/handlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/caralho", handlers.ListUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	return r
}
