package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateOrder cria um novo pedido
func CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string    `json:"customer_id" binding:"required"`
		Items      []string  `json:"items" binding:"required"`
		Total      float64   `json:"total" binding:"required"`
		OrderDate  time.Time `json:"order_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	order := OrderResponse{
		ID:         "new-order-123",
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Total:      req.Total,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder retorna um pedido específico
func GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order := OrderResponse{
		ID:         orderID,
		CustomerID: "customer-123",
		Items:      []string{"item1", "item2"},
		Total:      199.99,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus atualiza o status de um pedido
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         orderID,
		"status":     req.Status,
		"updated_at": time.Now(),
	})
}
