.PHONY: build test clean install run help

# Build the binary
build:
	@echo "Building json-to-go..."
	@go build -o json-to-go ./cmd/json-to-go
	@echo "Build complete! Run with: ./json-to-go"

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f json-to-go coverage.out coverage.html
	@echo "Clean complete!"

# Install the binary to $GOPATH/bin
install:
	@echo "Installing json-to-go..."
	@go install ./cmd/json-to-go
	@echo "Installed! You can now run: json-to-go"

# Run with example
run:
	@./json-to-go -type=Example example.json

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  run           - Run example"
	@echo "  help          - Show this help message"
