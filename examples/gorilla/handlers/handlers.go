package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: "1", Name: "User 1"},
		{ID: "2", Name: "User 2"},
	}
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := User{ID: vars["id"], Name: "User " + vars["id"]}
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
