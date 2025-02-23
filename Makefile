.PHONY: build test lint

build:
	@go build -o bin/ember ./cmd/ember

test:
	@go test ./ember_lang/...

lint:
	@golangci-lint run

run:
	@go run ./cmd/ember
