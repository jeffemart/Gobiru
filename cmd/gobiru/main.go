package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jeffemart/Gobiru/app"
	"github.com/jeffemart/Gobiru/app/gin"
	"github.com/jeffemart/Gobiru/app/models"
	"github.com/jeffemart/Gobiru/app/openapi"
	"github.com/jeffemart/Gobiru/app/server"
)

const version = "1.0.0"

func main() {
	// Command flags
	var (
		outputFile  = flag.String("output", "docs/routes.json", "Output file path for route documentation")
		openAPIFile = flag.String("openapi", "docs/openapi.json", "Output file path for OpenAPI specification")
		apiTitle    = flag.String("title", "API Documentation", "API title for OpenAPI spec")
		apiDesc     = flag.String("description", "", "API description for OpenAPI spec")
		apiVersion  = flag.String("api-version", "1.0.0", "API version for OpenAPI spec")
		framework   = flag.String("framework", "", "Framework to analyze (mux or gin)")
		routerFile  = flag.String("router", "", "Path to file containing router definition (defaults to routes.go)")
		serve       = flag.Bool("serve", false, "Start documentation server after generation")
		port        = flag.Int("port", 8081, "Port for documentation server")
		help        = flag.Bool("help", false, "Show help message")
		showVersion = flag.Bool("version", false, "Show version")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *showVersion {
		fmt.Printf("Gobiru version %s\n", version)
		return
	}

	if *framework != "gin" && *framework != "mux" {
		log.Fatal("Framework must be either 'gin' or 'mux'")
	}

	// Determinar o arquivo fonte
	sourceFile := *routerFile
	if sourceFile == "" {
		// Se não foi especificado, procurar por routes.go
		sourceFile = findRouterFile()
	}

	if sourceFile == "" {
		log.Fatal("Router file not found. Use -router flag or ensure routes.go exists in current directory")
	}

	// Create output directories
	ensureDir(*outputFile)
	if *openAPIFile != "" {
		ensureDir(*openAPIFile)
	}

	// Generate documentation
	routes, err := analyzeRoutes(*framework, sourceFile)
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Export documentation
	if err := exportDocs(routes, *outputFile, *openAPIFile, openapi.Info{
		Title:       *apiTitle,
		Description: *apiDesc,
		Version:     *apiVersion,
	}); err != nil {
		log.Fatalf("Failed to export documentation: %v", err)
	}

	fmt.Printf("Documentation generated successfully!\n")
	fmt.Printf("Routes JSON: %s\n", *outputFile)
	if *openAPIFile != "" {
		fmt.Printf("OpenAPI spec: %s\n", *openAPIFile)
	}

	// Copiar index.html se não existir
	indexPath := filepath.Join(filepath.Dir(*outputFile), "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI - Gobiru</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.11.0/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.11.0/favicon-16x16.png" sizes="16x16" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js" crossorigin></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "./openapi.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                plugins: [SwaggerUIBundle.plugins.DownloadUrl],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
		if err := ioutil.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
			log.Printf("Warning: Failed to create index.html: %v", err)
		}
	}

	// Start documentation server if requested
	if *serve {
		fmt.Printf("\nStarting documentation server...\n")
		server.Serve(*port, filepath.Dir(*outputFile))
	}
}

func showHelp() {
	fmt.Println(`Gobiru - API Documentation Generator

Usage:
  gobiru [options] <source-file>

Options:`)
	flag.PrintDefaults()
	fmt.Println(`
Examples:
  # Generate docs for a Gin application with specific router file
  gobiru -framework gin -router api/routes.go -output docs/routes.json

  # Generate docs using default routes.go
  gobiru -framework mux -output docs/routes.json

  # Full example with all options
  gobiru -framework gin \
         -router internal/routes/router.go \
         -output docs/routes.json \
         -openapi docs/openapi.json \
         -title "My API" \
         -description "My API description" \
         -version "1.0.0" \
         -serve \
         -port 8081`)
}

func ensureDir(filePath string) {
	dir := filepath.Dir(filePath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
}

func analyzeRoutes(framework, sourceFile string) ([]models.RouteInfo, error) {
	switch framework {
	case "gin":
		analyzer := gin.NewAnalyzer()
		return analyzer.AnalyzeFile(sourceFile)
	case "mux":
		analyzer := app.NewRouteAnalyzer()
		return analyzer.AnalyzeFile(sourceFile)
	default:
		return nil, fmt.Errorf("unsupported framework: %s", framework)
	}
}

func exportDocs(routes []models.RouteInfo, jsonFile, openAPIFile string, info openapi.Info) error {
	// Export routes JSON
	if err := exportJSON(routes, jsonFile); err != nil {
		return fmt.Errorf("failed to export routes JSON: %v", err)
	}

	// Export OpenAPI spec if requested
	if openAPIFile != "" {
		spec, err := openapi.ConvertToOpenAPI(routes, info)
		if err != nil {
			return fmt.Errorf("failed to convert to OpenAPI: %v", err)
		}

		if err := openapi.ExportOpenAPI(spec, openAPIFile); err != nil {
			return fmt.Errorf("failed to export OpenAPI spec: %v", err)
		}
	}

	return nil
}

func exportJSON(routes []models.RouteInfo, filepath string) error {
	data, err := json.MarshalIndent(routes, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %v", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// Função auxiliar para procurar o arquivo routes.go
func findRouterFile() string {
	// Primeiro procura no diretório atual
	if _, err := os.Stat("routes.go"); err == nil {
		return "routes.go"
	}

	// Procura em diretórios comuns
	commonPaths := []string{
		"internal/routes/routes.go",
		"pkg/routes/routes.go",
		"api/routes/routes.go",
		"src/routes/routes.go",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
