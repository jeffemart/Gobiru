// Package analyzer provides tools for analyzing and processing Go code for various frameworks.
package analyzer

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/gofiber/fiber/v2/log"
)

// Validation represents the validation structure for HTTP requests.
type Validation struct {
	Body   []Set
	Query  []Set
	Params []Set
}

// Set defines a key-value pair with validation rules.
type Set struct {
	Key      string
	Type     string
	Rules    string
	Required bool
}

// ImportsType maps import names to their paths.
type ImportsType map[string]string

// HttpMethods stores supported HTTP methods.
var HttpMethods = []string{"Get", "Post", "Put", "Patch", "Delete", "Head", "Options"}

// rootDir represents the root directory of the project.
var rootDir = ""

// paths stores the mapping of paths to their respective PathItems.
var paths = make(map[string]spec.PathItem)

// fset is used to hold the token file set.
var fset = token.NewFileSet()

// findRootDir initializes the rootDir by searching for the go.mod file.
func findRootDir(mainFilePath string) {
	path := filepath.Dir(mainFilePath)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}
	for _, v := range files {
		if v.Name() == "go.mod" {
			rootDir = path
			break
		}
	}
	if rootDir == "" {
		findRootDir(path)
	}
}

// traverseMain parses the main Go file to extract routing information.
func traverseMain(filePath string) error {
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return err
	}

	fiberName, mandragoraName, imports := trackImports(f)
	routerName := "app"

	for _, v := range f.Decls {
		switch x := v.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "main" {
				routerName = findRouterMainName(x, fiberName)
				findRoutes(x, routerName, mandragoraName, imports)
				break
			}
		}
	}
	return nil
}

// trackImports analyzes the imports in a Go file and returns relevant information.
func trackImports(f *ast.File) (string, string, ImportsType) {
	tempImports := ImportsType{}
	mandragoraName := "mandragora"
	fiberName := "fiber"
	for _, v := range f.Imports {
		name := ""
		ctx := build.Default
		ctx.Dir = rootDir
		pkg, err := ctx.Import(strings.ReplaceAll(v.Path.Value, "\"", ""), rootDir, 0)
		if err != nil {
			log.Errorf("Error importing package: %v\n", err)
		}
		if v.Name != nil {
			name = v.Name.Name
		} else {
			name = pkg.Name
		}
		tempImports[name] = pkg.Dir
		if v.Path.Value == "\"github.com/Camada8/mandragora\"" {
			mandragoraName = name
		}
		if v.Path.Value == "\"github.com/gofiber/fiber/v2\"" {
			fiberName = name
		}
	}
	return fiberName, mandragoraName, tempImports
}

// findRouterMainName identifies the main router name from the main function.
func findRouterMainName(mainFunc *ast.FuncDecl, fiberName string) (routerName string) {
	routerName = "app"
	ast.Inspect(mainFunc, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			call := x.Rhs[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
			if call.X.(*ast.Ident).Name == fiberName && call.Sel.Name == "New" {
				routerName = x.Lhs[0].(*ast.Ident).Name
			}
		}
		return true
	})
	return
}

// findRoutes inspects a function to identify route configurations and updates the global paths variable.
func findRoutes(routesFunc *ast.FuncDecl, routerName, mandragoraName string, imports ImportsType) {
	ast.Inspect(routesFunc, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			switch fun := x.Fun.(type) {
			case *ast.SelectorExpr:
				if fun.X.(*ast.Ident).Name == routerName && isValidHttpMethod(fun.Sel.Name) {
					parameters := getValidation(x.Args, mandragoraName, imports)
					path, err := strconv.Unquote(x.Args[0].(*ast.BasicLit).Value)
					if err != nil {
						log.Errorf("Error parsing path: %v\n", err)
					}
					var tempPath spec.PathItem
					var params []spec.Parameter
					if len(parameters) > 0 {
						for _, p := range parameters {
							params = append(params, *p)
						}
					}
					switch fun.Sel.Name {
					case "Get":
						tempPath.Get = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Post":
						tempPath.Post = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Put":
						tempPath.Put = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Patch":
						tempPath.Patch = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Delete":
						tempPath.Delete = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Head":
						tempPath.Head = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					case "Options":
						tempPath.Options = &spec.Operation{OperationProps: spec.OperationProps{Parameters: params}}
					}
					paths[path] = tempPath
				} else {
					for _, v := range x.Args {
						switch arg := v.(type) {
						case *ast.Ident:
							xName := fun.X.(*ast.Ident).Name
							selName := fun.Sel.Name
							if arg.Name == routerName && imports[xName] != "" {
								traverseMod(imports[xName], selName)
							}
						}
					}
				}
			}
		}
		return true
	})
}

