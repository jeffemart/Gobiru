package spec

// Documentation representa a documentação completa da API
type Documentation struct {
	Operations []*Operation
}

// Operation representa uma operação/rota da API
type Operation struct {
	Path        string
	Method      string
	Summary     string
	Parameters  []*Parameter
	RequestBody *RequestBody
	Responses   map[string]*Response
}

// Parameter representa um parâmetro da operação
type Parameter struct {
	Name        string
	In          string // path, query, header, cookie
	Description string
	Required    bool
	Schema      *Schema
}

// RequestBody representa o corpo da requisição
type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]*MediaType
}

// Response representa uma resposta da operação
type Response struct {
	Description string
	Content     map[string]*MediaType
}

// MediaType representa o tipo de mídia do conteúdo
type MediaType struct {
	Schema *Schema
}

// Schema representa a estrutura de dados
type Schema struct {
	Type       string             `json:"type"`
	Properties map[string]*Schema `json:"properties,omitempty"`
	Items      *Schema            `json:"items,omitempty"`
	Required   bool               `json:"required,omitempty"`
}
