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
	moduleName     string
}

func NewImportTracker(mainFile string) *ImportTracker {
	baseDir := filepath.Dir(mainFile)
	return &ImportTracker{
		processedFiles: make(map[string]bool),
		routeFiles:     make([]string, 0),
		handlerFiles:   make([]string, 0),
		baseDir:        baseDir,
		moduleName:     findModuleName(baseDir),
	}
}

// findModuleName tenta encontrar o nome do módulo no go.mod
func findModuleName(dir string) string {
	for {
		modFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			content, err := os.ReadFile(modFile)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "module ") {
						return strings.TrimSpace(strings.TrimPrefix(line, "module "))
					}
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func (t *ImportTracker) TrackImports(filePath string) error {
	if t.processedFiles[filePath] {
		return nil
	}
	t.processedFiles[filePath] = true

	// Se o arquivo é um diretório, processar todos os arquivos .go dentro dele
	fileInfo, err := os.Stat(filePath)
	if err == nil && fileInfo.IsDir() {
		files, err := os.ReadDir(filePath)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %v", filePath, err)
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
				fullPath := filepath.Join(filePath, file.Name())
				if err := t.TrackImports(fullPath); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Processar arquivo individual
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil // Ignorar arquivos que não podem ser parseados
	}

	// Verificar diretórios especiais
	dirName := filepath.Base(filepath.Dir(filePath))
	if dirName == "routes" {
		fmt.Printf("Found route file: %s\n", filePath)
		t.routeFiles = append(t.routeFiles, filePath)
	} else if dirName == "handlers" {
		fmt.Printf("Found handler file: %s\n", filePath)
		t.handlerFiles = append(t.handlerFiles, filePath)
	}

	// Processar imports
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")

		// Se for um import relativo ou local
		if strings.HasPrefix(importPath, ".") {
			dir := filepath.Dir(filePath)
			localPath := filepath.Join(dir, importPath)
			if err := t.TrackImports(localPath); err != nil {
				fmt.Printf("Warning: failed to process import %s: %v\n", importPath, err)
			}
			continue
		}

		// Tentar resolver no diretório do projeto
		if t.moduleName != "" {
			// Se o import começa com o nome do módulo, resolver relativamente ao diretório base
			if strings.HasPrefix(importPath, t.moduleName) {
				relativePath := strings.TrimPrefix(importPath, t.moduleName)
				localPath := filepath.Join(t.baseDir, relativePath)
				if err := t.TrackImports(localPath); err != nil {
					fmt.Printf("Warning: failed to process module import %s: %v\n", importPath, err)
				}
				continue
			}
		}

		// Verificar diretórios especiais no diretório atual
		routesPath := filepath.Join(t.baseDir, "routes")
		handlersPath := filepath.Join(t.baseDir, "handlers")

		if _, err := os.Stat(routesPath); err == nil {
			if err := t.TrackImports(routesPath); err != nil {
				fmt.Printf("Warning: failed to process routes directory: %v\n", err)
			}
		}
		if _, err := os.Stat(handlersPath); err == nil {
			if err := t.TrackImports(handlersPath); err != nil {
				fmt.Printf("Warning: failed to process handlers directory: %v\n", err)
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
	// case "fiber":
	// 	analyzer = NewFiberAnalyzer(config)
	default:
		return nil, fmt.Errorf("unsupported framework: %s", framework)
	}

	return analyzer, nil
}
