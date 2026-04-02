import api from './api';
import type { ApiResponse, Category } from '../types';

export const categoryService = {
  getAll: async () => {
    const res = await api.get<ApiResponse<Category[]>>('/categories');
    return res.data.data;
  },
  create: async (data: { name: string }) => {
    const res = await api.post<ApiResponse<Category>>('/categories', data);
    return res.data.data;
  },
  update: async (id: number, data: { name: string }) => {
    const res = await api.put<ApiResponse<Category>>(`/categories/${id}`, data);
    return res.data.data;
  },
  delete: async (id: number) => {
    await api.delete(`/categories/${id}`);
  },
};
