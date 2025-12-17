import apiClient from './client';
import type { ValidationResult } from '../types';

export const validationApi = {
  validate: async (response: unknown, expectedStatus?: number, schema?: unknown): Promise<ValidationResult> => {
    const result = await apiClient.post<ValidationResult>('/api/v1/validate', {
      response,
      expected_status: expectedStatus,
      schema,
    });
    return result.data;
  },
};

