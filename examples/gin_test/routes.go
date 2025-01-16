package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Post representa um post do blog
type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// SetupRouter configura as rotas da API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Grupo de rotas v1
	v1 := r.Group("/api/v1")
	{
		posts := v1.Group("/posts")
		{
			posts.GET("", listPosts)
			posts.GET("/:id", getPost)
			posts.POST("", createPost)
			posts.PUT("/:id", updatePost)
			posts.DELETE("/:id", deletePost)

			// Rotas adicionais
			posts.GET("/tags/:tag", getPostsByTag)
			posts.POST("/:id/publish", publishPost)
		}
	}

	// Rota de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func listPosts(c *gin.Context) {
	posts := []Post{
		{
			ID:      "1",
			Title:   "Introdução ao Gin",
			Content: "Gin é um framework web em Go...",
			Tags:    []string{"go", "gin", "web"},
		},
		{
			ID:      "2",
			Title:   "Documentação com Gobiru",
			Content: "Gobiru é uma ferramenta...",
			Tags:    []string{"go", "docs", "api"},
		},
	}
	c.JSON(http.StatusOK, posts)
}

func getPost(c *gin.Context) {
	id := c.Param("id")
	post := Post{
		ID:      id,
		Title:   "Post de Exemplo",
		Content: "Conteúdo do post de exemplo...",
		Tags:    []string{"exemplo", "teste"},
	}
	c.JSON(http.StatusOK, post)
}

func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.ID = "new-id"
	c.JSON(http.StatusCreated, post)
}

func updatePost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.ID = id
	c.JSON(http.StatusOK, post)
}

func deletePost(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func getPostsByTag(c *gin.Context) {
	tag := c.Param("tag")
	posts := []Post{
		{
			ID:      "1",
			Title:   "Post com tag " + tag,
			Content: "Conteúdo...",
			Tags:    []string{tag},
		},
	}
	c.JSON(http.StatusOK, posts)
}

func publishPost(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"status":  "published",
		"message": "Post publicado com sucesso",
	})
}