// getValidation takes a slice of arguments and a map of imports and returns a slice of
// spec.Parameter. It's used to parse the validation parameters of the
// fiber.ValidationWith middleware.
func getValidation(args []ast.Expr, mandragoraName string, imports ImportsType) []*spec.Parameter {
	var parameters []*spec.Parameter

	for _, arg := range args {
		if callExpr, ok := arg.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == mandragoraName && selectorExpr.Sel.Name == "WithValidation" {
					parameters = parseValidation(callExpr.Args[0].(*ast.CompositeLit), imports)
					break
				}
			}
		}
	}

	return parameters
}

// parseValidation takes a composite literal and a map of imports and returns a slice of
// spec.Parameter. It's used to parse the validation parameters of the
// fiber.ValidationWith middleware.
func parseValidation(compositeLit *ast.CompositeLit, imports ImportsType) []*spec.Parameter {
	var parameters []*spec.Parameter

	for _, elt := range compositeLit.Elts {
		if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			switch keyValueExpr.Key.(*ast.Ident).Name {
			case "Body":
				parameters = append(parameters, parseFields(keyValueExpr.Value.(*ast.CompositeLit), imports, "body"))
			case "Query":
				parameters = append(parameters, parseFields(keyValueExpr.Value.(*ast.CompositeLit), imports, "query"))
			case "Params":
				parameters = append(parameters, parseFields(keyValueExpr.Value.(*ast.CompositeLit), imports, "path"))
			}
		}
	}

	return parameters
}

// parseFields takes a composite literal, a map of imports, and a string representing the
// parameter location and returns a slice of spec.Parameter. It's used to parse the
// validation parameters of the fiber.ValidationWith middleware.
func parseFields(compositeLit *ast.CompositeLit, imports ImportsType, parIn string) *spec.Parameter {
	props := make(map[string]spec.Schema)
	requiredFields := []string{}

	switch t := compositeLit.Type.(type) {
	case *ast.Ident:
		for _, field := range t.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
			propName, prop, required := parseField(field)
			props[propName] = prop
			if required {
				requiredFields = append(requiredFields, propName)
			}
		}
	case *ast.SelectorExpr:
		typeFields := findTypeDeclaration(t, imports)
		for _, field := range typeFields {
			propName, prop, required := parseField(field)
			props[propName] = prop
			if required {
				requiredFields = append(requiredFields, propName)
			}
		}
	}

	return &spec.Parameter{
		ParamProps: spec.ParamProps{
			In:       parIn,
			Required: true,
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:       []string{"object"},
					Required:   requiredFields,
					Properties: props,
				},
			},
		},
	}
}

// parseField takes an ast.Field and returns a string representing the field name, a
// spec.Schema representing the field type, and a boolean indicating whether the field is
// required. It's used to parse the validation parameters of the fiber.ValidationWith
// middleware.
func parseField(field *ast.Field) (string, spec.Schema, bool) {
	tags := parseTags(field.Tag.Value)
	format := removeRequired(tags["validate"])
	return tags["json"], spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:   []string{parseFieldType(field.Type.(*ast.Ident).Name)},
			Format: format,
		},
	}, strings.Contains(field.Tag.Value, "required")
}

// removeRequired takes a string representing a tag and returns a string with the
// "required" keyword removed. It's used to parse the validation parameters of the
// fiber.ValidationWith middleware.
func removeRequired(tags string) string {
	tags = strings.ReplaceAll(tags, "required,", "")
	tags = strings.ReplaceAll(tags, ",required", "")
	tags = strings.ReplaceAll(tags, "required", "")
	tags = strings.ReplaceAll(tags, ",,", ",")
	return tags
}

