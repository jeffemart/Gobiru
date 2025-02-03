// internal/models/models_test.go
package models

import (
	"testing"
)

type YourModel struct {
	Field1 string
	Field2 string
}

func TestModelCreation(t *testing.T) {
	model := &YourModel{
		Field1: "value1",
		Field2: "value2",
	}

	if model.Field1 != "value1" {
		t.Errorf("Expected Field1 to be 'value1', got '%s'", model.Field1)
	}
	if model.Field2 != "value2" {
		t.Errorf("Expected Field2 to be 'value2', got '%s'", model.Field2)
	}
}