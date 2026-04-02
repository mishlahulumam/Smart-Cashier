import { useEffect, useState, useRef } from 'react';
import { productService } from '../services/products';
import { categoryService } from '../services/categories';
import { transactionService } from '../services/transactions';
import { customerService } from '../services/customers';
import { useCartStore } from '../store/cartStore';
import { formatCurrency } from '../utils/format';
import type { Product, Category, Customer, Transaction } from '../types';
import {
  Search,
  Plus,
  Minus,
  Trash2,
  ShoppingCart,
  CreditCard,
  X,
  Printer,
  CheckCircle,
} from 'lucide-react';
import toast from 'react-hot-toast';

export default function POS() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [search, setSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<number>(0);
  const [showPayment, setShowPayment] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');
  const [completedTx, setCompletedTx] = useState<Transaction | null>(null);
  const [loading, setLoading] = useState(false);
  const receiptRef = useRef<HTMLDivElement>(null);

  const {
    items, discount, customerId,
    addItem, removeItem, updateQuantity,
    setDiscount, setCustomerId, clearCart,
    getSubtotal, getTotal,
  } = useCartStore();

  useEffect(() => {
    loadProducts();
    categoryService.getAll().then(setCategories);
    customerService.getAll().then(setCustomers);
  }, []);

  useEffect(() => {
    loadProducts();
  }, [search, selectedCategory]);

  const loadProducts = async () => {
    try {
      const data = await productService.getAll({
        search: search || undefined,
        category_id: selectedCategory || undefined,
      });
      setProducts(data || []);
    } catch {
      setProducts([]);
    }
  };

  const handlePayment = async () => {
    const amount = parseFloat(paymentAmount);
    const total = getTotal();
    if (!amount || amount < total) {
      toast.error('Jumlah pembayaran kurang');
      return;
    }

    setLoading(true);
    try {
      const tx = await transactionService.create({
        customer_id: customerId,
        items: items.map((i) => ({
          product_id: i.product.id,
          quantity: i.quantity,
          discount: i.discount,
        })),
        discount,
        payment_amount: amount,
      });
      setCompletedTx(tx);
      clearCart();
      setShowPayment(false);
      setPaymentAmount('');
      toast.success('Transaksi berhasil!');
      loadProducts();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } };
      toast.error(error?.response?.data?.message || 'Gagal memproses transaksi');
    } finally {
      setLoading(false);
    }
  };

  const quickAmounts = [10000, 20000, 50000, 100000, 200000, 500000];

  return (
    <div className="flex flex-col lg:flex-row gap-4 h-[calc(100vh-7rem)]">
      {/* Product List */}
      <div className="flex-1 flex flex-col min-w-0">
        <div className="flex flex-col sm:flex-row gap-3 mb-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Cari produk..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none text-sm"
            />
          </div>
          <select
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(Number(e.target.value))}
            className="px-4 py-2.5 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
          >
            <option value={0}>Semua Kategori</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
        </div>

        <div className="flex-1 overflow-auto">
          <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3">
            {products.map((product) => {
              const inCart = items.find((i) => i.product.id === product.id);
              return (
                <button
                  key={product.id}
                  onClick={() => addItem(product)}
                  disabled={product.stock <= 0}
                  className={`bg-white border rounded-xl p-3 text-left transition hover:shadow-md hover:border-blue-300 disabled:opacity-50 disabled:cursor-not-allowed ${
                    inCart ? 'border-blue-400 ring-2 ring-blue-100' : 'border-gray-200'
                  }`}
                >
                  <div className="flex items-start justify-between mb-2">
                    <span className="text-xs px-2 py-0.5 bg-gray-100 rounded-full text-gray-600">
                      {product.category?.name}
                    </span>
                    {inCart && (
                      <span className="text-xs px-2 py-0.5 bg-blue-100 rounded-full text-blue-700 font-medium">
                        x{inCart.quantity}
                      </span>
                    )}
                  </div>
                  <h3 className="font-medium text-gray-900 text-sm truncate">{product.name}</h3>
                  <p className="text-blue-600 font-bold text-sm mt-1">{formatCurrency(product.price)}</p>
                  <p className={`text-xs mt-1 ${product.stock <= product.min_stock ? 'text-red-500' : 'text-gray-400'}`}>
                    Stok: {product.stock}
                  </p>
                </button>
              );
            })}
            {products.length === 0 && (
              <div className="col-span-full text-center py-12 text-gray-400">
                Tidak ada produk ditemukan
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Cart */}
      <div className="w-full lg:w-96 bg-white border border-gray-200 rounded-xl flex flex-col">
        <div className="px-4 py-3 border-b border-gray-200 flex items-center gap-2">
          <ShoppingCart className="w-5 h-5 text-blue-600" />
          <h2 className="font-semibold text-gray-900">Keranjang</h2>
          <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full ml-auto">
            {items.length} item
          </span>
        </div>

        <div className="flex-1 overflow-auto p-4 space-y-3">
          {items.length === 0 ? (
            <p className="text-gray-400 text-center py-8 text-sm">Keranjang kosong</p>
          ) : (
            items.map((item) => (
              <div key={item.product.id} className="flex gap-3 pb-3 border-b border-gray-100 last:border-0">
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 truncate">{item.product.name}</p>
                  <p className="text-xs text-gray-500">{formatCurrency(item.product.price)}</p>
                </div>
                <div className="flex items-center gap-1.5">
                  <button
                    onClick={() => updateQuantity(item.product.id, item.quantity - 1)}
                    className="w-7 h-7 rounded-lg border border-gray-300 flex items-center justify-center hover:bg-gray-100"
                  >
                    <Minus className="w-3 h-3" />
                  </button>
                  <span className="w-8 text-center text-sm font-medium">{item.quantity}</span>
                  <button
                    onClick={() => updateQuantity(item.product.id, item.quantity + 1)}
                    className="w-7 h-7 rounded-lg border border-gray-300 flex items-center justify-center hover:bg-gray-100"
                  >
                    <Plus className="w-3 h-3" />
                  </button>
                  <button
                    onClick={() => removeItem(item.product.id)}
                    className="w-7 h-7 rounded-lg text-red-400 hover:bg-red-50 hover:text-red-600 flex items-center justify-center"
                  >
                    <Trash2 className="w-3 h-3" />
                  </button>
                </div>
                <p className="text-sm font-semibold text-gray-900 w-24 text-right">
                  {formatCurrency(item.product.price * item.quantity - item.discount)}
                </p>
              </div>
            ))
          )}
        </div>

        <div className="p-4 border-t border-gray-200 space-y-3">
          <select
            value={customerId || ''}
            onChange={(e) => setCustomerId(e.target.value ? Number(e.target.value) : null)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
          >
            <option value="">Tanpa pelanggan</option>
            {customers.map((c) => (
              <option key={c.id} value={c.id}>{c.name} ({c.phone})</option>
            ))}
          </select>

          <div className="flex items-center gap-2">
            <label className="text-sm text-gray-600">Diskon:</label>
            <input
              type="number"
              value={discount || ''}
              onChange={(e) => setDiscount(Number(e.target.value) || 0)}
              className="flex-1 px-3 py-1.5 border border-gray-300 rounded-lg text-sm text-right focus:ring-2 focus:ring-blue-500 outline-none"
              placeholder="0"
            />
          </div>

          <div className="space-y-1 text-sm">
            <div className="flex justify-between text-gray-600">
              <span>Subtotal</span>
              <span>{formatCurrency(getSubtotal())}</span>
            </div>
            {discount > 0 && (
              <div className="flex justify-between text-red-600">
                <span>Diskon</span>
                <span>-{formatCurrency(discount)}</span>
              </div>
            )}
            <div className="flex justify-between text-lg font-bold text-gray-900 pt-1 border-t">
              <span>Total</span>
              <span>{formatCurrency(getTotal())}</span>
            </div>
          </div>

          <button
            onClick={() => setShowPayment(true)}
            disabled={items.length === 0}
            className="w-full bg-blue-600 text-white py-3 rounded-lg font-medium hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <CreditCard className="w-5 h-5" />
            Bayar
          </button>
        </div>
      </div>

      {/* Payment Modal */}
      {showPayment && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl w-full max-w-md p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900">Pembayaran</h2>
              <button onClick={() => setShowPayment(false)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="text-center mb-6">
              <p className="text-sm text-gray-500">Total yang harus dibayar</p>
              <p className="text-3xl font-bold text-blue-600">{formatCurrency(getTotal())}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Jumlah Bayar</label>
              <input
                type="number"
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg text-lg text-right font-semibold focus:ring-2 focus:ring-blue-500 outline-none"
                placeholder="0"
                autoFocus
              />
            </div>

            <div className="grid grid-cols-3 gap-2 mt-3">
              {quickAmounts.map((amount) => (
                <button
                  key={amount}
                  onClick={() => setPaymentAmount(String(amount))}
                  className="px-3 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 transition"
                >
                  {formatCurrency(amount)}
                </button>
              ))}
            </div>

            <button
              onClick={() => setPaymentAmount(String(getTotal()))}
              className="w-full mt-2 px-3 py-2 border border-blue-300 rounded-lg text-sm text-blue-600 hover:bg-blue-50 transition"
            >
              Uang Pas
            </button>

            {paymentAmount && parseFloat(paymentAmount) >= getTotal() && (
              <div className="mt-4 p-3 bg-green-50 rounded-lg text-center">
                <p className="text-sm text-gray-600">Kembalian</p>
                <p className="text-xl font-bold text-green-600">
                  {formatCurrency(parseFloat(paymentAmount) - getTotal())}
                </p>
              </div>
            )}

            <button
              onClick={handlePayment}
              disabled={loading || !paymentAmount || parseFloat(paymentAmount) < getTotal()}
              className="w-full mt-4 bg-green-600 text-white py-3 rounded-lg font-medium hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Memproses...' : 'Proses Pembayaran'}
            </button>
          </div>
        </div>
      )}

      {/* Receipt Modal */}
      {completedTx && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl w-full max-w-sm p-6">
            <div className="text-center mb-4">
              <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-2" />
              <h2 className="text-xl font-bold text-gray-900">Transaksi Berhasil!</h2>
            </div>

            <div ref={receiptRef} className="border border-gray-200 rounded-lg p-4 text-sm">
              <div className="text-center border-b pb-3 mb-3">
                <p className="font-bold text-lg">Smart Cashier</p>
                <p className="text-gray-500 text-xs">{completedTx.invoice_number}</p>
                <p className="text-gray-500 text-xs">
                  {new Date(completedTx.created_at).toLocaleString('id-ID')}
                </p>
              </div>

              {completedTx.items?.map((item, idx) => (
                <div key={idx} className="flex justify-between py-1">
                  <div>
                    <p>{item.product_name}</p>
                    <p className="text-xs text-gray-500">
                      {item.quantity} x {formatCurrency(item.price)}
                    </p>
                  </div>
                  <p className="font-medium">{formatCurrency(item.subtotal)}</p>
                </div>
              ))}

              <div className="border-t mt-3 pt-3 space-y-1">
                <div className="flex justify-between">
                  <span>Subtotal</span>
                  <span>{formatCurrency(completedTx.subtotal)}</span>
                </div>
                {completedTx.discount > 0 && (
                  <div className="flex justify-between text-red-600">
                    <span>Diskon</span>
                    <span>-{formatCurrency(completedTx.discount)}</span>
                  </div>
                )}
                <div className="flex justify-between font-bold text-base">
                  <span>Total</span>
                  <span>{formatCurrency(completedTx.total)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Bayar</span>
                  <span>{formatCurrency(completedTx.payment_amount)}</span>
                </div>
                <div className="flex justify-between font-bold text-green-600">
                  <span>Kembalian</span>
                  <span>{formatCurrency(completedTx.change_amount)}</span>
                </div>
              </div>

              {completedTx.customer && (
                <div className="border-t mt-3 pt-2 text-xs text-gray-500">
                  Pelanggan: {completedTx.customer.name}
                </div>
              )}
            </div>

            <div className="flex gap-3 mt-4">
              <button
                onClick={() => {
                  const printContent = receiptRef.current;
                  if (printContent) {
                    const w = window.open('', '', 'width=300,height=600');
                    w?.document.write('<html><head><title>Struk</title><style>body{font-family:monospace;font-size:12px;padding:10px}table{width:100%}.right{text-align:right}.bold{font-weight:bold}.center{text-align:center}.line{border-top:1px dashed #000;margin:8px 0}</style></head><body>');
                    w?.document.write(printContent.innerHTML);
                    w?.document.write('</body></html>');
                    w?.document.close();
                    w?.print();
                  }
                }}
                className="flex-1 border border-gray-300 py-2.5 rounded-lg font-medium text-gray-700 hover:bg-gray-50 transition flex items-center justify-center gap-2"
              >
                <Printer className="w-4 h-4" />
                Cetak
              </button>
              <button
                onClick={() => setCompletedTx(null)}
                className="flex-1 bg-blue-600 text-white py-2.5 rounded-lg font-medium hover:bg-blue-700 transition"
              >
                Transaksi Baru
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
