# Docker Usage Guide

## Table of Contents

- [Requirements](#requirements)
- [Environment Configuration](#environment-configuration)
- [Running the Application](#running-the-application)
  - [Development Mode](#development-mode)
  - [Production Mode](#production-mode)
- [Package Management](#package-management)
- [Database Migration](#database-migration)
- [Debug](#debug)
- [Useful Commands](#useful-commands)
- [Troubleshooting Common Errors](#troubleshooting-common-errors)

---

## Requirements

- Docker >= 20.10
- Docker Compose >= 2.0

## Environment Configuration

1. Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

2. Edit the environment variables in `.env` as needed:

```env
APP_ENV=development
APP_PORT=8000
LOG_LEVEL=info

DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=auth-db

REDIS_HOST=redis
REDIS_PORT=6379
```

> **Note**: When running in Docker, `DB_HOST` and `REDIS_HOST` must be the service names in docker-compose (postgres, redis).

---

## Running the Application

### Development Mode

Development mode uses **Air** for hot reload - automatically restarts when code changes.

```bash
# Start all services (postgres, redis, app)
docker compose -f docker-compose.dev.yml up -d

# View logs
docker compose -f docker-compose.dev.yml logs -f app

# Stop all services
docker compose -f docker-compose.dev.yml down
```

**Development Mode Features:**
- Hot reload with Air
- Source code is mounted into the container (code changes are automatically applied)
- Container name: `app-dev`

### Production Mode

```bash
# Build and start
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop all services
docker compose down
```

**Production Mode Features:**
- Multi-stage build optimizes image size
- Uses distroless base image (high security)
- Runs with non-root user
- Container name: `app`

---

## Package Management

### Adding New Packages

```bash
# Method 1: Run go get in the running container
docker exec -it app-dev go get github.com/package/name

# Method 2: Run go get and update go.mod/go.sum
docker exec -it app-dev go get -u github.com/package/name

# Method 3: Use a temporary container
docker compose -f docker-compose.dev.yml run --rm app go get github.com/package/name
```

### Update All Packages

```bash
docker exec -it app-dev go get -u ./...
```

### Clean Up Unused Packages

```bash
docker exec -it app-dev go mod tidy
```

### List All Packages

```bash
docker exec -it app-dev go list -m all
```

---

## Database Migration

### Run Migration

```bash
# Run migration tool
docker exec -it app-dev go run cmd/migrate/main.go up

# Or use a temporary container
docker compose -f docker-compose.dev.yml run --rm app go run cmd/migrate/main.go up
```

### Rollback Migration

```bash
docker exec -it app-dev go run cmd/migrate/main.go down
```

### Access PostgreSQL

```bash
# Connect to postgres container
docker exec -it postgres psql -U user -d auth-db

# Some useful SQL commands
\dt          # List tables
\d+ users    # Show table details for users
\q           # Exit
```

### Access Redis

```bash
docker exec -it redis redis-cli

# Some useful Redis commands
KEYS *       # List all keys
GET key      # Get value of key
FLUSHALL     # Delete all data (be careful!)
```

---

## Debug

### View Logs

```bash
# Logs of all services
docker compose -f docker-compose.dev.yml logs -f

# Logs of specific service
docker compose -f docker-compose.dev.yml logs -f app
docker compose -f docker-compose.dev.yml logs -f postgres
docker compose -f docker-compose.dev.yml logs -f redis

# Logs with limited lines
docker compose -f docker-compose.dev.yml logs --tail=100 app
```

### Access Container Shell

```bash
# Enter app container shell
docker exec -it app-dev sh

# Enter postgres container shell
docker exec -it postgres sh
```

### Check Services Status

```bash
# View containers status
docker compose -f docker-compose.dev.yml ps

# View resource usage
docker stats

# Check postgres health
docker exec -it postgres pg_isready -U user
```

### View Build Errors

When Air build fails, view the log file:

```bash
docker exec -it app-dev cat tmp/build-errors.log
```

---

## Useful Commands

### Docker Compose

```bash
# Rebuild container (when Dockerfile changes)
docker compose -f docker-compose.dev.yml up -d --build

# Restart a service
docker compose -f docker-compose.dev.yml restart app

# Remove volumes (reset database)
docker compose -f docker-compose.dev.yml down -v

# View network
docker network ls
docker network inspect go-auth-service_app-net
```

### Running Tests

```bash
# Run all tests
docker exec -it app-dev go test ./...

# Run tests with verbose
docker exec -it app-dev go test -v ./...

# Run tests for a specific package
docker exec -it app-dev go test -v ./test/unit/domain/entity/...

# Run tests with coverage
docker exec -it app-dev go test -cover ./...
```

### Running Linter

```bash
# Install golangci-lint
docker exec -it app-dev go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
docker exec -it app-dev golangci-lint run
```

### Generate Code

```bash
# If using go generate
docker exec -it app-dev go generate ./...
```

---

## Troubleshooting Common Errors

### 1. Port Already in Use

```
Error: bind: address already in use
```

**Solution:**

```bash
# Find process using the port
lsof -i :8000
lsof -i :5432

# Kill the process or change the port in .env
```

### 2. Database Connection Refused

```
Error: connection refused
```

**Solution:**

```bash
# Check if postgres container is running
docker compose -f docker-compose.dev.yml ps

# Check health status
docker inspect postgres | grep Health -A 10

# Wait for postgres to be healthy then restart app
docker compose -f docker-compose.dev.yml restart app
```

### 3. Go Module Download Failed

```
Error: go mod download timeout
```

**Solution:**

```bash
# Use proxy
docker exec -it app-dev sh -c "GOPROXY=https://proxy.golang.org go mod download"

# Or set proxy in docker-compose
# environment:
#   - GOPROXY=https://proxy.golang.org
```

### 4. Permission Denied

```
Error: permission denied
```

**Solution:**

```bash
# Fix directory permissions
sudo chown -R $(whoami):$(whoami) .

# Or rebuild container
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.dev.yml up -d --build
```

### 5. Container Keeps Restarting

**Solution:**

```bash
# View logs to find the cause
docker compose -f docker-compose.dev.yml logs app

# Check exit code
docker inspect app-dev --format='{{.State.ExitCode}}'
```

### 6. Disk Space Full

```bash
# Clean up Docker
docker system prune -a

# Remove unused volumes
docker volume prune

# Remove unused images
docker image prune -a
```

---

## Docker Files Structure

```
.
├── Dockerfile              # Production build (multi-stage, distroless)
├── Dockerfile.dev          # Development build (with Air hot reload)
├── docker-compose.yml      # Production compose
├── docker-compose.dev.yml  # Development compose
├── .air.toml               # Air configuration for hot reload
└── .env                    # Environment variables
```

---

## Tips

1. **Always use `docker-compose.dev.yml` during development** for hot reload.

2. **Do not commit `.env` file** to git, only commit `.env.example`.

3. **Use `docker compose down -v`** when you want to completely reset the database.

4. **Check logs frequently** when facing issues: `docker compose logs -f`.

5. **Rebuild image** after changing `go.mod` or `Dockerfile`:
   ```bash
   docker compose -f docker-compose.dev.yml up -d --build
   ```
