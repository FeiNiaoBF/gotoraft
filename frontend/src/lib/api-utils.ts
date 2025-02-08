import type { ApiResponse, ApiError } from './raft-types';

export function createApiResponse<T>(
  data?: T,
  error?: ApiError
): ApiResponse<T> {
  return { data, error };
}
