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
				Path:   "/api/v1/test/{id}",
				Method: "GET",
				Summary: "TestHandler lida com a rota de teste",
				Parameters: []*spec.Parameter{
					{
						Name:        "id",
						In:          "path",
						Required:    true,
						Description: "Path parameter: id",
						Schema: &spec.Schema{
							Type: "string",
						},
					},
				},
				RequestBody: &spec.RequestBody{
					Required: true,
					Content: map[string]*spec.MediaType{
						"application/json": {
							Schema: &spec.Schema{
								Type: "object",
								Properties: map[string]*spec.Schema{
									"name":  {Type: "string"},
									"email": {Type: "string"},
								},
							},
						},
					},
				},
				Responses: map[string]*spec.Response{
					"200": {
						Description: "Successful response",
						Content: map[string]*spec.MediaType{
							"application/json": {
								Schema: &spec.Schema{
									Type: "object",
								},
							},
						},
					},
				},
			},
		},
	}

	// Configurações para geração da documentação OpenAPI
	config := Config{
		OutputFile:  outputFile,
		Title:       "Test API",
		Description: "Test API Description",
		Version:     "1.0.0",
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