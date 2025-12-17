import apiClient from './client';
import type { ParseResult, ConstructedRequest } from '../types';

export const llmApi = {
  parse: async (naturalLanguage: string, provider?: string): Promise<ParseResult> => {
    const response = await apiClient.post<ParseResult>('/api/v1/parse', {
      natural_language: naturalLanguage,
      provider,
    });
    return response.data;
  },

  construct: async (parseResult: ParseResult, apiConfig?: unknown): Promise<ConstructedRequest> => {
    const response = await apiClient.post<ConstructedRequest>('/api/v1/construct', {
      parse_result: parseResult,
      api_config: apiConfig,
    });
    return response.data;
  },

  generateData: async (fieldName: string, fieldType: string, format?: string): Promise<unknown> => {
    const response = await apiClient.post('/api/v1/llm/generate-data', {
      field_name: fieldName,
      field_type: fieldType,
      format,
    });
    return response.data;
  },

  getProviders: async (): Promise<{ providers: string[]; default_provider: string }> => {
    const response = await apiClient.get('/api/v1/llm/providers');
    return response.data;
  },
};

