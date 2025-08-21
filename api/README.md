# Starterkit API

Go REST API with feature-oriented architecture.

## Tech Stack

- **Go 1.24** + PostgreSQL + pgx/v5
- **goose** (migrations) + **sqlc** (type-safe queries)
- **slog** (logging) + **OpenTelemetry** (tracing)

## Quick Start

```bash
task backend:install-tools  # Install dependencies
task db:start               # Start PostgreSQL
task backend:migrate        # Run migrations
task backend:dev            # Start with live-reload
```

API available at `http://localhost:8080`

## Commands

### Development
- `task backend:dev` - Live-reload server
- `task backend:build` - Production build
- `task dev:env` - Swagger UI at `:8082`

### Database
- `task backend:migrate` - Apply migrations
- `task backend:migrate:down` - Rollback
- `task backend:migrate:create -- <name>` - New migration

### Code Generation
- `task backend:generate` - All generation
- `task backend:generate:sqlc` - Database code from SQL

### Quality
- `task backend:test` - Run tests
- `task backend:lint` - Lint code
- `task backend:format` - Format code

## Structure

```
/api/
├── cmd/server/         # Entry point
├── internal/
│   ├── config/         # Config
│   ├── server/         # HTTP, routes, middleware
│   ├── platform/       # Database, telemetry
│   ├── user/           # User feature
│   └── <feature>/      # features
├── db/migrations/      # SQL migrations
└── sql/queries/        # SQL queries by feature
```

## Adding Features

1. Create package: `mkdir internal/product`
2. Add files: `models.go`, `repository.go`, `service.go`, `handler.go`
3. Write SQL in `/sql/queries/product.sql`
4. Generate: `task backend:generate:sqlc`
5. Register routes in `/internal/server/routes.go`

## SQL Queries

Write in `/sql/queries/<feature>.sql`:

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name)
VALUES ($1, $2) RETURNING *;
```

Generate code: `task backend:generate:sqlc`

## Configuration

Environment variables:

```env
SERVER_ADDRESS=:8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=starterkit
```
