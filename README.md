# Go Auth Service

A clean architecture authentication service built with Go, featuring user registration, validation, and comprehensive logging.

## Tech Stack

- **Go 1.24** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL** - Primary database
- **Redis** - Caching (planned)
- **Zap** - Structured logging
- **Prometheus** - Metrics collection
- **Docker** - Containerization

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   └── main.go              # Application entry point
│   └── migrate/
│       └── main.go              # Database migration tool
│
├── internal/
│   ├── app/                     # Application bootstrap
│   │   ├── app.go               # Main application setup
│   │   ├── database.go          # Database initialization
│   │   ├── server.go            # HTTP server setup
│   │   └── services.go          # Service container
│   │
│   ├── application/             # Application layer (Use Cases)
│   │   ├── input/               # Input DTOs
│   │   │   └── register.input.go
│   │   ├── output/              # Output DTOs
│   │   │   └── register.output.go
│   │   ├── service/             # Application services
│   │   │   └── audit.service.go
│   │   └── usecase/             # Use cases
│   │       └── register.usecase.go
│   │
│   ├── config/                  # Configuration
│   │   ├── config.go            # Main config loader
│   │   ├── db.config.go         # Database config
│   │   ├── redis.config.go      # Redis config
│   │   └── server.config.go     # Server config
│   │
│   ├── domain/                  # Domain layer (Core Business Logic)
│   │   ├── entity/              # Domain entities
│   │   │   ├── user.go
│   │   │   └── audit_log.go
│   │   ├── repository/          # Repository interfaces
│   │   │   ├── user.repository.go
│   │   │   └── audit.repository.go
│   │   ├── service/             # Domain services
│   │   │   └── uuid.service.go
│   │   ├── vo/                  # Value objects
│   │   │   ├── email.go
│   │   │   ├── user_id.go
│   │   │   └── username.go
│   │   └── exception/           # Domain exceptions
│   │       └── error.go
│   │
│   ├── infrastructure/          # Infrastructure layer
│   │   └── persistence/
│   │       └── postgres/        # PostgreSQL implementations
│   │           ├── db.connect.go
│   │           ├── user.repo.go
│   │           └── audit.repo.go
│   │
│   ├── pkg/                     # Shared packages
│   │   ├── correlationid/       # Request correlation ID
│   │   ├── logger/              # Logging utilities
│   │   └── metrics/             # Prometheus metrics
│   │
│   ├── presentation/            # Presentation layer (HTTP)
│   │   └── http/
│   │       ├── handler/         # HTTP handlers
│   │       │   └── auth.handler.go
│   │       ├── middleware/      # HTTP middleware
│   │       │   ├── correlation.go
│   │       │   ├── logging.go
│   │       │   ├── metrics.go
│   │       │   ├── ratelimit.go
│   │       │   └── recovery.go
│   │       ├── request/         # Request DTOs
│   │       │   └── register.request.go
│   │       └── router.go        # Route definitions
│   │
│   └── validation/              # Input validation
│       ├── messages.go          # Custom validation messages
│       ├── translator.go        # Error translation
│       └── validator.go         # Validator setup
│
├── migrations/                  # Database migrations
│
├── test/
│   └── unit/                    # Unit tests
│       ├── domain/
│       │   ├── entity/
│       │   └── valueobject/
│       └── validation/
│
├── Dockerfile                   # Production Dockerfile
├── Dockerfile.dev               # Development Dockerfile
├── docker-compose.yml           # Production compose
├── docker-compose.dev.yml       # Development compose
├── .air.toml                    # Air hot reload config
├── go.mod                       # Go modules
└── go.sum                       # Go dependencies lock
```

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│              (HTTP Handlers, Middleware, Router)             │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                    Application Layer                         │
│                (Use Cases, Input/Output DTOs)                │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                      Domain Layer                            │
│         (Entities, Value Objects, Repository Interfaces)     │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                   Infrastructure Layer                       │
│            (PostgreSQL Repositories, External Services)      │
└─────────────────────────────────────────────────────────────┘
```

## Getting Started

### Prerequisites

- Docker >= 20.10
- Docker Compose >= 2.0

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/thanhnamdk2710/auth-service.git
cd auth-service
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Start development environment:
```bash
docker compose -f docker-compose.dev.yml up -d
```

4. Run database migrations:
```bash
docker exec -it app-dev go run cmd/migrate/main.go up
```

5. The API will be available at `http://localhost:8000`

## API Endpoints

| Method | Endpoint           | Description         |
|--------|-------------------|---------------------|
| POST   | `/api/v1/register` | User registration   |
| GET    | `/metrics`         | Prometheus metrics  |

## Documentation

- [Docker Guide](GUIDE.md) - Detailed Docker usage instructions

## Running Tests

```bash
# Run all tests
docker exec -it app-dev go test ./...

# Run with coverage
docker exec -it app-dev go test -cover ./...
```

## License

MIT License
