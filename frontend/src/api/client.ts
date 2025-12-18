import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_GATEWAY_URL || '';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper to get token from either direct localStorage or Zustand persisted state
const getAuthToken = (): string | null => {
  // First try direct localStorage (set by authApi.login)
  const directToken = localStorage.getItem('token');
  if (directToken) return directToken;
  
  // Fallback: try to read from Zustand persisted storage
  try {
    const authStorage = localStorage.getItem('auth-storage');
    if (authStorage) {
      const parsed = JSON.parse(authStorage);
      const zustandToken = parsed?.state?.token;
      if (zustandToken) {
        // Sync to direct localStorage for future requests
        localStorage.setItem('token', zustandToken);
        return zustandToken;
      }
    }
  } catch {
    // Ignore JSON parse errors
  }
  return null;
};

// Request interceptor - add auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAuthToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle 401
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Clear all auth-related storage to prevent stale state after Zustand rehydration
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      localStorage.removeItem('auth-storage'); // Clear Zustand persisted state
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
