package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/models"
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

func (a *MuxAnalyzer) Analyze() ([]models.RouteInfo, error) {
	fset := token.NewFileSet()

	// Analisar arquivos
	routerFile, err := parser.ParseFile(fset, a.config.RouterFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse router file: %v", err)
	}

	handlersFile, err := parser.ParseFile(fset, a.config.HandlersFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse handlers file: %v", err)
	}

	var routes []models.RouteInfo

	// Encontrar estruturas para schemas
	schemas := make(map[string]interface{})
	ast.Inspect(handlersFile, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if structType, ok := typeSpec.Type.(*ast.StructType); ok {
			schemas[typeSpec.Name.Name] = extractSchema(structType)
		}
		return true
	})

	// Analisar rotas
	ast.Inspect(routerFile, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "HandleFunc" {
				route := models.RouteInfo{
					Version: "v1.0",
					RateLimit: models.RateLimitConfig{
						RequestsPerMinute: 100,
						TimeWindowSeconds: 60,
					},
				}

				// Extrair path
				if len(callExpr.Args) > 0 {
					if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
						route.Path = strings.Trim(lit.Value, "\"")
						route.Parameters = extractPathParams(route.Path)
					}
				}

				// Extrair handler
				if len(callExpr.Args) > 1 {
					if ident, ok := callExpr.Args[1].(*ast.Ident); ok {
						route.HandlerName = ident.Name

						// Analisar função handler
						if handlerFunc := findFunction(handlersFile, ident.Name); handlerFunc != nil {
							route.Description = extractDescription(handlerFunc)
							route.Request = extractRequestBody(handlerFunc, schemas)
							route.Responses = extractResponses(handlerFunc, schemas)
						}
					}
				}

				// Extrair método e query params
				parent := findParentCall(n)
				if parent != nil {
					if sel, ok := parent.Fun.(*ast.SelectorExpr); ok {
						if sel.Sel.Name == "Methods" {
							if len(parent.Args) > 0 {
								if lit, ok := parent.Args[0].(*ast.BasicLit); ok {
									route.Method = strings.Trim(lit.Value, "\"")
								}
							}
						}
					}
				}

				// Extrair query params
				queryParams := extractQueryParams(n)
				if len(queryParams) > 0 {
					route.QueryParams = queryParams
				}

				routes = append(routes, route)
			}
		}

		return true
	})

	return routes, nil
}

func extractPathParams(path string) []models.Parameter {
	var params []models.Parameter
	parts := strings.Split(path, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			name := strings.Trim(part, "{}")
			params = append(params, models.Parameter{
				Name:     name,
				Type:     "string",
				Required: true,
			})
		}
	}

	return params
}

func extractQueryParams(node ast.Node) []models.Parameter {
	var params []models.Parameter

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

									params = append(params, models.Parameter{
										Name:     paramName,
										Type:     "string",
										Required: required,
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

func extractDescription(fn *ast.FuncDecl) string {
	if fn.Doc != nil {
		return fn.Doc.Text()
	}
	return ""
}

func extractRequestBody(fn *ast.FuncDecl, schemas map[string]interface{}) models.RequestBody {
	// Procurar por decodificação de JSON no corpo da função
	var reqBody models.RequestBody
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Decode" || sel.Sel.Name == "ShouldBindJSON" {
					if len(call.Args) > 0 {
						if star, ok := call.Args[0].(*ast.UnaryExpr); ok {
							if ident, ok := star.X.(*ast.Ident); ok {
								reqBody.Type = ident.Name
								if schema, ok := schemas[ident.Name]; ok {
									reqBody.Schema = schema
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	return reqBody
}

func extractResponses(fn *ast.FuncDecl, schemas map[string]interface{}) []models.Response {
	var responses []models.Response

	// Procurar por WriteHeader e Encode no corpo da função
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				switch sel.Sel.Name {
				case "WriteHeader":
					if len(call.Args) > 0 {
						if ident, ok := call.Args[0].(*ast.Ident); ok {
							responses = append(responses, models.Response{
								StatusCode:  200, // valor padrão
								Description: ident.Name,
							})
						}
					}
				case "Encode":
					if len(call.Args) > 0 {
						if ident, ok := call.Args[0].(*ast.Ident); ok {
							if len(responses) > 0 {
								lastResponse := &responses[len(responses)-1]
								lastResponse.Type = ident.Name
								if schema, ok := schemas[ident.Name]; ok {
									lastResponse.Schema = schema
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	// Se nenhuma resposta foi encontrada, adicionar 200 OK como padrão
	if len(responses) == 0 {
		responses = append(responses, models.Response{
			StatusCode:  200,
			Description: "OK",
		})
	}

	return responses
}

// ... outras funções auxiliares
