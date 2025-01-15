package gin

import (
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

// NewGinAnalyzer creates a new GinAnalyzer instance
func NewGinAnalyzer() *GinAnalyzer {
	return &GinAnalyzer{
		routes: make([]models.RouteInfo, 0),
	}
}

// AnalyzeRoutes analyzes the given Gin engine and extracts route information
func (ga *GinAnalyzer) AnalyzeRoutes(engine *gin.Engine) error {
	routes := engine.Routes()

	for _, route := range routes {
		routeInfo := models.RouteInfo{
			Method:      route.Method,
			Path:        route.Path,
			HandlerName: getHandlerName(route.Handler),
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
