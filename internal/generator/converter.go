package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

// Funções de conversão compartilhadas
func convertSchema(schema *spec.Schema) map[string]interface{} {
	if schema == nil {
		return nil
	}

	result := make(map[string]interface{})
	if schema.Type != "" {
		result["type"] = schema.Type
	}
	if schema.Format != "" {
		result["format"] = schema.Format
	}
	if len(schema.Properties) > 0 {
		props := make(map[string]interface{})
		for name, prop := range schema.Properties {
			props[name] = convertSchema(prop)
		}
		result["properties"] = props
	}
	if schema.Items != nil {
		result["items"] = convertSchema(schema.Items)
	}
	return result
}

func convertParameters(params []*spec.Parameter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, p := range params {
		if p.Name == "" {
			continue
		}
		param := map[string]interface{}{
			"name":        p.Name,
			"in":          p.In,
			"required":    p.Required,
			"description": p.Description,
			"schema":      convertSchema(p.Schema),
		}
		result = append(result, param)
	}
	return result
}

func convertRequestBody(body *spec.RequestBody) map[string]interface{} {
	if body == nil {
		return nil
	}
	return map[string]interface{}{
		"required": body.Required,
		"content":  convertContent(body.Content),
	}
}

func convertContent(content map[string]*spec.MediaType) map[string]interface{} {
	result := make(map[string]interface{})
	for mediaType, mt := range content {
		result[mediaType] = map[string]interface{}{
			"schema": convertSchema(mt.Schema),
		}
	}
	return result
}

func convertResponses(responses map[string]*spec.Response) map[string]interface{} {
	result := make(map[string]interface{})
	for code, resp := range responses {
		result[code] = map[string]interface{}{
			"description": resp.Description,
			"content":     convertContent(resp.Content),
		}
	}
	return result
}

func extractTags(path string) []string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	tags := make([]string, 0)
	for _, segment := range segments {
		if segment != "api" && segment != "v1" && !strings.HasPrefix(segment, "{") {
			tags = append(tags, segment)
			break
		}
	}
	if len(tags) == 0 {
		tags = append(tags, "general")
	}
	return tags
}

func extractHandlerName(summary string) string {
	if idx := strings.Index(summary, " "); idx > 0 {
		return summary[:idx]
	}
	return summary
}

func convertMediaType(mediaType *spec.MediaType) map[string]interface{} {
	if mediaType == nil {
		return nil
	}

	result := make(map[string]interface{})
	if mediaType.Schema != nil {
		result["schema"] = convertSchema(mediaType.Schema)
	}
	return result
}

func convertPaths(operations []*spec.Operation) map[string]interface{} {
	paths := make(map[string]interface{})

	for _, op := range operations {
		if _, exists := paths[op.Path]; !exists {
			paths[op.Path] = make(map[string]interface{})
		}

		pathItem := paths[op.Path].(map[string]interface{})
		method := strings.ToLower(op.Method)

		pathItem[method] = map[string]interface{}{
			"summary":     op.Summary,
			"parameters":  convertParameters(op.Parameters),
			"requestBody": convertRequestBody(op.RequestBody),
			"responses":   convertResponses(op.Responses),
		}
	}

	return paths
}

// writeJSON escreve os dados em formato JSON no arquivo especificado
func writeJSON(filename string, data interface{}) error {
	// Criar diretório se não existir
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Converter para JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Escrever arquivo
	return os.WriteFile(filename, jsonData, 0644)
}
