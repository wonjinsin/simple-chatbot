# Go HTTP Boilerplate

A production-ready HTTP API boilerplate built with Go, following Clean Architecture principles with clear layer separation, custom error handling, and structured logging.

## ✨ Features

- **Clean Architecture** with clear layer separation (Domain, Repository, UseCase, Handler)
- **Custom Error Handling** system with 4-digit error codes aligned to HTTP status codes
- **Structured Logging** with TrID (Transaction ID) tracking using Zerolog
- **EntGo ORM** with PostgreSQL for type-safe database operations
- **Database Migrations** with go-migrate for version-controlled schema changes
- **Docker Compose** for local development infrastructure
- **HTTP Middleware Stack** (CORS, logging, request ID, recovery)
- **Standard JSON Response Format** with transaction ID and status codes
- **Environment-based Configuration** with fail-fast validation

## 🛠 Tech Stack

- **Go** 1.21+
- **Chi Router** - Lightweight, idiomatic HTTP router
- **EntGo** - Type-safe ORM with code generation
- **PostgreSQL** - Primary database
- **Zerolog** - Structured JSON logging
- **go-migrate** - Database migration tool
- **Docker Compose** - Local infrastructure management

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for convenience commands)

## 🚀 Quick Start

### Installation & Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd go-boilerplate
   ```

2. **Install Go dependencies**

   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure environment**

   The `.env.local` file is already included in the project with default settings:

   ```
   PORT=8080
   ENV=local
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=go_boilerplate
   DB_SSLMODE=disable
   ```

   You can modify these values if needed for your local environment.

4. **Start infrastructure (PostgreSQL)**

   ```bash
   make infra-up
   # Wait for PostgreSQL to be ready (5-10 seconds)
   ```

5. **Run database migrations**

   ```bash
   make migrate-up
   ```

6. **Build the application**

   ```bash
   make build
   ```

7. **Run the server**

   ```bash
   ./bin/server
   # or for development: go run cmd/server/main.go
   ```

   You should see:

   ```
   [ASCII art banner]
   Configuration loaded: ENV=local, PORT=8080, DB=postgres@localhost:5432/go_boilerplate
   HTTP server starting on :8080
   ```

### Testing the API

Once the server is running, test the endpoints using curl:

#### 1. Health Check

```bash
curl http://localhost:8080/healthz
```

Expected response:

```json
{
  "trid": "2025102616501424416161",
  "code": "0200",
  "result": {
    "status": "ok"
  }
}
```

#### 2. Create a User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

Expected response:

```json
{
  "trid": "2025102616501424416162",
  "code": "0201",
  "result": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-10-26T16:50:14.241Z"
  }
}
```

#### 3. Get User by ID

```bash
curl http://localhost:8080/users/1
```

Expected response:

```json
{
  "trid": "2025102616501424416163",
  "code": "0200",
  "result": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-10-26T16:50:14.241Z"
  }
}
```

#### 4. List Users (with pagination)

```bash
curl "http://localhost:8080/users?offset=0&limit=10"
```

Expected response:

```json
{
  "trid": "2025102616501424416164",
  "code": "0200",
  "result": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-10-26T16:50:14.241Z"
      }
    ],
    "total": 1,
    "offset": 0,
    "limit": 10
  }
}
```

#### 5. Test Error Handling

Try creating a duplicate user:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

Expected error response (409 Conflict):

```json
{
  "trid": "2025102616501424416165",
  "code": "0409",
  "result": {
    "msg": "duplicate email"
  }
}
```

Try getting a non-existent user:

```bash
curl http://localhost:8080/users/999
```

Expected error response (404 Not Found):

```json
{
  "trid": "2025102616501424416166",
  "code": "0404",
  "result": {
    "msg": "failed to get user: user not found: ent: user not found"
  }
}
```

## 📁 Project Structure

```
├── cmd/                    # Application entrypoints
│   ├── server/            # HTTP server
│   └── migrate/           # Migration runner
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── constants/        # Error codes and constants
│   ├── domain/           # Domain entities and business logic
│   ├── handler/          # HTTP handlers and middleware
│   │   └── http/
│   │       ├── dto/      # Data transfer objects
│   │       └── middleware/ # HTTP middleware
│   ├── repository/       # Data access layer
│   │   ├── repository.go # Repository interfaces
│   │   └── postgres/     # PostgreSQL implementation
│   │       ├── dao/      # EntGo generated code
│   │       └── user_mapper.go # Domain-DAO mapping
│   ├── usecase/          # Business use cases
│   │   ├── service.go    # Service interfaces
│   │   └── user_service.go # Service implementation
│   └── shared/           # Shared utilities
├── pkg/                   # Public reusable packages
│   ├── constants/        # Shared constants
│   ├── errors/           # Custom error system
│   ├── logger/           # Logging utilities
│   └── utils/            # Common utilities
├── migrations/            # Database migrations
├── test/                  # Tests and mocks
├── docker-compose.yml     # Local infrastructure
├── Makefile              # Development commands
└── .env.local            # Local configuration
```

## ⚙️ Configuration

Environment variables (all required):

| Variable      | Description                    | Example          |
| ------------- | ------------------------------ | ---------------- |
| `PORT`        | Server port                    | `8080`           |
| `ENV`         | Environment (local, dev, prod) | `local`          |
| `DB_HOST`     | PostgreSQL host                | `localhost`      |
| `DB_PORT`     | PostgreSQL port                | `5432`           |
| `DB_USER`     | Database user                  | `postgres`       |
| `DB_PASSWORD` | Database password              | `postgres`       |
| `DB_NAME`     | Database name                  | `go_boilerplate` |
| `DB_SSLMODE`  | SSL mode                       | `disable`        |

## 🔧 Development Guide

### Database Management

```bash
# Start PostgreSQL
make infra-up

