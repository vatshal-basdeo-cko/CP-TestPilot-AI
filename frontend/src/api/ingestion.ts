import apiClient from './client';
import type { APISpecification } from '../types';

export const ingestionApi = {
  listAPIs: async (): Promise<{ apis: APISpecification[]; count: number }> => {
    const response = await apiClient.get('/api/v1/apis');
    return response.data;
  },

  ingestFolder: async (folderPath: string): Promise<{
    message: string;
    ingested: number;
    skipped: number;
    failed: number;
    errors: string[] | null;
  }> => {
    const response = await apiClient.post('/api/v1/ingest/folder', {
      folder_path: folderPath,
    });
    return response.data;
  },

  getStatus: async (): Promise<unknown> => {
    const response = await apiClient.get('/api/v1/ingest/status');
    return response.data;
  },
};

