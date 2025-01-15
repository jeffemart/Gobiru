package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	analyzer := gobiru.NewRouteAnalyzer()

	// TODO: Implement logic to load and parse the routes file
	// This will require reading the user's routes.go file and getting the router instance

	// Export route information as JSON
	err := analyzer.ExportJSON(*outputFile)
	if err != nil {
		log.Fatalf("Failed to export routes: %v", err)
	}

	fmt.Printf("Route documentation generated successfully at: %s\n", *outputFile)

	// Export OpenAPI specification if requested
	if *openAPIFile != "" {
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
