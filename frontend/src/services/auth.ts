import api from './api';
import type { ApiResponse, User } from '../types';

interface LoginResponse {
  token: string;
  user: User;
}

export const authService = {
  login: async (username: string, password: string) => {
    const res = await api.post<ApiResponse<LoginResponse>>('/auth/login', { username, password });
    return res.data.data;
  },
  me: async () => {
    const res = await api.get<ApiResponse<User>>('/auth/me');
    return res.data.data;
  },
  register: async (data: { username: string; password: string; name: string; role: string }) => {
    const res = await api.post<ApiResponse<User>>('/auth/register', data);
    return res.data.data;
  },
};
