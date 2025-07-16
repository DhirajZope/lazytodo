# LazyTodo Makefile
# Common development and release tasks

.PHONY: build test clean install-deps release-test release-local help

# Default target
help: ## Show this help message
	@echo "LazyTodo Development Commands"
	@echo "============================="
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
build: ## Build the application
	go build -o lazytodo cmd/main.go

build-all: ## Build for all platforms (using GoReleaser)
	goreleaser build --snapshot --clean

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	rm -f lazytodo lazytodo.exe
	rm -rf dist/
	rm -f coverage.out coverage.html

# Dependencies
install-deps: ## Install Go dependencies
	go mod download

install-goreleaser: ## Install GoReleaser
	@which goreleaser > /dev/null || { \
		echo "Installing GoReleaser..."; \
		if command -v brew >/dev/null 2>&1; then \
			brew install goreleaser; \
		else \
			curl -sfL https://goreleaser.com/static/run | bash; \
		fi; \
	}

# Code quality
fmt: ## Format Go code
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter
	@which golangci-lint > /dev/null || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	}
	golangci-lint run

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all code quality checks

# Release testing
release-check: install-goreleaser ## Check GoReleaser configuration
	goreleaser check

release-test: install-goreleaser ## Test release build locally
	goreleaser release --snapshot --clean
	@echo
	@echo "Test release completed! Check dist/ folder for binaries."
	@echo "Example binaries created:"
	@find dist/ -name "lazytodo*" -type f | head -5

release-local: install-goreleaser ## Create a local test release
	goreleaser build --snapshot --clean
	@echo
	@echo "Local build completed! Test binaries:"
	@find dist/ -name "lazytodo" -type f | head -3

# Installation script testing
test-scripts: ## Test installation scripts syntax
	@echo "Testing installation scripts..."
	bash -n scripts/install-linux.sh
	bash -n scripts/install-mac.sh
	@echo "✅ All scripts have valid syntax"

# Docker
docker-build: ## Build Docker image locally
	docker build -t lazytodo:dev .

docker-run: docker-build ## Run Docker container
	docker run --rm -it lazytodo:dev --help

# Development setup
setup: install-deps install-goreleaser ## Set up development environment
	@echo "Development environment setup complete!"
	@echo
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make release-test   - Test full release build"
	@echo "  make check          - Run all quality checks"

# Version management
version: ## Show current version info
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD)"
	@echo "Git tag: $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'No tags')"

# Create release tag
tag: ## Create a new release tag (usage: make tag VERSION=v1.0.0)
ifndef VERSION
	@echo "Error: VERSION is required. Usage: make tag VERSION=v1.0.0"
	@exit 1
endif
	@echo "Creating tag $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Tag created. Push with: git push origin $(VERSION)"
	@echo "This will trigger the release workflow."

# Quick development cycle
dev: ## Quick development cycle: format, test, and build
	make fmt
	make test
	make build
	@echo "✅ Development cycle complete"

# Installation
install: build ## Install binary to local system
	@if [ -w /usr/local/bin ]; then \
		cp lazytodo /usr/local/bin/; \
		echo "✅ Installed to /usr/local/bin/lazytodo"; \
	elif [ -w $(HOME)/.local/bin ]; then \
		mkdir -p $(HOME)/.local/bin; \
		cp lazytodo $(HOME)/.local/bin/; \
		echo "✅ Installed to $(HOME)/.local/bin/lazytodo"; \
	else \
		echo "❌ No writable install directory found"; \
		echo "Try: sudo make install"; \
	fi

uninstall: ## Uninstall binary from local system
	@rm -f /usr/local/bin/lazytodo $(HOME)/.local/bin/lazytodo
	@echo "✅ LazyTodo uninstalled" 