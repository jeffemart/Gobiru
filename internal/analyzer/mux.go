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

type routeInfo struct {
	basePath    string
	path        string
	method      string
	handlerName string
}

type routeContext struct {
	paths  []string
	routes []routeInfo
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

	ctx := &routeContext{
		paths:  make([]string, 0),
		routes: make([]routeInfo, 0),
	}

	// Encontrar a função SetupRoutes
	var setupFound bool
	for _, decl := range routerFile.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == "SetupRoutes" {
				setupFound = true
				a.processRouterFunction(fn, ctx)
				break
			}
		}
	}

	if !setupFound {
		return nil, fmt.Errorf("SetupRoutes function not found in router file")
	}

	// Criar a documentação
	doc := &spec.Documentation{
		Operations: make([]*spec.Operation, 0),
	}

	for _, route := range ctx.routes {
		operation := &spec.Operation{
			Path:       route.path,
			Method:     route.method,
			Parameters: extractMuxParameters(route.path),
		}

		if handlerFunc := findFunction(handlersFile, route.handlerName); handlerFunc != nil {
			operation.Summary = extractSummaryFromComments(handlerFunc)
			operation.RequestBody = extractRequestBody(handlerFunc, a.config.HandlersFile)
			operation.Responses = extractResponses(handlerFunc, a.config.HandlersFile)
		}

		doc.Operations = append(doc.Operations, operation)
		fmt.Printf("Added operation: %s %s -> %s\n", operation.Method, operation.Path, route.handlerName)
	}

	if len(doc.Operations) == 0 {
		return nil, fmt.Errorf("no operations found in router file")
	}

	return doc, nil
}

func (a *MuxAnalyzer) processRouterFunction(fn *ast.FuncDecl, ctx *routeContext) {
	var currentPath []string
	var subrouterVars = make(map[string]string) // Mapear variáveis para seus caminhos

	// Primeira passagem: mapear todos os subrouters
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if expr, ok := n.(*ast.AssignStmt); ok {
			for i, rhs := range expr.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						if sel.Sel.Name == "Subrouter" {
							if pathCall, ok := sel.X.(*ast.CallExpr); ok {
								if pathSel, ok := pathCall.Fun.(*ast.SelectorExpr); ok {
									if pathSel.Sel.Name == "PathPrefix" && len(pathCall.Args) > 0 {
										if lit, ok := pathCall.Args[0].(*ast.BasicLit); ok {
											path := strings.Trim(lit.Value, "\"")
											currentPath = append(currentPath, path)

											if i < len(expr.Lhs) {
												if ident, ok := expr.Lhs[i].(*ast.Ident); ok {
													subrouterVars[ident.Name] = strings.Join(currentPath, "")
													fmt.Printf("Mapped subrouter %s to path: %s\n",
														ident.Name, subrouterVars[ident.Name])
												}
											}

											fmt.Printf("Added path prefix: %s (current path: %s)\n",
												path, strings.Join(currentPath, ""))
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	// Segunda passagem: encontrar todas as rotas
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if expr, ok := n.(*ast.CallExpr); ok {
			if sel, ok := expr.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "HandleFunc" {
					if x, ok := sel.X.(*ast.Ident); ok {
						basePath := subrouterVars[x.Name]
						fmt.Printf("Processing HandleFunc for subrouter %s with basePath: %s\n", x.Name, basePath)

						route := routeInfo{
							basePath: basePath,
							path:     "/", // Valor padrão para path vazio
						}

						// Extrair path
						if len(expr.Args) >= 1 {
							if lit, ok := expr.Args[0].(*ast.BasicLit); ok {
								pathValue := strings.Trim(lit.Value, "\"")
								if pathValue != "" {
									route.path = pathValue
								}
								fmt.Printf("Found path: %s\n", route.path)
							}
						}

						// Extrair handler
						if len(expr.Args) >= 2 {
							if ident, ok := expr.Args[1].(*ast.Ident); ok {
								route.handlerName = ident.Name
							} else if sel, ok := expr.Args[1].(*ast.SelectorExpr); ok {
								route.handlerName = sel.Sel.Name
							}
							fmt.Printf("Found handler: %s\n", route.handlerName)
						}

						// Procurar método HTTP
						ast.Inspect(expr, func(m ast.Node) bool {
							if call, ok := m.(*ast.CallExpr); ok {
								if methodSel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if methodSel.Sel.Name == "Methods" && len(call.Args) > 0 {
										if lit, ok := call.Args[0].(*ast.BasicLit); ok {
											route.method = strings.Trim(lit.Value, "\"")
											fmt.Printf("Found HTTP method: %s\n", route.method)
										}
									}
								}
							}
							return true
						})

						if route.handlerName != "" { // Removida a verificação de path vazio
							if route.method == "" {
								route.method = "GET"
							}
							route.path = route.basePath + route.path
							// Remover barras duplas se houverem
							route.path = strings.ReplaceAll(route.path, "//", "/")
							ctx.routes = append(ctx.routes, route)
							fmt.Printf("Added route: %s %s -> %s\n",
								route.method, route.path, route.handlerName)
						} else {
							fmt.Printf("Skipping incomplete route: path=%s, handler=%s, method=%s\n",
								route.path, route.handlerName, route.method)
						}
					}
				}
			}
		}
		return true
	})

	// Debug: imprimir todas as rotas encontradas
	fmt.Printf("\nDebug - Todas as rotas encontradas:\n")
	for _, route := range ctx.routes {
		fmt.Printf("Route: %s %s -> %s (basePath: %s)\n",
			route.method, route.path, route.handlerName, route.basePath)
	}

	if len(ctx.routes) == 0 {
		fmt.Printf("Warning: No routes were found in the router file\n")
	}
}

// Função auxiliar para verificar se uma chamada Methods está relacionada a um HandleFunc
func isRelatedMethod(methodCall, handleFunc *ast.CallExpr) bool {
	// Verificar se o Methods está na mesma cadeia de chamadas que o HandleFunc
	parent := methodCall
	for {
		if parent == nil {
			return false
		}
		if parent == handleFunc {
			return true
		}
		// Tentar encontrar o próximo pai na cadeia
		parent = findParentCall(parent)
	}
}

// Função auxiliar para encontrar a chamada pai
func findParentCall(node ast.Node) *ast.CallExpr {
	var parent *ast.CallExpr
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if call != node {
				parent = call
				return false // Parar após encontrar o primeiro pai
			}
		}
		return true
	})
	return parent
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
