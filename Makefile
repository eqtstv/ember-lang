.PHONY: build test lint install

build:
	@go build -o bin/ember ./cmd/ember

test:
	@go test ./ember_lang/...

lint:
	@golangci-lint run

run:
	@go run ./cmd/ember

install: build
	@sudo cp bin/ember /usr/local/bin/ember
	@echo "Ember installed to /usr/local/bin/ember"
