# Pickup Queue Management System

A modern web application for managing package pickup queues with real-time status tracking and intuitive user interface.

## Features

- **Package Management**: Create, view, update, and delete packages
- **Status Tracking**: Track packages through different states (WAITING, PICKED_UP, HANDED_OVER, EXPIRED)
- **Real-time Updates**: Live statistics and status updates
- **Filtering**: Filter packages by status
- **Responsive Design**: Modern, mobile-friendly interface
- **API Integration**: RESTful API with comprehensive endpoints

## Architecture

### Backend (Go)
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Architecture**: Clean architecture with domain, usecase, repository layers
- **API**: RESTful endpoints with proper error handling
- **Worker**: Background worker for handling expired packages

### Frontend (React)
- **Framework**: React 18 with TypeScript
- **State Management**: React Query for server state
- **UI Components**: Custom components with Tailwind CSS
- **Icons**: Lucide React icons
- **Notifications**: React Hot Toast

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository
2. Run the application:
   ```bash
   docker-compose up -d
   ```
3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health

### Manual Setup

#### Backend Setup
1. Navigate to backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```
4. Run migrations:
   ```bash
   make migrate-up
   ```
5. Start the API server:
   ```bash
   make run-api
   ```
6. Start the worker (optional):
   ```bash
   make run-worker
   ```

#### Frontend Setup
1. Navigate to frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Set up environment variables:
   ```bash
   cp env.example .env
   # Edit .env with your API URL
   ```
4. Start the development server:
   ```bash
   npm start
   ```

## API Endpoints

### Packages
- `GET /api/v1/packages` - List packages with optional filtering
- `POST /api/v1/packages` - Create a new package
- `GET /api/v1/packages/:id` - Get package by ID
- `GET /api/v1/packages/order/:orderRef` - Get package by order reference
- `PATCH /api/v1/packages/:id/status` - Update package status
- `DELETE /api/v1/packages/:id` - Delete package
- `GET /api/v1/packages/stats` - Get package statistics

### Health Check
- `GET /health` - API health check

## Package Status Flow

```
WAITING → PICKED_UP → HANDED_OVER
   ↓           ↓
EXPIRED    EXPIRED
```

### Status Descriptions
- **WAITING**: Package is waiting to be picked up
- **PICKED_UP**: Package has been picked up by driver
- **HANDED_OVER**: Package has been delivered to recipient
- **EXPIRED**: Package pickup/delivery has expired

## Environment Variables

### Backend (.env)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=pickup_queue
DB_SSL_MODE=disable
PORT=8080
GIN_MODE=release
```

### Frontend (.env)
```
REACT_APP_API_URL=http://localhost:8080/api/v1
```

## Development

### Backend Commands
```bash
# Run API server
make run-api

# Run worker
make run-worker

# Run tests
make test

# Database migrations
make migrate-up
make migrate-down

# Build
make build
```

### Frontend Commands
```bash
# Start development server
npm start

# Build for production
npm run build

# Run tests
npm test
```

## Testing

### API Testing
Use the provided Postman collection in `postman/Pickup_Queue_API.postman_collection.json` to test all API endpoints.

### Seed Data
Run the seed script to populate the database with sample data:
```bash
# Windows
scripts/seed-data.bat

# Linux/Mac
scripts/seed-data.sh

# PowerShell
scripts/seed-data.ps1
```

## Production Deployment

### Docker Deployment
1. Build and run with Docker Compose:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Manual Deployment
1. Build the backend:
   ```bash
   cd backend && make build
   ```
2. Build the frontend:
   ```bash
   cd frontend && npm run build
   ```
3. Deploy the built artifacts to your server

## Design Decisions

### Status Transition Rules
- Packages can only transition to valid next states
- Once HANDED_OVER or EXPIRED, no further transitions allowed
- Worker automatically expires packages based on business rules

### UI/UX Considerations
- Clean, modern interface matching the visual reference
- Real-time updates for better user experience
- Responsive design for mobile and desktop
- Clear status indicators with color coding
- Intuitive modal dialogs for actions

### Technical Choices
- **Go + Gin**: Fast, lightweight backend with excellent performance
- **PostgreSQL**: Reliable, ACID-compliant database
- **React + TypeScript**: Type-safe frontend development
- **React Query**: Efficient server state management
- **Tailwind CSS**: Utility-first CSS for rapid UI development
- **Docker**: Containerized deployment for consistency

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
