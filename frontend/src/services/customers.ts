import api from './api';
import type { ApiResponse, Customer } from '../types';

export const customerService = {
  getAll: async (search?: string) => {
    const res = await api.get<ApiResponse<Customer[]>>('/customers', { params: { search } });
    return res.data.data;
  },
  getById: async (id: number) => {
    const res = await api.get<ApiResponse<Customer>>(`/customers/${id}`);
    return res.data.data;
  },
  create: async (data: { name: string; phone: string; email: string }) => {
    const res = await api.post<ApiResponse<Customer>>('/customers', data);
    return res.data.data;
  },
  update: async (id: number, data: { name: string; phone: string; email: string }) => {
    const res = await api.put<ApiResponse<Customer>>(`/customers/${id}`, data);
    return res.data.data;
  },
  delete: async (id: number) => {
    await api.delete(`/customers/${id}`);
  },
};
