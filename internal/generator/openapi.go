package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeffemart/gobiru/internal/models"
)

type APIInfo struct {
	Title       string
	Description string
	Version     string
}

// GenerateOpenAPI gera a especificação OpenAPI
func GenerateOpenAPI(routes []models.RouteInfo, info APIInfo, outputFile string) error {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":       info.Title,
			"description": info.Description,
			"version":     info.Version,
		},
		"paths": generatePaths(routes),
	}

	// Criar diretório se não existir
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Gerar JSON
	data, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI spec: %v", err)
	}

	// Salvar arquivo
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func generatePaths(routes []models.RouteInfo) map[string]interface{} {
	paths := make(map[string]interface{})

	for _, route := range routes {
		method := route.Method
		path := route.Path

		operation := map[string]interface{}{
			"summary":    fmt.Sprintf("%s %s", method, path),
			"parameters": generateParameters(route.Parameters),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful operation",
				},
			},
		}

		if _, exists := paths[path]; !exists {
			paths[path] = make(map[string]interface{})
		}
		paths[path].(map[string]interface{})[method] = operation
	}

	return paths
}

func generateParameters(params []models.Parameter) []map[string]interface{} {
	var result []map[string]interface{}

	for _, param := range params {
		result = append(result, map[string]interface{}{
			"name":     param.Name,
			"in":       "path",
			"required": param.Required,
			"schema": map[string]string{
				"type": param.Type,
			},
		})
	}

	return result
}
