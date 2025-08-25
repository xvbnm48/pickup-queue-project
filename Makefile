# Pickup Queue Project Makefile
.PHONY: help setup clean build run test docker-up docker-down docker-logs compose-up compose-down compose-build compose-logs compose-clean compose-restart

# Default target
help: ## Show this help message
	@echo "Pickup Queue Management System"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# SETUP & DEPENDENCIES
# =============================================================================

setup: setup-backend setup-frontend ## Setup both backend and frontend dependencies
	@echo "✅ All dependencies installed successfully!"

setup-backend: ## Install backend dependencies
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	cd backend && go mod tidy

setup-frontend: ## Install frontend dependencies
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install

# =============================================================================
# DATABASE
# =============================================================================

docker-up: ## Start PostgreSQL database with Docker
	@echo "🐳 Starting PostgreSQL database..."
	docker network create pickup-network 2>/dev/null || true
	docker run -d \
		--name pickup-postgres \
		--network pickup-network \
		-e POSTGRES_DB=pickup_queue \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5432:5432 \
		-v pickup_postgres_data:/var/lib/postgresql/data \
		postgres:15-alpine || docker start pickup-postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@docker exec pickup-postgres pg_isready -U postgres -d pickup_queue || true
	@echo "✅ Database is ready!"

docker-down: ## Stop PostgreSQL database
	@echo "🛑 Stopping PostgreSQL database..."
	docker stop pickup-postgres 2>/dev/null || true
	docker rm pickup-postgres 2>/dev/null || true

docker-logs: ## View PostgreSQL logs
	docker logs -f pickup-postgres

db-migrate: ## Run database migrations
	@echo "🔄 Running database migrations..."
	docker exec -i pickup-postgres psql -U postgres -d pickup_queue < backend/migrations/001_create_packages_table.sql || true
	@echo "✅ Migrations completed!"

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "⚠️  Resetting database..."
	docker exec pickup-postgres psql -U postgres -d pickup_queue -c "DROP TABLE IF EXISTS packages CASCADE;"
	$(MAKE) db-migrate
	@echo "✅ Database reset completed!"

db-connect: ## Connect to PostgreSQL database
	docker exec -it pickup-postgres psql -U postgres -d pickup_queue

# =============================================================================
# BACKEND
# =============================================================================

run-backend: ## Run backend API server
	@echo "🚀 Starting backend API server..."
	cd backend && go run cmd/api/main.go

run-worker: ## Run background worker
	@echo "🔄 Starting background worker..."
	cd backend && go run cmd/worker/main.go

build-backend: ## Build backend binaries
	@echo "🔨 Building backend..."
	cd backend && go build -o bin/api cmd/api/main.go
	cd backend && go build -o bin/worker cmd/worker/main.go
	@echo "✅ Backend built successfully!"

test-backend: ## Run backend tests
	@echo "🧪 Running backend tests..."
	cd backend && go test ./... -v

clean-backend: ## Clean backend build artifacts
	@echo "🧹 Cleaning backend..."
	cd backend && rm -rf bin/
	cd backend && go clean

lint-backend: ## Lint backend code
	@echo "🔍 Linting backend code..."
	cd backend && go fmt ./...
	cd backend && go vet ./...

# =============================================================================
# FRONTEND
# =============================================================================

run-frontend: ## Run frontend development server
	@echo "🚀 Starting frontend development server..."
	cd frontend && npm run dev

build-frontend: ## Build frontend for production
	@echo "🔨 Building frontend..."
	cd frontend && npm run build
	@echo "✅ Frontend built successfully!"

test-frontend: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	cd frontend && npm test

clean-frontend: ## Clean frontend build artifacts
	@echo "🧹 Cleaning frontend..."
	cd frontend && rm -rf dist/
	cd frontend && rm -rf node_modules/.cache/

lint-frontend: ## Lint frontend code
	@echo "🔍 Linting frontend code..."
	cd frontend && npm run lint

