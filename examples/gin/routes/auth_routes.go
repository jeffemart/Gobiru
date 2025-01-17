package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/handlers"
)

// Renomeado de SetupAuthRoutes para SetupUserRoutes para corresponder ao main.go
func SetupUserRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	users := api.Group("/users")
	{
		users.GET("/:id", handlers.GetUser)
		users.POST("", handlers.CreateUser)
		users.PUT("/:id", handlers.UpdateUser)
	}
}
