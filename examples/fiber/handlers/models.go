package handlers

import "time"

// Respostas de Usuário
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Respostas de Produto
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories,omitempty"`
	SKU         string    `json:"sku,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Price       float64  `json:"price" validate:"required"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	SKU         string   `json:"sku"`
}

// Respostas de Pedido
type OrderResponse struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Items      []string  `json:"items"`
	Total      float64   `json:"total"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// Resposta de Erro
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
