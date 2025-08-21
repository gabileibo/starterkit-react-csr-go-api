export interface ApiResponse<T> {
  data: T;
  status: number;
}

export interface ApiError {
  status: number;
  message: string;
  details?: unknown;
}

export interface ErrorResponse {
  error: string;
}
