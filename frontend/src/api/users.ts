import apiClient from './client';

export interface User {
  id: string;
  username: string;
  role: 'admin' | 'user';
  created_at: string;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  role: 'admin' | 'user';
}

export interface UsersResponse {
  users: User[];
  count: number;
}

export interface CreateUserResponse {
  message: string;
  user: User;
}

export const usersApi = {
  // List all users (admin only)
  list: async (): Promise<User[]> => {
    const response = await apiClient.get<UsersResponse>('/api/v1/users');
    return response.data.users || [];
  },

  // Create a new user (admin only)
  create: async (data: CreateUserRequest): Promise<CreateUserResponse> => {
    const response = await apiClient.post<CreateUserResponse>('/api/v1/users', data);
    return response.data;
  },

  // Delete a user (admin only)
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/api/v1/users/${id}`);
  },
};

