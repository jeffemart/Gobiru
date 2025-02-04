.PHONY: build test lint clean docs

# Variáveis
BINARY_NAME=gobiru
MAIN_FILE=cmd/gobiru/main.go

# Build
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)

# Testes
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Limpar
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean -testcache

# Documentação
docs:
	@echo "Generating documentation..."
	@go run $(MAIN_FILE) --framework gin --base-dir examples/gin --output docs/gin-openapi.json
	@go run $(MAIN_FILE) --framework fiber --base-dir examples/fiber --output docs/fiber-openapi.json
	@go run $(MAIN_FILE) --framework mux --base-dir examples/gorilla --output docs/mux-openapi.json

# Instalar dependências de desenvolvimento
setup:
	@echo "Setting up development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod download

# Executar todos os checks
check: lint test

# Build para todas as plataformas
build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

# Exemplos
example-mux:
	@echo "Running Mux example..."
	@go run examples/mux/simple/main.go

example-gin:
	@echo "Running Gin example..."
	@go run examples/gin/simple/main.go

example-fiber:
	@echo "Running Fiber example..."
	@go run examples/fiber/simple/main.go