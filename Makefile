# Development targets
run:
	go run ./cmd/goutui

dev: ## Run in development mode with hot reloading (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular run..."; \
		$(MAKE) run; \
	fi

# Testing targets
test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race ./...

test-cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-cover-func: ## Show test coverage by function
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Linting and code quality
vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format Go code
	go fmt ./...

fmt-check: ## Check if code is formatted
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	fi

# Building targets
build: ## Build the application
	go build -o goutui ./cmd/goutui

build-race: ## Build with race detection
	go build -race -o goutui ./cmd/goutui

build-debug: ## Build with debug information
	go build -gcflags="all=-N -l" -o goutui ./cmd/goutui

# Benchmarking
bench: ## Run benchmarks
	go test -bench=. ./...

bench-mem: ## Run benchmarks with memory allocation info
	go test -bench=. -benchmem ./...

# Dependencies
deps-update: ## Update all dependencies
	go get -u ./...
	go mod tidy

deps-clean: ## Clean up dependencies
	go mod tidy
	go mod download

# Installation
install: ## Install the application to $GOPATH/bin
	go install ./cmd/goutui

install-dev: ## Install development dependencies
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Cleanup
clean: ## Clean build artifacts
	rm -f goutui
	rm -f coverage.out coverage.html

clean-all: clean ## Clean all artifacts including dependencies
	go clean -cache
	go clean -testcache
	go clean -modcache

# Development setup
setup-dev: ## Set up development environment
	$(MAKE) install-dev
	$(MAKE) deps-clean
	@echo "Development environment set up successfully!"

# CI targets
ci: fmt-check vet lint test ## Run all CI checks

# Help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: run dev test test-verbose test-race test-cover test-cover-func vet lint fmt fmt-check build build-race build-debug bench bench-mem deps-update deps-clean install install-dev clean clean-all setup-dev ci help
