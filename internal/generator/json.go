package generator

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jeffemart/gobiru/internal/spec"
)

type JSONGenerator struct{}

func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{}
}

func (g *JSONGenerator) Generate(doc *spec.Documentation, config Config) error {
	operations := make([]map[string]interface{}, 0)

	for _, op := range doc.Operations {
		operation := map[string]interface{}{
			"path":        op.Path,
			"method":      op.Method,
			"summary":     op.Summary,
			"parameters":  convertParameters(op.Parameters),
			"requestBody": convertRequestBody(op.RequestBody),
			"responses":   convertResponses(op.Responses),
		}
		operations = append(operations, operation)
	}

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
