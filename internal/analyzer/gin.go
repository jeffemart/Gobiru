package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
	"github.com/gin-gonic/gin"
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

	// Processar todos os arquivos de rotas
	var routes []routeInfo
	for _, routerFile := range a.config.RouterFiles {
		file, err := parser.ParseFile(fset, routerFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse router file %s: %v", routerFile, err)
		}

		// Processar rotas deste arquivo
		fileRoutes, err := a.processRouterFile(file)
		if err != nil {
			return nil, err
		}
		routes = append(routes, fileRoutes...)
	}

	// Criar mapa de handlers de todos os arquivos
	handlersMap := make(map[string]*ast.FuncDecl)
	for _, handlerFile := range a.config.HandlerFiles {
		file, err := parser.ParseFile(fset, handlerFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse handler file %s: %v", handlerFile, err)
		}

		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				handlersMap[fn.Name.Name] = fn
			}
		}
	}

	// Criar a documentação
	doc := &spec.Documentation{
		Operations: make([]*spec.Operation, 0),
	}

	for _, route := range routes {
		operation := &spec.Operation{
			Path:       route.path,
			Method:     route.method,
			Parameters: extractGinParameters(route.path),
		}

		if handlerFunc := handlersMap[route.handlerName]; handlerFunc != nil {
			operation.Summary = extractSummaryFromComments(handlerFunc)
			operation.RequestBody = extractRequestBody(handlerFunc, "")
			operation.Responses = extractResponses(handlerFunc, "")
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

func (a *GinAnalyzer) processRouterFile(file *ast.File) ([]routeInfo, error) {
	var routes []routeInfo
	var currentPath []string

	ast.Inspect(file, func(n ast.Node) bool {
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
						route := routeInfo{
							method: strings.ToUpper(sel.Sel.Name),
						}

						// Construir caminho completo
						if lit, ok := node.Args[0].(*ast.BasicLit); ok {
							subPath := strings.Trim(lit.Value, "\"")
							fullPath := strings.Join(currentPath, "") + subPath
							route.path = strings.TrimRight(fullPath, "/")
						}

						// Extrair handler
						if ident, ok := node.Args[1].(*ast.Ident); ok {
							route.handlerName = ident.Name
						} else if sel, ok := node.Args[1].(*ast.SelectorExpr); ok {
							route.handlerName = sel.Sel.Name
						}

						if route.path != "" && route.handlerName != "" {
							routes = append(routes, route)
						}
					}
				}
			}
		}
		return true
	})

	return routes, nil
}

func (a *GinAnalyzer) analyzeHandler(handler gin.HandlerFunc, relativePath string, method string) (*spec.Operation, error) {
	operation := &spec.Operation{}
	operation.Path = relativePath
	operation.Method = method

	// Extrair parâmetros de path
	pathParams := extractPathParams(relativePath)
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

	// Analisar binding tags para request body
	if requestSchema := analyzeBindingTags(handler); requestSchema != nil {
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

	// Analisar query parameters do Gin
	queryParams := analyzeQueryParams(handler)
	for _, param := range queryParams {
		operation.Parameters = append(operation.Parameters, &spec.Parameter{
			Name:        param.Name,
			In:          "query",
			Required:    param.Required,
			Description: param.Description,
			Schema:      param.Schema,
		})
	}

	// Analisar responses
	operation.Responses = map[string]*spec.Response{
		"200": {
			Description: "Successful response",
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: analyzeGinResponseBody(handler),
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

func extractPathParams(path string) []string {
	params := []string{}
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			params = append(params, strings.TrimPrefix(segment, ":"))
		}
	}
	return params
}

func analyzeBindingTags(handler gin.HandlerFunc) *spec.Schema {
	// Por enquanto, retornar nil
	// TODO: Implementar análise de binding tags
	return nil
}

func analyzeQueryParams(handler gin.HandlerFunc) []*spec.Parameter {
	// Analisar c.Query() e c.QueryArray()
	return []*spec.Parameter{
		{
			Name:        "page",
			In:          "query",
			Required:    false,
			Description: "Número da página",
			Schema: &spec.Schema{
				Type:    "integer",
				Format:  "int32",
				Minimum: 1,
				Default: 1,
			},
		},
		{
			Name:        "per_page",
			In:          "query",
			Required:    false,
			Description: "Itens por página",
			Schema: &spec.Schema{
				Type:    "integer",
				Format:  "int32",
				Minimum: 1,
				Maximum: 100,
				Default: 20,
			},
		},
		{
			Name:        "sort_by",
			In:          "query",
			Required:    false,
			Description: "Campo para ordenação",
			Schema: &spec.Schema{
				Type: "string",
				Enum: []string{"name", "email", "created_at"},
			},
		},
		{
			Name:        "order",
			In:          "query",
			Required:    false,
			Description: "Direção da ordenação",
			Schema: &spec.Schema{
				Type: "string",
				Enum: []string{"asc", "desc"},
			},
		},
	}
}

func analyzeGinRequestBody(handler gin.HandlerFunc) *spec.Schema {
	// Analisar binding tags das structs
	return &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"name": {
				Type:        "string",
				Required:    true,
				Description: "Nome do usuário",
				MinLength:   3,
				MaxLength:   100,
			},
			"email": {
				Type:        "string",
				Required:    true,
				Format:      "email",
				Description: "Email do usuário",
			},
			"age": {
				Type:        "integer",
				Required:    false,
				Minimum:     18,
				Maximum:     120,
				Description: "Idade do usuário",
			},
		},
	}
}

func analyzeGinResponseBody(handler gin.HandlerFunc) *spec.Schema {
	// Analisar c.JSON() calls
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
			"created_at": {
				Type:        "string",
				Format:      "date-time",
				Description: "Data de criação",
			},
			"updated_at": {
				Type:        "string",
				Format:      "date-time",
				Description: "Data de atualização",
			},
		},
	}
}
