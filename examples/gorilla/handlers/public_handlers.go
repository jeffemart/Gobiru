package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// Login autentica um usuário
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Simula autenticação
	response := LoginResponse{
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Register registra um novo usuário
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Simula criação de usuário
	user := UserResponse{
		ID:        "new-user-123",
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// ForgotPassword inicia o processo de recuperação de senha
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := decodeJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password reset instructions sent to your email",
	})
}

// ResetPassword altera a senha do usuário
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := decodeJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password successfully reset",
	})
}

// HealthCheck verifica a saúde da API
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// APIStatus retorna o status atual da API
func APIStatus(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status    string    `json:"status"`
		Version   string    `json:"version"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Status:    "operational",
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	respondWithJSON(w, http.StatusOK, status)
}

// GetAPIDocumentation retorna a documentação da API
func GetAPIDocumentation(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"version": "1.0.0",
		"title":   "API Documentation",
		"baseUrl": "/api/v1",
		"paths":   map[string]interface{}{},
	})
}

// GetSwaggerUI retorna a interface Swagger
func GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><h1>Swagger UI</h1></body></html>"))
}

// Funções auxiliares
func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, errorCode, message, details string) {
	respondWithJSON(w, code, ErrorResponse{
		Code:    errorCode,
		Message: message,
		Details: details,
	})
}
