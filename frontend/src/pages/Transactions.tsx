import { useEffect, useState } from 'react';
import { transactionService } from '../services/transactions';
import { formatCurrency, formatDate } from '../utils/format';
import type { Transaction } from '../types';
import { Search, Eye, X, ChevronLeft, ChevronRight } from 'lucide-react';

export default function Transactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [detail, setDetail] = useState<Transaction | null>(null);
  const limit = 20;

  useEffect(() => {
    loadTransactions();
  }, [page, startDate, endDate]);

  const loadTransactions = async () => {
    try {
      const res = await transactionService.getAll({
        page,
        limit,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
      });
      setTransactions(res.data || []);
      setTotal(res.total);
    } catch {
      setTransactions([]);
    }
  };

  const viewDetail = async (id: number) => {
    try {
      const tx = await transactionService.getById(id);
      setDetail(tx);
    } catch {
      // handle error
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Transaksi</h1>
        <p className="text-gray-500 text-sm">Histori semua transaksi</p>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <div className="flex items-center gap-2">
          <Search className="w-4 h-4 text-gray-400" />
          <input
            type="date"
            value={startDate}
            onChange={(e) => { setStartDate(e.target.value); setPage(1); }}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
          />
          <span className="text-gray-400">-</span>
          <input
            type="date"
            value={endDate}
            onChange={(e) => { setEndDate(e.target.value); setPage(1); }}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
          />
        </div>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="text-left px-4 py-3 font-medium text-gray-600">Invoice</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Tanggal</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Kasir</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Pelanggan</th>
                <th className="text-right px-4 py-3 font-medium text-gray-600">Total</th>
                <th className="text-right px-4 py-3 font-medium text-gray-600">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {transactions.map((tx) => (
                <tr key={tx.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono text-xs text-gray-900">{tx.invoice_number}</td>
                  <td className="px-4 py-3 text-gray-500">{formatDate(tx.created_at)}</td>
                  <td className="px-4 py-3 text-gray-700">{tx.user?.name}</td>
                  <td className="px-4 py-3 text-gray-500">{tx.customer?.name || '-'}</td>
                  <td className="px-4 py-3 text-right font-medium text-gray-900">{formatCurrency(tx.total)}</td>
                  <td className="px-4 py-3 text-right">
                    <button onClick={() => viewDetail(tx.id)} className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg">
                      <Eye className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
              {transactions.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-4 py-12 text-center text-gray-400">
                    Tidak ada transaksi
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200">
            <p className="text-sm text-gray-500">
              Menampilkan {(page - 1) * limit + 1}-{Math.min(page * limit, total)} dari {total}
            </p>
            <div className="flex items-center gap-1">
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page === 1}
                className="p-1.5 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <span className="px-3 py-1 text-sm font-medium">{page} / {totalPages}</span>
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page === totalPages}
                className="p-1.5 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Detail Modal */}
      {detail && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl w-full max-w-md p-6 max-h-[90vh] overflow-auto">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold text-gray-900">Detail Transaksi</h2>
              <button onClick={() => setDetail(null)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-3 text-sm">
              <div className="grid grid-cols-2 gap-2">
                <div><span className="text-gray-500">Invoice:</span></div>
                <div className="font-mono text-xs">{detail.invoice_number}</div>
                <div><span className="text-gray-500">Tanggal:</span></div>
                <div>{formatDate(detail.created_at)}</div>
                <div><span className="text-gray-500">Kasir:</span></div>
                <div>{detail.user?.name}</div>
                <div><span className="text-gray-500">Pelanggan:</span></div>
                <div>{detail.customer?.name || '-'}</div>
              </div>

              <div className="border-t pt-3">
                <p className="font-medium text-gray-900 mb-2">Item</p>
                {detail.items?.map((item, idx) => (
                  <div key={idx} className="flex justify-between py-1.5 border-b border-gray-100 last:border-0">
                    <div>
                      <p className="text-gray-900">{item.product_name}</p>
                      <p className="text-xs text-gray-500">{item.quantity} x {formatCurrency(item.price)}</p>
                    </div>
                    <p className="font-medium">{formatCurrency(item.subtotal)}</p>
                  </div>
                ))}
              </div>

              <div className="border-t pt-3 space-y-1">
                <div className="flex justify-between"><span className="text-gray-500">Subtotal</span><span>{formatCurrency(detail.subtotal)}</span></div>
                {detail.discount > 0 && (
                  <div className="flex justify-between text-red-600"><span>Diskon</span><span>-{formatCurrency(detail.discount)}</span></div>
                )}
                <div className="flex justify-between font-bold text-base"><span>Total</span><span>{formatCurrency(detail.total)}</span></div>
                <div className="flex justify-between"><span className="text-gray-500">Bayar</span><span>{formatCurrency(detail.payment_amount)}</span></div>
                <div className="flex justify-between font-bold text-green-600"><span>Kembalian</span><span>{formatCurrency(detail.change_amount)}</span></div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
