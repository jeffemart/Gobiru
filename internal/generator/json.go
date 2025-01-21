package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type JSONGenerator struct{}

func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{}
}

func (g *JSONGenerator) Generate(doc *spec.Documentation, config Config) error {
	if doc == nil {
		return fmt.Errorf("documentation is nil")
	}

	if len(doc.Operations) == 0 {
		return fmt.Errorf("no operations found to generate documentation")
	}

	operations := make([]map[string]interface{}, 0)

	for _, op := range doc.Operations {
		if op == nil {
			continue
		}

		// Limpar e normalizar o path
		path := strings.TrimPrefix(op.Path, "/api/v1")
		path = strings.TrimPrefix(path, "/api")
		path = "/api/v1" + path

		// Extrair o recurso base para as tags
		var tags []string
		segments := strings.Split(strings.Trim(path, "/"), "/")
		if len(segments) > 1 {
			// Usar o primeiro segmento significativo após api/v1
			for _, segment := range segments {
				if segment != "api" && segment != "v1" && !strings.HasPrefix(segment, "{") {
					tags = append(tags, segment)
					break
				}
			}
		}
		if len(tags) == 0 {
			tags = append(tags, "general")
		}

		// Extrair o nome do handler do Summary
		handlerName := op.Summary
		if idx := strings.Index(handlerName, " "); idx > 0 {
			handlerName = handlerName[:idx]
		}

		operation := map[string]interface{}{
			"method":       op.Method,
			"path":         path,
			"description":  op.Summary,
			"handler_name": handlerName,
			"tags":         tags,
		}

		// Adicionar parâmetros se existirem
		if len(op.Parameters) > 0 {
			operation["parameters"] = convertParameters(op.Parameters)
		}

		// Adicionar query parameters se existirem
		queryParams := extractQueryParameters(op.Parameters)
		if len(queryParams) > 0 {
			operation["query_parameters"] = queryParams
		}

		// Adicionar headers padrão
		operation["headers"] = []map[string]interface{}{
			{
				"name":     "Content-Type",
				"value":    "application/json",
				"required": true,
			},
		}

		// Adicionar request body se existir
		if op.RequestBody != nil {
			operation["request_body"] = map[string]interface{}{
				"type":   "object",
				"schema": convertRequestBodySchema(op.RequestBody),
			}
		}

		// Adicionar responses se existirem
		if len(op.Responses) > 0 {
			operation["responses"] = convertResponses(op.Responses)
		}

		// Adicionar campos adicionais
		operation["authentication"] = map[string]interface{}{
			"type":     "bearer",
			"required": true,
		}
		operation["estimated_time_ms"] = 100
		operation["permissions"] = []string{"read", "write"}
		operation["api_version"] = "v1.0"
		operation["deprecated"] = false
		operation["rate_limit"] = map[string]interface{}{
			"requests_per_minute": 100,
			"time_window_seconds": 60,
		}
		operation["notes"] = op.Summary

		operations = append(operations, operation)
	}

	// Ordenar operações por path e método
	sort.Slice(operations, func(i, j int) bool {
		pathI := operations[i]["path"].(string)
		pathJ := operations[j]["path"].(string)
		if pathI != pathJ {
			return pathI < pathJ
		}
		return operations[i]["method"].(string) < operations[j]["method"].(string)
	})

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
