package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

type MuxAnalyzer struct {
	BaseAnalyzer
}

func NewMuxAnalyzer(config Config) *MuxAnalyzer {
	return &MuxAnalyzer{
		BaseAnalyzer: BaseAnalyzer{
			config: config,
		},
	}
}

func (a *MuxAnalyzer) Analyze() (*spec.Documentation, error) {
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
			// Capturar PathPrefix e HandleFunc
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				switch sel.Sel.Name {
				case "PathPrefix":
					if len(node.Args) > 0 {
						if lit, ok := node.Args[0].(*ast.BasicLit); ok {
							prefix := strings.Trim(lit.Value, "\"")
							currentPath = append(currentPath, prefix)
						}
					}
				case "HandleFunc":
					if len(node.Args) >= 2 {
						route := struct {
							Path        string
							Method      string
							HandlerName string
						}{}

						// Construir caminho completo
						if lit, ok := node.Args[0].(*ast.BasicLit); ok {
							subPath := strings.Trim(lit.Value, "\"")
							fullPath := strings.Join(currentPath, "") + subPath
							route.Path = strings.TrimRight(fullPath, "/")
						}

						// Extrair handler
						if ident, ok := node.Args[1].(*ast.Ident); ok {
							route.HandlerName = ident.Name
						} else if sel, ok := node.Args[1].(*ast.SelectorExpr); ok {
							route.HandlerName = sel.Sel.Name
						}

						// Procurar pelo método HTTP
						ast.Inspect(n, func(m ast.Node) bool {
							if call, ok := m.(*ast.CallExpr); ok {
								if methodSel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if methodSel.Sel.Name == "Methods" && len(call.Args) > 0 {
										if lit, ok := call.Args[0].(*ast.BasicLit); ok {
											route.Method = strings.Trim(lit.Value, "\"")
										}
									}
								}
							}
							return true
						})

						if route.Path != "" && route.HandlerName != "" {
							if route.Method == "" {
								route.Method = "GET" // Método padrão
							}
							routes = append(routes, route)
						}
					}
				case "Subrouter":
					// Manter o caminho atual para o subrouter
					return true
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
			Method:     strings.ToUpper(route.Method),
			Parameters: extractMuxParameters(route.Path),
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

func extractMuxParameters(path string) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			paramName := strings.Trim(segment, "{}")
			params = append(params, &spec.Parameter{
				Name:        paramName,
				In:          "path",
				Required:    true,
				Description: fmt.Sprintf("Path parameter: %s", paramName),
				Schema: &spec.Schema{
					Type: "string",
				},
			})
		}
	}

	return params
}

func extractQueryParams(node ast.Node) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Queries" {
					for i := 0; i < len(call.Args); i += 2 {
						if i+1 < len(call.Args) {
							if name, ok := call.Args[i].(*ast.BasicLit); ok {
								if pattern, ok := call.Args[i+1].(*ast.BasicLit); ok {
									paramName := strings.Trim(name.Value, "\"")
									paramPattern := strings.Trim(pattern.Value, "\"")
									required := !strings.HasSuffix(paramPattern, "?")

									params = append(params, &spec.Parameter{
										Name:     paramName,
										In:       "query",
										Required: required,
										Schema: &spec.Schema{
											Type: "string",
										},
									})
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return params
}

// Funções auxiliares...
func findParentCall(node ast.Node) *ast.CallExpr {
	var parent *ast.CallExpr
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			parent = call
		}
		return true
	})
	return parent
}

func extractSchema(structType *ast.StructType) map[string]interface{} {
	schema := make(map[string]interface{})
	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 {
			fieldName := field.Names[0].Name
			fieldType := ""
			if ident, ok := field.Type.(*ast.Ident); ok {
				fieldType = ident.Name
			}
			schema[fieldName] = fieldType
		}
	}
	return schema
}

// Adicionar estas funções auxiliares
func findFunction(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == name {
				return fn
			}
		}
	}
	return nil
}
