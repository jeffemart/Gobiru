package gin

import (
	"github.com/jeffemart/Gobiru/app/analyzer"
	"github.com/jeffemart/Gobiru/app/models"
)

type GinAnalyzer struct {
	analyzer.BaseAnalyzer
}

func NewAnalyzer() *GinAnalyzer {
	return &GinAnalyzer{
		BaseAnalyzer: analyzer.BaseAnalyzer{
			FrameworkName: "gin",
		},
	}
}

func (ga *GinAnalyzer) GetTemplateMain() string {
	return `package main

import (
	"encoding/json"
	"fmt"
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
}`
}

func (ga *GinAnalyzer) GetDependencies() []string {
	return []string{
		"github.com/gin-gonic/gin v1.9.1",
		"github.com/bytedance/sonic v1.9.1",
		"github.com/go-playground/validator/v10 v10.14.0",
		"github.com/goccy/go-json v0.10.2",
		"github.com/json-iterator/go v1.1.12",
		"github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd",
		"github.com/modern-go/reflect2 v1.0.2",
		"github.com/pelletier/go-toml/v2 v2.0.8",
		"github.com/ugorji/go/codec v1.2.11",
		"golang.org/x/crypto v0.14.0",
		"golang.org/x/net v0.17.0",
		"golang.org/x/sys v0.15.0",
		"golang.org/x/text v0.13.0",
		"google.golang.org/protobuf v1.30.0",
		"gopkg.in/yaml.v3 v3.0.1",
	}
}

func (ga *GinAnalyzer) GetRoutes() []models.RouteInfo {
	return ga.Routes
}

func (ga *GinAnalyzer) GetFrameworkName() string {
	return ga.FrameworkName
}

func (ga *GinAnalyzer) AnalyzeFile(filePath string) ([]models.RouteInfo, error) {
	return analyzer.AnalyzeFile(ga, filePath)
}
