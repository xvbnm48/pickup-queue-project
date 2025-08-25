# Pickup Queue Management System

A full-stack application for managing package pickup queues with real-time status tracking.

## 🚀 Quick Start with Docker Compose

The easiest way to run the entire application (API + Database + Frontend + Worker) is using Docker Compose:

### Prerequisites
- Docker and Docker Compose installed
- Make (optional, for shortcuts)

### Start All Services
```bash
# Using Docker Compose directly
docker compose up -d

# Or using Make
make compose-up
```

This will start:
- 📊 **Frontend**: http://localhost:3000 (Next.js React app)
- 🔌 **Backend API**: http://localhost:8080 (Go REST API)
- 🗄️ **PostgreSQL Database**: localhost:5432
- ⚙️ **Background Worker**: (Package expiry automation)

### Other Docker Compose Commands
```bash
# Build all images
make compose-build

# View logs
make compose-logs

# Stop all services
make compose-down

# Clean up everything
make compose-clean

# Restart all services
make compose-restart

# Check service status
make compose-status
```

## 🏗️ Architecture

This project follows clean architecture principles with:

### Backend (Golang)

- **Clean Architecture**: Domain, UseCase, Repository, Handler layers
- **REST API**: RESTful endpoints with proper HTTP status codes
- **Database**: PostgreSQL with GORM ORM
- **Background Worker**: Automatic package expiry handling
- **Middleware**: CORS, logging, request ID tracking

### Frontend (TypeScript/React)

- **React 18**: Modern React with hooks
- **TypeScript**: Type safety throughout the application
- **Tailwind CSS**: Utility-first CSS framework
- **React Query**: Server state management and caching
- **Responsive Design**: Mobile-first approach

## 📁 Project Structure

```
pickup-project/
├── backend/
│   ├── cmd/
│   │   ├── api/           # API server entry point
│   │   └── worker/        # Background worker entry point
│   ├── internal/
│   │   ├── domain/        # Business entities and interfaces
│   │   ├── usecase/       # Business logic
│   │   ├── repository/    # Data access layer
│   │   ├── handler/       # HTTP handlers
│   │   └── middleware/    # HTTP middleware
│   ├── pkg/
│   │   ├── database/      # Database configuration
│   │   └── logger/        # Logging utilities
│   └── migrations/        # Database migrations
└── frontend/
    ├── src/
    │   ├── api/           # API client
    │   ├── components/    # Reusable components
    │   ├── pages/         # Page components
    │   ├── types/         # TypeScript type definitions
    │   └── utils/         # Utility functions
    └── public/            # Static assets
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 12+ (or Docker)
- Make (optional, for using Makefile commands)

### Quick Start with Makefile

The project includes comprehensive Makefiles for easy development:

```bash
# Complete setup (installs dependencies + starts database)
make install

# Start full development environment
make dev

# Or start individual services
make run-backend    # Start API server
make run-worker     # Start background worker
make run-frontend   # Start frontend dev server
```

### Manual Setup

If you prefer manual setup or don't have Make installed:

1. **Clone and navigate to backend:**

   ```bash
   cd backend
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Set up environment variables:**

   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Set up PostgreSQL database:**

   ```sql
   CREATE DATABASE pickup_queue;
   ```

5. **Run database migrations:**

   ```bash
   # Run the SQL migration manually or use your preferred migration tool
   psql -d pickup_queue -f migrations/001_create_packages_table.sql
   ```

6. **Start the API server:**

   ```bash
   go run cmd/api/main.go
   ```

7. **Start the background worker (in separate terminal):**

   ```bash
   # Option 1: Run directly with Go
   go run cmd/worker/main.go
   
   # Option 2: Use VS Code Task
   # Press Ctrl+Shift+P, type "Tasks: Run Task", select "Start Package Expiry Worker"
   
   # Option 3: Use provided scripts
   # On Windows:
   ./start-worker.bat
   
   # On Unix/Linux/Mac:
   ./start-worker.sh
   ```

   **Note:** The worker runs continuously and checks for expired packages every hour. Packages that have been in "WAITING" status for more than 24 hours are automatically marked as "EXPIRED".

### Frontend Setup

1. **Navigate to frontend:**

   ```bash
   cd frontend
   ```

2. **Install dependencies:**

   ```bash
   npm install
   ```

3. **Set up environment variables:**

   ```bash
   # Create .env file
   echo "VITE_API_BASE_URL=http://localhost:8080/api/v1" > .env
   ```

4. **Start the development server:**

   ```bash
   npm run dev
   ```

## 📡 API Endpoints

