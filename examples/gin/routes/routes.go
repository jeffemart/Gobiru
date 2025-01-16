package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/handlers"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")

	// Stores and Customers
	stores := v1.Group("/stores/:storeId")
	customers := stores.Group("/customers/:customerId")

	// Orders
	orders := customers.Group("/orders")
	orders.POST("/", handlers.CreateOrder)

	return r
}
