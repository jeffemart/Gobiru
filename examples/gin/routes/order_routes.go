package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/handlers"
)

// Rotas de pedidos e transações
func SetupOrderRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	orders := api.Group("/orders")
	{
		orders.POST("", handlers.CreateOrder)
		orders.GET("/:id", handlers.GetOrder)
		orders.PATCH("/:id/status", handlers.UpdateOrderStatus)
	}
}
