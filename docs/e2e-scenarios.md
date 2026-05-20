# FELO E2E Functional Test Scenarios

Dokumen ini merangkum skenario pengujian End-to-End (E2E) yang diimplementasikan untuk memastikan integritas sistem FELO lintas microservices.

## 1. Smoke Tests
Tujuan: Memastikan infrastruktur dasar dan konektivitas antar service berjalan.

| Skenario | Deskripsi | Service Terlibat |
|---|---|---|
| **Service Reachability** | Melakukan health check ke semua service pendukung (Ride, Matching, Location, Payment, Wallet, dll). | All Services |

## 2. Critical Flow (Happy Path)
Tujuan: Memastikan fitur utama berjalan mulus dari hulu ke hilir.

### 2.1 Regular Ride to Settlement
Alur lengkap perjalanan reguler FELO-City.
1. **Request Ride**: Customer meminta perjalanan dengan koordinat penjemputan & tujuan.
2. **Matching**: Memastikan `matching-service` memberikan driver yang sesuai.
3. **Location Tracking**: Driver mengirimkan koordinat GPS terbaru.
4. **Ride Lifecycle**: Perjalanan dimulai (`Started`) dan diselesaikan (`Completed`).
5. **Settlement**: Memastikan pembayaran diproses dan saldo driver dikreditkan.
6. **Event Validation**: Verifikasi event `ride.completed.v1` dan `payment.completed.v1` terbit di RabbitMQ.

### 2.2 FELO-Now QR Ride to Settlement
Alur instan menggunakan fitur QR.
1. **QR Generation**: Customer membuat QR code dengan tujuan yang sudah ditentukan.
2. **QR Scan**: Driver melakukan scan terhadap QR Customer.
3. **Acceptance**: Driver menyetujui perjalanan.
4. **Completion**: Perjalanan selesai dan pembayaran diproses secara otomatis.

## 3. Negative & Edge Case Flow
Tujuan: Memastikan sistem menangani kesalahan dengan benar (Resiliency).

| Skenario | Deskripsi | Ekspektasi |
|---|---|---|
| **No Driver Available** | Customer meminta ride di area yang tidak ada driver. | Sistem melakukan retry dan akhirnya memberikan status "No Match". |
| **Payment Failure** | Customer dengan saldo tidak mencukupi menyelesaikan perjalanan. | Sistem menerbitkan event `payment.failed.v1` untuk ditindaklanjuti oleh Fraud/Ops service. |
| **QR Expiry** | Driver mencoba melakukan scan terhadap QR yang sudah melewati batas waktu (10 menit). | Sistem menolak scan dan mengembalikan error "QR Expired". *(Status: In Development)* |
| **Duplicate Settlement** | Event penyelesaian yang sama dikirim dua kali ke `wallet-service`. | Idempotensi terjaga; saldo driver tidak bertambah dua kali. *(Status: In Development)* |

---
**Catatan Teknis:**
- Tests dijalankan menggunakan build tag `//go:build e2e`.
- Menggunakan `test-harness` untuk abstraksi komunikasi antar service (gRPC/REST/RabbitMQ).
- Semua test bersifat isolated menggunakan data seeder khusus.
