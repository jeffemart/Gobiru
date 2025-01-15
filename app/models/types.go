package models

// RouteInfo represents the complete information about an API route
type RouteInfo struct {
	Method          string         `json:"method"`
	Path            string         `json:"path"`
	Description     string         `json:"description"`
	HandlerName     string         `json:"handler_name"`
	Parameters      []Parameter    `json:"parameters"`
	QueryParameters []Parameter    `json:"query_parameters"`
	Headers         []Parameter    `json:"headers"`
	RequestBody     RequestBody    `json:"request_body"`
	Responses       []Response     `json:"responses"`
	Tags            []string       `json:"tags"`
	Authentication  Authentication `json:"authentication"`
	EstimatedTimeMs int            `json:"estimated_time_ms"`
	Permissions     []string       `json:"permissions"`
	APIVersion      string         `json:"api_version"`
	Deprecated      bool           `json:"deprecated"`
	RateLimit       RateLimit      `json:"rate_limit"`
	Notes           string         `json:"notes"`
}

// Parameter represents a route parameter, query parameter, or header
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// RequestBody represents the request body schema
type RequestBody struct {
	Type   string                 `json:"type"`
	Schema map[string]interface{} `json:"schema"`
}

// Response represents a possible API response
type Response struct {
	StatusCode  int         `json:"status_code"`
	Description string      `json:"description"`
	Example     interface{} `json:"example"`
}

// Authentication represents authentication requirements
type Authentication struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// RateLimit represents rate limiting configuration
type RateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	TimeWindowSeconds int `json:"time_window_seconds"`
}
