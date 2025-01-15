package gobiru

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
