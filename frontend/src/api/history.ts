import apiClient from './client';
import type { HistoryResponse, TestExecution } from '../types';

export interface HistoryFilters {
  status?: 'success' | 'failed' | 'error';
  api_id?: string;
  from_date?: string;
  to_date?: string;
  limit?: number;
  offset?: number;
}

export const historyApi = {
  list: async (filters?: HistoryFilters): Promise<HistoryResponse> => {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, String(value));
        }
      });
    }
    const response = await apiClient.get<HistoryResponse>(`/api/v1/history?${params.toString()}`);
    return response.data;
  },

  get: async (id: string): Promise<TestExecution> => {
    const response = await apiClient.get<TestExecution>(`/api/v1/history/${id}`);
    return response.data;
  },
};

