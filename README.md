# Smart Cashier

Aplikasi kasir pintar berbasis web untuk mengelola penjualan, produk, pelanggan, dan laporan bisnis secara real-time.

## Fitur

- **Autentikasi** — Login dengan JWT, role-based access (Admin & Kasir)
- **Point of Sale (POS)** — Transaksi cepat dengan pencarian produk, keranjang, diskon, dan struk
- **Manajemen Produk** — CRUD produk dengan kategori, harga, stok, dan foto
- **Manajemen Kategori** — Kelola kategori produk
- **Manajemen Pelanggan** — Data pelanggan dengan riwayat transaksi
- **Histori Transaksi** — Lihat detail setiap transaksi yang pernah dilakukan
- **Dashboard** — Ringkasan penjualan harian, produk terlaris, dan statistik bisnis
- **Laporan** — Laporan harian & bulanan dengan grafik interaktif

## Tech Stack

| Layer    | Teknologi                                          |
| -------- | -------------------------------------------------- |
| Frontend | React 19, TypeScript, Vite 5, Tailwind CSS 3       |
| Backend  | Go (Gin), GORM, JWT                                |
| Database | SQLite                                             |
| Infra    | Docker, Docker Compose, Nginx                      |

## Struktur Proyek

```
Smart-Cashier/
├── backend/
│   ├── cmd/main.go              # Entry point
│   ├── internal/
│   │   ├── config/              # Konfigurasi environment
│   │   ├── database/            # Koneksi DB & seeder
│   │   ├── handler/             # HTTP handlers
│   │   ├── middleware/          # Auth & role middleware
│   │   ├── model/               # Data models & DTOs
│   │   ├── repository/          # Database queries
│   │   └── service/             # Business logic
│   ├── pkg/                     # Shared packages (jwt, response)
│   ├── uploads/                 # File uploads
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/          # Layout, ProtectedRoute
│   │   ├── pages/               # Login, Dashboard, POS, dll
│   │   ├── services/            # API service layer
│   │   ├── store/               # Zustand state management
│   │   ├── types/               # TypeScript interfaces
│   │   └── utils/               # Helper functions
│   ├── nginx.conf               # Nginx config (Docker)
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
└── README.md
```

## Quick Start

### Dengan Docker (Rekomendasi)

Pastikan [Docker](https://docs.docker.com/get-docker/) sudah terinstal.

```bash
# Clone repository
git clone <repository-url>
cd Smart-Cashier

# Jalankan semua service
docker compose up --build

# Akses aplikasi di browser
# http://localhost:3000
```

Untuk menghentikan:

```bash
docker compose down
```

Untuk menghapus data (reset database):

```bash
docker compose down -v
```

### Manual (Development)

**Prasyarat:**
- Node.js >= 20
- Go >= 1.21
- Git

**Backend:**

```bash
cd backend
go mod download
go run ./cmd/main.go
# Server berjalan di http://localhost:8080
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
# Aplikasi berjalan di http://localhost:5173
```

## Akun Default

| Role  | Username | Password  |
| ----- | -------- | --------- |
| Admin | admin    | admin123  |
| Kasir | kasir    | kasir123  |

Akun default otomatis dibuat saat pertama kali menjalankan aplikasi.

## Environment Variables

| Variable     | Default                          | Deskripsi             |
| ------------ | -------------------------------- | --------------------- |
| `PORT`       | `8080`                           | Port backend API      |
| `JWT_SECRET` | `smart-cashier-secret-key-2026`  | Secret key untuk JWT  |
| `DB_PATH`    | `smart_cashier.db`               | Path file SQLite      |
| `GIN_MODE`   | `debug`                          | Mode Gin (`release`)  |

## API Endpoints

| Method | Endpoint                | Auth     | Deskripsi              |
| ------ | ----------------------- | -------- | ---------------------- |
| POST   | `/api/auth/login`       | -        | Login                  |
| GET    | `/api/auth/me`          | Token    | Profil user            |
| POST   | `/api/auth/register`    | Admin    | Registrasi user baru   |
| GET    | `/api/products`         | Token    | Daftar produk          |
| POST   | `/api/products`         | Admin    | Tambah produk          |
| PUT    | `/api/products/:id`     | Admin    | Edit produk            |
| DELETE | `/api/products/:id`     | Admin    | Hapus produk           |
| GET    | `/api/categories`       | Token    | Daftar kategori        |
| POST   | `/api/categories`       | Admin    | Tambah kategori        |
| PUT    | `/api/categories/:id`   | Admin    | Edit kategori          |
| DELETE | `/api/categories/:id`   | Admin    | Hapus kategori         |
| GET    | `/api/customers`        | Token    | Daftar pelanggan       |
| POST   | `/api/customers`        | Token    | Tambah pelanggan       |
| PUT    | `/api/customers/:id`    | Token    | Edit pelanggan         |
| DELETE | `/api/customers/:id`    | Admin    | Hapus pelanggan        |
| GET    | `/api/transactions`     | Token    | Daftar transaksi       |
| GET    | `/api/transactions/:id` | Token    | Detail transaksi       |
| POST   | `/api/transactions`     | Token    | Buat transaksi baru    |
| GET    | `/api/dashboard`        | Token    | Data dashboard         |
| GET    | `/api/reports/daily`    | Token    | Laporan harian         |
| GET    | `/api/reports/monthly`  | Token    | Laporan bulanan        |
| GET    | `/api/reports/top-products` | Token | Produk terlaris     |
| GET    | `/api/reports/summary`  | Token    | Ringkasan laporan      |
| GET    | `/api/health`           | -        | Health check           |

## License

MIT
