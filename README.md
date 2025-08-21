# Starterkit

An opinionated full-stack web application template built with Go and React.

## Prerequisites

- Go 1.24+
- Node.js 22 LTS+
- Docker & Docker Compose
- [Task](https://taskfile.dev)

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd starterkit
cp .env.example .env

# Install dependencies
task backend:install-tools
task frontend:install

# Start development
task dev
```

This starts:
- PostgreSQL (Docker)
- Backend (http://localhost:8080)
- Frontend (http://localhost:5173)

## Architecture

```
api/
├── cmd/server/        # Entry point
├── internal/          # Business logic
│   ├── config/        # Configuration
│   ├── server/        # HTTP server
│   ├── platform/      # Infrastructure
│   └── <feature>/     # Feature modules
├── db/migrations/     # Database migrations
└── sql/               # SQL queries (sqlc)

webapp/                # React frontend
├── src/
│   ├── components/    # Reusable UI
│   ├── features/      # Feature modules
│   ├── pages/         # Route components
│   └── store/         # State (Zustand)
└── dist/              # Production build
```

## Development

### Common Commands

```bash
task dev               # Start everything
task test              # Run all tests
task lint              # Run linters
task format            # Format code
task check             # Run all checks
task generate          # Code generation
```

### Database

```bash
task db:start          # Start PostgreSQL
task db:shell          # Connect to database
task db:clean          # Reset database

task backend:migrate   # Apply migrations
task backend:migrate:create -- <name>  # Create migration
```

### API Documentation

```bash
task dev:env           # Start with Swagger UI
# Visit http://localhost:8082
```

## Tech Stack

**Backend:** Go, PostgreSQL (pgx/sqlc), goose, OpenTelemetry, slog

**Frontend:** React 19, TypeScript, Vite, Tailwind CSS v4, Zustand, React Router

## Customization

This is a starter template. To adapt it for your project:

1. Update module name in `api/go.mod`
2. Update package name in `webapp/package.json`
3. Add your database migrations to `api/db/migrations/`
4. Define your SQL queries in `api/sql/queries/`
5. Implement your business logic in `api/internal/`
6. Build your UI in `webapp/src/`

## License

MIT