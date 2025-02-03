package analyzer

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestExtractSummaryFromComments(t *testing.T) {
	src := `
		package main

		// This is a test function
		func TestFunc() {}
	`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	summary := extractSummaryFromComments(node.Decls[0])
	expected := "This is a test function"
	if summary != expected {
		t.Errorf("Expected summary %q, got %q", expected, summary)
	}
}

func TestExtractRequestBody(t *testing.T) {
	src := `
		package main

		func TestFunc(c *fiber.Ctx) error {
			var req struct {
				Name  string ` + "`json:\"name\"`" + `
				Email string ` + "`json:\"email\"`" + `
			}
			c.BodyParser(&req)
			return nil
		}
	`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	reqBody := extractRequestBody(node.Decls[0], "test.go")
	if reqBody == nil {
		t.Error("Expected request body to be extracted, got nil")
	}
}

type Request struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
