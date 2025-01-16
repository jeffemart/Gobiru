package gin

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

// GinAnalyzer is responsible for analyzing Gin routes
type GinAnalyzer struct {
	routes []models.RouteInfo
}

// NewAnalyzer creates a new GinAnalyzer instance
func NewAnalyzer() *GinAnalyzer {
	return &GinAnalyzer{
		routes: make([]models.RouteInfo, 0),
	}
}

// AnalyzeFile analyzes a file containing Gin routes
func (ga *GinAnalyzer) AnalyzeFile(filePath string) ([]models.RouteInfo, error) {
	// Copiar o arquivo de rotas para o diretório temporário
	routesContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read routes file: %v", err)
	}

	// Extrair o conteúdo do arquivo routes.go, removendo package e imports
	content := string(routesContent)
	content = strings.Replace(content, "package main", "", 1)
	content = strings.Replace(content, `import (
	"net/http"

	"github.com/gin-gonic/gin"
)`, "", 1)

	// Criar um diretório temporário para compilar
	tmpDir, err := ioutil.TempDir("", "gobiru")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo temporário main
	tmpMain := filepath.Join(tmpDir, "main.go")
	mainContent := fmt.Sprintf(`
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
)

type RouteInfo struct {
	Method      string   %[1]sjson:"method"%[1]s
	Path        string   %[1]sjson:"path"%[1]s
	HandlerName string   %[1]sjson:"handler_name"%[1]s
	Parameters  []Parameter %[1]sjson:"parameters"%[1]s
}

type Parameter struct {
	Name     string %[1]sjson:"name"%[1]s
	Type     string %[1]sjson:"type"%[1]s
	Required bool   %[1]sjson:"required"%[1]s
}

%[2]s

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := SetupRouter()
	routes := router.Routes()
	var routeInfos []RouteInfo

	for _, route := range routes {
		info := RouteInfo{
			Method:      route.Method,
			Path:       route.Path,
			HandlerName: route.Handler,
		}

		// Extract path parameters
		parts := strings.Split(route.Path, "/")
		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				paramName := strings.TrimPrefix(part, ":")
				info.Parameters = append(info.Parameters, Parameter{
					Name:     paramName,
					Type:     "string",
					Required: true,
				})
			}
		}

		routeInfos = append(routeInfos, info)
	}

	data, err := json.Marshal(routeInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
`, "`", content)

	if err := ioutil.WriteFile(tmpMain, []byte(mainContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp main: %v", err)
	}

	// Criar go.mod
	modContent := `module temp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`
	if err := ioutil.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write go.mod: %v", err)
	}

	// Executar go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to tidy go module: %v", err)
	}

	// Executar o programa temporário
	cmd = exec.Command("go", "run", ".")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to run temporary program: %v\nStderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run temporary program: %v", err)
	}

	// Analisar a saída
	if err := json.Unmarshal(output, &ga.routes); err != nil {
		return nil, fmt.Errorf("failed to parse routes output: %v", err)
	}

	return ga.routes, nil
}

// GetRoutes returns the analyzed routes
func (ga *GinAnalyzer) GetRoutes() []models.RouteInfo {
	return ga.routes
}
