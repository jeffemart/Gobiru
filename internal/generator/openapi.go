package generator

import (
	"fmt"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type OpenAPIGenerator struct{}

func NewOpenAPIGenerator() *OpenAPIGenerator {
	return &OpenAPIGenerator{}
}

func (g *OpenAPIGenerator) Generate(doc *spec.Documentation, config Config) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	openapi := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       config.Title,
			"description": config.Description,
			"version":     config.Version,
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@example.com",
			},
			"license": map[string]interface{}{
				"name": "MIT",
				"url":  "https://opensource.org/licenses/MIT",
			},
		},
		"servers": []map[string]interface{}{
			{
				"url":         "{protocol}://{host}",
				"description": "API server",
				"variables": map[string]interface{}{
					"protocol": map[string]interface{}{
						"enum":    []string{"http", "https"},
						"default": "https",
					},
					"host": map[string]interface{}{
						"default": "api.example.com",
					},
				},
			},
		},
		"paths":      buildPaths(doc.Operations),
		"components": buildComponents(),
		"tags":       buildTags(doc.Operations),
		"security": []map[string][]string{
			{
				"bearerAuth": {},
			},
		},
	}

	return writeJSON(config.OutputFile, openapi)
}

func buildPaths(operations []*spec.Operation) map[string]interface{} {
	paths := make(map[string]interface{})

	for _, op := range operations {
		if _, exists := paths[op.Path]; !exists {
			paths[op.Path] = make(map[string]interface{})
		}

		pathItem := paths[op.Path].(map[string]interface{})
		method := strings.ToLower(op.Method)

		operation := map[string]interface{}{
			"tags":        extractTags(op.Path),
			"summary":     op.Summary,
			"operationId": extractHandlerName(op.Summary),
			"parameters":  convertParameters(op.Parameters),
			"responses":   convertResponses(op.Responses),
		}

		if op.RequestBody != nil {
			operation["requestBody"] = convertRequestBody(op.RequestBody)
		}

		pathItem[method] = operation
	}

	return paths
}

func buildComponents() map[string]interface{} {
	return map[string]interface{}{
		"securitySchemes": map[string]interface{}{
			"bearerAuth": map[string]interface{}{
				"type":         "http",
				"scheme":       "bearer",
				"bearerFormat": "JWT",
				"description":  "JWT Authorization header using the Bearer scheme",
			},
		},
		"schemas": map[string]interface{}{
			"Error": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":    "integer",
						"format":  "int32",
						"example": 400,
					},
					"message": map[string]interface{}{
						"type":    "string",
						"example": "Bad Request",
					},
				},
			},
		},
	}
}

func buildTags(operations []*spec.Operation) []map[string]interface{} {
	tagSet := make(map[string]bool)
	for _, op := range operations {
		tags := extractTags(op.Path)
		for _, tag := range tags {
			tagSet[tag] = true
		}
	}

	tags := make([]map[string]interface{}, 0)
	for tag := range tagSet {
		tags = append(tags, map[string]interface{}{
			"name":        tag,
			"description": fmt.Sprintf("Operations about %s", tag),
		})
	}

	return tags
}

func validateConfig(config Config) error {
	// Definir valores padrão se não fornecidos
	if config.OutputFile == "" {
		config.OutputFile = "docs/openapi.json"
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}
	if config.Title == "" {
		config.Title = "API Documentation"
	}
	if config.Description == "" {
		config.Description = "API documentation generated by Gobiru"
	}
	return nil
}
