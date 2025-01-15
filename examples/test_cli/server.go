package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8081, "Porta do servidor")
	flag.Parse()

	// Verificar se o diretório docs existe
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		log.Fatal("Diretório 'docs' não encontrado. Execute primeiro o comando para gerar a documentação.")
	}

	// Servir arquivos estáticos do diretório docs
	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Servidor rodando em http://localhost%s", addr)
	log.Printf("Documentação disponível em:")
	log.Printf("- Swagger UI: http://localhost%s/docs/index.html", addr)
	log.Printf("- OpenAPI JSON: http://localhost%s/docs/openapi.json", addr)
	log.Printf("- Routes JSON: http://localhost%s/docs/routes.json", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
