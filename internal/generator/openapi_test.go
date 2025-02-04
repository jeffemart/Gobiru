package generator

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jeffemart/gobiru/internal/spec"
)

func TestOpenAPIGeneration(t *testing.T) {
	doc := &spec.Documentation{
		Operations: []*spec.Operation{
			{
				Path:    "/api/v1/users/{id}",
				Method:  "GET",
				Summary: "GetUser retorna os dados do usuário",
				Parameters: []*spec.Parameter{
					{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &spec.Schema{Type: "string"},
					},
				},
				Responses: map[string]*spec.Response{
					"200": {
						Description: "Successful response",
						Content: map[string]*spec.MediaType{
							"application/json": {
								Schema: &spec.Schema{Type: "object"},
							},
						},
					},
				},
			},
		},
	}

	config := Config{
		OutputFile:  "test_output.json",
		Title:       "Test API",
		Description: "Test Description",
		Version:     "1.0.0",
	}

	gen := NewOpenAPIGenerator()
	err := gen.Generate(doc, config)
	if err != nil {
		t.Fatalf("Failed to generate OpenAPI: %v", err)
	}

	// Validar o JSON gerado
	data, err := os.ReadFile("test_output.json")
	if err != nil {
		t.Fatalf("Failed to read generated OpenAPI: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse generated OpenAPI: %v", err)
	}

	// Validar estrutura OpenAPI
	if result["openapi"] != "3.0.3" {
		t.Errorf("Expected OpenAPI version 3.0.3, got %v", result["openapi"])
	}

	info, ok := result["info"].(map[string]interface{})
	if !ok {
		t.Fatal("Info section not found")
	}

	if info["title"] != "Test API" {
		t.Errorf("Expected title 'Test API', got %v", info["title"])
	}

	paths, ok := result["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("Paths section not found")
	}

	userPath, ok := paths["/api/v1/users/{id}"].(map[string]interface{})
	if !ok {
		t.Fatal("User path not found")
	}

	get, ok := userPath["get"].(map[string]interface{})
	if !ok {
		t.Fatal("GET method not found")
	}

	if get["summary"] != "GetUser retorna os dados do usuário" {
		t.Errorf("Expected summary 'GetUser retorna os dados do usuário', got %v", get["summary"])
	}
} 