# Auth Service - System Flow

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Request                                  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Middleware Chain                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Correlation  │─▶│   Recovery   │─▶│   Logging    │─▶│   Metrics    │     │
│  │     ID       │  │  (Panic)     │  │   (Zap)      │  │ (Prometheus) │     │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                               │              │
│                                                               ▼              │
│                                                      ┌──────────────┐        │
│                                                      │ Rate Limiter │        │
│                                                      └──────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Presentation Layer                                 │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────┐         │
│  │                         Router (Gin)                           │         │
│  │  GET  /health     ─────────────────────▶ Health Check          │         │
│  │  GET  /metrics    ─────────────────────▶ Prometheus Metrics    │         │
│  │  POST /api/v1/auth/register ───────────▶ AuthHandler.Register  │         │
│  │  POST /api/v1/auth/login ──────────────▶ AuthHandler.Login     │         │
│  └────────────────────────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Application Layer                                  │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────┐         │
│  │                         Use Cases                              │         │
│  │  ┌──────────────────┐  ┌──────────────────┐                    │         │
│  │  │ RegisterUseCase  │  │   LoginUseCase   │  ...               │         │
│  │  └──────────────────┘  └──────────────────┘                    │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                           │                                                  │
│                           ▼                                                  │
│  ┌────────────────────────────────────────────────────────────────┐         │
│  │                      Audit Service                             │         │
│  │  (Async Writer - Buffered Channel - Non-blocking)              │         │
│  └────────────────────────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                             Domain Layer                                     │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │    Entities     │  │  Value Objects  │  │   Exceptions    │              │
│  │  ┌───────────┐  │  │  ┌───────────┐  │  │                 │              │
│  │  │   User    │  │  │  │  UserID   │  │  │  Domain Errors  │              │
│  │  │ AuditLog  │  │  │  │ Username  │  │  │                 │              │
│  │  └───────────┘  │  │  │   Email   │  │  │                 │              │
│  └─────────────────┘  │  └───────────┘  │  └─────────────────┘              │
│                       └─────────────────┘                                    │
│  ┌─────────────────────────────────────────────────────────────┐            │
│  │                    Repository Interfaces                     │            │
│  │  UserRepository, AuditRepository                            │            │
│  └─────────────────────────────────────────────────────────────┘            │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Infrastructure Layer                                │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────┐            │
│  │                   PostgreSQL Repositories                    │            │
│  │  PostgreUserRepo, AuditRepo                                 │            │
│  └─────────────────────────────────────────────────────────────┘            │
│                           │                                                  │
│                           ▼                                                  │
│  ┌─────────────────────────────────────────────────────────────┐            │
│  │                      Database (PostgreSQL)                   │            │
│  │  ┌─────────────┐  ┌─────────────────────┐                   │            │
│  │  │   users     │  │    audit_logs       │                   │            │
│  │  └─────────────┘  └─────────────────────┘                   │            │
│  └─────────────────────────────────────────────────────────────┘            │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Request Flow

### 1. User Registration Flow

```
Client                    Middleware              Handler             UseCase              Repository           Database
  │                          │                      │                    │                    │                    │
  │  POST /register          │                      │                    │                    │                    │
  │─────────────────────────▶│                      │                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                    ┌─────┴─────┐                │                    │                    │                    │
  │                    │ Generate  │                │                    │                    │                    │
  │                    │ Corr. ID  │                │                    │                    │                    │
  │                    └─────┬─────┘                │                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                    ┌─────┴─────┐                │                    │                    │                    │
  │                    │   Log     │                │                    │                    │                    │
  │                    │ Request   │                │                    │                    │                    │
  │                    └─────┬─────┘                │                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                    ┌─────┴─────┐                │                    │                    │                    │
  │                    │  Check    │                │                    │                    │                    │
  │                    │Rate Limit │                │                    │                    │                    │
  │                    └─────┬─────┘                │                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                          │  Bind & Validate     │                    │                    │                    │
  │                          │─────────────────────▶│                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │  Execute(ctx)      │                    │                    │
  │                          │                      │───────────────────▶│                    │                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │                    │ ExistsByUsername   │                    │
  │                          │                      │                    │───────────────────▶│                    │
  │                          │                      │                    │                    │    SELECT          │
  │                          │                      │                    │                    │───────────────────▶│
  │                          │                      │                    │                    │◀───────────────────│
  │                          │                      │                    │◀───────────────────│                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │                    │ ExistsByEmail      │                    │
  │                          │                      │                    │───────────────────▶│                    │
  │                          │                      │                    │                    │    SELECT          │
  │                          │                      │                    │                    │───────────────────▶│
  │                          │                      │                    │                    │◀───────────────────│
  │                          │                      │                    │◀───────────────────│                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │                    │ Create(user)       │                    │
  │                          │                      │                    │───────────────────▶│                    │
  │                          │                      │                    │                    │    INSERT          │
  │                          │                      │                    │                    │───────────────────▶│
  │                          │                      │                    │                    │◀───────────────────│
  │                          │                      │                    │◀───────────────────│                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │                    │ ┌───────────────┐  │                    │
  │                          │                      │                    │ │  Audit Log    │  │                    │
  │                          │                      │                    │ │  (Async)      │──┼───────────────────▶│
  │                          │                      │                    │ └───────────────┘  │                    │
  │                          │                      │                    │                    │                    │
  │                          │                      │◀───────────────────│                    │                    │
  │                          │◀─────────────────────│                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │                    ┌─────┴─────┐                │                    │                    │                    │
  │                    │   Log     │                │                    │                    │                    │
  │                    │ Response  │                │                    │                    │                    │
  │                    └─────┬─────┘                │                    │                    │                    │
  │                          │                      │                    │                    │                    │
  │◀─────────────────────────│                      │                    │                    │                    │
  │   201 Created            │                      │                    │                    │                    │
  │   X-Correlation-ID       │                      │                    │                    │                    │
```

