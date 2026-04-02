import { useEffect, useState } from 'react';
import { reportService } from '../services/reports';
import { formatCurrency } from '../utils/format';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
} from 'recharts';
import { Calendar, DollarSign, ShoppingCart, Package, TrendingUp } from 'lucide-react';

export default function Reports() {
  const [tab, setTab] = useState<'daily' | 'monthly'>('daily');
  const [monthlyData, setMonthlyData] = useState<Record<string, unknown>[]>([]);
  const [topProducts, setTopProducts] = useState<Record<string, unknown>[]>([]);
  const [summary, setSummary] = useState<Record<string, unknown> | null>(null);
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [year, setYear] = useState(new Date().getFullYear());
  const [month, setMonth] = useState(new Date().getMonth() + 1);

  useEffect(() => {
    loadReport();
  }, [tab, date, year, month]);

  const loadReport = async () => {
    try {
      if (tab === 'daily') {
        const [, top, sum] = await Promise.all([
          reportService.getDailySummary(date),
          reportService.getTopProducts({ start_date: date, end_date: date }),
          reportService.getSummary({ start_date: date, end_date: date }),
        ]);
        setTopProducts(top || []);
        setSummary(sum);
      } else {
        const startDate = `${year}-${String(month).padStart(2, '0')}-01`;
        const endDay = new Date(year, month, 0).getDate();
        const endDate = `${year}-${String(month).padStart(2, '0')}-${String(endDay).padStart(2, '0')}`;
        const [monthly, top, sum] = await Promise.all([
          reportService.getMonthlySummary(year, month),
          reportService.getTopProducts({ start_date: startDate, end_date: endDate }),
          reportService.getSummary({ start_date: startDate, end_date: endDate }),
        ]);
        setMonthlyData(monthly || []);
        setTopProducts(top || []);
        setSummary(sum);
      }
    } catch {
      // handle error
    }
  };

  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];

  const summaryStats = [
    {
      label: 'Total Pendapatan',
      value: formatCurrency(Number(summary?.total_revenue || 0)),
      icon: DollarSign,
      color: 'bg-green-100 text-green-600',
    },
    {
      label: 'Total Transaksi',
      value: Number(summary?.total_transactions || 0),
      icon: ShoppingCart,
      color: 'bg-blue-100 text-blue-600',
    },
    {
      label: 'Item Terjual',
      value: Number(summary?.total_items || 0),
      icon: Package,
      color: 'bg-purple-100 text-purple-600',
    },
    {
      label: 'Rata-rata Transaksi',
      value: formatCurrency(Number(summary?.avg_transaction || 0)),
      icon: TrendingUp,
      color: 'bg-orange-100 text-orange-600',
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Laporan</h1>
        <p className="text-gray-500 text-sm">Analisis penjualan toko Anda</p>
      </div>

      <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center">
        <div className="flex bg-gray-100 rounded-lg p-1">
          <button
            onClick={() => setTab('daily')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition ${
              tab === 'daily' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Harian
          </button>
          <button
            onClick={() => setTab('monthly')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition ${
              tab === 'monthly' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Bulanan
          </button>
        </div>

        {tab === 'daily' ? (
          <div className="flex items-center gap-2">
            <Calendar className="w-4 h-4 text-gray-400" />
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
            />
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <Calendar className="w-4 h-4 text-gray-400" />
            <select
              value={month}
              onChange={(e) => setMonth(Number(e.target.value))}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
            >
              {months.map((m, i) => (
                <option key={i} value={i + 1}>{m}</option>
              ))}
            </select>
            <input
              type="number"
              value={year}
              onChange={(e) => setYear(Number(e.target.value))}
              className="w-24 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
            />
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {summaryStats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="bg-white rounded-xl border border-gray-200 p-4">
              <div className="flex items-center gap-3">
                <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${stat.color}`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-gray-500">{stat.label}</p>
                  <p className="text-lg font-bold text-gray-900">{stat.value}</p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {tab === 'monthly' && monthlyData.length > 0 && (
          <div className="bg-white rounded-xl border border-gray-200 p-5 lg:col-span-2">
            <h3 className="font-semibold text-gray-900 mb-4">Tren Pendapatan</h3>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={monthlyData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} tickFormatter={(v) => `${(v / 1000).toFixed(0)}k`} />
                <Tooltip formatter={(value) => [formatCurrency(Number(value)), 'Pendapatan']} />
                <Line type="monotone" dataKey="revenue" stroke="#3b82f6" strokeWidth={2} dot={{ r: 3 }} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        )}

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="font-semibold text-gray-900 mb-4">Produk Terlaris</h3>
          {topProducts.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={topProducts} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis type="number" tick={{ fontSize: 11 }} />
                <YAxis dataKey="product_name" type="category" tick={{ fontSize: 11 }} width={100} />
                <Tooltip formatter={(value) => [String(value), 'Terjual']} />
                <Bar dataKey="quantity" fill="#3b82f6" radius={[0, 4, 4, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-400 text-center py-12">Tidak ada data</p>
          )}
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="font-semibold text-gray-900 mb-4">Pendapatan per Produk</h3>
          {topProducts.length > 0 ? (
            <div className="space-y-3">
              {topProducts.map((p, i) => (
                <div key={i} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                  <div className="flex items-center gap-3">
                    <span className="w-6 h-6 rounded-full bg-blue-100 flex items-center justify-center text-xs font-bold text-blue-700">
                      {i + 1}
                    </span>
                    <div>
                      <p className="text-sm font-medium text-gray-900">{String(p.product_name)}</p>
                      <p className="text-xs text-gray-500">{Number(p.quantity)} terjual</p>
                    </div>
                  </div>
                  <p className="text-sm font-semibold text-gray-900">{formatCurrency(Number(p.revenue))}</p>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-12">Tidak ada data</p>
          )}
        </div>
      </div>
    </div>
  );
}
