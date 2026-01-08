.PHONY: build build-all clean test run install deps help

# Build variables
BINARY_NAME=whatsapp-proxy
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS=-s -w \
	-X main.Version=$(VERSION) \
	-X main.BuildTime=$(BUILD_TIME) \
	-X main.GitCommit=$(GIT_COMMIT)

# Build output directory
DIST_DIR=dist

# Supported platforms
PLATFORMS=\
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	linux/arm \
	windows/amd64 \
	windows/386 \
	darwin/amd64 \
	darwin/arm64 \
	freebsd/amd64

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

build: deps ## Build for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(DIST_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME) ./cmd/whatsapp-proxy
	@echo "Build complete: $(DIST_DIR)/$(BINARY_NAME)"

build-all: deps ## Build for all platforms
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUTPUT="$(DIST_DIR)/$(BINARY_NAME)-$${OS}-$${ARCH}"; \
		if [ "$$OS" = "windows" ]; then \
			OUTPUT="$${OUTPUT}.exe"; \
		fi; \
		echo "Building $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -ldflags="$(LDFLAGS)" -o $$OUTPUT ./cmd/whatsapp-proxy || exit 1; \
		echo "  Created: $$OUTPUT"; \
	done
	@echo "All builds complete!"

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out

test-coverage: test ## Run tests and show coverage in browser
	go tool cover -html=coverage.out

run: build ## Build and run locally
	@echo "Running $(BINARY_NAME)..."
	./$(DIST_DIR)/$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(DIST_DIR) coverage.out
	@echo "Clean complete"

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install -ldflags="$(LDFLAGS)" ./cmd/whatsapp-proxy
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

lint: ## Run linters
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run --timeout 5m

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	go mod tidy

dev: ## Run in development mode with hot reload (requires air)
	@which air > /dev/null || (echo "air not installed. Run: go install github.com/cosmtrek/air@latest" && exit 1)
	air

.DEFAULT_GOAL := help
