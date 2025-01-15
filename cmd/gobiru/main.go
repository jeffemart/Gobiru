package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	if *framework == "" {
		log.Fatal("Framework is required. Use -framework flag with 'mux' or 'gin'")
	}

	if flag.NArg() < 1 {
		log.Fatal("Source file is required. Usage: gobiru [options] <source-file>")
	}

	sourceFile := flag.Arg(0)

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
  # Generate docs for a Gin application
  gobiru -framework gin -output docs/routes.json main.go

  # Generate docs and start server
  gobiru -framework mux -serve main.go

  # Full example with all options
  gobiru -framework gin \
         -output docs/routes.json \
         -openapi docs/openapi.json \
         -title "My API" \
         -description "My API description" \
         -version "1.0.0" \
         -serve \
         -port 8081 \
         main.go`)
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
