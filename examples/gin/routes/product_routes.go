package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/handlers"
)

// Rotas de produtos e catálogo
func SetupProductRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	products := api.Group("/products")
	{
		products.GET("", handlers.ListProducts)
		products.GET("/:id", handlers.GetProduct)
		products.POST("", handlers.CreateProduct)
	}
}
