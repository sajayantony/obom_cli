# Makefile for obom

# Variables
BINARY = obom
SRC = $(wildcard *.go)
GOFLAGS = -ldflags "-s -w"

# Build the binary
.PHONY: build
build: $(SRC)
	go build $(GOFLAGS) -o $(BINARY) $(SRC)

# Clean up build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out
