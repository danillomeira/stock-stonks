.PHONY: run build test lint fmt tidy

run:
	go run ./cmd/api

build:
	CGO_ENABLED=1 go build -o bin/api ./cmd/api

test:
	go test ./... -race -cover

lint:
	go vet ./...
	gofmt -l .

fmt:
	gofmt -w .

tidy:
	go mod tidy
