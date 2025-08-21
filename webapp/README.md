# Starterkit Web App

Modern React application with performance and maintainability focus.

## Tech Stack

- **React 19** + TypeScript + Vite
- **Tailwind CSS** - Utility-first styling
- **Zustand** - Lightweight state management
- **React Router** - Client-side routing with code splitting

## Quick Start

```bash
# Install dependencies
task frontend:install

# Start development server
task frontend:dev

# Visit http://localhost:5173
```

## Commands

### Development
- `task frontend:install` - Install dependencies
- `task frontend:dev` - Development server
- `task frontend:preview` - Preview production build
- `task frontend:build` - Production build

### Quality
- `task frontend:lint` - ESLint
- `task frontend:typecheck` - TypeScript checking
- `task frontend:format` - Format with Prettier
- `task frontend:test` - Run tests

### Utility
- `task frontend:clean` - Clean build artifacts

## Project Structure

```
src/
├── assets/         # Static assets
├── components/     # Reusable UI components
├── features/       # Feature modules (auth/, dashboard/, etc.)
├── hooks/          # Custom React hooks
├── lib/            # Utilities
├── pages/          # Route components (lazy-loaded)
├── services/       # API communication
├── store/          # Zustand state stores
├── styles/         # Global CSS
└── types/          # TypeScript definitions
```

## Deployment

```bash
# Build for production
task frontend:build

# Deploy webapp/dist directory to static hosting
```

## Development Guidelines

- Feature-based directory structure
- TypeScript strict mode
- Tailwind utilities over custom CSS
- Co-locate related functionality
