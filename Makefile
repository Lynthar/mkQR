# mkQR Makefile

# Build variables
BINARY_NAME=mkqr
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go variables
GOFLAGS=-ldflags "-X github.com/Lynthar/mkQR/internal/cli.Version=$(VERSION) \
                  -X github.com/Lynthar/mkQR/internal/cli.GitCommit=$(GIT_COMMIT) \
                  -X github.com/Lynthar/mkQR/internal/cli.BuildDate=$(BUILD_DATE)"

# Directories
BUILD_DIR=build
CMD_DIR=cmd/mkqr

.PHONY: all build clean test install uninstall cross help

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	go test -v ./...

# Install to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(GOFLAGS) ./$(CMD_DIR)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Uninstall from GOPATH/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Download dependencies
deps:
	go mod download
	go mod tidy

# Cross-compile for multiple platforms
cross:
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# Linux
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

	# macOS
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

	# Windows
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

	@echo "Cross-compilation complete. Binaries in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Show help
help:
	@echo "mkQR - QR Code Generator"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build for current platform"
	@echo "  make install  - Install to GOPATH/bin"
	@echo "  make test     - Run tests"
	@echo "  make cross    - Cross-compile for all platforms"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make deps     - Download dependencies"
	@echo "  make help     - Show this help"
