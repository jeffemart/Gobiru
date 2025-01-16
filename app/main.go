package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gobiru/app/models"
	"github.com/jeffemart/Gobiru/app/openapi"
)

// RouteAnalyzer is the main struct for analyzing routes
type RouteAnalyzer struct {
	routes []models.RouteInfo
}

// NewRouteAnalyzer creates a new RouteAnalyzer instance
func NewRouteAnalyzer() *RouteAnalyzer {
	return &RouteAnalyzer{
		routes: make([]models.RouteInfo, 0),
	}
}

// AnalyzeRoutes analyzes the given router and extracts route information
func (ra *RouteAnalyzer) AnalyzeRoutes(router *mux.Router) error {
	return router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		routeInfo := models.RouteInfo{}

		// Get path
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return fmt.Errorf("failed to get path template: %v", err)
		}
		routeInfo.Path = pathTemplate

		// Get HTTP method
		methods, err := route.GetMethods()
		if err == nil && len(methods) > 0 {
			routeInfo.Method = methods[0]
		}

		// Get handler name
		if handler := route.GetHandler(); handler != nil {
			routeInfo.HandlerName = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		}

		// Extract path parameters
		pathVars := extractPathParameters(pathTemplate)
		for _, param := range pathVars {
			routeInfo.Parameters = append(routeInfo.Parameters, models.Parameter{
				Name:     param,
				Type:     "string",
				Required: true,
			})
		}

		// Set default values
		routeInfo.APIVersion = "v1.0"
		routeInfo.Deprecated = false
		routeInfo.RateLimit = models.RateLimit{
			RequestsPerMinute: 100,
			TimeWindowSeconds: 60,
		}

		ra.routes = append(ra.routes, routeInfo)
		return nil
	})
}

// GetRoutes returns the analyzed routes
func (ra *RouteAnalyzer) GetRoutes() []models.RouteInfo {
	return ra.routes
}

// extractPathParameters extracts path parameters from the route template
func extractPathParameters(pathTemplate string) []string {
	var params []string
	parts := strings.Split(pathTemplate, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			param := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
			params = append(params, param)
		}
	}
	return params
}

// ExportJSON exports the route information as JSON
func (ra *RouteAnalyzer) ExportJSON(filepath string) error {
	data, err := json.MarshalIndent(ra.routes, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %v", err)
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// ExportOpenAPI exports the route information as OpenAPI specification
func (ra *RouteAnalyzer) ExportOpenAPI(filepath string, info openapi.Info) error {
	spec, err := openapi.ConvertToOpenAPI(ra.routes, info)
	if err != nil {
		return fmt.Errorf("failed to convert to OpenAPI: %v", err)
	}

	data, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI spec: %v", err)
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %v", err)
	}

	return nil
}

// AnalyzeFile analyzes a file containing Mux routes
func (ra *RouteAnalyzer) AnalyzeFile(filePath string) ([]models.RouteInfo, error) {
	// Copiar o arquivo de rotas para o diretório temporário
	routesContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read routes file: %v", err)
	}

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
	"os"
	"strings"
	"github.com/gorilla/mux"
	"runtime"
	"reflect"
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

func main() {
	router := SetupRouter()
	var routes []RouteInfo

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		info := RouteInfo{}
		
		if path, err := route.GetPathTemplate(); err == nil {
			info.Path = path
			
			// Extract path parameters
			for _, part := range splitPath(path) {
				if isParameter(part) {
					info.Parameters = append(info.Parameters, Parameter{
						Name:     trimParameter(part),
						Type:     "string",
						Required: true,
					})
				}
			}
		}
		
		if methods, err := route.GetMethods(); err == nil && len(methods) > 0 {
			info.Method = methods[0]
		}

		if handler := route.GetHandler(); handler != nil {
			info.HandlerName = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		}

		routes = append(routes, info)
		return nil
	})
	
	data, err := json.Marshal(routes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func isParameter(part string) bool {
	return strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")
}

func trimParameter(part string) string {
	return strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
}
`, "`")

	if err := ioutil.WriteFile(tmpMain, []byte(mainContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp main: %v", err)
	}

	tmpRoutes := filepath.Join(tmpDir, "routes.go")
	if err := ioutil.WriteFile(tmpRoutes, routesContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp routes: %v", err)
	}

	// Criar go.mod
	modContent := `module temp

go 1.21

require (
	github.com/gorilla/mux v1.8.1
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
	if err := json.Unmarshal(output, &ra.routes); err != nil {
		return nil, fmt.Errorf("failed to parse routes output: %v", err)
	}

	return ra.routes, nil
}
