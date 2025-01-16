package main

import (
	"flag"
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
		outputFile   = flag.String("output", "", "Path to output JSON file")
		openAPIFile  = flag.String("openapi", "", "Path to output OpenAPI file")
		title        = flag.String("title", "", "Title for OpenAPI documentation")
		description  = flag.String("description", "", "Description for OpenAPI documentation")
		version      = flag.String("version", "", "Version for OpenAPI documentation")
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
	analyzer, err := analyzer.New(*framework, *mainFile, *routerFile, *handlersFile)
	if err != nil {
		log.Fatalf("Failed to create analyzer: %v", err)
	}

	doc, err := analyzer.Analyze()
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Configurações para geração da documentação
	jsonConfig := generator.Config{
		OutputFile: *outputFile,
	}

	openapiConfig := generator.Config{
		OutputFile:  *openAPIFile,
		Title:       *title,
		Description: *description,
		Version:     *version,
	}

	// Gerar documentação
	jsonGen := generator.NewJSONGenerator()
	if err := jsonGen.Generate(doc, jsonConfig); err != nil {
		log.Fatalf("Failed to generate JSON documentation: %v", err)
	}

	openapiGen := generator.NewOpenAPIGenerator()
	if err := openapiGen.Generate(doc, openapiConfig); err != nil {
		log.Fatalf("Failed to generate OpenAPI documentation: %v", err)
	}

	log.Printf("Documentation generated successfully!")
	log.Printf("JSON documentation: %s", *outputFile)
	log.Printf("OpenAPI documentation: %s", *openAPIFile)
}
