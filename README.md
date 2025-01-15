# Gobiru 

<div align="left">
       <img src="https://res.cloudinary.com/dx70wyorg/image/upload/v1736953035/photo_2025-01-15_11-40-32_esheqe.jpg" width="200" alt="Gobiru Mascot">
</div>

Gobiru é um gerador de documentação automático para APIs Go, com suporte para Gin e Gorilla Mux.

## 🚀 Instalação

```bash
# Instalar o CLI globalmente
go install github.com/jeffemart/Gobiru/cmd/gobiru@latest

# Ou adicionar como dependência em seu projeto
go get github.com/jeffemart/Gobiru
```

## 📖 Uso do CLI

### Comandos Básicos

```bash
# Gerar documentação para API Gin
gobiru -framework gin -output docs/routes.json main.go

# Gerar documentação para API Mux
gobiru -framework mux -output docs/routes.json main.go

# Gerar documentação e iniciar servidor Swagger UI
gobiru -framework gin -output docs/routes.json -openapi docs/openapi.json -serve main.go
```

### Opções Disponíveis

```bash
Opções:
  -framework string     Framework a ser analisado (gin ou mux)
  -output string       Caminho do arquivo JSON de rotas (default "docs/routes.json")
  -openapi string      Caminho do arquivo OpenAPI/Swagger (default "docs/openapi.json")
  -title string        Título da API para spec OpenAPI (default "API Documentation")
  -description string  Descrição da API para spec OpenAPI
  -api-version string  Versão da API para spec OpenAPI (default "1.0.0")
  -serve              Iniciar servidor de documentação após geração
  -port int           Porta do servidor de documentação (default 8081)
  -help              Mostrar mensagem de ajuda
  -version           Mostrar versão do Gobiru
```

### Exemplos Completos

#### Para APIs usando Gin

```bash
# 1. Instalar o Gobiru
go install github.com/jeffemart/Gobiru/cmd/gobiru@latest

# 2. Em seu projeto Gin, gerar a documentação
gobiru -framework gin \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API Gin" \
       -description "API de exemplo usando Gin" \
       -api-version "1.0.0" \
       -serve \
       main.go

# 3. Acessar a documentação
# - Swagger UI: http://localhost:8081/docs/index.html
# - OpenAPI JSON: http://localhost:8081/docs/openapi.json
# - Routes JSON: http://localhost:8081/docs/routes.json
```

#### Para APIs usando Gorilla Mux

```bash
# 1. Instalar o Gobiru
go install github.com/jeffemart/Gobiru/cmd/gobiru@latest

# 2. Em seu projeto Mux, gerar a documentação
gobiru -framework mux \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API Mux" \
       -description "API de exemplo usando Gorilla Mux" \
       -api-version "1.0.0" \
       -serve \
       main.go

# 3. Acessar a documentação
# - Swagger UI: http://localhost:8081/docs/index.html
# - OpenAPI JSON: http://localhost:8081/docs/openapi.json
# - Routes JSON: http://localhost:8081/docs/routes.json
```

## 💡 Uso como Biblioteca

```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    gobiru "github.com/jeffemart/Gobiru/app/gin"
    "github.com/jeffemart/Gobiru/app/openapi"
)

func main() {
    // Criar router Gin
    router := gin.Default()
    
    // Definir rotas
    router.GET("/users", getUsers)
    router.POST("/users", createUser)
    
    // Criar analisador
    analyzer := gobiru.NewAnalyzer()
    
    // Analisar rotas
    err := analyzer.AnalyzeRoutes(router)
    if err != nil {
        log.Fatal(err)
    }
    
    // Exportar documentação
    info := openapi.Info{
        Title: "Minha API",
        Description: "Descrição da minha API",
        Version: "1.0.0",
    }
    
    err = analyzer.ExportOpenAPI("docs/openapi.json", info)
    if err != nil {
        log.Fatal(err)
    }
    
    // Iniciar servidor
    router.Run(":8080")
}
```

## 🔄 Workflow Recomendado

1. Instale o Gobiru globalmente
2. Desenvolva sua API normalmente usando Gin ou Mux
3. Use o comando `gobiru` para gerar a documentação
4. Inicie o servidor de documentação com a flag `-serve`
5. Acesse a documentação via Swagger UI
6. Atualize a documentação sempre que modificar as rotas

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, leia nosso guia de contribuição antes de enviar um PR.

## 📝 Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📞 Suporte

- Abra uma issue no GitHub
- Entre em contato via [LinkedIn](https://www.linkedin.com/in/jefferson-martins-dev/)
- Email: jefferson.developers@gmail.com
```