### Packages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/health` | Health check |
| `POST` | `/api/v1/packages` | Create new package |
| `GET` | `/api/v1/packages` | List packages (with pagination and filtering) |
| `GET` | `/api/v1/packages/{id}` | Get package by ID |
| `GET` | `/api/v1/packages/order/{orderRef}` | Get package by order reference |
| `PATCH` | `/api/v1/packages/{id}/status` | Update package status |
| `DELETE` | `/api/v1/packages/{id}` | Delete package |
| `GET` | `/api/v1/packages/stats` | Get package statistics |

### API Examples

#### 1. Create Package

**Request:**

```bash
curl -X POST http://localhost:8080/api/v1/packages \
  -H "Content-Type: application/json" \
  -d '{
    "order_ref": "ORD-20250824-001",
    "driver_code": "DRV-JAKARTA-01"
  }'
```

**Request Body:**

```json
{
  "order_ref": "ORD-20250824-001",
  "driver_code": "DRV-JAKARTA-01"
}
```

**Response (201 Created):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "order_ref": "ORD-20250824-001",
  "driver_code": "DRV-JAKARTA-01",
  "status": "WAITING",
  "created_at": "2025-08-24T15:30:45Z",
  "updated_at": "2025-08-24T15:30:45Z",
  "picked_up_at": null,
  "handed_over_at": null,
  "expired_at": null
}
```

#### 2. List Packages

**Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/packages?limit=10&offset=0&status=WAITING"
```

**Response (200 OK):**

```json
{
  "packages": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "order_ref": "ORD-20250824-001",
      "driver_code": "DRV-JAKARTA-01",
      "status": "WAITING",
      "created_at": "2025-08-24T15:30:45Z",
      "updated_at": "2025-08-24T15:30:45Z",
      "picked_up_at": null,
      "handed_over_at": null,
      "expired_at": null
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "order_ref": "ORD-20250824-002",
      "driver_code": "DRV-BANDUNG-01",
      "status": "WAITING",
      "created_at": "2025-08-24T14:20:30Z",
      "updated_at": "2025-08-24T14:20:30Z",
      "picked_up_at": null,
      "handed_over_at": null,
      "expired_at": null
    }
  ],
  "limit": 10,
  "offset": 0,
  "count": 2
}
```

#### 3. Update Package Status

**Request:**

```bash
curl -X PATCH http://localhost:8080/api/v1/packages/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PICKED_UP"
  }'
```

**Request Body:**

```json
{
  "status": "PICKED_UP"
}
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "order_ref": "ORD-20250824-001",
  "driver_code": "DRV-JAKARTA-01",
  "status": "PICKED_UP",
  "created_at": "2025-08-24T15:30:45Z",
  "updated_at": "2025-08-24T16:15:20Z",
  "picked_up_at": "2025-08-24T16:15:20Z",
  "handed_over_at": null,
  "expired_at": null
}
```

#### 4. Get Package Statistics

**Request:**

```bash
curl -X GET http://localhost:8080/api/v1/packages/stats
```

**Response (200 OK):**

```json
{
  "total": 150,
  "waiting": 45,
  "picked_up": 30,
  "handed_over": 70,
  "expired": 5
}
```

### Sample Data for Testing

Here are some sample package data you can use for testing:

```json
[
  {
    "order_ref": "ORD-20250824-001",
    "driver_code": "DRV-JAKARTA-01"
  },
  {
    "order_ref": "ORD-20250824-002",
    "driver_code": "DRV-BANDUNG-01"
  },
  {
    "order_ref": "ORD-20250824-003",
    "driver_code": "DRV-SURABAYA-01"
  },
  {
    "order_ref": "ORD-20250824-004",
    "driver_code": "DRV-MEDAN-01"
  },
  {
    "order_ref": "ORD-20250824-005",
    "driver_code": "DRV-JAKARTA-02"
  }
]
```

#### Quick Data Seeding

Use the provided scripts to quickly populate your database with sample data:

```bash
# Windows (Command Prompt)
scripts\seed-data.bat

# Windows (PowerShell)
scripts\seed-data.ps1

# Unix/Linux/Mac
scripts/seed-data.sh
```

#### Postman Collection

Import the Postman collection for easy API testing:

1. Open Postman
2. Import `postman/Pickup_Queue_API.postman_collection.json`
3. Set the `baseUrl` variable to `http://localhost:8080/api/v1`
4. Start testing the API endpoints

### Error Responses

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:

