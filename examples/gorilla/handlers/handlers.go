package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Categories  []string `json:"categories" validate:"required,min=1"`
	SKU         string   `json:"sku" validate:"required"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	StoreID     string    `json:"store_id"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Categories  []string  `json:"categories"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["storeId"]
	categoryID := vars["categoryId"]

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Simular criação do produto
	product := ProductResponse{
		ID:          uuid.New().String(),
		StoreID:     storeID,
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  req.Categories,
		SKU:         req.SKU,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
