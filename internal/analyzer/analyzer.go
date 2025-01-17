package analyzer

import (
	"fmt"
	"os"

	"github.com/jeffemart/gobiru/internal/spec"
)

// Analyzer define a interface para análise de rotas
type Analyzer interface {
	Analyze() (*spec.Documentation, error)
}

// Config contém as configurações para análise
type Config struct {
	MainFile     string
	RouterFile   string
	HandlersFile string
}

// BaseAnalyzer contém a implementação comum para todos os analisadores
type BaseAnalyzer struct {
	config Config
}

// New cria um novo analisador baseado no framework
func New(framework string, mainFile, routerFile, handlersFile string) (Analyzer, error) {
	// Add validation for file paths
	if routerFile == "" || handlersFile == "" {
		return nil, fmt.Errorf("router and handlers files are required")
	}

	config := Config{
		MainFile:     mainFile,
		RouterFile:   routerFile,
		HandlersFile: handlersFile,
	}

	var analyzer Analyzer
	switch framework {
	case "gin":
		analyzer = NewGinAnalyzer(config)
	case "mux":
		analyzer = NewMuxAnalyzer(config)
	case "fiber":
		analyzer = NewFiberAnalyzer(config)
	default:
		return nil, fmt.Errorf("unsupported framework: %s", framework)
	}

	// Validate that files exist
	if err := validateFiles(config); err != nil {
		return nil, err
	}

	return analyzer, nil
}

// validateFiles checks if the required files exist
func validateFiles(config Config) error {
	files := []string{config.RouterFile, config.HandlersFile}
	if config.MainFile != "" {
		files = append(files, config.MainFile)
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", file)
		}
	}
	return nil
}
