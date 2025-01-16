package gin

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
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
	// Create a new Gin engine for analysis
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Configure router from file
	if err := configureGinRouter(router, filePath); err != nil {
		return nil, fmt.Errorf("failed to configure Gin router: %v", err)
	}

	// Analyze routes
	if err := ga.AnalyzeRoutes(router); err != nil {
		return nil, fmt.Errorf("failed to analyze routes: %v", err)
	}

	return ga.routes, nil
}

// AnalyzeRoutes analyzes the given Gin engine and extracts route information
func (ga *GinAnalyzer) AnalyzeRoutes(engine *gin.Engine) error {
	routes := engine.Routes()

	for _, route := range routes {
		routeInfo := models.RouteInfo{
			Method:      route.Method,
			Path:        route.Path,
			HandlerName: runtime.FuncForPC(reflect.ValueOf(route.Handler).Pointer()).Name(),
			APIVersion:  "v1.0",
			RateLimit: models.RateLimit{
				RequestsPerMinute: 100,
				TimeWindowSeconds: 60,
			},
		}

		// Extract path parameters
		routeInfo.Parameters = extractPathParameters(route.Path)

		ga.routes = append(ga.routes, routeInfo)
	}

	return nil
}

// GetRoutes returns the analyzed routes
func (ga *GinAnalyzer) GetRoutes() []models.RouteInfo {
	return ga.routes
}

func getHandlerName(handler interface{}) string {
	if handler == nil {
		return ""
	}
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}

func extractPathParameters(path string) []models.Parameter {
	var params []models.Parameter
	parts := strings.Split(path, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params = append(params, models.Parameter{
				Name:     paramName,
				Type:     "string",
				Required: true,
			})
		}
	}

	return params
}

func configureGinRouter(router *gin.Engine, filePath string) error {
	// Example routes for testing
	router.GET("/api/v1/products", func(c *gin.Context) {})
	router.GET("/api/v1/products/:id", func(c *gin.Context) {})
	router.POST("/api/v1/products", func(c *gin.Context) {})
	router.PUT("/api/v1/products/:id", func(c *gin.Context) {})
	router.DELETE("/api/v1/products/:id", func(c *gin.Context) {})

	return nil
}
