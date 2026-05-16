.PHONY: run build tidy test lint

BIN     := imghat
CMD     := ./cmd/server

run:
	go run $(CMD)/main.go

build:
	go build -o bin/$(BIN) $(CMD)/main.go

tidy:
	go mod tidy

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...