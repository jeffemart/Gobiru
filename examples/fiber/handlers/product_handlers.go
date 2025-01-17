package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ListProducts retorna a lista de produtos
func ListProducts(c *fiber.Ctx) error {
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
	return c.JSON(products)
}

// GetProduct retorna um produto específico
func GetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	product := ProductResponse{
		ID:          productID,
		Name:        "Sample Product",
		Price:       99.99,
		Description: "Sample Description",
		CreatedAt:   time.Now(),
	}
	return c.JSON(product)
}

// CreateProduct cria um novo produto
func CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
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

	return c.Status(fiber.StatusCreated).JSON(product)
}