preview-frontend: ## Preview frontend production build
	@echo "👀 Previewing frontend..."
	cd frontend && npm run preview

# =============================================================================
# DEVELOPMENT
# =============================================================================

dev: docker-up ## Start full development environment
	@echo "🚀 Starting full development environment..."
	@sleep 2
	$(MAKE) db-migrate
	@echo "Starting all services..."
	@echo "1. Database: ✅ Ready"
	@echo "2. Starting Backend API..."
	@echo "3. Starting Worker..."
	@echo "4. Starting Frontend..."
	@echo ""
	@echo "🎯 Run the following in separate terminals:"
	@echo "   make run-backend"
	@echo "   make run-worker"
	@echo "   make run-frontend"

dev-api: docker-up db-migrate run-backend ## Start database and API only

dev-full: ## Start all services in background (experimental)
	@echo "🚀 Starting all services..."
	$(MAKE) docker-up
	$(MAKE) db-migrate
	@echo "Starting services in background..."
	cd backend && go run cmd/api/main.go &
	cd backend && go run cmd/worker/main.go &
	cd frontend && npm run dev &
	@echo "✅ All services started!"
	@echo "API: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"

stop-dev: ## Stop all development services
	@echo "🛑 Stopping all services..."
	@pkill -f "go run cmd/api/main.go" 2>/dev/null || true
	@pkill -f "go run cmd/worker/main.go" 2>/dev/null || true
	@pkill -f "npm run dev" 2>/dev/null || true
	$(MAKE) docker-down
	@echo "✅ All services stopped!"

# =============================================================================
# BUILDING & TESTING
# =============================================================================

build: build-backend build-frontend ## Build both backend and frontend

test: test-backend test-frontend ## Run all tests

lint: lint-backend lint-frontend ## Lint all code

clean: clean-backend clean-frontend ## Clean all build artifacts
	@echo "✅ All artifacts cleaned!"

# =============================================================================
# DOCKER PRODUCTION
# =============================================================================

docker-build: ## Build production Docker images
	@echo "🐳 Building Docker images..."
	docker build -t pickup-backend -f backend/Dockerfile backend/
	docker build -t pickup-frontend -f frontend/Dockerfile frontend/
	@echo "✅ Docker images built!"

docker-run-prod: ## Run production Docker containers
	@echo "🚀 Starting production containers..."
	docker network create pickup-network 2>/dev/null || true
	docker run -d --name pickup-postgres --network pickup-network \
		-e POSTGRES_DB=pickup_queue \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5432:5432 \
		postgres:15-alpine
	@sleep 5
	docker run -d --name pickup-backend --network pickup-network \
		-p 8080:8080 \
		-e DB_HOST=pickup-postgres \
		pickup-backend
	docker run -d --name pickup-frontend \
		-p 3000:80 \
		pickup-frontend
	@echo "✅ Production environment started!"

# =============================================================================
# UTILITIES
# =============================================================================

status: ## Show status of all services
	@echo "📊 Service Status:"
	@echo "Database:"
	@docker ps --filter "name=pickup-postgres" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || echo "  Not running"
	@echo ""
	@echo "API Server:"
	@pgrep -f "go run cmd/api/main.go" >/dev/null && echo "  ✅ Running" || echo "  ❌ Not running"
	@echo ""
	@echo "Worker:"
	@pgrep -f "go run cmd/worker/main.go" >/dev/null && echo "  ✅ Running" || echo "  ❌ Not running"
	@echo ""
	@echo "Frontend:"
	@pgrep -f "npm run dev" >/dev/null && echo "  ✅ Running" || echo "  ❌ Not running"

logs-backend: ## Show backend logs (if running)
	@echo "📝 Backend logs:"
	@ps aux | grep "go run cmd/api/main.go" | grep -v grep || echo "Backend not running"

