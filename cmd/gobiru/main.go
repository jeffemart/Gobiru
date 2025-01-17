package main

import (
	"flag"
	"log"
	"strings"

	"github.com/jeffemart/gobiru/internal/analyzer"
	"github.com/jeffemart/gobiru/internal/generator"
)

func main() {
	var (
		framework   string
		mainFile    string
		output      string
		openapi     string
		title       string
		description string
		version     string
	)

	flag.StringVar(&framework, "framework", "", "Framework usado (gin, fiber, mux)")
	flag.StringVar(&mainFile, "main", "", "Arquivo principal da aplicação")
	flag.StringVar(&output, "output", "", "Path to output JSON file")
	flag.StringVar(&openapi, "openapi", "", "Path to output OpenAPI file")
	flag.StringVar(&title, "title", "", "Title for OpenAPI documentation")
	flag.StringVar(&description, "description", "", "Description for OpenAPI documentation")
	flag.StringVar(&version, "version", "", "Version for OpenAPI documentation")

	flag.Parse()

	// Validar parâmetros obrigatórios
	if framework == "" {
		log.Fatal("Framework is required (-framework)")
	}
	if mainFile == "" {
		log.Fatal("Path to main.go is required (-main)")
	}

	config := analyzer.Config{
		MainFile: mainFile,
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

	// Configurações para geração da documentação
	jsonConfig := generator.Config{
		OutputFile: output,
	}

	openapiConfig := generator.Config{
		OutputFile:  openapi,
		Title:       title,
		Description: description,
		Version:     version,
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
	log.Printf("JSON documentation: %s", output)
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
