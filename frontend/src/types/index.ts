export interface User {
  id: number;
  username: string;
  name: string;
  role: 'admin' | 'cashier';
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  name: string;
  sku: string;
  price: number;
  stock: number;
  min_stock: number;
  image: string;
  category_id: number;
  category: Category;
  created_at: string;
  updated_at: string;
}

export interface Customer {
  id: number;
  name: string;
  phone: string;
  email: string;
  points: number;
  created_at: string;
  updated_at: string;
}

export interface TransactionItem {
  id: number;
  transaction_id: number;
  product_id: number;
  product_name: string;
  price: number;
  quantity: number;
  discount: number;
  subtotal: number;
}

export interface Transaction {
  id: number;
  invoice_number: string;
  user_id: number;
  user: User;
  customer_id: number | null;
  customer: Customer | null;
  items: TransactionItem[];
  subtotal: number;
  discount: number;
  total: number;
  payment_amount: number;
  change_amount: number;
  created_at: string;
}

export interface CartItem {
  product: Product;
  quantity: number;
  discount: number;
}

export interface DashboardData {
  today_revenue: number;
  today_transactions: number;
  total_products: number;
  low_stock_count: number;
  recent_transactions: Transaction[];
  top_products: TopProduct[];
}

export interface TopProduct {
  product_id: number;
  product_name: string;
  quantity: number;
  revenue: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> extends ApiResponse<T> {
  total: number;
  page: number;
  limit: number;
}
