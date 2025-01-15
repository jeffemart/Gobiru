.PHONY: build test run docker-build docker-run clean

build:
	go build -o gobiru ./cmd/gobiru/main.go

test:
	go test -v ./...

run: build
	./gobiru -output docs/routes.json -openapi docs/openapi.json examples/test_cli/routes.go
	go run examples/test_cli/server.go

docker-build:
	docker build -t gobiru .

docker-run: docker-build
	docker run -p 8081:8081 gobiru

clean:
	rm -f gobiru
	rm -rf docs/ 