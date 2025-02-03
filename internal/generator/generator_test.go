package generator

import (
	"testing"
	"os"
	"path/filepath"

	"github.com/jeffemart/gobiru/internal/spec"
)

func TestGenerateOpenAPIDocumentation(t *testing.T) {
	// Crie um diretório temporário para o teste
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "openapi.json")

	// Crie um exemplo de documentação
	doc := &spec.Documentation{
		Operations: []*spec.Operation{
			{
				Path:   "/api/v1/products",
				Method: "GET",
				Summary: "List all products",
			},
		},
	}

	// Configurações para geração da documentação OpenAPI
	config := Config{
		OutputFile: outputFile,
		Title:      "Test API",
		Description: "Test API Description",
		Version:    "1.0.0",
	}

	// Gerar documentação OpenAPI
	gen := NewOpenAPIGenerator()
	if err := gen.Generate(doc, config); err != nil {
		t.Fatalf("Failed to generate OpenAPI documentation: %v", err)
	}

	// Verifique se o arquivo foi criado
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected OpenAPI documentation file to be created: %s", outputFile)
	}
} 