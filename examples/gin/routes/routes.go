package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/users", handlers.ListUsers)
	r.GET("/peste/:id", handlers.GetUser)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.GET("/users/:userId/posts/:postId", handlers.GetUserPost)
	r.GET("/search", handlers.SearchUsers)

	return r
}
