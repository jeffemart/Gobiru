# Gobiru - API Documentation Generator

Gobiru é uma ferramenta de linha de comando para gerar documentação de APIs Go automaticamente, suportando os frameworks Gin, Gorilla Mux e Fiber.

## Instalação

```bash
go install github.com/jeffemart/gobiru/cmd/gobiru@latest
```

## Uso

O Gobiru requer três arquivos principais do seu projeto:
- main.go: Arquivo principal da aplicação
- routes.go: Arquivo com as definições das rotas
- handlers.go: Arquivo com os handlers da API

```bash
gobiru -framework [gin|mux|fiber] \
       -main path/to/main.go \
       -router path/to/routes.go \
       -handlers path/to/handlers.go \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API" \
       -description "Descrição da minha API" \
       -version "1.0.0"
```

### Flags Disponíveis

- `-framework`: Framework utilizado (gin, mux, ou fiber)
- `-main`: Caminho para o arquivo main.go
- `-router`: Caminho para o arquivo routes.go
- `-handlers`: Caminho para o arquivo handlers.go
- `-output`: Caminho para o arquivo JSON de saída (default: docs/routes.json)
- `-openapi`: Caminho para o arquivo OpenAPI (default: docs/openapi.json)
- `-title`: Título da API para documentação OpenAPI
- `-description`: Descrição da API
- `-version`: Versão da API

## Exemplos

O repositório inclui exemplos para cada framework suportado:

### Gorilla Mux
```bash
gobiru -framework mux \
       -main examples/mux/simple/main.go \
       -router examples/mux/simple/main.go \
       -handlers examples/mux/simple/main.go
```

### Gin
```bash
gobiru -framework gin \
       -main examples/gin/simple/main.go \
       -router examples/gin/simple/main.go \
       -handlers examples/gin/simple/main.go
```

### Fiber
```bash
gobiru -framework fiber \
       -main examples/fiber/simple/main.go \
       -router examples/fiber/simple/main.go \
       -handlers examples/fiber/simple/main.go
```

## Desenvolvimento

Para contribuir com o projeto:

1. Clone o repositório
2. Instale as dependências: `go mod download`
3. Execute os testes: `make test`
4. Faça o build: `make build`

## Licença

MIT License