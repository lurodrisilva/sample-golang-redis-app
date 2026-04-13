# sample-golang-redis-app

[![CI](https://github.com/lurodrisilva/sample-golang-redis-app/actions/workflows/ci.yml/badge.svg)](https://github.com/lurodrisilva/sample-golang-redis-app/actions/workflows/ci.yml)

A sample Go REST API demonstrating **Hexagonal (Clean) Architecture** with Redis as the persistence layer.

## Architecture

The project follows strict hexagonal layering with dependency inversion:

```
HTTP Handlers (Inbound Adapters)
        ↓
  Application Use Cases (CQRS)
        ↓
    Domain Layer (Business Logic)
        ↓
  Repository Interface (Port)
        ↓
Redis Repository (Outbound Adapter)
```

### Project Structure

```
├── cmd/api/                          # Application entry point
├── internal/
│   ├── domain/item/                  # Aggregate root, value objects, errors, repository port
│   ├── application/itemapp/          # Use cases: CreateItemHandler, GetItemHandler
│   ├── adapter/
│   │   ├── inbound/http/             # HTTP handlers, router, middleware
│   │   └── outbound/persistence/     # Redis repository implementation
│   └── infrastructure/
│       ├── config/                   # Environment-based configuration
│       └── server/                   # HTTP server with graceful shutdown
├── Dockerfile                        # Multi-stage build (build → test → distroless)
├── Makefile                          # Build, test, and dev automation
└── .github/workflows/ci.yml         # CI pipeline
```

## API Endpoints

| Method | Path | Description | Status Codes |
|--------|------|-------------|--------------|
| `POST` | `/items` | Create a new item | `201`, `400` |
| `GET` | `/items/{id}` | Retrieve an item by ID | `200`, `404`, `500` |
| `GET` | `/health/live` | Liveness probe | `200` |

### Examples

**Create an item:**

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name": "Widget", "description": "A useful widget"}'
```

```json
{"id": "550e8400-e29b-41d4-a716-446655440000"}
```

**Get an item:**

```bash
curl http://localhost:8080/items/550e8400-e29b-41d4-a716-446655440000
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Widget",
  "description": "A useful widget",
  "created_at": "2025-01-15T10:00:00Z"
}
```

## Getting Started

### Prerequisites

- Go 1.26+
- Redis (or Docker)

### Running Locally

```bash
# Start Redis (if not running)
docker run -d -p 6379:6379 redis:alpine

# Run the application
make run
```

The server starts on `:8080` by default.

### Configuration

All settings are loaded from environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | HTTP server port |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | _(empty)_ | Redis password |
| `REDIS_DB` | `0` | Redis database number |

### Docker

```bash
# Build the image
make docker/build

# Run the container
make docker/run
```

The Docker image uses a multi-stage build: compile with `golang:1.24`, run tests during build, and produce a minimal `distroless` runtime image (~50 MB) running as a non-root user.

## Development

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Compile binary to `bin/` |
| `make run` | Run the application |
| `make test` | Run all tests |
| `make test-race` | Run tests with race detector |
| `make test-coverage` | Run tests and enforce 95% coverage |
| `make test-bench` | Run benchmarks |
| `make lint` | Run `go vet` and `staticcheck` |
| `make fmt` | Format code with `gofmt` and `goimports` |
| `make clean` | Remove build artifacts |
| `make docker/build` | Build Docker image |
| `make docker/run` | Run Docker container |

### Pre-commit Hooks

This project uses [pre-commit](https://pre-commit.com/) with [pre-commit-golang](https://github.com/tekwizely/pre-commit-golang) hooks for formatting, linting, security scanning, testing, and build verification.

```bash
# Install pre-commit (macOS)
brew install pre-commit

# Install hooks (one-time)
make precommit/install

# Run all hooks manually
make precommit/run
```

### Testing

Tests exist at every layer with a minimum coverage threshold of **95%**.

- **Domain**: Entity validation, value object parsing
- **Application**: Use case success/failure paths
- **Adapter**: HTTP handler responses, Redis round-trip with [miniredis](https://github.com/alicebob/miniredis)
- **Infrastructure**: Config loading, graceful shutdown

```bash
make test-coverage
```

## CI/CD

GitHub Actions runs on every push and PR to `master`:

1. **Lint** -- `go vet` and `gofmt` checks
2. **Test** -- Full test suite with race detector
3. **Coverage** -- Enforces 95% minimum
4. **Build** -- Compiles the binary
5. **Docker** -- Builds and pushes to `ghcr.io` (master only, tagged with commit SHA and `latest`)
