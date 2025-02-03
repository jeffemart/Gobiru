.PHONY: build test clean

# Variáveis
BINARY_NAME=gobiru
BUILD_DIR=bin

# Comandos
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/gobiru/main.go

test:
	@echo "Running tests..."
	@go test ./internal/... ./cmd/...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

install:
	@echo "Installing..."
	@go install ./cmd/gobiru

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