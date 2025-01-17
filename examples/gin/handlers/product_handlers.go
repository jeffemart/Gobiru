package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListProducts retorna a lista de produtos
func ListProducts(c *gin.Context) {
	products := []ProductResponse{
		{
			ID:          "1",
			Name:        "Product 1",
			Price:       99.99,
			Description: "Description 1",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Product 2",
			Price:       149.99,
			Description: "Description 2",
			CreatedAt:   time.Now(),
		},
	}
	c.JSON(http.StatusOK, products)
}

// GetProduct retorna um produto específico
func GetProduct(c *gin.Context) {
	productID := c.Param("id")
	product := ProductResponse{
		ID:          productID,
		Name:        "Sample Product",
		Price:       99.99,
		Description: "Sample Description",
		CreatedAt:   time.Now(),
	}
	c.JSON(http.StatusOK, product)
}

// CreateProduct cria um novo produto
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
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

	c.JSON(http.StatusCreated, product)
}
