package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Grupo de rotas da API
	api := r.Group("/api")
	{
		// Rotas de produtos
		products := api.Group("/products")
		{
			products.GET("", listProducts)
			products.GET("/:id", getProduct)
			products.POST("", createProduct)
			products.PUT("/:id", updateProduct)
			products.DELETE("/:id", deleteProduct)
		}
	}

	return r
}

func listProducts(c *gin.Context) {
	products := []Product{
		{ID: "1", Name: "Produto 1", Price: 19.99},
		{ID: "2", Name: "Produto 2", Price: 29.99},
	}
	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	product := Product{ID: id, Name: "Produto Exemplo", Price: 99.99}
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
