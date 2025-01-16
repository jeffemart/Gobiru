package generator

import (
	"strings"

	"github.com/jeffemart/gobiru/internal/spec"
)

// Funções auxiliares para conversão
func convertSchema(schema *spec.Schema) map[string]interface{} {
	if schema == nil {
		return nil
	}

	result := map[string]interface{}{
		"type": schema.Type,
	}

	if schema.Properties != nil {
		props := make(map[string]interface{})
		for name, prop := range schema.Properties {
			props[name] = convertSchema(prop)
		}
		result["properties"] = props
	}

	if schema.Items != nil {
		result["items"] = convertSchema(schema.Items)
	}

	if schema.Required {
		result["required"] = true
	}

	return result
}

func convertParameters(params []*spec.Parameter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, p := range params {
		param := map[string]interface{}{
			"name":        p.Name,
			"in":          p.In,
			"required":    p.Required,
			"description": p.Description,
		}
		if p.Schema != nil {
			param["schema"] = convertSchema(p.Schema)
		}
		result = append(result, param)
	}
	return result
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

func convertRequestBody(body *spec.RequestBody) map[string]interface{} {
	if body == nil {
		return nil
	}
	return map[string]interface{}{
		"description": body.Description,
		"required":    body.Required,
		"content":     body.Content,
	}
}

func convertResponses(responses map[string]*spec.Response) map[string]interface{} {
	result := make(map[string]interface{})
	for code, resp := range responses {
		content := make(map[string]interface{})
		for mediaType, mt := range resp.Content {
			content[mediaType] = convertMediaType(mt)
		}

		result[code] = map[string]interface{}{
			"description": resp.Description,
			"content":     content,
		}
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
