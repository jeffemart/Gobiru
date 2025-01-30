package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffemart/gobiru/internal/analyzer"
	"github.com/jeffemart/gobiru/internal/generator"

	swag "github.com/go-openapi/spec"
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

	// Resolver caminho absoluto do main.go
	absMainPath, err := filepath.Abs(mainFile)
	if err != nil {
		log.Fatalf("Failed to resolve main file path: %v", err)
	}

	if framework == "fiber" {
		doc, err := analyzer.Fiber(mainFile)
		if err != nil {
			log.Fatalf("Failed to analyze routes: %v", err)
		}
		doc.Info = &swag.Info{
			InfoProps: swag.InfoProps{
				Title:       title,
				Description: description,
				Version:     version,
			},
		}
		json, err := doc.MarshalJSON()
		jsonTemp := string(json)
		jsonTemp = strings.ReplaceAll(jsonTemp, `{"info":`, `{"openapi":"3.1.0","info":`)
		json = []byte(jsonTemp)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		err = os.WriteFile(openapi, json, 0644)
		if err != nil {
			log.Fatalf("Failed to write JSON file: %v", err)
		}

		log.Printf("Documentation generated successfully!")
		log.Printf("OpenAPI documentation: %s", openapi)
	} else {

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
}
