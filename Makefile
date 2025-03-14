.PHONY: build test clean run

# Build the chess game
build:
	go build -o bin/chess cmd/chess/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the chess game
run:
	go run cmd/chess/main.go

# Run with ASCII mode
run-ascii:
	go run cmd/chess/main.go -ascii

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Build and run tests with coverage
coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the chess game"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  run        - Run the chess game"
	@echo "  run-ascii  - Run with ASCII display mode"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  coverage   - Generate test coverage report"
	@echo "  help       - Show this help message" 