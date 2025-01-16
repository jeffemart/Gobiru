package models

// RouteInfo representa uma rota da API
type RouteInfo struct {
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Description string          `json:"description"`
	HandlerName string          `json:"handler_name"`
	Parameters  []Parameter     `json:"parameters"`
	QueryParams []Parameter     `json:"query_parameters"`
	Headers     []Parameter     `json:"headers"`
	Request     RequestBody     `json:"request_body"`
	Responses   []Response      `json:"responses"`
	Tags        []string        `json:"tags"`
	Auth        Authentication  `json:"authentication"`
	TimeMs      int             `json:"estimated_time_ms"`
	Permissions []string        `json:"permissions"`
	Version     string          `json:"api_version"`
	Deprecated  bool            `json:"deprecated"`
	RateLimit   RateLimitConfig `json:"rate_limit"`
	Notes       string          `json:"notes"`
}

// Parameter representa um parâmetro da rota
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type RequestBody struct {
	Type   string      `json:"type"`
	Schema interface{} `json:"schema"`
}

type Response struct {
	StatusCode  int         `json:"status_code"`
	Type        string      `json:"type"`
	Schema      interface{} `json:"schema"`
	Description string      `json:"description"`
}

type Authentication struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	TimeWindowSeconds int `json:"time_window_seconds"`
}
