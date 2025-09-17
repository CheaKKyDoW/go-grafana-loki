# Go Clean Architecture with Fiber & Loki

A production-ready Go mono repository template using Clean Architecture principles, Fiber framework, and integrated with Grafana Loki for centralized logging.

## 🚀 Features

- **Clean Architecture** implementation with proper dependency inversion
- **Fiber Framework** for high-performance HTTP APIs
- **Grafana Loki Integration** for centralized logging and observability
- **Mono Repository** structure supporting multiple services
- **PostgreSQL** database with migrations
- **JWT Authentication** with middleware
- **Docker Compose** orchestration with your existing Loki stack
- **Structured Logging** with automatic Loki forwarding
- **Graceful Shutdown** and health checks

## 📁 Project Structure

```
├── cmd/api/                    # Main API service entrypoint
├── services/
│   ├── user-service/          # Example microservice
│   └── ...                    # Additional services
├── internal/
│   ├── user/                  # User domain
│   │   ├── domain/
│   │   │   ├── entity/       # User entities and DTOs
│   │   │   └── repository/   # User repository interfaces
│   │   ├── usecase/          # User business logic
│   │   ├── repository/       # User repository implementations
│   │   └── delivery/http/    # User HTTP handlers
│   ├── cart/                  # Cart domain (example)
│   │   ├── domain/
│   │   ├── usecase/
│   │   ├── repository/
│   │   └── delivery/http/
│   ├── shared/                # Shared infrastructure
│   │   ├── config/           # Application configuration
│   │   ├── infrastructure/   # Database, logging, etc.
│   │   └── middleware/       # Shared middleware
│   └── delivery/http/router/ # Main router orchestrating all domains
├── pkg/                      # Shared utilities
├── migrations/               # Database migrations
├── configs/                  # Configuration files
├── promtail-config.yml       # Promtail configuration for log collection
└── docker-compose.yaml       # Multi-service orchestration
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25.1+
- Docker & Docker Compose
- Your existing Grafana Loki setup (integrated automatically)

### Quick Start

1. **Clone and setup**:
   ```bash
   cp .env.example .env
   make deps
   ```

2. **Start all services** (integrates with your existing Loki):
   ```bash
   make docker-run
   ```

3. **Access services**:
   - Main API: http://localhost:8080
   - User Service: http://localhost:8081
   - Grafana: http://localhost:3000 (your existing setup)
   - Loki: http://localhost:3100 (your existing setup)

### Local Development

```bash
# Run main API only
make run

# Run user service only
make run-user-service

# Run all services locally
make run-all-local

# View logs in Loki
make logs-loki
```

## 🔌 Loki Integration

The template automatically integrates with your existing Grafana Loki setup:

- **Structured Logs**: All services send JSON logs to Loki
- **Service Labels**: Automatic labeling by service name, environment, version
- **Log Aggregation**: Promtail collects and forwards container logs
- **Grafana Dashboards**: Pre-configured Loki datasource

### Log Labels Applied:
```yaml
service: "go-clean-arch-api" | "go-clean-arch-user-service"
env: "development" | "production"
version: "1.0.0"
container: "container-name"
```

## 📡 API Endpoints

### Main API (Port 8080)
- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/profile` - Get user profile (protected)

### User Service (Port 8081)
- `GET /health` - Health check
- `GET /users` - List users

## 🐳 Docker Services

The `docker-compose.yaml` integrates with your existing Loki stack and adds:

- **go-clean-arch-api**: Main API service
- **go-clean-arch-user-service**: Example microservice
- **postgres**: Shared database
- Connects to your existing **loki**, **promtail**, and **grafana**

## 🔧 Mono Repo Commands

```bash
# Build all services
make build-all

# Test all services
make test-all

# Docker operations
make docker-run          # Start all services
make docker-stop         # Stop all services
make docker-logs         # View all logs

# Database
make migrate-up          # Run migrations
make migrate-down        # Rollback migrations

# Development
make dev-setup          # Initial setup
make fmt                # Format all code
make lint               # Lint all services
```

## 📊 Monitoring & Observability

### Grafana Queries for Your Services:
```logql
# All API logs
{service="go-clean-arch-api"}

# Error logs only
{service="go-clean-arch-api"} |= "ERROR"

# User registration events
{service="go-clean-arch-api"} |= "User created successfully"

# Performance monitoring
rate({service="go-clean-arch-api"}[5m])
```

### Health Checks:
```bash
# Check all services
curl http://localhost:8080/health
curl http://localhost:8081/health
```

## 🏗️ Adding New Services

1. **Create service directory**:
   ```bash
   mkdir -p services/new-service
   cp -r services/user-service/* services/new-service/
   ```

2. **Update docker-compose.yaml**:
   ```yaml
   new-service:
     build: ./services/new-service
     ports:
       - "8082:8080"
     environment:
       - LOKI_URL=http://loki:3100
     networks:
       - app-network
   ```

3. **Add Makefile targets**:
   ```makefile
   run-new-service:
       cd services/new-service && go run cmd/api/main.go
   ```

## 🚀 Production Deployment

1. **Environment Configuration**:
   ```bash
   export LOKI_URL=https://your-loki-instance.com
   export DATABASE_HOST=your-db-host
   export JWT_SECRET=your-production-secret
   ```

2. **Deploy with Docker**:
   ```bash
   docker-compose -f docker-compose.prod.yaml up -d
   ```

3. **Monitor in Grafana**:
   - Service health dashboards
   - Error rate monitoring
   - Performance metrics
   - Log aggregation views

## 🏛️ Architecture Benefits

- **Domain-Driven Design**: Code organized by business domains (user, cart, order, etc.)
- **Clean Architecture**: Testable, maintainable, framework-independent
- **Domain Isolation**: Each domain is self-contained and independent
- **Fiber Performance**: High-throughput HTTP handling
- **Centralized Logging**: All services log to your existing Loki
- **Scalable Structure**: Easy to add new domains without affecting existing ones
- **Team Organization**: Different teams can work on different domains independently

## 🏗️ Domain-Based Architecture

The project follows a domain-driven approach where each business domain has its own isolated structure:

### Current Domains
- **User Domain** (`internal/user/`): Authentication, user management
- **Cart Domain** (`internal/cart/`): Shopping cart functionality (example)

### Adding New Domains
```bash
# Create new domain structure
mkdir -p internal/order/{domain/{entity,repository},usecase,repository,delivery/http/handler}

# Add routes in internal/delivery/http/router/router.go
func setupOrderRoutes(api fiber.Router, orderHandler *orderHandler.OrderHandler, authMiddleware *middleware.AuthMiddleware) {
    orders := api.Group("/orders")
    orders.Use(authMiddleware.RequireAuth())
    orders.Post("/", orderHandler.CreateOrder)
    orders.Get("/:id", orderHandler.GetOrder)
}
```

### Domain Benefits
- **Isolation**: Each domain is self-contained
- **Scalability**: Add new domains without affecting existing ones
- **Team Collaboration**: Different teams can work on different domains
- **Clear Boundaries**: Business logic is organized by domain context

This template gives you a production-ready foundation for building scalable microservices with proper domain organization and observability integration into your existing Grafana Loki stack.