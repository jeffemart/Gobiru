package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUser retorna os detalhes de um usuário
func GetUser(c *gin.Context) {
	userID := c.Param("id")
	user := UserResponse{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser cria um novo usuário
func CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	user := UserResponse{
		ID:        "new-user-123",
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser atualiza um usuário existente
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	user := UserResponse{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, user)
}
