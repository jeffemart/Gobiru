package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/jeffemart/gobiru/internal/analyzer"
	"github.com/jeffemart/gobiru/internal/generator"
)

func main() {
	var (
		framework   string
		mainFile    string
		title       string
		description string
		version     string
	)

	flag.StringVar(&framework, "framework", "", "Framework usado (gin, fiber, mux)")
	flag.StringVar(&mainFile, "main", "", "Arquivo principal da aplicação")
	// Removendo a flag -openapi
	// flag.StringVar(&openapi, "openapi", "", "Path to output OpenAPI file")

	// Definindo um caminho padrão para o arquivo OpenAPI
	openapi := "docs/openapi.json" // ou "examples/gorilla/docs/openapi.json"

	flag.StringVar(&title, "title", "", "Title for OpenAPI documentation")
	flag.StringVar(&description, "description", "", "Description for OpenAPI documentation")
	flag.StringVar(&version, "version", "", "Version for OpenAPI documentation")

	flag.Parse()

	// Se a flag -main não for passada, tentar encontrar o main.go automaticamente
	if mainFile == "" {
		var err error
		mainFile, err = analyzer.FindMainFile(".") // Passando o diretório atual
		if err != nil {
			log.Fatalf("Error finding main.go: %v", err)
		}
	}

	// Resolver caminho absoluto do main.go
	absMainPath, err := filepath.Abs(mainFile)
	if err != nil {
		log.Fatalf("Failed to resolve main file path: %v", err)
	}

	config := analyzer.Config{
		MainFile: absMainPath,
	}

	// Criar analisador baseado no framework
	analyzer, err := analyzer.New(framework, config)
	if err != nil {
		log.Fatalf("Failed to create analyzer: %v", err)
	}

	doc, err := analyzer.Analyze()
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Configurações para geração da documentação OpenAPI
	openapiConfig := generator.Config{
		OutputFile:  openapi,
		Title:       title,
		Description: description,
		Version:     version,
	}

	// Gerar documentação OpenAPI
	openapiGen := generator.NewOpenAPIGenerator()
	if err := openapiGen.Generate(doc, openapiConfig); err != nil {
		log.Fatalf("Failed to generate OpenAPI documentation: %v", err)
	}

	log.Printf("OpenAPI documentation generated successfully!")
	log.Printf("OpenAPI documentation: %s", openapi)
}

// stringSliceFlag implementa a interface flag.Value para aceitar múltiplos valores
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
