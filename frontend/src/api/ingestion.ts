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

  uploadPostman: async (file: File): Promise<{
    message: string;
    api_id: string;
    name: string;
    endpoints: number;
  }> => {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.post('/api/v1/ingest/postman', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

