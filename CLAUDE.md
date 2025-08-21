# Starterkit Architecture Guide

Full-stack template: Go REST API + React SPA

## Core Principles

### Architecture Philosophy
- **Idiomatic Go & React**: Follow language/framework conventions
- **Start Simple, Scale Gracefully**: Minimal initial complexity, clear growth path
- **Feature-Oriented (Vertical Slice)**: Organize by business domain, not technical layer
- **Performance by Default**: Optimized from the start, not as afterthought

## Tech Stack

### Backend (Go)
- **Core**: Go 1.22+ with net/http, no heavy frameworks
- **Database**: PostgreSQL + pgx/v5 (connection pooling)
- **Migrations**: goose (version-controlled SQL)
- **SQL**: sqlc (compile-time safe, NO ORMs)
- **Observability**: OpenTelemetry + slog
- **DI Pattern**: Constructor-based injection

### Frontend (React)
- **Core**: React 19 + TypeScript + Vite
- **Styling**: Tailwind CSS only (JIT mode, utility-first)
- **State**: Zustand (selector-based, no prop drilling)
- **Routing**: React Router + lazy loading
- **Components**: shadcn/ui (copy-paste, not npm)
- **Structure**: Feature modules with co-location

## Project Structure

```
/api/
  /cmd/server/          # Entry point, dependency wiring
  /internal/
    /config/            # App configuration
    /server/            # HTTP routes, middleware
    /platform/          # Shared infra (db, telemetry)
    /<feature>/         # Business features
      handler.go        # HTTP handlers
      service.go        # Business logic (pure)
      repository.go     # Data access
      models.go         # Domain types
  /db/migrations/       # Timestamped SQL files
  /sql/queries/         # sqlc source files

/webapp/
  /src/
    /components/        # Global UI components
    /features/          # Feature modules
      /<feature>/
        components/     # Feature-specific UI
        hooks/          # Feature logic
        services/       # API calls
        stores/         # Feature state
        types.ts        # TypeScript types
        index.ts        # Public exports
    /pages/             # Route components (lazy-loaded)
    /services/          # API client setup
    /store/             # Global app state
```

## Key Patterns

### Backend Patterns
```go
// Dependency injection chain (in main.go)
dbPool := database.Connect(cfg)
userRepo := user.NewRepository(dbPool)
userService := user.NewService(userRepo)
userHandler := user.NewHandler(userService)

// HTTP routing (Go 1.22+)
mux.HandleFunc("GET /users/{id}", handler.GetUser)
```

### Frontend Patterns
```typescript
// Feature-based structure
/features/auth/
  components/LoginForm.tsx
  hooks/useAuth.ts
  stores/authStore.ts
  types.ts
  index.ts  // Public API

// Zustand store (no context providers)
const useAppStore = create<State>((set) => ({
  count: 0,
  increment: () => set(state => ({ count: state.count + 1 }))
}))

// Route-based code splitting
const HomePage = lazy(() => import('./pages/Home'))
```

## Development Workflow

### Commands (via Taskfile)
```bash
task dev                           # Start with hot reload
task build                         # Production build
task db:migrate                    # Run migrations
task db:migrate:create NAME        # New migration
task task backend:generate:sqlc    # Generate Go from SQL
task test                          # Run all tests
```

### Quality Gates
- Backend: golangci-lint, go vet, tests
- Frontend: ESLint, TypeScript strict, Prettier
- Both: Pre-commit hooks, CI/CD checks

## Best Practices

### DO
- Write raw SQL (sqlc generates types)
- Use feature folders for business logic
- Implement lazy loading for routes
- Use Tailwind utilities exclusively
- Keep handlers thin, services pure
- Version control migrations

### DON'T
- Use ORMs or query builders
- Create unnecessary abstractions
- Mix concerns across features
- Write custom CSS
- Use global state for local data
- Store secrets in code

## Quick Start

1. Clone and setup `.env`
2. `task db:migrate` - Setup database
3. `task dev` - Start development
4. Implement features in `/api/internal/<feature>` and `/webapp/src/features/<feature>`