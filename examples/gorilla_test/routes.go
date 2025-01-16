package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Task representa uma tarefa
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// SetupRouter configura as rotas da API
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Rotas para tarefas
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	// Rota de health check
	r.HandleFunc("/health", healthCheck).Methods("GET")

	return r
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{
		{ID: "1", Title: "Implementar API", Description: "Criar API REST com Gorilla Mux", Done: true},
		{ID: "2", Title: "Documentar API", Description: "Usar Gobiru para documentação", Done: false},
	}
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	task := Task{
		ID:          vars["id"],
		Title:       "Tarefa Exemplo",
		Description: "Descrição da tarefa de exemplo",
		Done:        false,
	}
	json.NewEncoder(w).Encode(task)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = "new-id"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var task Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = vars["id"]
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
