.PHONY: build install clean test run deps

# Build the Base Framework MCP server
build:
	go build -o base-mcp .

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install the MCP server to GOPATH/bin
install: build
	go install .

# Clean build artifacts
clean:
	rm -f base-mcp
	go clean

# Run tests
test:
	go test -v ./...

# Run the MCP server
run:
	go run .

# Development build with race detection
dev:
	go build -race -o base-mcp .

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o base-mcp-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o base-mcp-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o base-mcp-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o base-mcp-windows-amd64.exe .

# Show help
help:
	@echo "Available commands:"
	@echo "  build      - Build the MCP server"
	@echo "  install    - Install the MCP server"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  run        - Run the MCP server"
	@echo "  dev        - Build with race detection"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  deps       - Install dependencies"