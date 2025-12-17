import apiClient from './client';
import type { ExecuteRequest, ExecuteResponse, Environment } from '../types';

export const executionApi = {
  execute: async (request: ExecuteRequest): Promise<ExecuteResponse> => {
    const response = await apiClient.post<ExecuteResponse>('/api/v1/execute', request);
    return response.data;
  },

  getEnvironments: async (): Promise<{ environments: Environment[]; count: number }> => {
    const response = await apiClient.get('/api/v1/environments');
    return response.data;
  },

  createEnvironment: async (data: Omit<Environment, 'id' | 'created_at' | 'updated_at'>): Promise<Environment> => {
    const response = await apiClient.post<Environment>('/api/v1/environments', data);
    return response.data;
  },
};

