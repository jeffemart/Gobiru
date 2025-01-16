package mux

import (
	"github.com/jeffemart/Gobiru/app/analyzer"
	"github.com/jeffemart/Gobiru/app/models"
)

type MuxAnalyzer struct {
	analyzer.BaseAnalyzer
}

func NewAnalyzer() *MuxAnalyzer {
	return &MuxAnalyzer{
		BaseAnalyzer: analyzer.BaseAnalyzer{
			FrameworkName: "mux",
		},
	}
}

func (ma *MuxAnalyzer) GetTemplateMain() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"runtime"
	"reflect"
	"net/http"
	"github.com/gorilla/mux"
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
	router := SetupRouter()
	var routeInfos []RouteInfo

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		info := RouteInfo{}
		
		if path, err := route.GetPathTemplate(); err == nil {
			info.Path = path
			
			// Extract path parameters
			for _, part := range strings.Split(path, "/") {
				if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
					paramName := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
					info.Parameters = append(info.Parameters, Parameter{
						Name:     paramName,
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

		routeInfos = append(routeInfos, info)
		return nil
	})
	
	data, err := json.Marshal(routeInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}`
}

func (ma *MuxAnalyzer) GetDependencies() []string {
	return []string{
		"github.com/gorilla/mux v1.8.1",
	}
}

func (ma *MuxAnalyzer) AnalyzeFile(filePath string) ([]models.RouteInfo, error) {
	return analyzer.AnalyzeFile(ma, filePath)
}

func (ma *MuxAnalyzer) GetRoutes() []models.RouteInfo {
	return ma.Routes
}

func (ma *MuxAnalyzer) GetFrameworkName() string {
	return ma.FrameworkName
}
