import apiClient from './client';
import type { LoginRequest, LoginResponse, User } from '../types';

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<{ token: string }>('/api/v1/auth/login', data);
    const token = response.data.token;
    
    // Get user info
    localStorage.setItem('token', token);
    const userResponse = await apiClient.get<User>('/api/v1/auth/me');
    
    return {
      token,
      user: userResponse.data,
    };
  },

  me: async (): Promise<User> => {
    const response = await apiClient.get<User>('/api/v1/auth/me');
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },
};

