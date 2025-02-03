package spec

import (
	"testing"
)

func TestSchemaCreation(t *testing.T) {
	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}

	if schema.Type != "object" {
		t.Errorf("Expected schema type to be 'object', got '%s'", schema.Type)
	}
	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}
}

func TestOperationCreation(t *testing.T) {
	op := &Operation{
		Path:    "/api/v1/test/{id}",
		Method:  "GET",
		Summary: "TestHandler lida com a rota de teste",
		Parameters: []*Parameter{
			{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "Path parameter: id",
				Schema: &Schema{
					Type: "string",
				},
			},
		},
		RequestBody: &RequestBody{
			Required: true,
			Content: map[string]*MediaType{
				"application/json": {
					Schema: &Schema{
						Type: "object",
						Properties: map[string]*Schema{
							"name":  {Type: "string"},
							"email": {Type: "string"},
						},
					},
				},
			},
		},
		Responses: map[string]*Response{
			"200": {
				Description: "Successful response",
				Content: map[string]*MediaType{
					"application/json": {
						Schema: &Schema{
							Type: "object",
						},
					},
				},
			},
		},
	}

	if op.Path != "/api/v1/test/{id}" {
		t.Errorf("Expected path %q, got %q", "/api/v1/test/{id}", op.Path)
	}
	if op.Method != "GET" {
		t.Errorf("Expected method %q, got %q", "GET", op.Method)
	}
}
