.PHONY: build test lint

build:
	@go build -o bin/ember ./cmd/ember

test:
	@go test ./internal/...

lint:
	@golangci-lint run

run:
	@go run ./cmd/ember
