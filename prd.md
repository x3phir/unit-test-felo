# FELO Mobile App — Product Requirements Document (REVISED)

| Field | Value |
|---|---|
| Version | 1.1 — MVP (Revised) |
| Status | Final Draft |
| Last Updated | May 2024 |
| Platform | iOS (Swift) + Android (Kotlin) |
| Audience | PM, Backend Engineer, Mobile Engineer, Stakeholders |

---

## 1. Overview
FELO adalah platform on-demand yang menyediakan layanan transportasi (FELO-City), pesan-antar makanan (FELO-Food), dan pengiriman barang (FELO-Send). MVP ini fokus pada efisiensi operasional, biaya rendah, dan sistem mitigasi kecurangan (fraud) yang inovatif.

### 1.1 Product Suite
1.  **FELO-City:** Layanan ride-hailing reguler dan fitur instan (FELO-Now).
2.  **FELO-Food:** Pesan-antar makanan dengan sistem Merchant POS terintegrasi.
3.  **FELO-Send:** Pengiriman barang (single-stop).

### 1.2 Tech Stack & Infra
- **Architecture:** 18 Microservices.
- **Messaging:** RabbitMQ (Asynchronous & Consistency).
- **Maps:** OpenStreetMap (OSM) / MapLibre (Free tier/Open Source).
- **Communication:** REST API & WebSockets.

---

## 2. Goals
- **G1:** Menyediakan layanan transportasi dengan matching < 5 detik.
- **G2:** Mitigasi order fiktif pada layanan Food hingga 95% melalui validasi lokasi.
- **G3:** Pengelolaan keuangan transparan menggunakan sistem Escrow.

---

## 3. User Stories & Functional Requirements

### 3.1 FELO-Now (QR Instant Ride)
- **Pricing:** Tarif FELO-Now sama dengan tarif Regular Ride.
- **Prioritas Matching:** Jika driver mendapatkan tawaran dari sistem matching, driver tidak bisa di-scan oleh penumpang FELO-Now. Sistem matching memiliki prioritas lebih tinggi.
- **Flow:** Penumpang scan QR Driver → Validasi Driver Availability → Konfirmasi → Perjalanan dimulai.

### 3.2 FELO Wallet & Escrow System
- **Escrow Ledger:** Dana dari transaksi (Wallet/E-wallet) masuk ke rekening penampung (Escrow) selama perjalanan berlangsung.
- **Settlement:** Dana dipindahkan dari Escrow ke Driver Earning setelah status perjalanan "Completed".
- **Refund Logic (SAGA Pattern):** Jika terjadi kegagalan sistem setelah saldo terpotong, RabbitMQ akan memicu event kompensasi untuk mengembalikan saldo secara otomatis.
- **Automated Reversal:** Jika pesanan Food tidak selesai dalam 1x24 jam, dana di Escrow otomatis dikembalikan ke User.

### 3.3 Merchant POS App
- **Unified App:** Aplikasi FELO untuk user memiliki mode "Merchant POS".
- **Role:** User yang terdaftar sebagai Owner Merchant dapat beralih ke tampilan POS untuk manajemen menu, pesanan, dan laporan penjualan.

### 3.4 Mitigasi Order Fiktif (Khusus FELO-Food)
Sistem memiliki dua opsi saat checkout:
1.  **Kirim untuk Saya:**
    - User wajib mengaktifkan *Live Shared Location* saat memesan hingga makanan sampai.
    - Driver dapat melacak posisi user secara real-time.
2.  **Kirim untuk Teman:**
    - Pembayaran wajib E-Wallet (untuk keamanan) atau COD dengan syarat tambahan.
    - Penerima wajib melakukan validasi melalui notifikasi App atau WhatsApp/SMS (jika belum punya App).
    - Penerima wajib mengirimkan *Shared Location* satu kali (one-time consent) untuk divalidasi oleh sistem terhadap titik koordinat pengiriman.

---

## 4. Business Rules & Penalties
- **Penalty Score:** Pembatalan sepihak oleh User (terutama pada metode Cash) tidak memotong saldo, namun mengurangi *Reputation Score*. User dengan score rendah akan dibatasi fiturnya (misal: dilarang COD).
- **Account Security:** 1 Nomor HP = 1 User = 1 Device. Sistem melakukan pengecekan Device ID/IMEI untuk mencegah login ganda.
- **Commission:** 0% Profit-taking (Perusahaan belum mengambil keuntungan dari pendapatan driver di fase MVP).

---

## 5. Non-Functional Requirements (NFR)

### 5.1 Adaptive GPS Polling
- **Mode On-Ride:** Polling GPS setiap 3-5 detik (Akurasi tinggi).
- **Mode Idle (Searching):** Polling GPS turun ke 15-30 detik untuk menghemat baterai.
- **Low Battery:** Jika baterai < 15%, polling turun ke 60 detik.

### 5.2 Zombie Trip Protection
- Jika aplikasi Driver offline > 10 menit saat trip aktif, sistem mengirimkan notifikasi konfirmasi ke Customer.
- Customer dapat menutup trip secara manual jika driver benar-benar hilang.
- Data koordinat selama offline disimpan di *Local Cache* HP Driver dan di-sync saat internet kembali tersedia.

### 5.3 KYC (Know Your Customer)
- SLA Verifikasi: Admin wajib memproses dokumen (KTP/STNK) dalam waktu maksimal 1x24 jam.

---

## 6. Architecture Appendix (The 18 Services)
Layanan inti meliputi: `user-service`, `driver-service`, `merchant-service`, `order-service`, `ride-service`, `shipment-service`, `payment-service`, `wallet-service`, `pricing-service`, `matching-service`, `tracking-service`, `location-service`, `notification-service`, `feedback-service`, `invoice-service`, `fraud-service`, `promotion-service`, `chat-service` (internal system use).

---

## 7. Cross-Context Rules
1. Komunikasi antar-service WAJIB menggunakan RabbitMQ untuk menjaga integritas data (Event-driven).
2. Dilarang keras *Database Sharing* antar microservice.
3. Semua perhitungan tarif terpusat di `pricing-service`.

---