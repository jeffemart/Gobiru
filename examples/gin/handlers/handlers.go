package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
	Posts []Post `json:"posts,omitempty"`
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func ListUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	users := []User{
		{ID: "1", Name: "User 1", Email: "user1@example.com", Type: "user"},
		{ID: "2", Name: "User 2", Email: "user2@example.com", Type: "admin"},
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	user := User{
		ID:    id,
		Name:  "User " + id,
		Email: "user" + id + "@example.com",
		Type:  "user",
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func GetUserPost(c *gin.Context) {
	userId := c.Param("userId")
	postId := c.Param("postId")

	post := Post{
		ID:      postId,
		Title:   "Post " + postId + " from User " + userId,
		Content: "Content here...",
	}
	c.JSON(http.StatusOK, post)
}

func SearchUsers(c *gin.Context) {
	query := c.Query("query")
	sort := c.DefaultQuery("sort", "asc")
	userType := c.Query("type")

	users := []User{
		{ID: "1", Name: "Found User 1", Type: userType},
		{ID: "2", Name: "Found User 2", Type: userType},
	}

	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"sort":  sort,
		"type":  userType,
		"users": users,
	})
}
