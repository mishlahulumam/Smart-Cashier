import api from './api';
import type { ApiResponse, Product } from '../types';

export const productService = {
  getAll: async (params?: { category_id?: number; search?: string; low_stock?: boolean }) => {
    const res = await api.get<ApiResponse<Product[]>>('/products', { params });
    return res.data.data;
  },
  getById: async (id: number) => {
    const res = await api.get<ApiResponse<Product>>(`/products/${id}`);
    return res.data.data;
  },
  create: async (data: Partial<Product>) => {
    const res = await api.post<ApiResponse<Product>>('/products', data);
    return res.data.data;
  },
  update: async (id: number, data: Partial<Product>) => {
    const res = await api.put<ApiResponse<Product>>(`/products/${id}`, data);
    return res.data.data;
  },
  delete: async (id: number) => {
    await api.delete(`/products/${id}`);
  },
};
