# Makefile for obom

# Variables

PROJECT_PKG = github.com/sajayantony/obom
BINARY = obom
SRC = $(wildcard *.go)
LDFLAGS = -s -w

GIT_COMMIT  = $(shell git rev-parse HEAD)
GIT_TAG     = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY   = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

LDFLAGS += -X '$(PROJECT_PKG)/internal/version.GitCommit=${GIT_COMMIT}'
LDFLAGS += -X '$(PROJECT_PKG)/internal/version.GitTreeState=${GIT_DIRTY}'

# Build the binary
.PHONY: build
build: $(SRC)
	go build -v --ldflags="$(LDFLAGS)" -o $(BINARY) . 

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
