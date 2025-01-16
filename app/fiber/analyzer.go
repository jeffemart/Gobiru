package fiber

import (
	"github.com/jeffemart/Gobiru/app/analyzer"
	"github.com/jeffemart/Gobiru/app/models"
)

type FiberAnalyzer struct {
	analyzer.BaseAnalyzer
}

func NewAnalyzer() *FiberAnalyzer {
	return &FiberAnalyzer{
		BaseAnalyzer: analyzer.BaseAnalyzer{
			FrameworkName: "fiber",
		},
	}
}

func (fa *FiberAnalyzer) GetTemplateMain() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"reflect"
	"github.com/gofiber/fiber/v2"
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

// Conteúdo do arquivo routes.go
%[2]s

func main() {
	app := SetupRouter()
	
	var routeInfos []RouteInfo
	stack := app.Stack()
	
	for _, routes := range stack {
		for _, route := range routes {
			info := RouteInfo{
				Method:      route.Method,
				Path:       route.Path,
			}

			// Get handler name from the first handler in the chain
			if len(route.Handlers) > 0 {
				info.HandlerName = getFunctionName(route.Handlers[0])
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
	}

	data, err := json.Marshal(routeInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func getFunctionName(i interface{}) string {
	return reflect.TypeOf(i).String()
}`
}

func (fa *FiberAnalyzer) GetDependencies() []string {
	return []string{
		"github.com/gofiber/fiber/v2 v2.52.0",
	}
}

func (fa *FiberAnalyzer) AnalyzeFile(filePath string) ([]models.RouteInfo, error) {
	return analyzer.AnalyzeFile(fa, filePath)
}

func (fa *FiberAnalyzer) GetRoutes() []models.RouteInfo {
	return fa.Routes
}

func (fa *FiberAnalyzer) GetFrameworkName() string {
	return fa.FrameworkName
}
