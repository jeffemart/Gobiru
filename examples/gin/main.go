package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/routes"
)

func main() {
	r := gin.Default()

	// Setup all routes
	routes.SetupUserRoutes(r)
	routes.SetupProductRoutes(r)
	routes.SetupOrderRoutes(r)

	r.Run(":8080")
}
