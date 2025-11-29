# Copilot Instructions for Weeate

## Architecture Overview

This is a Go backend service following **clean architecture** with strict separation of concerns across four distinct layers:

- **Domain** (`internal/domain/`): Core business entities, interfaces, and domain errors
- **Application** (`internal/app/`): Business logic and command handlers
- **API** (`internal/api/`): HTTP endpoints and request/response handling
- **Infrastructure** (`internal/infrastructure/`): External dependencies (DB, config)

## Key Patterns & Conventions

### Layer Dependencies

**Critical**: Dependencies flow inward only. Infrastructure implements domain interfaces, never the reverse.

```go
// Domain defines interfaces
type UserRepository interface {
    CreateUser(user *User) error
}

// Infrastructure implements them
func NewUserRepository(db *gorm.DB) domain_auth.UserRepository
```

### Import Aliasing Convention

Always use package aliases to distinguish between layers:

```go
api_auth "github.com/SirNacou/weeate/backend/internal/api/auth"
app_auth "github.com/SirNacou/weeate/backend/internal/app/auth"
domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
infra_auth "github.com/SirNacou/weeate/backend/internal/infrastructure/repositories/auth"
```

### Command Handler Pattern

Application layer uses command handlers for business operations:

```go
type RegisterCommand struct { Username, Password, Fullname string }
func (h *RegisterCommandHandler) Handle(ctx context.Context, cmd RegisterCommand) error
```

### Domain Error Handling

Define domain-specific errors in domain layer (`domain/auth/user.go`):

```go
var (
    ErrUserNotFound = errors.New("User not found")
    ErrUsernameAlreadyExist = errors.New("Username already exists")
)
```

## Development Workflow

### Local Development

```bash
# Start services
docker-compose up

# Direct Go development
go run ./backend/cmd/server/main.go
```

### Database

- PostgreSQL with GORM for ORM
- Auto-migrations run on startup: `db.AutoMigrate(&domain_auth.User{})`
- Repository pattern with transaction support via `WithTx(tx *gorm.DB)`

### Configuration

Environment-based config using `github.com/caarlos0/env/v11`:

- DB settings default to Docker Compose values
- See `internal/infrastructure/configs/config.go` for all options

## When Adding New Features

1. **Start with Domain**: Define entities and interfaces in `internal/domain/`
2. **Add Application Logic**: Create command handlers in `internal/app/`
3. **Implement Infrastructure**: Add repository implementations in `internal/infrastructure/repositories/`
4. **Wire up API**: Create endpoints in `internal/api/` and register routes
5. **Update main.go**: Wire dependencies and register endpoints

### Module Structure Example

For a new "posts" feature, create:

- `internal/domain/posts/` - Post entity and PostRepository interface
- `internal/app/posts/` - CreatePostCommandHandler, etc.
- `internal/infrastructure/repositories/posts/` - GORM implementation
- `internal/api/posts/` - HTTP handlers and route registration

The dependency injection happens in `main.go` following the existing auth pattern.
