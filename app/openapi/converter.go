package openapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/jeffemart/Gobiru/app/models"
)

// OpenAPISpec represents the root OpenAPI specification
type OpenAPISpec struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // path, query, header
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required"`
	Schema      *Schema `json:"schema"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema *Schema `json:"schema"`
}

type Schema struct {
	Type       string            `json:"type,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Required   []string          `json:"required,omitempty"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

// ConvertToOpenAPI converts Gobiru route information to OpenAPI specification
func ConvertToOpenAPI(routes []models.RouteInfo, info Info) (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info:    info,
		Paths:   make(map[string]PathItem),
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"bearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		},
	}

	for _, route := range routes {
		pathItem, exists := spec.Paths[route.Path]
		if !exists {
			pathItem = PathItem{}
		}

		operation := &Operation{
			Tags:        route.Tags,
			Summary:     route.Description,
			Description: route.Notes,
			Deprecated:  route.Deprecated,
			Responses:   make(map[string]Response),
		}

		// Convert parameters
		operation.Parameters = convertParameters(route)

		// Convert request body
		if route.RequestBody.Type != "" {
			operation.RequestBody = convertRequestBody(route.RequestBody)
		}

		// Convert responses
		for _, resp := range route.Responses {
			statusCode := fmt.Sprintf("%d", resp.StatusCode)
			operation.Responses[statusCode] = Response{
				Description: resp.Description,
				Content: map[string]MediaType{
					"application/json": {
						Schema: convertToSchema(resp.Example),
					},
				},
			}
		}

		// Add security if required
		if route.Authentication.Required {
			operation.Security = []map[string][]string{
				{"bearerAuth": {}},
			}
		}

		// Set the operation based on the HTTP method
		switch strings.ToUpper(route.Method) {
		case "GET":
			pathItem.Get = operation
		case "POST":
			pathItem.Post = operation
		case "PUT":
			pathItem.Put = operation
		case "DELETE":
			pathItem.Delete = operation
		case "PATCH":
			pathItem.Patch = operation
		}

		spec.Paths[route.Path] = pathItem
	}

	return spec, nil
}

func convertParameters(route models.RouteInfo) []Parameter {
	var parameters []Parameter

	// Path parameters
	for _, param := range route.Parameters {
		parameters = append(parameters, Parameter{
			Name:        param.Name,
			In:          "path",
			Description: param.Description,
			Required:    param.Required,
			Schema: &Schema{
				Type: param.Type,
			},
		})
	}

	// Query parameters
	for _, param := range route.QueryParameters {
		parameters = append(parameters, Parameter{
			Name:        param.Name,
			In:          "query",
			Description: param.Description,
			Required:    param.Required,
			Schema: &Schema{
				Type: param.Type,
			},
		})
	}

	// Headers
	for _, header := range route.Headers {
		parameters = append(parameters, Parameter{
			Name:        header.Name,
			In:          "header",
			Description: header.Description,
			Required:    header.Required,
			Schema: &Schema{
				Type: header.Type,
			},
		})
	}

	return parameters
}

func convertRequestBody(reqBody models.RequestBody) *RequestBody {
	return &RequestBody{
		Required: true,
		Content: map[string]MediaType{
			reqBody.Type: {
				Schema: convertToSchema(reqBody.Schema),
			},
		},
	}
}

func convertToSchema(data interface{}) *Schema {
	if data == nil {
		return &Schema{Type: "object"}
	}

	schema := &Schema{}
	switch v := data.(type) {
	case map[string]interface{}:
		schema.Type = "object"
		schema.Properties = make(map[string]Schema)
		for key, value := range v {
			schema.Properties[key] = *convertToSchema(value)
		}
	case []interface{}:
		schema.Type = "array"
		if len(v) > 0 {
			schema.Items = convertToSchema(v[0])
		}
	case string:
		schema.Type = "string"
	case float64:
		schema.Type = "number"
	case bool:
		schema.Type = "boolean"
	case int:
		schema.Type = "integer"
	}
	return schema
}

// ExportOpenAPI exports the OpenAPI specification to a file
func ExportOpenAPI(spec *OpenAPISpec, outputPath string) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI spec: %v", err)
	}

	ext := filepath.Ext(outputPath)
	if ext == "" {
		outputPath += ".json"
	}

	err = ioutil.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
