# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar arquivos necessários
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar o CLI
RUN CGO_ENABLED=0 GOOS=linux go build -o /gobiru ./cmd/gobiru/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /gobiru /usr/local/bin/gobiru

# Copiar arquivos estáticos
COPY examples/test_cli/docs/index.html /app/docs/index.html

# Expor a porta do servidor
EXPOSE 8081

# Criar um script de entrada
COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"] 