- `400` - Bad Request (validation error)
- `404` - Not Found (resource doesn't exist)
- `409` - Conflict (duplicate order reference)
- `500` - Internal Server Error

### Package Status Flow

```
WAITING → PICKED_UP → HANDED_OVER
    ↓         ↓
  EXPIRED   EXPIRED
```

## 🛠️ Makefile Commands

The project includes comprehensive Makefiles for streamlined development:

### Root Makefile Commands

```bash
# Setup & Dependencies
make setup              # Install all dependencies (backend + frontend)
make install            # Quick install (setup + database)

# Development
make dev                # Start full development environment
make dev-api            # Start database and API only
make start              # Alias for dev
make restart            # Restart all services
make stop-dev           # Stop all services

# Database
make docker-up          # Start PostgreSQL with Docker
make docker-down        # Stop PostgreSQL
make db-migrate         # Run database migrations
make db-reset           # Reset database (WARNING: deletes data)
make db-connect         # Connect to database

# Building & Testing
make build              # Build both backend and frontend
make test               # Run all tests
make lint               # Lint all code
make clean              # Clean all build artifacts

# Utilities
make status             # Show status of all services
make health             # Check health of all services
make info               # Show environment information
make help               # Show all available commands
```

### Backend Makefile Commands (backend/)

```bash
# Development
make dev                # Quick start development
make run-api            # Run API server
make run-worker         # Run background worker
make watch              # Watch for changes and restart

# Building & Testing
make build              # Build both API and worker binaries
make test               # Run tests
make test-coverage      # Run tests with coverage
make lint               # Run linting and formatting

# Database
make migrate            # Run database migrations
make migrate-create     # Create new migration file

# Utilities
make install-tools      # Install development tools
make security           # Run security analysis
make env-check          # Check environment setup
```

### Frontend Makefile Commands (frontend/)

```bash
# Development
make dev                # Start development server
make dev-host           # Start with host access
make dev-open           # Start and open browser

# Building & Testing
make build              # Build for production
make preview            # Preview production build
make test               # Run tests
make test-coverage      # Run tests with coverage

# Code Quality
make lint               # Run ESLint
make lint-fix           # Run ESLint with auto-fix
make format             # Format code with Prettier
make type-check         # Run TypeScript type checking

# Analysis
make analyze            # Analyze bundle size
make lighthouse         # Run Lighthouse audit

# Deployment
make deploy-netlify     # Deploy to Netlify
make deploy-vercel      # Deploy to Vercel
```

## 🎯 Features

### Core Features

- ✅ Create packages with order reference and driver code
- ✅ Real-time status tracking (WAITING → PICKED_UP → HANDED_OVER)
- ✅ Automatic package expiry (24-hour rule)
- ✅ Search and filter packages
- ✅ Package statistics dashboard
- ✅ Responsive design for mobile and desktop

### Technical Features

- ✅ Clean architecture with separation of concerns
- ✅ Type-safe API with TypeScript
- ✅ Real-time updates with React Query
- ✅ Comprehensive error handling
- ✅ Input validation and sanitization
- ✅ Database indexing for performance
- ✅ CORS support for cross-origin requests

## 🔧 Configuration

### Backend Configuration (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=pickup_queue
DB_SSL_MODE=disable
PORT=8080
GIN_MODE=debug
WORKER_INTERVAL=1h
```

### Frontend Configuration (.env)

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 🧪 Testing

### Backend Testing

```bash
cd backend
go test ./...
```

### Frontend Testing

```bash
cd frontend
npm test
```

## 📦 Deployment

### Backend Deployment

1. Build the binary:

   ```bash
   go build -o pickup-api cmd/api/main.go
   go build -o pickup-worker cmd/worker/main.go
   ```

2. Run with environment variables:

   ```bash
   ./pickup-api
   ./pickup-worker
   ```

### Frontend Deployment

1. Build the application:

   ```bash
   npm run build
   ```

2. Serve the `dist` folder with any static file server.

## 🔨 Development

### Code Style

- **Backend**: Follow Go conventions with `gofmt` and `golint`
- **Frontend**: ESLint and Prettier for consistent code style

### Adding New Features

1. **Backend**: Add domain entities → repository interfaces → use cases → handlers
2. **Frontend**: Add types → API functions → components → pages

## 🐛 Troubleshooting

### Common Issues

1. **Database connection error:**
   - Verify PostgreSQL is running
   - Check database credentials in `.env`
   - Ensure database exists

2. **Frontend API errors:**
   - Verify backend server is running on correct port
   - Check CORS configuration
   - Verify API base URL in frontend `.env`

3. **Build errors:**
   - Run `go mod tidy` for backend dependencies
   - Run `npm install` for frontend dependencies

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

For questions or issues, please create an issue in the repository.
