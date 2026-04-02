import api from './api';
import type { ApiResponse, PaginatedResponse, Transaction } from '../types';

interface CreateTransactionData {
  customer_id?: number | null;
  items: { product_id: number; quantity: number; discount: number }[];
  discount: number;
  payment_amount: number;
}

export const transactionService = {
  getAll: async (params?: { start_date?: string; end_date?: string; page?: number; limit?: number }) => {
    const res = await api.get<PaginatedResponse<Transaction[]>>('/transactions', { params });
    return res.data;
  },
  getById: async (id: number) => {
    const res = await api.get<ApiResponse<Transaction>>(`/transactions/${id}`);
    return res.data.data;
  },
  create: async (data: CreateTransactionData) => {
    const res = await api.post<ApiResponse<Transaction>>('/transactions', data);
    return res.data.data;
  },
};
