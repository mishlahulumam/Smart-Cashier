import { create } from 'zustand';
import type { CartItem, Product } from '../types';

interface CartState {
  items: CartItem[];
  discount: number;
  customerId: number | null;
  addItem: (product: Product) => void;
  removeItem: (productId: number) => void;
  updateQuantity: (productId: number, quantity: number) => void;
  updateItemDiscount: (productId: number, discount: number) => void;
  setDiscount: (discount: number) => void;
  setCustomerId: (id: number | null) => void;
  clearCart: () => void;
  getSubtotal: () => number;
  getTotal: () => number;
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  discount: 0,
  customerId: null,

  addItem: (product) => {
    const items = get().items;
    const existing = items.find((i) => i.product.id === product.id);
    if (existing) {
      if (existing.quantity < product.stock) {
        set({
          items: items.map((i) =>
            i.product.id === product.id ? { ...i, quantity: i.quantity + 1 } : i
          ),
        });
      }
    } else {
      if (product.stock > 0) {
        set({ items: [...items, { product, quantity: 1, discount: 0 }] });
      }
    }
  },

  removeItem: (productId) => {
    set({ items: get().items.filter((i) => i.product.id !== productId) });
  },

  updateQuantity: (productId, quantity) => {
    if (quantity <= 0) {
      get().removeItem(productId);
      return;
    }
    set({
      items: get().items.map((i) =>
        i.product.id === productId ? { ...i, quantity: Math.min(quantity, i.product.stock) } : i
      ),
    });
  },

  updateItemDiscount: (productId, discount) => {
    set({
      items: get().items.map((i) =>
        i.product.id === productId ? { ...i, discount } : i
      ),
    });
  },

  setDiscount: (discount) => set({ discount }),
  setCustomerId: (id) => set({ customerId: id }),

  clearCart: () => set({ items: [], discount: 0, customerId: null }),

  getSubtotal: () => {
    return get().items.reduce(
      (sum, item) => sum + item.product.price * item.quantity - item.discount,
      0
    );
  },

  getTotal: () => {
    const subtotal = get().getSubtotal();
    return Math.max(0, subtotal - get().discount);
  },
}));
