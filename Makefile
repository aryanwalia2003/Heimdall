.PHONY: build test clean lint

# Default target
all: build

# Build the shared object library
build:
	go build -v -buildmode=c-shared -o bin/libheimdall.so core/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Lint the codebase
lint:
	golangci-lint run ./...
