package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	gobiru "github.com/jeffemart/Gobiru/app"
	"github.com/jeffemart/Gobiru/app/openapi"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Criar router
	router := mux.NewRouter()

	// Configurar rotas
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Criar analisador de rotas
	analyzer := gobiru.NewRouteAnalyzer()

	// Analisar rotas
	err := analyzer.AnalyzeRoutes(router)
	if err != nil {
		log.Fatalf("Erro ao analisar rotas: %v", err)
	}

	// Exportar documentação JSON
	err = analyzer.ExportJSON("routes.json")
	if err != nil {
		log.Fatalf("Erro ao exportar JSON: %v", err)
	}

	// Exportar documentação OpenAPI
	info := openapi.Info{
		Title:       "API de Exemplo",
		Description: "API de exemplo para testar o Gobiru",
		Version:     "1.0.0",
	}

	err = analyzer.ExportOpenAPI("openapi.json", info)
	if err != nil {
		log.Fatalf("Erro ao exportar OpenAPI: %v", err)
	}

	log.Println("Documentação gerada com sucesso!")
	log.Println("JSON: routes.json")
	log.Println("OpenAPI: openapi.json")

	// Iniciar servidor (opcional, apenas para teste)
	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: "1", Name: "João", Email: "joao@example.com"},
		{ID: "2", Name: "Maria", Email: "maria@example.com"},
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := User{ID: vars["id"], Name: "João", Email: "joao@example.com"}
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = vars["id"]
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
