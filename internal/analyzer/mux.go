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
	paths  []string
	routes []routeInfo
}

func (a *MuxAnalyzer) Analyze() (*spec.Documentation, error) {
	operations := make([]*spec.Operation, 0)
	basePath := ""

	// Processar arquivos de rota
	for _, routeFile := range a.config.RouterFiles {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, routeFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route file %s: %v", routeFile, err)
		}

		// Primeiro encontrar o PathPrefix se houver
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "PathPrefix" && len(call.Args) > 0 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							basePath = strings.Trim(lit.Value, "\"")
						}
					}
				}
			}
			return true
		})

		// Depois analisar as rotas
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					method := strings.ToUpper(sel.Sel.Name)
					if isHTTPMethod(method) && len(call.Args) >= 2 {
						operation := &spec.Operation{
							Method: method,
						}

						// Extrair path
						if pathLit, ok := call.Args[0].(*ast.BasicLit); ok {
							path := strings.Trim(pathLit.Value, "\"")
							if basePath != "" {
								operation.Path = basePath + path
							} else {
								operation.Path = path
							}
						}

						// Extrair handler name
						var handlerName string
						if ident, ok := call.Args[1].(*ast.Ident); ok {
							handlerName = ident.Name
							operation.OperationID = handlerName
						}

						// Extrair tags do path
						pathSegments := strings.Split(operation.Path, "/")
						if len(pathSegments) > 1 {
							operation.Tags = []string{pathSegments[1]} // Usar primeiro segmento após / como tag
						}

						// Analisar handler
						if handlerFunc := findHandler(handlerName, a.config.HandlerFiles); handlerFunc != nil {
							// Extrair comentários
							operation.Summary = extractHandlerComments(handlerFunc)

							// Extrair parâmetros de path
							operation.Parameters = append(operation.Parameters,
								extractPathParameters(operation.Path)...)

							// Extrair query parameters
							operation.Parameters = append(operation.Parameters,
								extractQueryParameters(handlerFunc)...)

							// Extrair request body
							if reqBody := a.extractRequestBody(handlerFunc, handlerName); reqBody != nil {
								operation.RequestBody = reqBody
							}

							// Extrair responses
							operation.Responses = map[string]*spec.Response{
								"200": {
									Description: "Successful response",
									Content: map[string]*spec.MediaType{
										"application/json": {
											Schema: &spec.Schema{
												Type: "object",
												Properties: map[string]*spec.Schema{
													"data": {
														Type: "object",
													},
												},
											},
										},
									},
								},
								"400": {
									Description: "Bad Request",
									Content: map[string]*spec.MediaType{
										"application/json": {
											Schema: &spec.Schema{
												Type: "object",
												Properties: map[string]*spec.Schema{
													"error": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							}
						}

						operations = append(operations, operation)
					}
				}
			}
			return true
		})
	}

	return &spec.Documentation{Operations: operations}, nil
}

func extractPathParameters(path string) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			name := strings.Trim(segment, "{}")
			params = append(params, &spec.Parameter{
				Name:        name,
				In:          "path",
				Required:    true,
				Description: fmt.Sprintf("ID do %s", name),
				Schema: &spec.Schema{
					Type: "string",
				},
			})
		}
	}

	return params
}

func extractQueryParameters(handlerFunc *ast.FuncDecl) []*spec.Parameter {
	params := make([]*spec.Parameter, 0)

	// Analisar o corpo da função para encontrar r.URL.Query()
	ast.Inspect(handlerFunc.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Query" {
					// Encontrou uma chamada Query, extrair parâmetro
					if parent := findParentCall(n); parent != nil {
						if lit, ok := parent.Args[0].(*ast.BasicLit); ok {
							paramName := strings.Trim(lit.Value, "\"")
							params = append(params, &spec.Parameter{
								Name:        paramName,
								In:          "query",
								Required:    false,
								Description: fmt.Sprintf("Parâmetro de consulta: %s", paramName),
								Schema: &spec.Schema{
									Type: "string",
								},
							})
						}
					}
				}
			}
		}
		return true
	})

	return params
}

func (a *MuxAnalyzer) extractRequestBody(handlerFunc *ast.FuncDecl, handlerName string) *spec.RequestBody {
	// Procurar por json.NewDecoder(r.Body).Decode(&req)
	var reqType *ast.TypeSpec
	ast.Inspect(handlerFunc.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Decode" {
					// Encontrou Decode, procurar o tipo
					if len(call.Args) > 0 {
						if unary, ok := call.Args[0].(*ast.UnaryExpr); ok {
							if ident, ok := unary.X.(*ast.Ident); ok {
								// Encontrar a definição do tipo no arquivo do handler
								for _, handlerFile := range a.config.HandlerFiles {
									if file, err := parser.ParseFile(token.NewFileSet(), handlerFile, nil, parser.ParseComments); err == nil {
										if typeSpec := findTypeSpec(file, ident.Name); typeSpec != nil {
											reqType = typeSpec
											break
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

	if reqType != nil {
		return &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: convertASTTypeToSchema(reqType),
				},
			},
		}
	}

	return nil
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

func findHandler(handlerName string, handlerFiles []string) *ast.FuncDecl {
	for _, file := range handlerFiles {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name.Name == handlerName {
					return fn
				}
			}
		}
	}
	return nil
}

func extractHandlerComments(handler *ast.FuncDecl) string {
	if handler.Doc != nil {
		return strings.TrimSpace(handler.Doc.Text())
	}
	return ""
}

func convertASTTypeToSchema(typeSpec *ast.TypeSpec) *spec.Schema {
	schema := &spec.Schema{
		Type:       "object",
		Properties: make(map[string]*spec.Schema),
	}

	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				fieldName := field.Names[0].Name
				fieldSchema := &spec.Schema{}

				// Determinar tipo do campo
				switch fieldType := field.Type.(type) {
				case *ast.Ident:
					fieldSchema.Type = strings.ToLower(fieldType.Name)
				case *ast.StarExpr:
					if ident, ok := fieldType.X.(*ast.Ident); ok {
						fieldSchema.Type = strings.ToLower(ident.Name)
					}
				}

				// Processar tags
				if field.Tag != nil {
					tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
					if jsonTag := tag.Get("json"); jsonTag != "" {
						parts := strings.Split(jsonTag, ",")
						fieldName = parts[0]
					}
					if validateTag := tag.Get("validate"); validateTag != "" {
						// Processar regras de validação
						rules := strings.Split(validateTag, ",")
						for _, rule := range rules {
							switch {
							case rule == "required":
								fieldSchema.Required = true
							case strings.HasPrefix(rule, "min="):
								// Processar min
							case strings.HasPrefix(rule, "max="):
								// Processar max
							}
						}
					}
				}

				schema.Properties[fieldName] = fieldSchema
			}
		}
	}

	return schema
}
