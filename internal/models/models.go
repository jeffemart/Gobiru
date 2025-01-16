package models

// Config contém configurações gerais da aplicação
type Config struct {
	RateLimit RateLimitConfig
}

// RateLimitConfig contém configurações de rate limiting
type RateLimitConfig struct {
	RequestsPerMinute int
	TimeWindowSeconds int
}

// Outros tipos específicos que não fazem parte da especificação OpenAPI
