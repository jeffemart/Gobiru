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

type routeContext struct {
	routes []routeInfo
}

func (a *MuxAnalyzer) Analyze() (*spec.Documentation, error) {
	operations := make([]*spec.Operation, 0)

	for _, routeFile := range a.config.RouterFiles {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, routeFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route file %s: %v", routeFile, err)
		}

		// Contexto para rastrear informações durante a análise
		ctx := &routeContext{
			routes: make([]routeInfo, 0),
		}

		// Primeira passagem: encontrar todos os subrouters
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				a.processRouterFunction(fn, ctx)
			}
		}

		// Converter rotas em operações
		for _, route := range ctx.routes {
			// Encontrar handler correspondente
			handlerDoc := ""
			handlerName := route.handlerName
			for _, handlerFile := range a.config.HandlerFiles {
				if doc := findHandlerDoc(handlerFile, handlerName); doc != "" {
					handlerDoc = doc
					break
				}
			}

			// Determinar método HTTP
			method := route.method
			if method == "" {
				// Inferir método do nome do handler
				switch {
				case strings.HasPrefix(handlerName, "Get"):
					method = "GET"
				case strings.HasPrefix(handlerName, "Create"):
					method = "POST"
				case strings.HasPrefix(handlerName, "Update"):
					method = "PUT"
				case strings.HasPrefix(handlerName, "Delete"):
					method = "DELETE"
				case strings.HasPrefix(handlerName, "Patch"):
					method = "PATCH"
				default:
					method = "GET"
				}
			}

			operation := &spec.Operation{
				Path:    route.path,
				Method:  method,
				Summary: handlerDoc,
			}

			// Extrair parâmetros do path
			operation.Parameters = extractMuxParameters(route.path)

			// Adicionar query parameters se existirem
			queryParams := extractQueryParams(route.node)
			operation.Parameters = append(operation.Parameters, queryParams...)

			// Adicionar response padrão
			operation.Responses = map[string]*spec.Response{
				"200": {
					Description: "Successful response",
					Content: map[string]*spec.MediaType{
						"application/json": {
							Schema: &spec.Schema{
								Type: "object",
							},
						},
					},
				},
			}

			operations = append(operations, operation)
		}
	}

	return &spec.Documentation{
		Operations: operations,
	}, nil
}

// findHandlerDoc procura a documentação de um handler específico
func findHandlerDoc(handlerFile string, handlerName string) string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, handlerFile, nil, parser.ParseComments)
	if err != nil {
		return ""
	}

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == handlerName && fn.Doc != nil {
				return strings.TrimSpace(fn.Doc.Text())
			}
		}
	}

	return ""
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
							path:     "/",  // Valor padrão para path vazio
							node:     expr, // Adicionar o nó AST
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
						route.method = a.extractMethod(expr)

						if route.handlerName != "" {
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

	// Remover rotas incompletas ou inválidas antes de adicionar ao contexto
	validRoutes := make([]routeInfo, 0)
	for _, route := range ctx.routes {
		if route.path != "" {
			// Normalizar o path
			route.path = strings.TrimSpace(route.path)
			route.path = strings.TrimSuffix(route.path, "/")
			if !strings.HasPrefix(route.path, "/") {
				route.path = "/" + route.path
			}

			// Inferir método HTTP baseado no path e nome do handler
			if route.method == "" || route.method == "GET" {
				// Primeiro tentar pelo nome do handler
				switch {
				case strings.HasPrefix(route.handlerName, "Create"):
					route.method = "POST"
				case strings.HasPrefix(route.handlerName, "Update"):
					route.method = "PUT"
				case strings.HasPrefix(route.handlerName, "Delete"):
					route.method = "DELETE"
				case strings.HasPrefix(route.handlerName, "Patch"):
					route.method = "PATCH"
				case strings.HasPrefix(route.handlerName, "Get"):
					route.method = "GET"
				default:
					// Se não encontrou pelo handler, tentar pelo path
					pathLower := strings.ToLower(route.path)
					switch {
					case strings.Contains(pathLower, "login") ||
						strings.Contains(pathLower, "register") ||
						strings.Contains(pathLower, "forgot-password") ||
						strings.Contains(pathLower, "reset-password"):
						route.method = "POST"
					case strings.Contains(pathLower, "employees"):
						if strings.HasSuffix(pathLower, "employees") {
							route.method = "GET" // Listar employees
							if route.handlerName == "" {
								route.handlerName = "ListEmployees"
								route.description = "Lista todos os funcionários"
							}
						} else if strings.Contains(pathLower, "status") {
							route.method = "PUT" // Atualizar status
							if route.handlerName == "" {
								route.handlerName = "UpdateEmployeeStatus"
								route.description = "Atualiza o status do funcionário"
							}
						} else {
							route.method = "GET" // Get employee by ID
							if route.handlerName == "" {
								route.handlerName = "GetEmployee"
								route.description = "Retorna os dados de um funcionário específico"
							}
						}
					}
				}
			}

			validRoutes = append(validRoutes, route)
		}
	}

	// Remover duplicatas mantendo a versão mais completa
	seenPaths := make(map[string]routeInfo)
	for _, route := range validRoutes {
		key := fmt.Sprintf("%s %s", route.method, route.path)
		if existing, exists := seenPaths[key]; exists {
			// Manter a versão com mais informações
			if route.handlerName != "" && existing.handlerName == "" {
				seenPaths[key] = route
			}
		} else {
			seenPaths[key] = route
		}
	}

	// Converter mapa de volta para slice
	uniqueRoutes := make([]routeInfo, 0, len(seenPaths))
	for _, route := range seenPaths {
		uniqueRoutes = append(uniqueRoutes, route)
	}

	ctx.routes = uniqueRoutes
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
	if node == nil {
		return params
	}

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

func (a *MuxAnalyzer) extractMethod(node *ast.CallExpr) string {
	var method string

	// Procurar por chamadas ao método Methods() em toda a cadeia
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Methods" && len(call.Args) > 0 {
					// Verificar se o Methods está relacionado ao HandleFunc atual
					if isRelatedMethod(call, node) {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							method = strings.Trim(lit.Value, "\"")
							return false // Parar após encontrar o método
						}
					}
				}
			}
		}
		return true
	})

	// Se não encontrou Methods(), usar outras estratégias...
	if method == "" {
		// ... resto do código permanece igual ...
	}

	return method
}
