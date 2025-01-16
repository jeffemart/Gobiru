package main

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin application routes
func SetupRouter() *gin.Engine {
	router := gin.New()

	// API routes group
	api := router.Group("/api")

	// Users routes
	users := api.Group("/users")
	users.GET("/", listUsers)
	users.GET("/:id", getUser)
	users.POST("/", createUser)

	// Health check
	router.GET("/health", healthCheck)

	return router
}

func listUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"users": []string{"user1", "user2"},
	})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"id":   id,
		"name": "Example User",
	})
}

func createUser(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "User created",
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
