package main

import (
	"log"
	"net/http"
)

func main() {
	router := SetupRouter()

	// Servir arquivos estáticos do diretório docs
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	log.Printf("Servidor rodando em http://localhost:8080")
	log.Printf("Documentação disponível em:")
	log.Printf("- Swagger UI: http://localhost:8080/docs/index.html")
	log.Printf("- OpenAPI JSON: http://localhost:8080/docs/openapi.json")
	log.Printf("- Routes JSON: http://localhost:8080/docs/routes.json")

	log.Fatal(http.ListenAndServe(":8080", router))
}
