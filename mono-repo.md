# Mono Repo Structure

This project is structured as a mono repository supporting multiple services with shared infrastructure.

## Structure

```
├── cmd/api/                    # Main API service
├── services/
│   ├── user-service/          # User management service
│   ├── auth-service/          # Authentication service (future)
│   └── notification-service/  # Notification service (future)
├── shared/                    # Shared libraries and utilities
│   ├── pkg/                   # Shared packages
│   ├── proto/                 # Protocol buffer definitions
│   └── configs/               # Shared configurations
├── deployments/               # Deployment configurations
├── scripts/                   # Build and deployment scripts
└── docker-compose.yaml        # Multi-service orchestration
```

## Services

### Main API (Port 8080)
- Clean Architecture implementation
- User authentication and management
- Integrated with Loki logging

### User Service (Port 8081)
- Dedicated user management service
- Microservice example
- Separate database schema

## Development

### Running All Services
```bash
make docker-run
```

### Running Individual Services
```bash
# Main API
make run

# User Service
cd services/user-service && go run cmd/api/main.go
```

### Adding New Services

1. Create service directory under `services/`
2. Copy template structure from existing service
3. Add service to `docker-compose.yaml`
4. Update Makefile with service-specific commands

## Logging Integration

All services are configured to send logs to your existing Loki instance:
- Loki URL: http://localhost:3100
- Grafana Dashboard: http://localhost:3000
- Structured logging with service labels
- Automatic log aggregation via Promtail

## Service Discovery

Services communicate via:
- Docker network (`app-network`)
- Environment-based configuration
- Health check endpoints (`/health`)

## Shared Dependencies

Common dependencies are managed at the root level:
- Database connections
- Logging configuration
- Authentication middleware
- Response utilities