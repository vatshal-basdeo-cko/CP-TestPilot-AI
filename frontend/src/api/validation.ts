import apiClient from './client';
import type { ValidationResult } from '../types';

interface ValidateParams {
  status_code: number;
  body: unknown;
}

export const validationApi = {
  validate: async (response: ValidateParams, expectedStatus?: number, schema?: unknown): Promise<ValidationResult> => {
    const result = await apiClient.post<ValidationResult>('/api/v1/validate', {
      response: response.body,
      status_code: response.status_code,
      expected_status: expectedStatus,
      expected_schema: schema,
    });
    return result.data;
  },
};

