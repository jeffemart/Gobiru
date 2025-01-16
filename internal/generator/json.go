package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeffemart/gobiru/internal/models"
)

// GenerateJSON gera a documentação em formato JSON
func GenerateJSON(routes []models.RouteInfo, outputFile string) error {
	// Criar diretório se não existir
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Gerar JSON
	data, err := json.MarshalIndent(routes, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %v", err)
	}

	// Salvar arquivo
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
