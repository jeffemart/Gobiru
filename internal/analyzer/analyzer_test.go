package analyzer

import (
	"testing"
	"os"
	"path/filepath"
)

func TestFindMainFile(t *testing.T) {
	// Crie um diretório temporário para o teste
	tempDir := t.TempDir()
	mainFilePath := filepath.Join(tempDir, "main.go")

	// Crie um arquivo main.go para o teste
	if err := os.WriteFile(mainFilePath, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Teste a função FindMainFile
	foundFile, err := FindMainFile(tempDir)
	if err != nil {
		t.Fatalf("Expected to find main.go, got error: %v", err)
	}

	if foundFile != mainFilePath {
		t.Errorf("Expected %s, got %s", mainFilePath, foundFile)
	}
}

func TestAnalyzeRoutes(t *testing.T) {
	// Aqui você pode adicionar um teste para a função Analyze
	// Isso pode incluir a criação de arquivos de rotas e handlers
	// e verificar se as operações são analisadas corretamente.
} 