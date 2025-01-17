package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

// ImportTracker rastreia os arquivos e seus imports
type ImportTracker struct {
	processedFiles map[string]bool
	routeFiles     []string
	handlerFiles   []string
	baseDir        string
}

func NewImportTracker(mainFile string) *ImportTracker {
	// Encontrar o diretório base do projeto
	baseDir := filepath.Dir(mainFile)
	for !strings.HasSuffix(baseDir, "examples") && baseDir != "/" && baseDir != "." {
		baseDir = filepath.Dir(baseDir)
	}

	return &ImportTracker{
		processedFiles: make(map[string]bool),
		routeFiles:     make([]string, 0),
		handlerFiles:   make([]string, 0),
		baseDir:        baseDir,
	}
}

func (t *ImportTracker) TrackImports(filePath string) error {
	if t.processedFiles[filePath] {
		return nil
	}
	t.processedFiles[filePath] = true

	// Verificar se é um diretório
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %v", filePath, err)
	}

	// Se for um diretório, processar todos os arquivos .go dentro dele
	if fileInfo.IsDir() {
		files, err := os.ReadDir(filePath)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %v", filePath, err)
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
				fullPath := filepath.Join(filePath, file.Name())
				if err := t.processFile(fullPath); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return t.processFile(filePath)
}

func (t *ImportTracker) processFile(filePath string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	// Verificar se é um arquivo de rotas ou handlers
	if strings.Contains(filePath, "routes") {
		t.routeFiles = append(t.routeFiles, filePath)
	} else if strings.Contains(filePath, "handlers") {
		t.handlerFiles = append(t.handlerFiles, filePath)
	}

	// Processar imports
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")

		// Se for um import local do projeto
		if strings.Contains(importPath, "examples/") {
			// Extrair o caminho relativo após "examples/"
			parts := strings.Split(importPath, "examples/")
			if len(parts) > 1 {
				localPath := filepath.Join(t.baseDir, parts[1])
				if err := t.TrackImports(localPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Analyzer define a interface para análise de rotas
type Analyzer interface {
	Analyze() (*spec.Documentation, error)
}

// Config contém as configurações para análise
type Config struct {
	MainFile     string
	RouterFiles  []string
	HandlerFiles []string
}

// BaseAnalyzer contém a implementação comum para todos os analisadores
type BaseAnalyzer struct {
	config Config
}

// New cria um novo analisador baseado no framework
func New(framework string, config Config) (Analyzer, error) {
	if config.MainFile == "" {
		return nil, fmt.Errorf("main file is required")
	}

	// Rastrear imports a partir do main.go
	tracker := NewImportTracker(config.MainFile)
	if err := tracker.TrackImports(config.MainFile); err != nil {
		return nil, fmt.Errorf("failed to track imports: %v", err)
	}

	// Atualizar config com os arquivos encontrados
	config.RouterFiles = tracker.routeFiles
	config.HandlerFiles = tracker.handlerFiles

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

	return analyzer, nil
}
