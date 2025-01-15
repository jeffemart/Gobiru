# Gobiru 

<div align="left">
       <img src="https://res.cloudinary.com/dx70wyorg/image/upload/v1736953035/photo_2025-01-15_11-40-32_esheqe.jpg" width="200" alt="Gobiru Mascot">
</div>

## 🚀 Funcionalidades

- Análise automática de rotas do Gorilla Mux
- Geração de documentação no formato JSON personalizado
- Geração de especificação OpenAPI (Swagger)
- Servidor de documentação integrado com Swagger UI
- Suporte a Docker
- Documentação interativa com Swagger UI

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Docker (opcional)
- Make (opcional)
- Gorilla Mux

## 🔧 Instalação

```bash
# Clonar o repositório
git clone https://github.com/jeffemart/Gobiru.git
cd Gobiru

# Instalar dependências
go mod download

# Compilar
make build
```

## 💻 Uso

### Estrutura do Projeto
```
.
├── app/
│   ├── main.go              # Core da aplicação
│   ├── models/              # Modelos de dados
│   └── openapi/            # Conversor OpenAPI
├── cmd/
│   └── gobiru/             # CLI da aplicação
├── examples/               # Exemplos de uso
│   ├── simple/            # Exemplo básico
│   └── test_cli/          # Exemplo com servidor
└── scripts/               # Scripts auxiliares
```

### CLI

```bash
# Gerar documentação
gobiru -output docs/routes.json \
       -openapi docs/openapi.json \
       -title "Minha API" \
       -description "Documentação da minha API" \
       path/to/routes.go

# Iniciar servidor de documentação
go run examples/test_cli/server.go
```

### Docker

```bash
# Construir e executar com Docker
make docker-run

# Ou usando docker-compose
docker-compose up
```

### Acessar a Documentação

Após iniciar o servidor, acesse:
- Swagger UI: http://localhost:8081/docs/index.html
- OpenAPI JSON: http://localhost:8081/docs/openapi.json
- Routes JSON: http://localhost:8081/docs/routes.json

## ⚙️ Configuração

### Opções do CLI

```bash
gobiru [options] <path-to-routes-file>

Options:
  -output string
        Caminho do arquivo de saída JSON (default "routes.json")
  -openapi string
        Caminho do arquivo de saída OpenAPI
  -title string
        Título da API para documentação OpenAPI (default "API Documentation")
  -description string
        Descrição da API para documentação OpenAPI
  -version string
        Versão da API (default "1.0.0")
```

## 🌟 Exemplo

```go
// examples/test_cli/routes.go
package main

import (
    "net/http"
    "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users/{id}", getUser).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")

    return router
}
```

### Exemplo de Saída JSON
```json
{
    "method": "GET",
    "path": "/users/{id}",
    "description": "",
    "handler_name": "main.getUser",
    "parameters": [
        {
            "name": "id",
            "type": "string",
            "required": true,
            "description": ""
        }
    ],
    "api_version": "v1.0"
}
```

## 🛠️ Desenvolvimento

```bash
# Executar testes
make test

# Executar localmente
make run

# Limpar arquivos gerados
make clean
```

## 🐳 Docker

O projeto inclui configurações Docker para facilitar o deployment:

```bash
# Construir imagem
docker build -t gobiru .

# Executar container
docker run -p 8081:8081 gobiru
```

## 🔄 CI/CD

O projeto utiliza GitHub Actions para:
- Executar testes automaticamente
- Construir e publicar imagem Docker
- Verificar qualidade do código
- Deploy automático em tags

## 👥 Contribuição

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## ✨ Próximos Passos

- [ ] Suporte para análise de comentários do código
- [ ] Detecção automática de request/response bodies
- [ ] Suporte para outros frameworks de roteamento
- [ ] Melhor análise de tipos Go
- [ ] Interface web para gerenciamento da documentação
- [ ] Suporte para autenticação e autorização
- [ ] Geração de documentação em outros formatos (PDF, Markdown)
- [ ] Integração com mais ferramentas de documentação

## 🤝 Contribuidores

- [@jeffemart](https://github.com/jeffemart) - Criador e mantenedor

## 📞 Suporte

Para suporte:
- Abra uma issue no GitHub
- Entre em contato via [LinkedIn](https://www.linkedin.com/in/jefferson-martins-dev/)
- Email: jefferson.martins.dev@gmail.com

## 📄 Licença

Este projeto está sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.
```