# Stop PostgreSQL
make infra-down

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration version
make migrate-version
```

### Build and Run

```bash
# Build the application
make build

# Run the server
./bin/server
# or
make start

# Development mode (hot reload)
go run cmd/server/main.go
```

### Testing

```bash
# Run tests
make test

# Generate mocks
make build-mocks

# Run full test suite (tests + vet + fmt + lint)
make test-all
```

### Code Quality

```bash
# Format code
make fmt

# Vet code
make vet

# Lint code
make lint
```

### Development Tools

```bash
# Install all development tools
make tool

# Clean build artifacts
make clean
```

## 📡 API Documentation

### Endpoints

| Method | Path          | Description            | Auth |
| ------ | ------------- | ---------------------- | ---- |
| `GET`  | `/healthz`    | Health check           | No   |
| `POST` | `/users`      | Create user            | No   |
| `GET`  | `/users`      | List users (paginated) | No   |
| `GET`  | `/users/{id}` | Get user by ID         | No   |

### Request/Response Format

All responses follow a standard format:

```json
{
  "trid": "string", // Transaction ID for request tracing
  "code": "string", // 4-digit status code (e.g., "0200", "0404")
  "result": {} // Response data or error message
}
```

#### Success Response Example

```json
{
  "trid": "2025102616501424416161",
  "code": "0200",
  "result": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

#### Error Response Example

```json
{
  "trid": "2025102616501424416162",
  "code": "0404",
  "result": {
    "msg": "failed to get user: user not found"
  }
}
```

## 🏗 Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│           Handler Layer                 │
│  (HTTP Controllers, Middleware, DTOs)   │
│  - Request/Response handling            │
│  - Logging and error handling           │
└──────────────┬──────────────────────────┘
               │ depends on
               ▼
┌─────────────────────────────────────────┐
│           UseCase Layer                 │
│     (Business Logic, Services)          │
│  - Application workflows                │
│  - Business orchestration               │
└──────────────┬──────────────────────────┘
               │ depends on
               ▼
┌─────────────────────────────────────────┐
│         Repository Layer                │
│    (Data Access, EntGo, Mappers)        │
│  - Database operations                  │
│  - Domain-DAO transformation            │
└──────────────┬──────────────────────────┘
               │ depends on
               ▼
┌─────────────────────────────────────────┐
│           Domain Layer                  │
│   (Business Entities, Validation)       │
│  - Domain models and rules              │
│  - Pure business logic                  │
└─────────────────────────────────────────┘
```

### Dependency Flow

```
Handler → UseCase → Repository → Domain
  (HTTP)    (Service)  (Data Access)  (Entity)
    ↓           ↓           ↓            ↓
  [Log]    [Orchestrate] [Transform]  [Validate]
```

### Error Handling Strategy

1. **Domain Layer**: Creates errors with validation messages
2. **Repository Layer**: Wraps database errors with context
3. **UseCase Layer**: Wraps errors from repository
4. **Handler Layer**: Extracts error codes and determines HTTP status

```go
// Domain Layer
errors.New(constants.InvalidParameter, "invalid email", nil)

// Repository Layer
errors.Wrap(err, "failed to find user")

// UseCase Layer
errors.Wrap(err, "failed to get user")

// Handler Layer
code := errors.GetCode(err)  // "0404"
httpStatus := http.StatusNotFound  // 404
```

### Logging Strategy

- **Structured JSON logging** with Zerolog
- **TrID (Transaction ID)** in all logs for request correlation
- **Log only at handler layer** (controllers, middleware)
- **No logging in usecase or repository layers**

Log format:

```json
{
  "level": "info",
  "trid": "2025102616501424416161",
  "time": "2025/01/01 01:01:01.333",
  "message": "user created successfully"
}
```

## 📐 Conventions

Key development conventions:

1. **Custom Collection Types**: Always use named slice types

   ```go
   // Define in domain
   type Users []*User

   // Use everywhere
   func GetUsers() (Users, error)
   ```

2. **Error Handling**: Use custom error system

   ```go
   // Create error
   errors.New(constants.NotFound, "user not found", err)

   // Wrap error
   errors.Wrap(err, "failed to get user")
   ```

3. **UTC Timezone**: Set globally in main.go

   ```go
   time.Local = time.UTC
   ```

4. **Interface Definitions**:

   - Repository interfaces in `internal/repository/`
   - Service interfaces in `internal/usecase/`

5. **DTO Mapping**: Separate mapper files
   - `user_mapper.go` for domain → DTO transformations
   - Keep DTOs in `internal/handler/http/dto/`

## 🔨 Development Commands

| Command                | Description                              |
| ---------------------- | ---------------------------------------- |
| `make infra-up`        | Start Docker infrastructure (PostgreSQL) |
| `make infra-down`      | Stop Docker infrastructure               |
| `make migrate-up`      | Run database migrations                  |
| `make migrate-down`    | Rollback database migrations             |
| `make migrate-version` | Check current migration version          |
| `make build`           | Build the application                    |
| `make start`           | Run the built binary                     |
| `make test`            | Run unit tests                           |
| `make build-mocks`     | Generate mock implementations            |
| `make test-all`        | Run tests + vet + fmt + lint             |
| `make fmt`             | Format code with go fmt                  |
| `make vet`             | Vet code with go vet                     |
| `make lint`            | Lint code with golangci-lint             |
| `make tool`            | Install development tools                |
| `make clean`           | Clean build artifacts                    |

## 🤝 Contributing

1. Follow the conventions in `.cursor/rules/convention.mdc`
2. Write tests for new features
3. Run `make test-all` before committing
4. Use descriptive commit messages
5. Keep functions short and focused
6. Document public APIs with GoDoc comments

---

**Note**: This boilerplate is designed for production use with best practices for maintainability, testability, and observability.
