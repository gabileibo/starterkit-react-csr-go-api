# Feature-Based Architecture

This directory contains feature modules that encapsulate all related functionality. Each feature is self-contained with its own components, hooks, services, and types.

## Structure Pattern

Each feature should follow this structure:

```
/features/
├── /auth/
│   ├── components/
│   │   ├── LoginForm.tsx
│   │   ├── RegisterForm.tsx
│   │   └── ProtectedRoute.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   └── usePermissions.ts
│   ├── services/
│   │   └── authService.ts
│   ├── stores/
│   │   └── authStore.ts
│   ├── types.ts
│   └── index.ts
├── /dashboard/
│   ├── components/
│   │   ├── DashboardLayout.tsx
│   │   ├── StatsCard.tsx
│   │   └── ActivityFeed.tsx
│   ├── hooks/
│   │   └── useDashboardData.ts
│   ├── services/
│   │   └── dashboardService.ts
│   └── types.ts
└── /user-profile/
    ├── components/
    │   ├── ProfileForm.tsx
    │   ├── AvatarUpload.tsx
    │   └── PreferencesPanel.tsx
    ├── hooks/
    │   └── useProfile.ts
    ├── services/
    │   └── profileService.ts
    └── types.ts
```

## Guidelines

### 1. Modularity
Each feature is completely self-contained. All code related to a specific functionality lives together.

### 2. Co-location
- **components/**: UI components specific to this feature
- **hooks/**: Custom React hooks for this feature's logic
- **services/**: API calls and external communication
- **stores/**: Feature-specific Zustand stores (if needed)
- **types.ts**: TypeScript interfaces and types
- **index.ts**: Public API exports

### 3. Example Feature Module

```typescript
// /features/auth/types.ts
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'user' | 'admin';
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
```

```typescript
// /features/auth/stores/authStore.ts
import { create } from 'zustand';
import type { AuthState, User } from '../types';

interface AuthStore extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  updateUser: (user: User) => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  
  login: async (email, password) => {
    set({ isLoading: true });
    // API call logic
    set({ isLoading: false });
  },
  
  logout: () => {
    set({ user: null, isAuthenticated: false });
  },
  
  updateUser: (user) => {
    set({ user });
  },
}));
```

```typescript
// /features/auth/hooks/useAuth.ts
import { useAuthStore } from '../stores/authStore';

export function useAuth() {
  const { user, isAuthenticated, login, logout } = useAuthStore();
  
  return {
    user,
    isAuthenticated,
    login,
    logout,
  };
}
```

```typescript
// /features/auth/index.ts
// Public API - only export what other features need
export { LoginForm } from './components/LoginForm';
export { ProtectedRoute } from './components/ProtectedRoute';
export { useAuth } from './hooks/useAuth';
export type { User, AuthState } from './types';
```

## When to Create a New Feature

Create a new feature directory when:
- The functionality represents a distinct business domain
- Multiple components and hooks work together for a specific purpose
- The code would otherwise be scattered across type-based folders
- You want to enable independent development and testing

## Shared Code

For truly global/shared code, use the root-level directories:
- `/components/` - Reusable UI components (Button, Modal, etc.)
- `/hooks/` - Global hooks (useLocalStorage, useDebounce, etc.)
- `/services/` - Core API client setup
- `/store/` - Global application state
- `/types/` - Shared TypeScript definitions