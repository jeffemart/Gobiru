package handlers

import (
	"github.com/gofiber/fiber/v2"
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

func ListUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	users := []User{
		{ID: "1", Name: "User 1", Email: "user1@example.com", Type: "user"},
		{ID: "2", Name: "User 2", Email: "user2@example.com", Type: "admin"},
	}
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := User{
		ID:    id,
		Name:  "User " + id,
		Email: "user" + id + "@example.com",
		Type:  "user",
	}
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user.ID = id
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func GetUserPost(c *fiber.Ctx) error {
	userId := c.Params("userId")
	postId := c.Params("postId")

	post := Post{
		ID:      postId,
		Title:   "Post " + postId + " from User " + userId,
		Content: "Content here...",
	}
	return c.JSON(post)
}

func SearchUsers(c *fiber.Ctx) error {
	query := c.Query("query")
	sort := c.Query("sort", "asc")
	userType := c.Query("type")

	users := []User{
		{ID: "1", Name: "Found User 1", Type: userType},
		{ID: "2", Name: "Found User 2", Type: userType},
	}

	return c.JSON(fiber.Map{
		"query": query,
		"sort":  sort,
		"type":  userType,
		"users": users,
	})
}
