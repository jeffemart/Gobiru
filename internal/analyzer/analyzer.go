package analyzer

import (
	"fmt"

	"github.com/jeffemart/gobiru/internal/models"
)

// Analyzer define a interface para análise de rotas
type Analyzer interface {
	Analyze() ([]models.RouteInfo, error)
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
	routes []models.RouteInfo
}

// New cria um novo analisador baseado no framework
func New(framework string, mainFile, routerFile, handlersFile string) (Analyzer, error) {
	config := Config{
		MainFile:     mainFile,
		RouterFile:   routerFile,
		HandlersFile: handlersFile,
	}

	switch framework {
	case "gin":
		return NewGinAnalyzer(config), nil
	case "mux":
		return NewMuxAnalyzer(config), nil
	case "fiber":
		return NewFiberAnalyzer(config), nil
	default:
		return nil, fmt.Errorf("unsupported framework: %s", framework)
	}
}
