package models

// RouteInfo contém informações sobre uma rota da API
type RouteInfo struct {
	Path        string
	Method      string
	HandlerName string
	Version     string
	Description string
	Parameters  []Parameter
	QueryParams []Parameter
	Request     RequestBody
	Responses   []Response
}

// Parameter representa um parâmetro da rota
type Parameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

// RequestBody representa o corpo da requisição
type RequestBody struct {
	Type   string
	Schema interface{}
}

// Response representa uma resposta da API
type Response struct {
	StatusCode  int
	Type        string
	Description string
	Schema      interface{}
}
