export const API_ERROR = '@birdseye/API_ERROR';

export function apiError(error) {
  return {
    type: API_ERROR,
    error,
  };
}

export function resetApiError() {
  return apiError(null);
}
