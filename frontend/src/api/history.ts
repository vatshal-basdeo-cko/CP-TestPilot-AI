import apiClient from './client';
import type { HistoryResponse, TestExecution, ValidationResult } from '../types';

export interface HistoryFilters {
  status?: 'success' | 'failed' | 'error';
  api_id?: string;
  from_date?: string;
  to_date?: string;
  search?: string;
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

  delete: async (id: string): Promise<{ message: string; id: string }> => {
    const response = await apiClient.delete<{ message: string; id: string }>(`/api/v1/history/${id}`);
    return response.data;
  },

  updateValidation: async (id: string, validationResult: ValidationResult): Promise<{ message: string; id: string }> => {
    const response = await apiClient.patch<{ message: string; id: string }>(
      `/api/v1/history/${id}/validation`,
      { validation_result: validationResult }
    );
    return response.data;
  },
};

