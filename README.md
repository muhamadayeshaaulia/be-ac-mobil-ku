# 🚗 Backend Sistem Rekomendasi Layanan Service AC Mobil

Sistem backend ini dikembangkan menggunakan bahasa pemrograman **Go (Golang)** dengan **Gin Web Framework** menggunakan **Clean Architecture** (domain, repository, usecase, delivery).

Sistem ini memiliki fitur rekomendasi hybrid: **Collaborative Filtering (User-Based Cosine Similarity)** untuk preferensi pengguna dan **Constraint Filtering (Haversine Distance)** untuk jarak lokasi spasial GPS secara real-time.

---

## 🏗️ Struktur Arsitektur (Clean Architecture)

```text
├── domain/          # Model Entitas & Kontrak Interface 
│   ├── user.go
│   ├── bengkel.go
│   ├── layanan.go
│   ├── booking.go
│   ├── rating.go
│   └── recommendation.go
├── repository/      # Implementasi Database Layer (GORM)
│   ├── db.go        # Koneksi GORM & Automigrasi (PostgreSQL/SQLite)
│   ├── gorm_user.go
│   ├── gorm_bengkel.go
│   ├── gorm_layanan.go
│   ├── gorm_booking.go
│   └── gorm_rating.go
├── usecase/         # Logika Bisnis Utama (CF, Haversine, Booking)
│   ├── user_usecase.go
│   ├── bengkel_usecase.go
│   ├── layanan_usecase.go
│   ├── booking_usecase.go
│   ├── rating_usecase.go
│   └── recommendation_usecase.go
├── delivery/        # HTTP Layer (Gin Handlers & RESTful Routes)
│   ├── user_handler.go
│   ├── bengkel_handler.go
│   ├── layanan_handler.go
│   ├── booking_handler.go
│   ├── rating_handler.go
│   └── recommendation_handler.go
├── middleware/      # Keamanan & Validasi Token Firebase Admin SDK
│   └── firebase_auth.go
├── firebase-service-account.json # Kredensial Firebase Admin SDK
├── main.go          # Dependency Injection & Bootstrapping
└── go.mod           # Daftar Dependensi Modul Go
```

---

## ⚙️ Variabel Lingkungan (Environment Variables)

Buat variabel lingkungan berikut untuk konfigurasi aplikasi:

| Nama Variabel | Tipe | Default / Keterangan |
|---|---|---|
| `APP_ENV` | `string` | `development` (aktifkan bypass Firebase Auth untuk testing) |
| `DB_DRIVER` | `string` | `sqlite` atau `postgres` |
| `DB_HOST` | `string` | `localhost` |
| `DB_PORT` | `string` | `5420` / `5432` |
| `DB_USER` | `string` | `postgres` |
| `DB_PASSWORD`| `string` | Sandi database PostgreSQL Anda |
| `DB_NAME` | `string` | `ac_mobil_ku` |
| `DB_SQLITE_PATH` | `string`| `ac_mobil_ku.db` (Jika menggunakan SQLite) |
| `FIREBASE_SERVICE_ACCOUNT_PATH` | `string` | `firebase-service-account.json` |
| `PORT` | `string` | `8080` (Port HTTP Server) |

---

## 🚀 Cara Menjalankan Server

1. **Unduh Dependensi:**
   ```bash
   go mod tidy
   ```

2. **Jalankan Aplikasi:**
   ```bash
   go run main.go
   ```

---

## 🧪 Metode Pengujian Offline / Dev Mode
Jika Anda menjalankan aplikasi dalam `APP_ENV=development`, Anda dapat menggunakan bypass token Firebase untuk mempermudah testing backend tanpa integrasi Google Sign-In langsung:
- Kirim Header: `Authorization: Bearer dev-token-[nama-pengguna]`
- Jika token mengandung kata `pengelola`, role otomatis diset sebagai `pengelola_bengkel`. Contoh: `dev-token-pengelola-dimas`.
- Jika tidak, role diset sebagai `pelanggan`. Contoh: `dev-token-pelanggan-budi`.

---

## 📌 Endpoint API Utama

Semua endpoint dilindungi oleh `Authorization: Bearer <token>` kecuali `/api/health`.

### 1. Autentikasi & Pengguna (User)
* **GET** `/api/user/profile` - Mengambil profil user yang masuk.
* **POST** `/api/user/profile` - Registrasi awal atau update profil (Pelanggan / Pengelola Bengkel).
* **POST** `/api/user/location` - Sinkronisasi koordinat GPS terkini pelanggan (`latitude`, `longitude`).

### 2. Bengkel
* **GET** `/api/bengkel` - Melihat daftar semua bengkel beserta ulasan rata-ratanya.
* **GET** `/api/bengkel/detail/:id` - Melihat detail lengkap satu bengkel.
* **GET** `/api/bengkel/my` - Mengambil profil bengkel yang dikelola oleh Pengelola Bengkel yang sedang login.
* **POST** `/api/bengkel` - Mendaftarkan profil bengkel baru.
* **PUT** `/api/bengkel` - Mengubah detail profil bengkel.

### 3. Katalog Layanan AC (Layanan)
* **GET** `/api/layanan/bengkel/:bengkel_id` - Melihat daftar layanan AC dari bengkel tertentu.
* **POST** `/api/layanan` - Menambahkan layanan baru ke katalog bengkel.
* **PUT** `/api/layanan/:id` - Memperbarui detail layanan (Nama, Harga, Ketersediaan).
* **DELETE** `/api/layanan/:id` - Menghapus layanan dari katalog.

### 4. Reservasi & Antrean (Booking)
* **POST** `/api/booking` - Membuat reservasi baru untuk layanan AC (Pelanggan).
* **GET** `/api/booking/history` - Melihat riwayat booking service (Pelanggan).
* **GET** `/api/booking/queue` - Melihat daftar antrean booking masuk ke bengkel (Pengelola Bengkel).
* **PUT** `/api/booking/:id/status` - Mengubah status antrean booking (`dikonfirmasi`, `selesai`, `dibatalkan`) (Pengelola Bengkel).

### 5. Rating & Ulasan
* **POST** `/api/rating` - Memberikan rating kualitas (1-5), rating harga (1-5), dan ulasan ulasan tertulis setelah booking berstatus `selesai`.
* **GET** `/api/rating/bengkel/:bengkel_id` - Melihat seluruh ulasan pelanggan pada bengkel terkait.

### 6. Rekomendasi Pintar (Hybrid Filtering)
* **GET** `/api/recommendations?limit=5` - Menghasilkan Top-N bengkel terdekat dan terbaik untuk pelanggan.
  * **Algoritma:** Memfilter bengkel dengan **Haversine Distance** (Maksimal 50 km) lalu memprediksi preferensi rating menggunakan **Cosine Similarity (Collaborative Filtering)**.