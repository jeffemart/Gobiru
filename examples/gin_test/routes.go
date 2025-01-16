package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// User representa um usuário do sistema
type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	Active   bool     `json:"active"`
	Skills   []string `json:"skills"`
	Password string   `json:"-"` // não será exibido no JSON
}

// SetupRouter configura as rotas da API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Grupo de rotas v1
	v1 := r.Group("/api/v1")
	{
		// Autenticação
		auth := v1.Group("/auth")
		{
			auth.POST("/login", login)
			auth.POST("/register", register)
			auth.POST("/forgot-password", forgotPassword)
			auth.POST("/reset-password", resetPassword)
		}

		// Usuários
		users := v1.Group("/users")
		{
			users.GET("", listUsers)
			users.GET("/:id", getUser)
			users.POST("", createUser)
			users.PUT("/:id", updateUser)
			users.DELETE("/:id", deleteUser)

			// Rotas adicionais
			users.GET("/search", searchUsers)
			users.PUT("/:id/activate", activateUser)
			users.PUT("/:id/deactivate", deactivateUser)
			users.GET("/:id/skills", getUserSkills)
			users.POST("/:id/skills", addUserSkill)
			users.DELETE("/:id/skills/:skill", removeUserSkill)
		}
	}

	// Health check
	r.GET("/health", healthCheck)

	return r
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": "jwt-token-example",
	})
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = "new-id"
	c.JSON(http.StatusCreated, user)
}

func forgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
	})
}

func resetPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

func listUsers(c *gin.Context) {
	users := []User{
		{ID: "1", Name: "John Doe", Email: "john@example.com", Role: "admin", Active: true},
		{ID: "2", Name: "Jane Doe", Email: "jane@example.com", Role: "user", Active: true},
	}
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user := User{
		ID:     id,
		Name:   "John Doe",
		Email:  "john@example.com",
		Role:   "admin",
		Active: true,
		Skills: []string{"golang", "docker", "kubernetes"},
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = "new-id"
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

func searchUsers(c *gin.Context) {
	query := c.Query("q")
	users := []User{
		{ID: "1", Name: "John Doe", Email: "john@example.com", Role: "admin", Active: true},
	}
	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"users": users,
	})
}

func activateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"active":  true,
		"message": "User activated successfully",
	})
}

func deactivateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"active":  false,
		"message": "User deactivated successfully",
	})
}

func getUserSkills(c *gin.Context) {
	id := c.Param("id")
	skills := []string{"golang", "docker", "kubernetes"}
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"skills": skills,
	})
}

func addUserSkill(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Skill string `json:"skill"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"skill":   body.Skill,
		"message": "Skill added successfully",
	})
}

func removeUserSkill(c *gin.Context) {
	id := c.Param("id")
	skill := c.Param("skill")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"skill":   skill,
		"message": "Skill removed successfully",
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
