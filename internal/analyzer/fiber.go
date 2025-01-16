package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/models"
)

type FiberAnalyzer struct {
	config Config
}

func NewFiberAnalyzer(config Config) *FiberAnalyzer {
	return &FiberAnalyzer{config: config}
}

func (a *FiberAnalyzer) Analyze() ([]models.RouteInfo, error) {
	// Analisar o arquivo de rotas
	fset := token.NewFileSet()
	routerFile, err := parser.ParseFile(fset, a.config.RouterFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse router file: %v", err)
	}

	var routes []models.RouteInfo

	// Encontrar as definições de rotas
	ast.Inspect(routerFile, func(n ast.Node) bool {
		// Procurar chamadas de métodos HTTP
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Verificar se é uma chamada de método HTTP
		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			method := strings.ToUpper(sel.Sel.Name)
			if isHTTPMethod(method) {
				route := models.RouteInfo{
					Method: method,
				}

				// Extrair path
				if len(callExpr.Args) > 0 {
					if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
						route.Path = strings.Trim(lit.Value, "\"")
					}
				}

				// Extrair handler
				if len(callExpr.Args) > 1 {
					if ident, ok := callExpr.Args[1].(*ast.Ident); ok {
						route.HandlerName = ident.Name
					}
				}

				// Extrair parâmetros da rota
				if route.Path != "" {
					route.Parameters = extractFiberParameters(route.Path)
				}

				routes = append(routes, route)
			}
		}

		return true
	})

	return routes, nil
}

func extractFiberParameters(path string) []models.Parameter {
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
