package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetUser retorna os detalhes de um usuário
func GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	user := UserResponse{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}
	return c.JSON(user)
}

// CreateUser cria um novo usuário
func CreateUser(c *fiber.Ctx) error {
	var req struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	user := UserResponse{
		ID:        "new-user-123",
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser atualiza um usuário existente
func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	user := UserResponse{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	return c.JSON(user)
}
