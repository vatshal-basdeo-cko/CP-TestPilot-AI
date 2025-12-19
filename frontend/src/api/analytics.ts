import apiClient from './client';
import type { AnalyticsOverview, APIStats } from '../types';

export interface AnalyticsFilters {
  start_date?: string;
  end_date?: string;
}

export const analyticsApi = {
  getOverview: async (filters?: AnalyticsFilters): Promise<AnalyticsOverview> => {
    const params = new URLSearchParams();
    if (filters) {
      if (filters.start_date) {
        params.append('start_date', filters.start_date);
      }
      if (filters.end_date) {
        params.append('end_date', filters.end_date);
      }
    }
    const queryString = params.toString();
    const url = queryString ? `/api/v1/analytics/overview?${queryString}` : '/api/v1/analytics/overview';
    const response = await apiClient.get<AnalyticsOverview>(url);
    return response.data;
  },

  getByAPI: async (apiId: string): Promise<APIStats> => {
    const response = await apiClient.get<APIStats>(`/api/v1/analytics/by-api/${apiId}`);
    return response.data;
  },
};


