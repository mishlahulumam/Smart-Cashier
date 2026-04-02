import { useEffect, useState } from 'react';
import { reportService } from '../services/reports';
import { formatCurrency, formatDate } from '../utils/format';
import type { DashboardData } from '../types';
import {
  DollarSign,
  ShoppingCart,
  Package,
  AlertTriangle,
  TrendingUp,
} from 'lucide-react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';

export default function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const result = await reportService.getDashboard();
      setData(result);
    } catch {
      // silently handle
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  const stats = [
    {
      label: 'Pendapatan Hari Ini',
      value: formatCurrency(data?.today_revenue || 0),
      icon: DollarSign,
      color: 'bg-green-100 text-green-600',
    },
    {
      label: 'Transaksi Hari Ini',
      value: data?.today_transactions || 0,
      icon: ShoppingCart,
      color: 'bg-blue-100 text-blue-600',
    },
    {
      label: 'Total Produk',
      value: data?.total_products || 0,
      icon: Package,
      color: 'bg-purple-100 text-purple-600',
    },
    {
      label: 'Stok Menipis',
      value: data?.low_stock_count || 0,
      icon: AlertTriangle,
      color: 'bg-red-100 text-red-600',
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 mt-1">Ringkasan bisnis Anda hari ini</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="bg-white rounded-xl border border-gray-200 p-5">
              <div className="flex items-center gap-4">
                <div className={`w-12 h-12 rounded-lg flex items-center justify-center ${stat.color}`}>
                  <Icon className="w-6 h-6" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">{stat.label}</p>
                  <p className="text-xl font-bold text-gray-900">{stat.value}</p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="flex items-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-blue-600" />
            <h2 className="text-lg font-semibold text-gray-900">Produk Terlaris</h2>
          </div>
          {data?.top_products && data.top_products.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={data.top_products}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="product_name" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip
                  formatter={(value) => [String(value), 'Terjual']}
                  labelStyle={{ fontWeight: 600 }}
                />
                <Bar dataKey="quantity" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-400 text-center py-12">Belum ada data</p>
          )}
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Transaksi Terakhir</h2>
          {data?.recent_transactions && data.recent_transactions.length > 0 ? (
            <div className="space-y-3">
              {data.recent_transactions.map((tx) => (
                <div
                  key={tx.id}
                  className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0"
                >
                  <div>
                    <p className="text-sm font-medium text-gray-900">{tx.invoice_number}</p>
                    <p className="text-xs text-gray-500">{formatDate(tx.created_at)}</p>
                  </div>
                  <p className="text-sm font-semibold text-gray-900">
                    {formatCurrency(tx.total)}
                  </p>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-12">Belum ada transaksi</p>
          )}
        </div>
      </div>
    </div>
  );
}
