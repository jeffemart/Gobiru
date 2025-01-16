# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o /gobiru ./cmd/gobiru/main.go

# Final stage
FROM alpine:latest

# Copiar o binário do estágio de build
COPY --from=builder /gobiru /usr/local/bin/gobiru

# Expor a porta do servidor
EXPOSE 8081

# Definir o diretório de trabalho padrão
VOLUME /work
WORKDIR /work

# Comando padrão
ENTRYPOINT ["gobiru"] 