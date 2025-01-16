package generator

import (
	"github.com/jeffemart/gobiru/internal/spec"
)

type OpenAPIGenerator struct{}

func NewOpenAPIGenerator() *OpenAPIGenerator {
	return &OpenAPIGenerator{}
}

func (g *OpenAPIGenerator) Generate(doc *spec.Documentation, config Config) error {
	openapi := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       config.Title,
			"description": config.Description,
			"version":     config.Version,
		},
		"paths": convertPaths(doc.Operations),
	}

	return writeJSON(config.OutputFile, openapi)
}
