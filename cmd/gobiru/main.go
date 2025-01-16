package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jeffemart/gobiru/internal/analyzer"
	"github.com/jeffemart/gobiru/internal/generator"
)

func main() {
	var (
		framework    = flag.String("framework", "", "Framework to analyze (gin, mux, or fiber)")
		mainFile     = flag.String("main", "", "Path to main.go file")
		routerFile   = flag.String("router", "", "Path to routes.go file")
		handlersFile = flag.String("handlers", "", "Path to handlers.go file")
		outputFile   = flag.String("output", "docs/routes.json", "Output path for routes JSON")
		openAPIFile  = flag.String("openapi", "docs/openapi.json", "Output path for OpenAPI spec")
		title        = flag.String("title", "API Documentation", "API title")
		description  = flag.String("description", "", "API description")
		version      = flag.String("version", "1.0.0", "API version")
	)

	flag.Parse()

	// Validar parâmetros obrigatórios
	if *framework == "" {
		log.Fatal("Framework is required (-framework)")
	}
	if *mainFile == "" {
		log.Fatal("Path to main.go is required (-main)")
	}
	if *routerFile == "" {
		log.Fatal("Path to routes.go is required (-router)")
	}
	if *handlersFile == "" {
		log.Fatal("Path to handlers.go is required (-handlers)")
	}

	// Criar analisador baseado no framework
	a, err := analyzer.New(*framework, *mainFile, *routerFile, *handlersFile)
	if err != nil {
		log.Fatalf("Failed to create analyzer: %v", err)
	}

	// Analisar rotas
	routes, err := a.Analyze()
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Gerar documentação JSON
	if err := generator.GenerateJSON(routes, *outputFile); err != nil {
		log.Fatalf("Failed to generate JSON: %v", err)
	}

	// Gerar documentação OpenAPI
	if err := generator.GenerateOpenAPI(routes, generator.APIInfo{
		Title:       *title,
		Description: *description,
		Version:     *version,
	}, *openAPIFile); err != nil {
		log.Fatalf("Failed to generate OpenAPI: %v", err)
	}

	fmt.Println("Documentation generated successfully!")
	fmt.Printf("Routes JSON: %s\n", *outputFile)
	fmt.Printf("OpenAPI spec: %s\n", *openAPIFile)
}
