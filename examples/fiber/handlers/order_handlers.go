package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateOrder cria um novo pedido
func CreateOrder(c *fiber.Ctx) error {
	var req struct {
		CustomerID string    `json:"customerId" validate:"required"`
		Items      []string  `json:"items" validate:"required"`
		Total      float64   `json:"total" validate:"required"`
		OrderDate  time.Time `json:"orderDate"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	order := OrderResponse{
		ID:         "new-order-123",
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Total:      req.Total,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetOrder retorna um pedido específico
func GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	order := OrderResponse{
		ID:         orderID,
		CustomerID: "customer-123",
		Items:      []string{"item1", "item2"},
		Total:      199.99,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}
	return c.JSON(order)
}

// UpdateOrderStatus atualiza o status de um pedido
func UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var req struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":         orderID,
		"status":     req.Status,
		"updated_at": time.Now(),
	})
}
