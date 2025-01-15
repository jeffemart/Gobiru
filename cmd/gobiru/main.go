package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	gobiru "github.com/jeffemart/Gobiru/app"
	"github.com/jeffemart/Gobiru/app/openapi"
)

func main() {
	outputFile := flag.String("output", "routes.json", "Output file path for the route documentation")
	openAPIFile := flag.String("openapi", "", "Output file path for OpenAPI specification")
	apiTitle := flag.String("title", "API Documentation", "API title for OpenAPI spec")
	apiDesc := flag.String("description", "", "API description for OpenAPI spec")
	apiVersion := flag.String("version", "1.0.0", "API version for OpenAPI spec")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: gobiru [options] <path-to-routes-file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	routesFile := flag.Arg(0)
	router, err := parseRouterFromFile(routesFile)
	if err != nil {
		log.Fatalf("Failed to parse router from file: %v", err)
	}

	analyzer := gobiru.NewRouteAnalyzer()
	err = analyzer.AnalyzeRoutes(router)
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Create output directory if it doesn't exist
	if dir := filepath.Dir(*outputFile); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Export route information as JSON
	err = analyzer.ExportJSON(*outputFile)
	if err != nil {
		log.Fatalf("Failed to export routes: %v", err)
	}

	fmt.Printf("Route documentation generated successfully at: %s\n", *outputFile)

	// Export OpenAPI specification if requested
	if *openAPIFile != "" {
		if dir := filepath.Dir(*openAPIFile); dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Failed to create OpenAPI output directory: %v", err)
			}
		}

		info := openapi.Info{
			Title:       *apiTitle,
			Description: *apiDesc,
			Version:     *apiVersion,
		}

		err = analyzer.ExportOpenAPI(*openAPIFile, info)
		if err != nil {
			log.Fatalf("Failed to export OpenAPI specification: %v", err)
		}
		fmt.Printf("OpenAPI specification generated successfully at: %s\n", *openAPIFile)
	}
}

func parseRouterFromFile(filePath string) (*mux.Router, error) {
	router := mux.NewRouter()

	// Rotas de exemplo para teste
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")

	return router, nil
}
