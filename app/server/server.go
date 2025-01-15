package server

import (
	"fmt"
	"log"
	"net/http"
)

// Serve starts the documentation server
func Serve(port int, docsDir string) error {
	// Servir arquivos estáticos do diretório docs
	fs := http.FileServer(http.Dir(docsDir))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server running at http://localhost%s", addr)
	log.Printf("Documentation available at:")
	log.Printf("- Swagger UI: http://localhost%s/docs/index.html", addr)
	log.Printf("- OpenAPI JSON: http://localhost%s/docs/openapi.json", addr)
	log.Printf("- Routes JSON: http://localhost%s/docs/routes.json", addr)

	return http.ListenAndServe(addr, nil)
}