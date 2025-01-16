package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID   string       `json:"customer_id" validate:"required"`
	Items        []OrderItem  `json:"items" validate:"required,min=1"`
	ShippingInfo ShippingInfo `json:"shipping_info" validate:"required"`
	PaymentInfo  PaymentInfo  `json:"payment_info" validate:"required"`
	Notes        string       `json:"notes"`
	Status       string       `json:"status" validate:"required,oneof=pending processing confirmed"`
}

type OrderItem struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type ShippingInfo struct {
	Address    string  `json:"address" validate:"required"`
	City       string  `json:"city" validate:"required"`
	PostalCode string  `json:"postal_code" validate:"required"`
	Country    string  `json:"country" validate:"required"`
	Cost       float64 `json:"cost" validate:"required,gte=0"`
	Method     string  `json:"method" validate:"required"`
}

type PaymentInfo struct {
	Method string  `json:"method" validate:"required,oneof=credit_card debit_card bank_transfer"`
	Total  float64 `json:"total" validate:"required,gt=0"`
	Status string  `json:"status" validate:"required"`
}

type OrderResponse struct {
	ID           string       `json:"id"`
	StoreID      string       `json:"store_id"`
	CustomerID   string       `json:"customer_id"`
	Items        []OrderItem  `json:"items"`
	ShippingInfo ShippingInfo `json:"shipping_info"`
	PaymentInfo  PaymentInfo  `json:"payment_info"`
	Notes        string       `json:"notes"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func CreateOrder(c *gin.Context) {
	storeID := c.Param("storeId")
	customerID := c.Param("customerId")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Simular criação do pedido
	order := OrderResponse{
		ID:           uuid.New().String(),
		StoreID:      storeID,    // Usando o parâmetro da rota
		CustomerID:   customerID, // Usando o parâmetro da rota
		Items:        req.Items,
		ShippingInfo: req.ShippingInfo,
		PaymentInfo:  req.PaymentInfo,
		Notes:        req.Notes,
		Status:       req.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	c.JSON(http.StatusCreated, order)
}
