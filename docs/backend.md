# A Canonical Repository Structure for Production-Grade Go Web Services

## Table of Contents

1. [Introduction: A Philosophy of Structure](#introduction-a-philosophy-of-structure)
2. [The Anatomy of the Repository](#the-anatomy-of-the-repository)
3. [The Application Core: Services, Handlers, and Routing](#the-application-core-services-handlers-and-routing)
4. [The Data Layer: PostgreSQL with goose and sqlc](#the-data-layer-postgresql-with-goose-and-sqlc)
5. [Full-Stack Observability: OpenTelemetry (otel) and slog](#full-stack-observability-opentelemetry-otel-and-slog)
6. [The Developer Workflow: Taskfile and Air](#the-developer-workflow-taskfile-and-air)
7. [Conclusion](#conclusion)

---

## Introduction: A Philosophy of Structure

This document presents a canonical standard for structuring a production-grade REST web service in Go. The foundational philosophy guiding this architecture is that the optimal structure for a Go application is not a rigid, one-size-fits-all template but an evolvable framework that balances initial simplicity with the capacity to manage growing complexity.

The approach is deeply rooted in idiomatic Go, which favors clarity, convention, and leveraging the power of the standard library over complex frameworks or configuration. It is designed to serve as an authoritative guide for building maintainable, scalable, and observable systems.

### Guiding Principles

This architecture is built upon a set of core principles that inform every structural and technological choice:

#### 1. Idiomatic Go
- The structure adheres strictly to Go's unique properties and conventions
- Uses packages for encapsulation rather than classes
- Prefers explicit dependency injection over global state
- Leverages the standard library wherever feasible
- Avoids patterns common in other object-oriented languages that don't translate well to Go's design philosophy

#### 2. Start Simple, Scale Gracefully
- The structure must be approachable for a new project yet robust enough for a large, multi-team application
- Provides clear boundaries from day one without imposing unnecessary complexity
- Balances the debate between minimal flat structure and comprehensive layout from the outset

#### 3. Separation of Concerns
- Clear and enforced delineation between business logic, infrastructure code, and application delivery mechanisms
- Paramount for achieving high testability, enhancing maintainability, and enabling parallel development efforts

#### 4. Feature-Oriented (Vertical Slice) Architecture
- Code organized by business domain or feature rather than by technical layer
- Reduces coupling between features, improves context locality for developers
- Significantly simplifies future refactoring, such as decomposing the monolith into microservices

### The Structure Debate

The ongoing debate surrounding Go project structure, exemplified by the controversy around the `golang-standards/project-layout` repository, is fundamentally a conflict between two philosophies:

- **"You Ain't Gonna Need It" (YAGNI)**: Advocates for minimal structure until pain is explicitly felt
- **"Scalability First"**: Argues for establishing strong architectural boundaries early to prevent the project from devolving into a "big ball of mud"

This canonical structure synthesizes these ideas into a more effective model:
- Employs the standard `cmd/internal` pattern to establish high-level separation of concerns
- Within the `internal` directory, code is organized by feature (vertical slicing)
- Encapsulates a domain's handler, service logic, and repository access within a single, cohesive package
- Provides a clear and scalable path for evolution from single-feature service to well-organized application

---

## The Anatomy of the Repository

This section provides a detailed blueprint of the repository layout. Each directory's purpose is explained and justified with references to community conventions and Go language features.

### The Root Directory

The root of the project serves as the primary location for project-level metadata, configuration, and developer tooling. It is intentionally kept clean of Go source files to avoid clutter from non-Go assets and build artifacts.

#### Root-Level Files

| File | Purpose |
|------|---------|
| `go.mod` & `go.sum` | Define the Go module and manage its dependencies |
| `Taskfile.yml` | Entry point for the developer workflow, automating common tasks |
| `.air.toml` | Configuration for the air live-reloading tool |
| `.gitignore` | Specifies files and directories to be ignored by version control |
| `README.md` | Provides an overview of the project, setup instructions, and documentation |
| `.env` | Contains local environment variables (should be in `.gitignore`) |

### `/cmd/server/main.go`

This directory contains the single, unambiguous entry point for the application binary. The name of the subdirectory, `server`, directly corresponds to the name of the executable that will be produced.

#### Responsibilities

The main package is designed to be minimal, acting as an orchestrator that wires together the various components:

1. **Loading application configuration** from files and environment variables
2. **Initializing core dependencies** (logger, database connection pool, OpenTelemetry SDK)
3. **Instantiating services and handlers** by injecting their dependencies
4. **Defining the HTTP server**, registering routes, and applying global middleware
5. **Starting the server** and managing graceful shutdown

This structure enforces a clean separation between the application's startup and configuration logic and its core business logic.

### `/internal/`: The Application's Core

This directory houses all private application and library code. The Go compiler enforces the privacy of any package within an `internal` directory, preventing it from being imported by projects outside of the current module.

#### `/internal/config/`

This package manages all application configuration:

- Defines a `Config` struct that holds all configuration values
- Provides a `Load()` function that reads configuration from sources like `.env` files
- Centralizes all configuration logic into a single, predictable location

#### `/internal/server/`

This package encapsulates all logic related to the HTTP server:

- Defines a `Server` struct that holds dependencies (router, logger, application services)
- Contains methods for registering routes (`routes.go`)
- Defines and applies middleware (`middleware.go`)
- Manages the server lifecycle

#### `/internal/platform/`

This directory contains packages that provide foundational capabilities or interact with external systems:

##### `/internal/platform/database/`
- Handles all PostgreSQL connection logic
- Contains a `Connect()` function that initializes a `pgxpool` connection pool
- Ensures the connection pool is created only once at application startup

##### `/internal/platform/telemetry/`
- Responsible for initializing the OpenTelemetry SDK
- Contains functions like `InitTracerProvider()` and `InitMeterProvider()`
- Configures OTLP exporter, defines service resource, and registers global providers

#### `/internal/<feature>/` (e.g., `/internal/user/`)

This is the heart of the feature-oriented architecture. Each subdirectory represents a "vertical slice" or "domain," containing all the logic for a specific business feature.

**Example structure for `/internal/user/`:**

| File | Purpose |
|------|---------|
| `handler.go` | `net/http` handlers for the user feature |
| `service.go` | Core business logic, decoupled from HTTP and database concerns |
| `repository.go` | Data access layer interface and implementation |
| `models.go` | Core Go structs for the feature's domain |

### `/db/migrations/`

This directory stores all database schema migrations as timestamped SQL files:

- Managed by the `goose` tool
- Provides a version-controlled, auditable history of the database schema
- Example file name: `20240101120000_create_users_table.sql`

### `/sql/`

This directory provides a dedicated home for all SQL source code that `sqlc` uses for code generation:

- **`/sql/schema.sql`**: Single file containing `CREATE TABLE`, `CREATE TYPE`, and other DDL statements
- **`/sql/queries/user.sql`**: Contains all SQL queries related to the user feature, each annotated with `sqlc` comments

### The `/pkg` Directory: An Explicit Omission

This project structure deliberately omits the `/pkg` directory. The `/pkg` convention is intended for code that will be imported and used by external Go projects. For a self-contained web service, all code is internal to the application. The compiler-enforced privacy of the `/internal` directory is the more precise and idiomatic choice.

---

## The Application Core: Services, Handlers, and Routing

The core of the application is defined by how its components are instantiated, how they communicate, and how incoming requests are directed to the appropriate logic.

### Dependency Injection via Constructors

A clean, explicit dependency injection pattern is used throughout the application, with the main function in `/cmd/server/main.go` acting as the central "injector" or composer.

#### Instantiation Flow

```go
// 1. Initialize platform-level dependencies
dbPool := database.Connect(config.Database)
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// 2. Construct data access layer
userRepo := user.NewRepository(dbPool)

// 3. Construct business logic layer
userService := user.NewService(userRepo)

// 4. Construct presentation layer
userHandler := user.NewHandler(userService)
```

This chain of constructor functions creates a clear, compile-time-verified dependency graph.

### HTTP Routing with net/http (Go 1.22+)

This architecture leverages the enhanced `net/http.ServeMux` introduced in Go 1.22, which supports method-based routing and path wildcards.

#### Route Registration Example

```go
// Located in internal/server/routes.go
package server

import "net/http"

// registerRoutes sets up the routing for the application.
func (s *Server) registerRoutes() {
    mux := http.NewServeMux()

    // Health check endpoint
    mux.HandleFunc("GET /health", s.handleHealthCheck())

    // API v1 routes
    v1Mux := http.NewServeMux()
    userHandler := s.userHandler 
    v1Mux.HandleFunc("GET /users/{id}", userHandler.HandleGetUser())
    v1Mux.HandleFunc("POST /users", userHandler.HandleCreateUser())

    // Mount the v1 sub-router under the /api/v1/ prefix
    mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))

    s.router = mux
}
```

### Idiomatic HTTP Handlers

Handlers are implemented following idiomatic Go patterns:

#### Handler Structs
- For each feature, a `Handler` struct is defined (e.g., `user.Handler`)
- Holds the feature's service as a dependency, injected via constructor

#### Closure-Based Methods
- Each handler action (e.g., `HandleGetUser`) is a method on the `Handler` struct
- Returns an `http.HandlerFunc`
- Uses closure to capture handler dependencies

#### Responsibilities
- Decode JSON request body into Go structs
- Validate input
- Call appropriate service-layer method
- Write JSON response with correct HTTP status code

---

## The Data Layer: PostgreSQL with goose and sqlc

The data layer is engineered for performance, type safety, and maintainability by combining best-in-class open-source tools.

### Database Connection Management

#### Tooling
- Uses the high-performance `pgx/v5` driver for PostgreSQL
- Specifically uses `pgxpool` package for robust and efficient connection pooling

#### Initialization
- Database connection pool is initialized once at application startup
- Logic resides in `/internal/platform/database/` package
- Resulting `*pgxpool.Pool` object is passed down via dependency injection

#### Configuration Parameters

```go
// Essential connection pool parameters
pool.SetMaxOpenConns(25)        // Total connections to database
pool.SetMaxIdleConns(5)         // Idle connections to keep open
pool.SetConnMaxLifetime(5*time.Minute)  // Max reuse time
pool.SetConnMaxIdleTime(1*time.Minute)  // Max idle time
```

### Schema Migrations with goose

Database schema changes are managed declaratively and version-controlled using `goose`.

#### Workflow

1. **Install goose CLI:**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. **Create migration:**
   ```bash
   goose -dir db/migrations create create_users_table sql
   ```

3. **Define schema changes:**
   ```sql
   -- +goose Up
   CREATE TABLE users (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       email VARCHAR(255) UNIQUE NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   );

   -- +goose Down
   DROP TABLE users;
   ```

4. **Apply migrations:**
   ```bash
   goose -dir db/migrations postgres "user=... password=... dbname=... sslmode=disable" up
   ```

#### Best Practices
- Migrations should be small, atomic, and reversible
- For complex schema changes, use multi-step migrations
- Always test both up and down migrations

### Type-Safe Data Access with sqlc

`sqlc` generates fully type-safe, idiomatic Go code from raw SQL queries.

#### Configuration (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/queries/"
    schema: "sql/schema.sql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
```

#### Workflow

1. **Define schema** in `/sql/schema.sql`
2. **Write queries** in feature-specific files (e.g., `/sql/queries/user.sql`)
3. **Generate code:**
   ```bash
   sqlc generate
   ```

#### Query Example

```sql
-- name: GetUserByID :one
SELECT id, email, created_at 
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email) 
VALUES ($1) 
RETURNING id, email, created_at;
```

#### Integration

```go
// Repository implementation
type Repository struct {
    *db.Queries
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
    return &Repository{
        Queries: db.New(db),
        db:      db,
    }
}
```

### Data Access Technology Comparison

| Method | Pros | Cons | Why sqlc is Chosen |
|--------|------|------|-------------------|
| Raw `database/sql` | Maximum performance and control; no dependencies | Verbose, error-prone due to manual row scanning; no compile-time type safety | sqlc automates boilerplate while adding compile-time type safety |
| sqlx | Light wrapper over `database/sql` that simplifies scanning | Still relies on runtime reflection; SQL errors discovered at runtime | sqlc provides compile-time validation of SQL queries |
| ORMs (e.g., GORM) | Fast for simple CRUD by abstracting away SQL | Performance overhead, "magic" behavior, less idiomatic Go | sqlc embraces SQL rather than hiding it |
| **sqlc (Chosen)** | **Compile-time type safety, high performance, explicit SQL control** | **Requires developers to write raw SQL** | **Perfect alignment with Go's philosophy** |

---

## Full-Stack Observability: OpenTelemetry (otel) and slog

A production-grade service requires robust observability. This architecture integrates a modern, best-in-class stack for logging, tracing, and metrics.

### Structured Logging with log/slog

As of Go 1.21, the `log/slog` package is the official standard for structured logging.

#### Implementation

```go
// In main.go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Create base logger with common attributes
baseLogger := logger.With(
    "service.name", "my-service",
    "service.version", "1.0.0",
)

// Inject logger into components
server := server.New(baseLogger, userHandler)
```

#### Usage in Components

```go
func (s *Service) CreateUser(ctx context.Context, user *User) error {
    logger := slog.Ctx(ctx)
    logger.Info("creating new user", "email", user.Email)
    
    // ... business logic ...
    
    logger.Info("user created successfully", "user_id", user.ID)
    return nil
}
```

### OpenTelemetry SDK Initialization

All OpenTelemetry SDK setup is centralized in `/internal/platform/telemetry/`.

#### Components

```go
// Exporter configuration
exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"))

// Resource definition
resource := resource.NewWithAttributes(
    semconv.SchemaURL,
    semconv.ServiceNameKey.String("my-service"),
    semconv.ServiceVersionKey.String("1.0.0"),
)

// Provider registration
tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource),
)
otel.SetTracerProvider(tp)
```

### Automated Instrumentation

#### HTTP Tracing

```go
// Wrap main router with otelhttp
mux := http.NewServeMux()
// ... register routes ...
otelHandler := otelhttp.NewHandler(mux, "http-server")
```

#### Database Tracing

```go
// Wrap PostgreSQL driver
import _ "github.com/jackc/pgx/v5/stdlib"
import "go.opentelemetry.io/contrib/instrumentation/database/sql/otelsql"

// In database connection setup
db, err := otelsql.Open("postgres", dsn, otelsql.WithAttributes(
    semconv.DBSystemPostgreSQL,
))
```

### Correlating Logs and Traces

The true value lies in correlating different telemetry signals.

#### Custom Logging Middleware

```go
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract trace context
            span := trace.SpanFromContext(r.Context())
            traceID := span.SpanContext().TraceID().String()
            spanID := span.SpanContext().SpanID().String()
            
            // Create request-specific logger
            requestLogger := logger.With("trace_id", traceID, "span_id", spanID)
            
            // Inject logger into request context
            ctx := context.WithValue(r.Context(), "logger", requestLogger)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

#### Usage in Handlers

```go
func (h *Handler) HandleGetUser() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger := r.Context().Value("logger").(*slog.Logger)
        logger.Info("handling get user request", "user_id", userID)
        
        // ... handler logic ...
    }
}
```

---

## The Developer Workflow: Taskfile and Air

A smooth, efficient, and reproducible developer workflow is critical for team productivity.

### Task Automation with Taskfile.yml

Taskfile is chosen as a modern, dependency-free alternative to GNU make.

#### Example Taskfile.yml

```yaml
version: '3'

# Load environment variables from a local .env file
env:
  dotenv: ['.env']

tasks:
  # The default task that runs when you just type 'task'
  default:
    cmds:
      - task: dev

  # Run the development server with live reload
  dev:
    desc: "Starts the server with live-reloading using Air"
    cmds:
      - air

  # Build a production binary
  build:
    desc: "Builds the application binary for production"
    cmds:
      - CGO_ENABLED=0 go build -o ./bin/server ./cmd/server

  # Run all tests
  test:
    desc: "Runs all Go tests"
    cmds:
      - go test -v -race ./...

  # Generate Go code from SQL queries using sqlc
  sqlc:generate:
    desc: "Generates Go code from SQL queries via sqlc"
    cmds:
      - sqlc generate

  # Database migration tasks
  db:migrate:
    desc: "Applies all pending database migrations"
    cmds:
      - goose -dir db/migrations postgres "{{.POSTGRES_DSN}}" up
  
  db:migrate:down:
    desc: "Rolls back the last applied database migration"
    cmds:
      - goose -dir db/migrations postgres "{{.POSTGRES_DSN}}" down

  db:migrate:create:
    desc: "Creates a new SQL migration file"
    cmds:
      - goose -dir db/migrations create {{.CLI_ARGS}} sql
```

### Live Reloading with air

Air provides a fast and efficient feedback loop during development.

#### Configuration (.air.toml)

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "./tmp/main"
include_ext = ["go", "sql", "yml", "toml"]
exclude_dir = ["tmp", "vendor", "bin", "docs"]
log = "air-build-errors.log"

[misc]
clean_on_exit = true

[screen]
clear_on_rebuild = true
```

### Developer Onboarding

The combination of Taskfile and air creates a fully declarative, reproducible development environment:

1. **Clone the repository**
2. **Create local `.env` file** for database credentials
3. **Run `task db:migrate`** to set up the database schema
4. **Run `task dev`** (or simply `air`)

This workflow ensures a fully functional, database-connected, live-reloading server is running in minutes.

---

## Conclusion

This document has detailed a canonical repository structure and architectural standard for building production-grade REST web services in Go. The presented structure is not an arbitrary collection of files and directories but a carefully considered system designed to optimize for specific outcomes critical to modern software development.

### Key Benefits

| Aspect | Benefit |
|--------|---------|
| **Maintainability** | Feature-oriented architecture with clear separation of concerns results in easier understanding, modification, and debugging |
| **Scalability** | Structure scales gracefully in terms of code complexity and team size, with strong encapsulation allowing concurrent development |
| **Testability** | Decoupled business logic from external dependencies enables isolated unit testing |
| **Observability** | Deep integration of `slog` and OpenTelemetry provides comprehensive view of application health and performance |
| **Developer Experience** | Streamlined workflow with Taskfile and air provides fast, consistent, and frictionless development environment |

### Architectural Decisions Summary

By adhering to idiomatic Go principles, leveraging the power of the modern standard library, and integrating best-in-class tooling for the data and observability layers, this standard provides a robust, defensible, and production-ready foundation for any new Go web service. It embodies the principles of clarity, simplicity, and performance that define modern Go development.

---

## Repository Structure Summary

| Directory/File | Purpose & Justification |
|----------------|------------------------|
| `/cmd/server/` | Executable entry point. Wires dependencies and starts the server. Keeps main minimal and focused on orchestration. |
| `/internal/` | All private application code. Compiler-enforced privacy prevents external imports, ensuring encapsulation. |
| `/internal/config/` | Centralized application configuration loading and parsing from files and environment variables. |
| `/internal/server/` | HTTP server setup, routing, and middleware. Decouples server logic from the main package. |
| `/internal/platform/` | Infrastructure-level code (e.g., database connection, telemetry SDK) shared across features. |
| `/internal/<feature>/` | Self-contained business feature modules (e.g., user). The core of the vertical slice architecture. |
| `/db/migrations/` | goose SQL migration files. Provides a version-controlled history of the database schema. |
| `/sql/` | sqlc source SQL for schema and queries. Separates SQL from Go for clarity and tooling support. |
| `Taskfile.yml` | Automation for build, test, run, and migration tasks. Codifies the developer workflow into simple commands. |
| `.air.toml` | Configuration for air live-reloading, enhancing the local development feedback loop. |