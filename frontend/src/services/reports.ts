import api from './api';
import type { ApiResponse, DashboardData } from '../types';

export const reportService = {
  getDashboard: async () => {
    const res = await api.get<ApiResponse<DashboardData>>('/dashboard');
    return res.data.data;
  },
  getDailySummary: async (date?: string) => {
    const res = await api.get<ApiResponse<Record<string, unknown>>>('/reports/daily', { params: { date } });
    return res.data.data;
  },
  getMonthlySummary: async (year?: number, month?: number) => {
    const res = await api.get<ApiResponse<Record<string, unknown>[]>>('/reports/monthly', { params: { year, month } });
    return res.data.data;
  },
  getTopProducts: async (params?: { start_date?: string; end_date?: string; limit?: number }) => {
    const res = await api.get<ApiResponse<Record<string, unknown>[]>>('/reports/top-products', { params });
    return res.data.data;
  },
  getSummary: async (params?: { start_date?: string; end_date?: string }) => {
    const res = await api.get<ApiResponse<Record<string, unknown>>>('/reports/summary', { params });
    return res.data.data;
  },
};
