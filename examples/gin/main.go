package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffemart/gobiru/examples/gin/routes"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := routes.SetupRoutes()
	r.Run(":3000")
}
