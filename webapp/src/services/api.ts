// API client for backend communication
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface ApiResponse<T> {
  data: T;
  error?: string;
}

class ApiClient {
  private baseURL: string;
  private headers: Record<string, string>;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
    this.headers = {
      'Content-Type': 'application/json',
    };
  }

  private async request<T>(
    method: string,
    path: string,
    options?: {
      body?: unknown;
      params?: Record<string, string | number | boolean>;
      headers?: Record<string, string>;
    }
  ): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseURL}${path}`);

    if (options?.params) {
      Object.entries(options.params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    try {
      const response = await fetch(url.toString(), {
        method,
        headers: {
          ...this.headers,
          ...options?.headers,
        },
        body: options?.body ? JSON.stringify(options.body) : undefined,
      });

      const data = await response.json();

      if (!response.ok) {
        return {
          data: null as T,
          error: data.error || `Request failed with status ${response.status}`,
        };
      }

      return { data };
    } catch (error) {
      return {
        data: null as T,
        error:
          error instanceof Error ? error.message : 'An unknown error occurred',
      };
    }
  }

  get<T>(path: string, params?: Record<string, string | number | boolean>) {
    return this.request<T>('GET', path, { params });
  }

  post<T>(path: string, body?: unknown) {
    return this.request<T>('POST', path, { body });
  }

  put<T>(path: string, body?: unknown) {
    return this.request<T>('PUT', path, { body });
  }

  delete<T>(path: string) {
    return this.request<T>('DELETE', path);
  }
}

export const apiClient = new ApiClient();

// User types
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface UsersListResponse {
  users: User[];
  limit: number;
  offset: number;
}

// API functions
export const api = {
  health: () => apiClient.get<{ status: string }>('/health'),
  
  users: {
    list: (params?: { limit?: number; offset?: number }) => 
      apiClient.get<UsersListResponse>('/api/v1/users', params),
    
    getById: (id: string) => 
      apiClient.get<User>(`/api/v1/users/${id}`),
  },
};
