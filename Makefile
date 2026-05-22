.PHONY: build run test lint

build:
	go build -o bin/gendiff ./cmd/gendiff

ARGS ?=

run: build
	./bin/gendiff $(ARGS)

test:
	go test -v ./...

lint:
	golangci-lint run ./...

run-server:
	go run main.go
