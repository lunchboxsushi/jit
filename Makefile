.PHONY: help build test lint clean install

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the jit binary"
	@echo "  test      - Run tests with coverage"
	@echo "  lint      - Run golangci-lint"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install jit binary"
	@echo "  ci        - Run full CI pipeline locally"

# Build the binary
build:
	go build -v -o jit ./cmd/jit

# Run tests with coverage
test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	go vet ./...
	go fmt ./...

# Clean build artifacts
clean:
	rm -f jit coverage.out coverage.html
	go clean -cache

# Install binary
install:
	go install ./cmd/jit

# Run full CI pipeline locally
ci: lint test build
	@echo "CI pipeline completed successfully!"

# Run tests for specific package
test-pkg:
	@read -p "Enter package path: " pkg; \
	go test -v ./$$pkg

# Generate test coverage badge
coverage-badge:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' > coverage.txt
	@echo "Coverage percentage saved to coverage.txt" 