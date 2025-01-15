package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Product represents a product in our API
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func main() {
	router := gin.Default()

	// Product routes
	v1 := router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", listProducts)
			products.POST("", createProduct)
			products.GET("/:id", getProduct)
			products.PUT("/:id", updateProduct)
			products.DELETE("/:id", deleteProduct)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Fatal(router.Run(":8080"))
}

func listProducts(c *gin.Context) {
	products := []Product{
		{ID: "1", Name: "Product 1", Description: "Description 1", Price: 19.99},
		{ID: "2", Name: "Product 2", Description: "Description 2", Price: 29.99},
	}
	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	product := Product{
		ID:          id,
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       99.99,
	}
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simulate product creation
	product.ID = "new-id"
	c.JSON(http.StatusCreated, product)
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = id
	c.JSON(http.StatusOK, product)
}

func deleteProduct(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
