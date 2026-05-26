.PHONY: build test lint run-server

build:
	go build -o bin/server .

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

run-server:
	go run .
