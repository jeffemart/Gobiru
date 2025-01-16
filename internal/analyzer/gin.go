package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type GinAnalyzer struct {
	BaseAnalyzer
}

func NewGinAnalyzer(config Config) *GinAnalyzer {
	return &GinAnalyzer{
		BaseAnalyzer: BaseAnalyzer{
			config: config,
		},
	}
}

func (a *GinAnalyzer) Analyze() (*spec.Documentation, error) {
	fset := token.NewFileSet()

	routerFile, err := parser.ParseFile(fset, a.config.RouterFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse router file: %v", err)
	}

	handlersFile, err := parser.ParseFile(fset, a.config.HandlersFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse handlers file: %v", err)
	}

	var routes []struct {
		Path        string
		Method      string
		HandlerName string
	}

	// Função para construir o caminho completo da rota
	var currentPath []string

	ast.Inspect(routerFile, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				switch sel.Sel.Name {
				case "Group":
					if len(node.Args) > 0 {
						if lit, ok := node.Args[0].(*ast.BasicLit); ok {
							prefix := strings.Trim(lit.Value, "\"")
							currentPath = append(currentPath, prefix)
						}
					}
				case "POST", "GET", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
					if len(node.Args) >= 2 {
						route := struct {
							Path        string
							Method      string
							HandlerName string
						}{
							Method: strings.ToUpper(sel.Sel.Name),
						}

						// Construir caminho completo
						if lit, ok := node.Args[0].(*ast.BasicLit); ok {
							subPath := strings.Trim(lit.Value, "\"")
							fullPath := strings.Join(currentPath, "") + subPath
							route.Path = strings.TrimRight(fullPath, "/")
						}

						// Extrair handler
						if len(node.Args) > 1 {
							if ident, ok := node.Args[1].(*ast.Ident); ok {
								route.HandlerName = ident.Name
							} else if sel, ok := node.Args[1].(*ast.SelectorExpr); ok {
								route.HandlerName = sel.Sel.Name
							}
						}

						if route.Path != "" && route.HandlerName != "" {
							routes = append(routes, route)
						}
					}
				}
			}
		}
		return true
	})

	doc := &spec.Documentation{
		Operations: make([]*spec.Operation, 0),
	}

	for _, route := range routes {
		operation := &spec.Operation{
			Path:       route.Path,
			Method:     route.Method,
			Parameters: extractGinParameters(route.Path),
			Responses:  make(map[string]*spec.Response),
		}

		if handlerFunc := findFunction(handlersFile, route.HandlerName); handlerFunc != nil {
			operation.Summary = extractSummaryFromComments(handlerFunc)
			operation.RequestBody = extractRequestBody(handlerFunc, a.config.HandlersFile)
			operation.Responses = extractResponses(handlerFunc, a.config.HandlersFile)
		}

		doc.Operations = append(doc.Operations, operation)
	}

	return doc, nil
}

func extractGinParameters(path string) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			paramName := strings.TrimPrefix(segment, ":")
			params = append(params, &spec.Parameter{
				Name:        paramName,
				In:          "path",
				Required:    true,
				Description: fmt.Sprintf("Path parameter: %s", paramName),
				Schema: &spec.Schema{
					Type: "string",
				},
			})
		} else if strings.HasPrefix(segment, "*") {
			paramName := strings.TrimPrefix(segment, "*")
			params = append(params, &spec.Parameter{
				Name:        paramName,
				In:          "path",
				Required:    true,
				Description: fmt.Sprintf("Wildcard parameter: %s", paramName),
				Schema: &spec.Schema{
					Type: "string",
				},
			})
		}
	}

	return params
}
