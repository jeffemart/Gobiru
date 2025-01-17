package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type JSONGenerator struct{}

func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{}
}

func (g *JSONGenerator) Generate(doc *spec.Documentation, config Config) error {
	if doc == nil || len(doc.Operations) == 0 {
		return fmt.Errorf("no operations found to generate documentation")
	}

	operations := make([]map[string]interface{}, 0)

	for _, op := range doc.Operations {
		operation := map[string]interface{}{
			"method":           op.Method,
			"path":             op.Path,
			"description":      op.Summary,
			"handler_name":     strings.Split(op.Path, "/")[len(strings.Split(op.Path, "/"))-1],
			"parameters":       convertParameters(op.Parameters),
			"query_parameters": extractQueryParameters(op.Parameters),
			"headers": []map[string]interface{}{
				{
					"name":     "Content-Type",
					"value":    "application/json",
					"required": true,
				},
			},
			"request_body": map[string]interface{}{
				"type":   "object",
				"schema": convertRequestBodySchema(op.RequestBody),
			},
			"responses": convertResponses(op.Responses),
			"tags": []string{
				strings.Split(op.Path, "/")[1], // Primeiro segmento após a raiz
			},
			"authentication": map[string]interface{}{
				"type":     "bearer",
				"required": true,
			},
			"estimated_time_ms": 100,
			"permissions":       []string{"read", "write"},
			"api_version":       "v1.0",
			"deprecated":        false,
			"rate_limit": map[string]interface{}{
				"requests_per_minute": 100,
				"time_window_seconds": 60,
			},
			"notes": op.Summary,
		}
		operations = append(operations, operation)
	}

	fmt.Printf("Generating documentation for %d operations\n", len(operations))

	return writeJSON(config.OutputFile, operations)
}

func writeJSON(filename string, data interface{}) error {
	// Criar diretório se não existir
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Converter para JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Escrever arquivo
	return os.WriteFile(filename, jsonData, 0644)
}

func extractQueryParameters(params []*spec.Parameter) []map[string]interface{} {
	queryParams := make([]map[string]interface{}, 0)
	for _, p := range params {
		if p.In == "query" {
			queryParams = append(queryParams, map[string]interface{}{
				"name":        p.Name,
				"type":        "string",
				"required":    p.Required,
				"description": p.Description,
			})
		}
	}
	return queryParams
}

func convertRequestBodySchema(body *spec.RequestBody) map[string]interface{} {
	if body == nil || len(body.Content) == 0 {
		return nil
	}

	for _, mediaType := range body.Content {
		if mediaType.Schema != nil {
			return convertSchema(mediaType.Schema)
		}
	}
	return nil
}