---

## Middleware Chain

| Order | Middleware      | Purpose                                                    |
|-------|-----------------|-----------------------------------------------------------|
| 1     | Correlation ID  | Generate/extract X-Correlation-ID for request tracing    |
| 2     | Recovery        | Catch panics, log stack trace, return 500                 |
| 3     | Logging         | Log request start/end with duration, status, correlation  |
| 4     | Metrics         | Record HTTP metrics (duration, count, status)             |
| 5     | Rate Limiter    | IP-based rate limiting (10 req/s, burst 20)              |

---

## Graceful Shutdown Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                     Application Running                          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ HTTP Server  │  │ Audit Svc    │  │  Database    │          │
│  │  (Running)   │  │  (Running)   │  │ (Connected)  │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ SIGTERM / SIGINT
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Shutdown Initiated                          │
│                                                                  │
│  1. Stop accepting new HTTP connections                         │
│  2. Wait for in-flight requests (30s timeout)                   │
│  3. Stop Audit Service (drain buffer)                           │
│  4. Close Database connections                                   │
│  5. Flush Logger                                                 │
│  6. Exit                                                         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Audit Log Flow (Async)

```
┌─────────────┐      ┌─────────────────────┐      ┌─────────────┐
│   UseCase   │      │    Audit Service    │      │  PostgreSQL │
│             │      │                     │      │             │
│  Create     │      │  ┌───────────────┐  │      │             │
│  AuditLog   │─────▶│  │ Buffered Chan │  │      │             │
│  (Non-block)│      │  │  (size=1000)  │  │      │             │
│             │      │  └───────┬───────┘  │      │             │
└─────────────┘      │          │          │      │             │
                     │          ▼          │      │             │
                     │  ┌───────────────┐  │      │             │
                     │  │   Worker      │  │      │             │
                     │  │  Goroutine    │──┼─────▶│   INSERT    │
                     │  └───────────────┘  │      │             │
                     │                     │      │             │
                     └─────────────────────┘      └─────────────┘
```

---

## Project Structure

```
go-auth-service/
├── cmd/
│   ├── api/
│   │   └── main.go              # Application entry point
│   └── migrate/
│       └── main.go              # Migration CLI
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_audit_logs_table.up.sql
│   └── 000002_create_audit_logs_table.down.sql
├── internal/
│   ├── config/                  # Configuration management
│   ├── pkg/                     # Shared packages
│   │   ├── logger/              # Zap logger wrapper
│   │   ├── correlationid/       # Correlation ID utilities
│   │   └── metrics/             # Prometheus metrics
│   ├── presentation/            # HTTP layer
│   │   └── http/
│   │       ├── handler/         # Request handlers
│   │       ├── middleware/      # HTTP middlewares
│   │       ├── request/         # Request DTOs
│   │       └── router.go        # Route definitions
│   ├── application/             # Business logic
│   │   ├── usecase/             # Use cases
│   │   ├── service/             # Application services
│   │   ├── input/               # Input DTOs
│   │   └── output/              # Output DTOs
│   ├── domain/                  # Domain layer
│   │   ├── entity/              # Entities
│   │   ├── vo/                  # Value Objects
│   │   ├── repository/          # Repository interfaces
│   │   ├── service/             # Domain services
│   │   └── exception/           # Domain exceptions
│   ├── infrastructure/          # Infrastructure layer
│   │   └── persistence/
│   │       └── postgres/        # PostgreSQL implementations
│   └── validation/              # Validation logic
├── docker-compose.yml           # Production Docker setup
├── docker-compose.dev.yml       # Development with hot reload
├── Dockerfile                   # Production Dockerfile
├── Dockerfile.dev               # Development Dockerfile
└── .air.toml                    # Air hot reload config
```

---

## Environment Variables

| Variable                  | Default      | Description                    |
|---------------------------|--------------|--------------------------------|
| `APP_ENV`                 | `local`      | Environment (local/production) |
| `APP_PORT`                | `8000`       | HTTP server port               |
| `LOG_LEVEL`               | `info`       | Log level (debug/info/warn)    |
| `DB_HOST`                 | `postgres`   | PostgreSQL host                |
| `DB_PORT`                 | `5432`       | PostgreSQL port                |
| `DB_USER`                 | `user`       | PostgreSQL user                |
| `DB_PASSWORD`             | `password`   | PostgreSQL password            |
| `DB_NAME`                 | `auth-db`    | PostgreSQL database name       |
| `DB_MAX_OPEN_CONNS`       | `25`         | Max open connections           |
| `DB_MAX_IDLE_CONNS`       | `5`          | Max idle connections           |
| `REDIS_HOST`              | `redis`      | Redis host                     |
| `REDIS_PORT`              | `6379`       | Redis port                     |

---

## Endpoints

| Method | Endpoint                  | Description            | Rate Limited |
|--------|---------------------------|------------------------|--------------|
| GET    | `/health`                 | Health check           | No           |
| GET    | `/metrics`                | Prometheus metrics     | No           |
| POST   | `/api/v1/auth/register`   | User registration      | Yes          |
| POST   | `/api/v1/auth/login`      | User login             | Yes          |
| POST   | `/api/v1/auth/forgot-password` | Password reset    | Yes          |

---

## Development Commands

```bash
# Start development environment with hot reload
docker-compose -f docker-compose.dev.yml up --build

# Run migrations
go run cmd/migrate/main.go -direction up

# Run migrations (rollback)
go run cmd/migrate/main.go -direction down -steps 1

# Check migration version
go run cmd/migrate/main.go -direction version

# Build for production
docker-compose up --build
```
