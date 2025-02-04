package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

// Config contém as configurações para análise
type Config struct {
	MainFile     string
	BaseDir      string
	RouterFiles  []string
	HandlerFiles []string
}

// Funções comuns utilizadas por múltiplos analyzers
func extractSummaryFromComments(node ast.Node) string {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		if funcDecl.Doc != nil {
			return strings.TrimSpace(funcDecl.Doc.Text())
		}
	}
	return ""
}

func extractRequestBody(node ast.Node, filename string) *spec.RequestBody {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		// Procurar por c.BodyParser no código
		var reqBody *spec.RequestBody
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selExpr.Sel.Name == "BodyParser" {
						if len(callExpr.Args) > 0 {
							if unary, ok := callExpr.Args[0].(*ast.UnaryExpr); ok {
								if ident, ok := unary.X.(*ast.Ident); ok {
									reqBody = &spec.RequestBody{
										Required: true,
										Content: map[string]*spec.MediaType{
											"application/json": {
												Schema: &spec.Schema{
													Type:       ident.Name,
													Properties: extractStructProperties(filename, ident.Name),
												},
											},
										},
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
	return nil
}

func extractResponses(node ast.Node, filename string) map[string]*spec.Response {
	responses := make(map[string]*spec.Response)

	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selExpr.Sel.Name == "Status" {
						if len(callExpr.Args) > 0 {
							if statusLit, ok := callExpr.Args[0].(*ast.SelectorExpr); ok {
								statusCode := statusLit.Sel.Name
								// Procurar pela chamada JSON no mesmo bloco
								ast.Inspect(funcDecl.Body, func(m ast.Node) bool {
									if jsonCall, ok := m.(*ast.CallExpr); ok {
										if jsonSel, ok := jsonCall.Fun.(*ast.SelectorExpr); ok {
											if jsonSel.Sel.Name == "JSON" {
												if len(jsonCall.Args) > 0 {
													if respArg := jsonCall.Args[0]; respArg != nil {
														var respType string
														switch t := respArg.(type) {
														case *ast.Ident:
															respType = t.Name
														case *ast.CompositeLit:
															if ident, ok := t.Type.(*ast.Ident); ok {
																respType = ident.Name
															}
														}

														if respType != "" {
															responses[statusCode] = &spec.Response{
																Description: fmt.Sprintf("%s Response", statusCode),
																Content: map[string]*spec.MediaType{
																	"application/json": {
																		Schema: &spec.Schema{
																			Type:       respType,
																			Properties: extractStructProperties(filename, respType),
																		},
																	},
																},
															}
														}
													}
												}
											}
										}
									}
									return true
								})
							}
						}
					}
				}
			}
			return true
		})
	}

	// Se nenhuma resposta foi encontrada, adicionar 200 OK como padrão
	if len(responses) == 0 {
		responses["200"] = &spec.Response{
			Description: "Successful response",
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
					},
				},
			},
		}
	}

	return responses
}

func isHTTPMethod(method string) bool {
	methods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	return methods[method]
}

func extractStructProperties(filename, structName string) map[string]*spec.Schema {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	properties := make(map[string]*spec.Schema)

	ast.Inspect(file, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if typeSpec.Name.Name == structName {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if len(field.Names) > 0 {
							fieldName := field.Names[0].Name

							// Extrair tags json e validações
							var required bool
							if field.Tag != nil {
								tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
								if jsonTag := tag.Get("json"); jsonTag != "" {
									parts := strings.Split(jsonTag, ",")
									if len(parts) > 0 {
										fieldName = parts[0]
									}
								}
								if validateTag := tag.Get("validate"); validateTag != "" {
									required = strings.Contains(validateTag, "required")
								}
							}

							// Determinar o tipo do campo
							var schema *spec.Schema
							switch t := field.Type.(type) {
							case *ast.Ident:
								schema = &spec.Schema{
									Type:     t.Name,
									Required: required,
								}
							case *ast.ArrayType:
								if ident, ok := t.Elt.(*ast.Ident); ok {
									schema = &spec.Schema{
										Type:     "array",
										Required: required,
										Items: &spec.Schema{
											Type: ident.Name,
										},
									}
								}
							case *ast.SelectorExpr:
								schema = &spec.Schema{
									Type:     t.Sel.Name,
									Required: required,
								}
							}

							if schema != nil {
								properties[fieldName] = schema
							}
						}
					}
				}
			}
		}
		return true
	})

	return properties
}

func extractFiberParameters(path string) []*spec.Parameter {
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

// routeInfo representa uma rota da API
type routeInfo struct {
	path        string
	method      string
	handlerName string
	basePath    string   // Usado pelo Mux para subrouters
	node        ast.Node // Usado para análise adicional
	description string   // Descrição da rota
}

// NewAnalyzer cria um novo analisador baseado no framework
func NewAnalyzer(framework string, config Config) (Analyzer, error) {
	if config.MainFile == "" {
		// Tentar encontrar o main.go se não foi especificado
		mainFile, err := FindMainFile(config.BaseDir)
		if err != nil {
			return nil, err
		}
		config.MainFile = mainFile
	}

	// Rastrear imports a partir do main.go
	tracker := NewImportTracker(config.MainFile)
	if err := tracker.TrackImports(config.MainFile); err != nil {
		return nil, fmt.Errorf("failed to track imports: %v", err)
	}

	// Atualizar config com os arquivos encontrados
	config.RouterFiles = tracker.routeFiles
	config.HandlerFiles = tracker.handlerFiles

	var analyzer Analyzer
	switch framework {
	case "gin":
		analyzer = NewGinAnalyzer(config)
	case "mux":
		analyzer = NewMuxAnalyzer(config)
	case "fiber":
		analyzer = NewFiberAnalyzer(config)
	default:
		return nil, fmt.Errorf("unsupported framework: %s", framework)
	}

	return analyzer, nil
}

func findTypeSpec(file *ast.File, typeName string) *ast.TypeSpec {
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == typeName {
						return typeSpec
					}
				}
			}
		}
	}
	return nil
}
