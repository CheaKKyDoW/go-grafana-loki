# Domain-Based Internal Structure

This project follows a domain-driven design approach where each business domain has its own isolated structure within the `internal` directory.

## Structure Overview

```
internal/
├── user/                           # User domain
│   ├── domain/
│   │   ├── entity/                # User entities and DTOs
│   │   └── repository/            # User repository interfaces
│   ├── usecase/                   # User business logic
│   ├── repository/                # User repository implementations
│   └── delivery/
│       └── http/
│           └── handler/           # User HTTP handlers
├── cart/                          # Cart domain (example)
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   ├── usecase/
│   ├── repository/
│   └── delivery/
│       └── http/
│           └── handler/
├── shared/                        # Shared infrastructure
│   ├── config/                    # Application configuration
│   ├── infrastructure/            # Database, logging, etc.
│   │   ├── database/
│   │   └── logger/
│   └── middleware/                # Shared middleware
└── delivery/
    └── http/
        └── router/                # Main router that orchestrates all domains
```

## Benefits

1. **Domain Isolation**: Each domain is self-contained with its own entities, use cases, and repositories
2. **Easy to Navigate**: Developers can quickly find domain-specific code
3. **Scalable**: New domains can be added without affecting existing ones
4. **Clean Dependencies**: Each domain only depends on its own interfaces and shared infrastructure
5. **Team Organization**: Different teams can work on different domains independently

## Adding a New Domain

To add a new domain (e.g., `order`):

1. **Create domain structure**:
   ```bash
   mkdir -p internal/order/{domain/{entity,repository},usecase,repository,delivery/http/handler}
   ```

2. **Define entities** in `internal/order/domain/entity/`
3. **Define repository interfaces** in `internal/order/domain/repository/`
4. **Implement use cases** in `internal/order/usecase/`
5. **Implement repositories** in `internal/order/repository/`
6. **Create HTTP handlers** in `internal/order/delivery/http/handler/`
7. **Add routes** to `internal/delivery/http/router/router.go`

## Domain Dependencies

- **Domains should NOT depend on each other directly**
- **Use events or shared interfaces for cross-domain communication**
- **All domains can depend on `shared/` infrastructure**
- **External dependencies (pkg/) are available to all domains**

## Example Domain Communication

```go
// Good: Using events for cross-domain communication
type UserCreatedEvent struct {
    UserID uuid.UUID
    Email  string
}

// Good: Using shared interfaces
type NotificationService interface {
    SendEmail(to, subject, body string) error
}

// Bad: Direct domain dependency
// import "go-clean-arch/internal/user/domain/entity"
```

## Migration from Old Structure

The old structure has been reorganized as follows:

- `internal/domain/entity/user.go` → `internal/user/domain/entity/user.go`
- `internal/usecase/user_usecase.go` → `internal/user/usecase/user_usecase.go`
- `internal/repository/user_repository.go` → `internal/user/repository/user_repository.go`
- `internal/delivery/http/handler/user_handler.go` → `internal/user/delivery/http/handler/user_handler.go`
- `internal/config/` → `internal/shared/config/`
- `internal/infrastructure/` → `internal/shared/infrastructure/`
- `internal/delivery/http/middleware/` → `internal/shared/middleware/`