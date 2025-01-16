# Gobiru 

<div align="left">
       <img src="https://res.cloudinary.com/dx70wyorg/image/upload/v1736953035/photo_2025-01-15_11-40-32_esheqe.jpg" width="200" alt="Gobiru Mascot">
</div>

Gobiru é um gerador de documentação automático para APIs Go, com suporte para Gin, Gorilla Mux e Fiber.

## 🚀 Instalação

```bash
# Instalar o CLI globalmente
go install github.com/jeffemart/Gobiru/cmd/gobiru@latest

# Ou adicionar como dependência em seu projeto
go get github.com/jeffemart/Gobiru
```

## 🚀 Recursos

- Suporte para múltiplos frameworks:
  - Gin
  - Gorilla Mux
  - Fiber
- Geração automática de documentação OpenAPI/Swagger
- Interface Swagger UI embutida
- Detecção automática de rotas e parâmetros
- Servidor de documentação integrado
- Personalização via flags de comando

## 📝 Exemplos de Uso

### Gin Framework
```bash
gobiru -framework gin \
       -router routes.go \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API Gin" \
       -serve
```

### Gorilla Mux
```bash
gobiru -framework mux \
       -router routes.go \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API Mux" \
       -serve
```

### Fiber Framework
```bash
gobiru -framework fiber \
       -router routes.go \
       -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API Fiber" \
       -serve
```

## 🔧 Opções do CLI

```bash
Opções:
  -framework string    Framework a ser analisado (gin, mux ou fiber)
  -router string      Caminho do arquivo com definição do router (padrão: routes.go)
  -output string      Caminho do arquivo JSON de rotas (default "docs/routes.json")
  -openapi string     Caminho do arquivo OpenAPI/Swagger (default "docs/openapi.json")
  -title string       Título da API para spec OpenAPI (default "API Documentation")
  -description string Descrição da API para spec OpenAPI
  -api-version string Versão da API para spec OpenAPI (default "1.0.0")
  -serve             Iniciar servidor de documentação após geração
  -port int          Porta do servidor de documentação (default 8081)
  -help             Mostrar mensagem de ajuda
  -version          Mostrar versão do Gobiru
```

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, leia nosso guia de contribuição antes de enviar um PR.

## 📝 Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📞 Suporte

- Abra uma issue no GitHub
- Entre em contato via [LinkedIn](https://www.linkedin.com/in/jefferson-martins-dev/)
- Email: jefferson.developers@gmail.com
```