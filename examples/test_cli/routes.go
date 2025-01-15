package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")

	return router
}

func getUsers(w http.ResponseWriter, r *http.Request)   {}
func getUser(w http.ResponseWriter, r *http.Request)    {}
func createUser(w http.ResponseWriter, r *http.Request) {}
