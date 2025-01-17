package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetProduct retorna um produto específico
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	product := ProductResponse{
		ID:          productID,
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       99.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct atualiza um produto existente
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

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

	product := ProductResponse{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  req.Categories,
		SKU:         req.SKU,
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct remove um produto
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// CreateOrder cria um novo pedido
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string    `json:"customer_id"`
		Items      []string  `json:"items"`
		Total      float64   `json:"total"`
		OrderDate  time.Time `json:"order_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	order := struct {
		ID         string    `json:"id"`
		CustomerID string    `json:"customer_id"`
		Items      []string  `json:"items"`
		Total      float64   `json:"total"`
		Status     string    `json:"status"`
		OrderDate  time.Time `json:"order_date"`
		CreatedAt  time.Time `json:"created_at"`
	}{
		ID:         uuid.New().String(),
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Total:      req.Total,
		Status:     "pending",
		OrderDate:  req.OrderDate,
		CreatedAt:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder retorna um pedido específico
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	order := struct {
		ID         string    `json:"id"`
		CustomerID string    `json:"customer_id"`
		Items      []string  `json:"items"`
		Total      float64   `json:"total"`
		Status     string    `json:"status"`
		OrderDate  time.Time `json:"order_date"`
		CreatedAt  time.Time `json:"created_at"`
	}{
		ID:         orderID,
		CustomerID: "customer123",
		Items:      []string{"item1", "item2"},
		Total:      199.99,
		Status:     "completed",
		OrderDate:  time.Now(),
		CreatedAt:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// UpdateOrderStatus atualiza o status de um pedido
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         orderID,
		"status":     req.Status,
		"updated_at": time.Now(),
	})
}

// Funções relacionadas a funcionários
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req Employee
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	req.ID = uuid.New().String()
	req.HireDate = time.Now()
	req.LastModified = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employeeId"]

	employee := Employee{
		ID:           employeeID,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Role:         "Developer",
		Department:   "Engineering",
		Status:       "Active",
		Skills:       []string{"Go", "Docker", "Kubernetes"},
		HireDate:     time.Now().AddDate(-1, 0, 0),
		LastModified: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employeeId"]

	var req Employee
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	req.ID = employeeID
	req.LastModified = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func UpdateEmployeeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employeeId"]

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            employeeID,
		"status":        req.Status,
		"last_modified": time.Now(),
	})
}

// CreateProduct cria um novo produto
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	product := ProductResponse{
		ID:          "new-product-123",
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Categories:  req.Categories,
		SKU:         req.SKU,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	respondWithJSON(w, http.StatusCreated, product)
}
