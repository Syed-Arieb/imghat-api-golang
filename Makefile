.PHONY: run build build-linux tidy test lint

BIN     := imghat
CMD     := ./cmd/server

run:
	go run $(CMD)/main.go

build:
	go build -o bin/$(BIN) $(CMD)/main.go

build-linux:
	cmd /c "set GOOS=linux&set GOARCH=amd64&go build -o bin/$(BIN)-linux $(CMD)/main.go"

tidy:
	go mod tidy

test:
	go test ./... -v

lint:
	golangci-lint run ./...