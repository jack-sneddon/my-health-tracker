BINARY=bin/tracker

all: fmt tidy vet test build

fmt:
	go fmt ./...

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./...

build:
	go build -o $(BINARY) ./cmd/tracker
