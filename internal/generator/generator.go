package generator

import (
	"github.com/jeffemart/gobiru/internal/spec"
)

// Config contém as configurações para geração da documentação
type Config struct {
	OutputFile  string
	Title       string
	Description string
	Version     string
}

// Generator define a interface para geração de documentação
type Generator interface {
	Generate(doc *spec.Documentation, config Config) error
}

// New cria um novo gerador baseado no formato
func New(format string) Generator {
	switch format {
	case "json":
		return NewJSONGenerator()
	case "openapi":
		return NewOpenAPIGenerator()
	default:
		return NewJSONGenerator()
	}
}
