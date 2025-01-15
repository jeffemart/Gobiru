package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gobiru "github.com/jeffemart/Gobiru/app/gin"
	"github.com/jeffemart/Gobiru/app/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Define routes
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Create Gin analyzer
	analyzer := gobiru.NewAnalyzer()

	// Analyze routes
	err := analyzer.AnalyzeRoutes(router)
	if err != nil {
		log.Fatalf("Failed to analyze routes: %v", err)
	}

	// Export OpenAPI specification
	info := openapi.Info{
		Title:       "Gin API Example",
		Description: "Example API using Gin framework",
		Version:     "1.0.0",
	}

	spec, err := openapi.ConvertToOpenAPI(analyzer.GetRoutes(), info)
	if err != nil {
		log.Fatalf("Failed to convert to OpenAPI: %v", err)
	}

	err = openapi.ExportOpenAPI(spec, "openapi.json")
	if err != nil {
		log.Fatalf("Failed to export OpenAPI spec: %v", err)
	}

	log.Println("Documentation generated successfully!")
	log.Fatal(router.Run(":8080"))
}

func getUsers(c *gin.Context) {
	users := []User{
		{ID: "1", Name: "John"},
		{ID: "2", Name: "Jane"},
	}
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user := User{ID: id, Name: "John"}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
