# Gobiru - Gerador de Documentação de API

## Estrutura do Projeto

gobiru/
├── cmd/
│   └── gobiru/
│       └── main.go
├── internal/
│   ├── analyzer/
│   │   ├── analyzer.go
│   │   ├── common.go
│   │   ├── fiber.go
│   │   └── mux.go
│   └── generator/
│       ├── json.go
│       └── openapi.go
├── examples/
│   ├── fiber/
│   ├── gin/
│   └── gorilla/
├── go.mod
├── go.sum
├── LICENSE
└── README.md

## Testes

O projeto inclui testes para garantir a funcionalidade correta. Os testes estão localizados nos seguintes diretórios:

- `internal/analyzer/`: Testes para o analisador, incluindo a função `FindMainFile`.
- `examples/fiber/handlers/`: Testes para os handlers do Fiber.
- `examples/fiber/routes/`: Testes para as rotas do Fiber.

### Executando os Testes

Para executar todos os testes do projeto, use o seguinte comando:

```bash
go test ./...
```

Os testes verificarão se as funções e rotas estão funcionando conforme o esperado.