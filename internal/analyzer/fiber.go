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
	operations := make([]*spec.Operation, 0)

	// Processar arquivos de rota
	for _, routeFile := range a.config.RouterFiles {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, routeFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route file %s: %v", routeFile, err)
		}

		// Encontrar todas as definições de rota
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					method := strings.ToUpper(sel.Sel.Name)

					// Verificar se é um método HTTP
					if isFiberHTTPMethod(method) && len(call.Args) >= 2 {
						operation := &spec.Operation{
							Method: method,
						}

						// Extrair path
						if pathLit, ok := call.Args[0].(*ast.BasicLit); ok {
							operation.Path = strings.Trim(pathLit.Value, "\"")
						}

						// Extrair handler
						if ident, ok := call.Args[1].(*ast.Ident); ok {
							operation.Summary = findHandlerComments(a.config.HandlerFiles, ident.Name)
							operations = append(operations, operation)
						} else if sel, ok := call.Args[1].(*ast.SelectorExpr); ok {
							operation.Summary = findHandlerComments(a.config.HandlerFiles, sel.Sel.Name)
							operations = append(operations, operation)
						}
					}
				}
			}
			return true
		})
	}

	if len(operations) == 0 {
		fmt.Println("Warning: No operations found in route files")
		fmt.Println("Route files found:", a.config.RouterFiles)
		fmt.Println("Handler files found:", a.config.HandlerFiles)
	}

	return &spec.Documentation{
		Operations: operations,
	}, nil
}

func isFiberHTTPMethod(method string) bool {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	method = strings.ToUpper(method)
	for _, m := range methods {
		if method == m {
			return true
		}
	}
	return false
}

func findHandlerComments(handlerFiles []string, handlerName string) string {
	for _, file := range handlerFiles {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name.Name == handlerName && fn.Doc != nil {
					return strings.TrimSpace(fn.Doc.Text())
				}
			}
		}
	}
	return ""
}

// Função auxiliar para extrair parâmetros da rota
func extractRouteParams(path string) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			name := strings.TrimPrefix(segment, ":")
			params = append(params, &spec.Parameter{
				Name:        name,
				In:          "path",
				Required:    true,
				Description: fmt.Sprintf("Path parameter: %s", name),
				Schema: &spec.Schema{
					Type: "string",
				},
			})
		}
	}

	return params
}
