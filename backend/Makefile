# Backend Makefile for Pickup Queue API
.PHONY: help build run test clean deps lint migrate

# Variables
BINARY_NAME=pickup-api
WORKER_BINARY=pickup-worker
BUILD_DIR=bin
API_CMD=cmd/api/main.go
WORKER_CMD=cmd/worker/main.go

# Default target
help: ## Show help message
	@echo "Backend Makefile for Pickup Queue API"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# DEPENDENCIES
# =============================================================================

deps: ## Download and tidy dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

deps-update: ## Update all dependencies
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy

# =============================================================================
# BUILD
# =============================================================================

build: build-api build-worker ## Build both API and worker binaries

build-api: ## Build API server binary
	@echo "🔨 Building API server..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(API_CMD)
	@echo "✅ API server built: $(BUILD_DIR)/$(BINARY_NAME)"

build-worker: ## Build worker binary
	@echo "🔨 Building worker..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(WORKER_BINARY) $(WORKER_CMD)
	@echo "✅ Worker built: $(BUILD_DIR)/$(WORKER_BINARY)"

build-prod: ## Build optimized production binaries
	@echo "🔨 Building production binaries..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BUILD_DIR)/$(BINARY_NAME) $(API_CMD)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BUILD_DIR)/$(WORKER_BINARY) $(WORKER_CMD)
	@echo "✅ Production binaries built!"

# =============================================================================
# RUN
# =============================================================================

run: run-api ## Run API server (alias)

run-api: ## Run API server in development mode
	@echo "🚀 Starting API server..."
	go run $(API_CMD)

run-worker: ## Run background worker
	@echo "🔄 Starting background worker..."
	go run $(WORKER_CMD)

run-prod: build ## Run production binaries
	@echo "🚀 Starting production API server..."
	./$(BUILD_DIR)/$(BINARY_NAME) &
	@echo "🔄 Starting production worker..."
	./$(BUILD_DIR)/$(WORKER_BINARY) &

# =============================================================================
# TESTING
# =============================================================================

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "🧪 Running tests with race detection..."
	go test ./... -race -v

benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	go test ./... -bench=. -benchmem

# =============================================================================
# LINTING & FORMATTING
# =============================================================================

lint: ## Run linting
	@echo "🔍 Running linters..."
	go fmt ./...
	go vet ./...
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

format: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...
	@which goimports > /dev/null 2>&1 || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .

# =============================================================================
# DATABASE
# =============================================================================

migrate: ## Run database migrations
	@echo "🔄 Running database migrations..."
	@if [ ! -f migrations/001_create_packages_table.sql ]; then \
		echo "❌ Migration file not found!"; \
		exit 1; \
	fi
	@echo "Applying migration to database..."
	@# Check if we can connect to the database
	@go run -c 'package main; import ("database/sql"; _ "github.com/lib/pq"); func main() { db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=pickup_queue sslmode=disable"); if err != nil { panic(err) }; defer db.Close(); if err = db.Ping(); err != nil { panic(err) } }' 2>/dev/null || (echo "❌ Cannot connect to database. Make sure PostgreSQL is running." && exit 1)

migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d_%H%M%S); \
	filename="migrations/$${timestamp}_$${name}.sql"; \
	touch $$filename; \
	echo "✅ Created migration file: $$filename"

# =============================================================================
# DEVELOPMENT UTILITIES
# =============================================================================

watch: ## Watch for changes and restart API server
	@echo "👀 Watching for changes..."
	@which air > /dev/null 2>&1 || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

hot-reload: watch ## Alias for watch

generate: ## Run go generate
	@echo "🔄 Running go generate..."
	go generate ./...

mod-graph: ## Show module dependency graph
	@echo "📊 Module dependency graph:"
	go mod graph

mod-why: ## Show why packages are needed
	@echo "❓ Analyzing package dependencies..."
	@read -p "Enter package name: " pkg; \
	go mod why $$pkg

# =============================================================================
# SECURITY & ANALYSIS
# =============================================================================

security: ## Run security analysis
	@echo "🔒 Running security analysis..."
	@which gosec > /dev/null 2>&1 || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

vuln-check: ## Check for vulnerabilities
	@echo "🛡️ Checking for vulnerabilities..."
	@which govulncheck > /dev/null 2>&1 || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)
	govulncheck ./...

# =============================================================================
# DOCKER
# =============================================================================

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t pickup-backend .

docker-run: ## Run Docker container
	@echo "🚀 Running Docker container..."
	docker run -d \
		--name pickup-backend \
		-p 8080:8080 \
		-e DB_HOST=host.docker.internal \
		pickup-backend

docker-stop: ## Stop Docker container
	@echo "🛑 Stopping Docker container..."
	docker stop pickup-backend || true
	docker rm pickup-backend || true

# =============================================================================
# CLEANUP
# =============================================================================

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)/
	rm -f coverage.out coverage.html
	go clean -cache
	go clean -testcache

clean-deps: ## Clean dependency cache
	@echo "🧹 Cleaning dependency cache..."
	go clean -modcache

# =============================================================================
# UTILITIES
# =============================================================================

env-check: ## Check environment setup
	@echo "🔍 Checking environment..."
	@echo "Go version: $$(go version)"
	@echo "GOPATH: $$GOPATH"
	@echo "GOROOT: $$GOROOT"
	@echo "Current directory: $$(pwd)"
	@echo ""
	@echo "Environment variables:"
	@grep -E '^[A-Z_]+=' .env 2>/dev/null || echo "No .env file found"

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development tools installed!"

server-info: ## Show server information
	@echo "📋 Server Information:"
	@echo "API Server: http://localhost:8080"
	@echo "Health Check: http://localhost:8080/health"
	@echo ""
	@echo "Available endpoints:"
	@echo "  GET    /health"
	@echo "  POST   /api/v1/packages"
	@echo "  GET    /api/v1/packages"
	@echo "  GET    /api/v1/packages/:id"
	@echo "  PATCH  /api/v1/packages/:id/status"
	@echo "  DELETE /api/v1/packages/:id"
	@echo "  GET    /api/v1/packages/stats"

# =============================================================================
# QUICK COMMANDS
# =============================================================================

dev: deps run-api ## Quick start development

all: clean deps lint test build ## Run full build pipeline

ci: deps lint test-coverage security ## Run CI pipeline

.DEFAULT_GOAL := help
