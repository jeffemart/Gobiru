package analyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jeffemart/Gobiru/app/models"
)

// RouteAnalyzer define a interface comum para todos os analisadores
type RouteAnalyzer interface {
	GetRoutes() []models.RouteInfo
	GetFrameworkName() string
	GetTemplateMain() string
	GetDependencies() []string
	AnalyzeFile(string) ([]models.RouteInfo, error)
}

// BaseAnalyzer contém a implementação comum para todos os analisadores
type BaseAnalyzer struct {
	Routes        []models.RouteInfo
	FrameworkName string
}

// AnalyzeFile é o método comum para analisar arquivos de rotas
func AnalyzeFile(analyzer RouteAnalyzer, filePath string) ([]models.RouteInfo, error) {
	// Read the routes file
	routesContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read routes file: %v", err)
	}

	// Limpar o conteúdo do arquivo routes.go
	content := string(routesContent)
	content = strings.Replace(content, "package main", "", 1)
	content = strings.TrimSpace(content)

	// Se houver bloco de import, remova-o
	if idx := strings.Index(content, "import ("); idx >= 0 {
		if endIdx := strings.Index(content[idx:], ")"); endIdx >= 0 {
			content = content[idx+endIdx+1:]
		}
	}
	content = strings.TrimSpace(content)

	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "gobiru")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create temporary main file
	tmpMain := filepath.Join(tmpDir, "main.go")
	mainContent := fmt.Sprintf(analyzer.GetTemplateMain(), "`", content)

	if err := ioutil.WriteFile(tmpMain, []byte(mainContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp main: %v", err)
	}

	// Create go.mod
	modContent := fmt.Sprintf(`module temp

go 1.21

require (
	%s
)
`, strings.Join(analyzer.GetDependencies(), "\n\t"))

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write go.mod: %v", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to tidy go module: %v", err)
	}

	// Run the temporary program
	cmd = exec.Command("go", "run", ".")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to run temporary program: %v\nStderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run temporary program: %v", err)
	}

	var routes []models.RouteInfo
	if err := json.Unmarshal(output, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse routes output: %v", err)
	}

	return routes, nil
}