logs-worker: ## Show worker logs (if running)
	@echo "📝 Worker logs:"
	@ps aux | grep "go run cmd/worker/main.go" | grep -v grep || echo "Worker not running"

health: ## Check health of all services
	@echo "🏥 Health Check:"
	@echo "Database:"
	@docker exec pickup-postgres pg_isready -U postgres -d pickup_queue 2>/dev/null && echo "  ✅ Healthy" || echo "  ❌ Unhealthy"
	@echo ""
	@echo "API:"
	@curl -s http://localhost:8080/health >/dev/null && echo "  ✅ Healthy" || echo "  ❌ Unhealthy"
	@echo ""
	@echo "Frontend:"
	@curl -s http://localhost:5173 >/dev/null && echo "  ✅ Healthy" || echo "  ❌ Unhealthy"

# =============================================================================
# QUICK START
# =============================================================================

install: setup docker-up db-migrate ## Quick install (setup + database)
	@echo ""
	@echo "🎉 Installation complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Start backend:  make run-backend"
	@echo "2. Start worker:   make run-worker"
	@echo "3. Start frontend: make run-frontend"
	@echo ""
	@echo "Or use: make dev"

start: dev ## Alias for dev

restart: stop-dev dev ## Restart all services

# =============================================================================
# ENVIRONMENT INFO
# =============================================================================

info: ## Show environment information
	@echo "📋 Environment Information:"
	@echo ""
	@echo "Backend:"
	@echo "  Go version: $(shell go version 2>/dev/null || echo 'Not installed')"
	@echo "  Working directory: backend/"
	@echo ""
	@echo "Frontend:"
	@echo "  Node version: $(shell node --version 2>/dev/null || echo 'Not installed')"
	@echo "  NPM version: $(shell npm --version 2>/dev/null || echo 'Not installed')"
	@echo "  Working directory: frontend/"
	@echo ""
	@echo "Database:"
	@echo "  PostgreSQL container: pickup-postgres"
	@echo "  Host: localhost:5432"
	@echo "  Database: pickup_queue"
	@echo ""
	@echo "URLs:"
	@echo "  API: http://localhost:8080"
	@echo "  API Health: http://localhost:8080/health"
	@echo "  Frontend: http://localhost:3000"

# =============================================================================
# DOCKER COMPOSE COMMANDS
# =============================================================================

.PHONY: compose-up compose-down compose-build compose-logs compose-clean compose-restart

compose-build: ## Build all Docker Compose services
	@echo "🐳 Building all Docker Compose services..."
	docker compose build

compose-up: ## Start all services with Docker Compose (API + DB + Frontend + Worker)
	@echo "🚀 Starting all services with Docker Compose..."
	docker compose up -d
	@echo "✅ All services started!"
	@echo ""
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔌 API: http://localhost:8080"
	@echo "🗄️  Database: localhost:5432"
	@echo "📊 API Health: http://localhost:8080/health"

compose-down: ## Stop all Docker Compose services
	@echo "🛑 Stopping all Docker Compose services..."
	docker compose down

compose-logs: ## View logs from all Docker Compose services
	docker compose logs -f

compose-clean: ## Clean up all Docker Compose resources
	@echo "🧹 Cleaning up all Docker Compose resources..."
	docker compose down -v --rmi all
	docker system prune -f

compose-restart: compose-down compose-up ## Restart all Docker Compose services

compose-status: ## Check status of all services
	docker compose ps

# Start only specific services
compose-db: ## Start only database service
	docker compose up -d postgres

compose-backend: ## Start only backend services (database + API + worker)
	docker compose up -d postgres backend worker

# Database access via Docker Compose
compose-db-shell: ## Access PostgreSQL shell via Docker Compose
	docker compose exec postgres psql -U postgres -d pickup_queue

compose-db-migrate: ## Run database migrations via Docker Compose
	docker compose exec postgres psql -U postgres -d pickup_queue -f /docker-entrypoint-initdb.d/001_create_packages_table.sql