// parseFieldType takes a string representing a field type and returns a string
// representing the corresponding OpenAPI type. It's used to parse the validation
// parameters of the fiber.ValidationWith middleware.
func parseFieldType(fieldType string) string {
	switch fieldType {
	case "string":
		return "string"
	case "uint":
		return "integer"
	case "int":
		return "integer"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}

// findTypeDeclaration takes an ast.SelectorExpr and a map of imports and returns a slice of
// ast.Field. It's used to parse the validation parameters of the
// fiber.ValidationWith middleware.
func findTypeDeclaration(sel *ast.SelectorExpr, imports ImportsType) []*ast.Field {
	var result []*ast.Field
	modPath := imports[sel.X.(*ast.Ident).Name]
	typeName := sel.Sel.Name
	pkgs, err := parser.ParseDir(fset, modPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Error parsing directory: %v\n", err)
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			// fiberName, mandragoraName, imports := trackImports(f)
			for _, decl := range f.Decls {
				switch decl.(type) {
				case *ast.GenDecl:
					if decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name == typeName {
						result = decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
					}
				}
			}
		}
	}
	return result
}

// traverseMod takes a string representing a module path and a string representing a
// function name and traverses the module to find the function and its routes. It's used
// to generate the OpenAPI documentation for the routes.
func traverseMod(modPath string, funcName string) {
	pkgs, err := parser.ParseDir(fset, modPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Error parsing directory: %v\n", err)
		return
	}

	for _, v := range pkgs {
		for _, f := range v.Files {
			fiberName, mandragoraName, imports := trackImports(f)
			for _, v := range f.Decls {
				switch x := v.(type) {
				case *ast.FuncDecl:
					if x.Name.Name == funcName {
						routerName := findRouterName(x, fiberName)
						findRoutes(x, routerName, mandragoraName, imports)
					}
				}
			}
		}
	}
}

// findRouterName takes an ast.FuncDecl and a string representing the fiber package name
// and returns a string representing the name of the fiber router. It's used to generate
// the OpenAPI documentation for the routes.
func findRouterName(routesFunc *ast.FuncDecl, fiberName string) (routerName string) {
	routerName = "app"
	for _, param := range routesFunc.Type.Params.List {
		switch t := param.Type.(type) {
		case *ast.StarExpr:
			if t.X.(*ast.SelectorExpr).X.(*ast.Ident).Name == fiberName && t.X.(*ast.SelectorExpr).Sel.Name == "App" {
				routerName = param.Names[0].Name
			}
		}
	}
	return
}

// isValidHttpMethod takes a string representing a HTTP method and returns a boolean
// indicating whether the method is valid. It's used to generate the OpenAPI
// documentation for the routes.
func isValidHttpMethod(str string) bool {
	for _, s := range HttpMethods {
		if s == str {
			return true
		}
	}
	return false
}

// parseTags takes a string representing a tag and returns a map of string to string.
// It's used to parse the validation parameters of the fiber.ValidationWith middleware.
func parseTags(tagString string) map[string]string {
	tags := make(map[string]string)

	// Split the tag string by spaces
	tagParts := strings.Split(strings.ReplaceAll(tagString, "`", ""), " ")

	// Iterate through the tag parts
	for _, part := range tagParts {
		// Split the tag part by the colon
		keyValue := strings.Split(part, ":")

		// If the tag part has a key and value
		if len(keyValue) == 2 {
			key := strings.Trim(keyValue[0], `"`)
			value := strings.Trim(keyValue[1], `"`)
			tags[key] = value
		}
	}

	return tags
}

// Fiber takes a string representing a main file path and returns a pointer to spec.Swagger.
// It's used to generate the OpenAPI documentation for the routes.
func Fiber(mainFilePath string) (*spec.Swagger, error) {
	absPath, err := filepath.Abs(mainFilePath)
	if err != nil {
		return nil, err
	}
	findRootDir(absPath)
	err = traverseMain(mainFilePath)
	if err != nil {
		return nil, err
	}
	return &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: paths,
			},
		},
	}, nil
}
