package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/gofiber/fiber/v2"
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
							operation.Parameters = extractFiberParameters(operation.Path)
						}

						// Extrair handler
						var handlerName string
						if ident, ok := call.Args[1].(*ast.Ident); ok {
							handlerName = ident.Name
						} else if sel, ok := call.Args[1].(*ast.SelectorExpr); ok {
							handlerName = sel.Sel.Name
						}

						// Processar handler para mais detalhes
						if handlerName != "" {
							handlerFile := findHandlerFile(a.config.HandlerFiles, handlerName)
							if handlerFile != "" {
								// Adicionar summary do comentário
								operation.Summary = findHandlerComments(a.config.HandlerFiles, handlerName)

								// Extrair informações adicionais do handler
								handler := findHandlerNode(handlerFile, handlerName)
								if handler != nil {
									// Extrair corpo da requisição
									operation.RequestBody = extractRequestBody(handler, handlerFile)

									// Extrair respostas
									operation.Responses = extractResponses(handler, handlerFile)
								}
							}
						}

						operations = append(operations, operation)
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

// Encontra o arquivo do handler para um determinado nome de handler
func findHandlerFile(handlerFiles []string, handlerName string) string {
	for _, file := range handlerFiles {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name.Name == handlerName {
					return file
				}
			}
		}
	}
	return ""
}

// Encontra o nó AST do handler
func findHandlerNode(filename, handlerName string) ast.Node {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	var handlerNode ast.Node
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == handlerName {
				handlerNode = n
				return false
			}
		}
		return true
	})

	return handlerNode
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

func (a *FiberAnalyzer) analyzeHandler(handler fiber.Handler, path string, method string) (*spec.Operation, error) {
	operation := &spec.Operation{}
	operation.Path = path
	operation.Method = method

	// Extrair parâmetros de path
	pathParams := extractFiberPathParams(path)
	for _, param := range pathParams {
		operation.Parameters = append(operation.Parameters, &spec.Parameter{
			Name:        param,
			In:          "path",
			Required:    true,
			Description: fmt.Sprintf("Parameter %s", param),
			Schema: &spec.Schema{
				Type: "string",
			},
		})
	}

	// Analisar query parameters
	queryParams := analyzeFiberQueryParams(handler)
	for _, param := range queryParams {
		operation.Parameters = append(operation.Parameters, &spec.Parameter{
			Name:        param.Name,
			In:          "query",
			Required:    param.Required,
			Description: param.Description,
			Schema:      param.Schema,
		})
	}

	// Analisar request body
	if requestSchema := analyzeFiberRequestBody(handler); requestSchema != nil {
		operation.RequestBody = &spec.RequestBody{
			Description: "Request body",
			Required:    true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: requestSchema,
				},
			},
		}
	}

	// Analisar responses
	operation.Responses = map[string]*spec.Response{
		"200": {
			Description: "Successful response",
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: analyzeFiberResponseBody(handler),
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
							"error": {Type: "string"},
						},
					},
				},
			},
		},
	}

	return operation, nil
}

func extractFiberPathParams(path string) []string {
	params := []string{}
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			params = append(params, strings.TrimPrefix(segment, ":"))
		}
	}
	return params
}

func analyzeFiberQueryParams(handler fiber.Handler) []*spec.Parameter {
	// Analisar c.Query() e c.QueryParser()
	return []*spec.Parameter{
		{
			Name:        "search",
			In:          "query",
			Required:    false,
			Description: "Termo de busca",
			Schema: &spec.Schema{
				Type: "string",
			},
		},
		{
			Name:        "status",
			In:          "query",
			Required:    false,
			Description: "Status do usuário",
			Schema: &spec.Schema{
				Type: "string",
				Enum: []string{"active", "inactive", "pending"},
			},
		},
		{
			Name:        "filters",
			In:          "query",
			Required:    false,
			Description: "Filtros adicionais",
			Schema: &spec.Schema{
				Type: "array",
				Items: &spec.Schema{
					Type: "string",
				},
			},
		},
	}
}

func analyzeFiberRequestBody(handler fiber.Handler) *spec.Schema {
	// Analisar c.BodyParser()
	return &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"name": {
				Type:        "string",
				Required:    true,
				Description: "Nome do usuário",
			},
			"email": {
				Type:        "string",
				Required:    true,
				Format:      "email",
				Description: "Email do usuário",
			},
			"profile": {
				Type: "object",
				Properties: map[string]*spec.Schema{
					"avatar": {
						Type:        "string",
						Format:      "uri",
						Description: "URL do avatar",
					},
					"bio": {
						Type:        "string",
						Description: "Biografia do usuário",
					},
				},
			},
		},
	}
}

func analyzeFiberResponseBody(handler fiber.Handler) *spec.Schema {
	// Analisar c.JSON()
	return &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"id": {
				Type:        "string",
				Format:      "uuid",
				Description: "ID do usuário",
			},
			"name": {
				Type:        "string",
				Description: "Nome do usuário",
			},
			"email": {
				Type:        "string",
				Format:      "email",
				Description: "Email do usuário",
			},
			"profile": {
				Type: "object",
				Properties: map[string]*spec.Schema{
					"avatar": {
						Type:        "string",
						Format:      "uri",
						Description: "URL do avatar",
					},
					"bio": {
						Type:        "string",
						Description: "Biografia do usuário",
					},
				},
			},
			"created_at": {
				Type:        "string",
				Format:      "date-time",
				Description: "Data de criação",
			},
		},
	}
}
