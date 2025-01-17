# Gobiru - Gerador de Documentação de API

![Gobiru Logo](https://res.cloudinary.com/dx70wyorg/image/upload/v1736953035/photo_2025-01-15_11-40-32_esheqe.jpg)

Gobiru é uma ferramenta para gerar automaticamente documentação de APIs Go, suportando os frameworks Gin, Gorilla Mux e Fiber.

## Instalação

```bash
go install github.com/jeffemart/gobiru/cmd/gobiru@latest
```

## Uso

O Gobiru analisa seu código fonte e gera documentação em formato JSON e OpenAPI (Swagger).

### Parâmetros

- `-framework`: Framework usado (gin, mux, fiber)
- `-main`: Arquivo principal da aplicação
- `-output`: Caminho para o arquivo JSON de saída
- `-openapi`: Caminho para o arquivo OpenAPI de saída
- `-title`: Título da documentação
- `-description`: Descrição da API
- `-version`: Versão da API

## Exemplos por Framework

### Gin
```bash
./gobiru -framework gin \
       -main examples/gin/main.go \
       -output examples/gin/docs/routes.json \
       -openapi examples/gin/docs/openapi.json \
       -title "API Gin" \
       -description "API de exemplo usando Gin" \
       -version "1.0.0"
```

Exemplo de rota com Gin:
```go
func SetupAuthRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    users := api.Group("/users")
    {
        users.GET("/:id", handlers.GetUser)
        users.POST("", handlers.CreateUser)
        users.PUT("/:id", handlers.UpdateUser)
    }
}
```

### Gorilla Mux
```bash
./gobiru -framework mux \
       -main examples/gorilla/main.go \
       -output examples/gorilla/docs/routes.json \
       -openapi examples/gorilla/docs/openapi.json \
       -title "API Mux" \
       -description "API de exemplo usando Gorilla Mux" \
       -version "1.0.0"
```

Exemplo de rota com Mux:
```go
func SetupPublicRoutes(r *mux.Router) {
    api := r.PathPrefix("/api/v1").Subrouter()
    
    // Autenticação
    auth := api.PathPrefix("/auth").Subrouter()
    auth.HandleFunc("/login", handlers.Login).Methods("POST")
    auth.HandleFunc("/register", handlers.Register).Methods("POST")
}
```

### Fiber
```bash
./gobiru -framework fiber \
       -main examples/fiber/main.go \
       -output examples/fiber/docs/routes.json \
       -openapi examples/fiber/docs/openapi.json \
       -title "API Fiber" \
       -description "API de exemplo usando Fiber" \
       -version "1.0.0"
```

Exemplo de rota com Fiber:
```go
func SetupRoutes(app *fiber.App) {
    api := app.Group("/api/v1")
    
    // Produtos
    products := api.Group("/products")
    products.Get("/", handlers.ListProducts)
    products.Post("/", handlers.CreateProduct)
    products.Get("/:id", handlers.GetProduct)
}
```

## Estrutura do Projeto

O Gobiru funciona analisando a estrutura do seu projeto. Ele:

1. Começa pelo arquivo main.go
2. Segue os imports para encontrar arquivos de rotas e handlers
3. Analisa a definição das rotas
4. Gera documentação completa

### Exemplo de Estrutura

```
seu-projeto/
├── main.go
├── routes/
│   ├── auth_routes.go    # Rotas de autenticação
│   ├── product_routes.go # Rotas de produtos
│   └── order_routes.go   # Rotas de pedidos
└── handlers/
    ├── auth_handlers.go    # Handlers de autenticação
    ├── product_handlers.go # Handlers de produtos
    ├── order_handlers.go   # Handlers de pedidos
    └── models.go          # Definições de estruturas
```

## Características

- Detecção automática de arquivos de rotas e handlers
- Suporte a múltiplos arquivos
- Geração de documentação JSON e OpenAPI
- Suporte a comentários para descrição das rotas
- Análise de parâmetros de rota, query e body
- Documentação de respostas e códigos de status

## Contribuindo

Contribuições são bem-vindas! Por favor, sinta-se à vontade para submeter pull requests.

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.