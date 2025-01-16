package main

import (
	"github.com/gofiber/fiber/v2"
)

// Product representa um produto no sistema
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

// SetupRouter configures the Fiber application routes
func SetupRouter() *fiber.App {
	app := fiber.New()

	// API routes group
	api := app.Group("/api")

	// Products routes
	products := api.Group("/products")
	products.Get("/", listProducts)                            // Lista todos os produtos
	products.Post("/", createProduct)                          // Cria um novo produto
	products.Get("/:id", getProduct)                           // Obtém um produto específico
	products.Put("/:id", updateProduct)                        // Atualiza um produto
	products.Delete("/:id", deleteProduct)                     // Remove um produto
	products.Get("/category/:category", getProductsByCategory) // Produtos por categoria
	products.Post("/:id/reviews", addProductReview)            // Adiciona uma review
	products.Get("/search", searchProducts)                    // Busca produtos com query params

	// Bulk operations
	products.Post("/bulk", bulkCreateProducts)   // Cria múltiplos produtos
	products.Delete("/bulk", bulkDeleteProducts) // Remove múltiplos produtos

	// Health check
	app.Get("/health", healthCheck)

	return app
}

func listProducts(c *fiber.Ctx) error {
	return c.JSON([]Product{
		{ID: "1", Name: "Produto 1", Price: 19.99},
		{ID: "2", Name: "Produto 2", Price: 29.99},
	})
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	return c.Status(201).JSON(product)
}

func getProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(Product{
		ID:          id,
		Name:        "Exemplo de Produto",
		Price:       99.99,
		Description: "Descrição detalhada do produto",
		Categories:  []string{"eletrônicos", "gadgets"},
	})
}

func updateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	product.ID = id
	return c.JSON(product)
}

func deleteProduct(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func getProductsByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	return c.JSON([]Product{
		{ID: "1", Name: "Produto da Categoria", Price: 19.99, Categories: []string{category}},
	})
}

func addProductReview(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"product_id": id,
		"message":    "Review added successfully",
	})
}

func searchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	category := c.Query("category")

	return c.JSON(fiber.Map{
		"query":     query,
		"min_price": minPrice,
		"max_price": maxPrice,
		"category":  category,
		"results":   []Product{},
	})
}

func bulkCreateProducts(c *fiber.Ctx) error {
	var products []Product
	if err := c.BodyParser(&products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	return c.Status(201).JSON(products)
}

func bulkDeleteProducts(c *fiber.Ctx) error {
	var ids []string
	if err := c.BodyParser(&ids); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	return c.SendStatus(204)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
