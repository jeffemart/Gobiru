package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type FiberAnalyzer struct {
	BaseAnalyzer
}

func NewFiberAnalyzer(config Config) *FiberAnalyzer {
	return &FiberAnalyzer{
		BaseAnalyzer: BaseAnalyzer{
			config: config,
		},
	}
}

func (a *FiberAnalyzer) Analyze() (*spec.Documentation, error) {
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
		Parameters  []*spec.Parameter
	}

	// Manter controle dos grupos e seus caminhos
	type groupInfo struct {
		path     string
		children []*groupInfo
	}

	var groups []*groupInfo
	currentGroup := &groupInfo{path: ""}
	groups = append(groups, currentGroup)

	// Primeira passagem: encontrar grupos
	ast.Inspect(routerFile, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Group" && len(callExpr.Args) > 0 {
					if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
						groupPath := strings.Trim(lit.Value, "\"")
						newGroup := &groupInfo{
							path: currentGroup.path + groupPath,
						}
						currentGroup.children = append(currentGroup.children, newGroup)
						currentGroup = newGroup
					}
				}
			}
		}
		return true
	})

	// Segunda passagem: encontrar rotas
	ast.Inspect(routerFile, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				method := strings.ToUpper(sel.Sel.Name)
				if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
					if len(callExpr.Args) >= 2 {
						route := struct {
							Path        string
							Method      string
							HandlerName string
							Parameters  []*spec.Parameter
						}{
							Method:     method,
							Parameters: make([]*spec.Parameter, 0),
						}

						// Extrair path
						if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
							subPath := strings.Trim(lit.Value, "\"")
							route.Path = currentGroup.path + subPath
						}

						// Extrair handler
						if ident, ok := callExpr.Args[1].(*ast.Ident); ok {
							route.HandlerName = ident.Name
						} else if sel, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
							route.HandlerName = sel.Sel.Name
						}

						// Extrair parâmetros
						if route.Path != "" {
							route.Parameters = extractFiberParameters(route.Path)
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
			Parameters: route.Parameters,
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
