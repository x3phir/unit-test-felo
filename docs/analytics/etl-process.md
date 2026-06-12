# FELO Data Warehouse ETL/ELT Process

Dokumen ini menjelaskan rancangan ekstraksi, transformasi, loading, rekonsiliasi, dan data quality untuk Data Warehouse FELO.

## Arsitektur Layer

| Layer | Fungsi | Contoh Isi |
|---|---|---|
| Raw Layer | Menyimpan hasil ekstraksi apa adanya dari source service | Snapshot `trips`, `orders`, `payment_transactions`, `wallet_transactions` |
| Staging Layer | Membersihkan data, standardisasi tipe, deduplication, timezone conversion, dan validasi awal | `stg_trips`, `stg_payment_transactions`, `stg_users` |
| Warehouse Layer | Menyimpan model dimensional final | `dim_*`, `fact_*` |
| Mart Layer | Menyimpan agregasi untuk dashboard | `mart_daily_revenue`, `mart_driver_performance`, `mart_food_merchant_performance` |

## Pola Ekstraksi

- Ekstraksi awal memakai batch incremental setiap 15 menit, 1 jam, atau harian sesuai kebutuhan domain.
- Incremental key memakai `created_at`, `updated_at`, atau event timestamp.
- Watermark disimpan per source table, misalnya `last_successful_extracted_at`.
- Source finance seperti payment, invoice, dan wallet memakai window overlap untuk menghindari transaksi terlambat.
- Semua timestamp dikonversi ke `Asia/Jakarta` di staging.
- Ekstraksi tidak melakukan query analitik berat ke database OLTP.

Contoh incremental load:

```text
Extract trips where updated_at > last_watermark_trip
Extract trip_state_logs where event_timestamp > last_watermark_trip_state
Extract payment_transactions where updated_at > last_watermark_payment
Extract wallet_transactions where updated_at > last_watermark_wallet
```

## Source to Target Mapping

| Source Service | Tabel/Collection Sumber | Data yang Diambil | Target DWH |
|---|---|---|---|
| User Service | `users` | `user_id`, status, registered timestamp, city, verification status | `dim_user` |
| Auth Service | `auth_users`, `login_sessions` | created timestamp, phone/email verification, login activity | Enrichment `dim_user`, optional login mart |
| Driver Service | `driver_profiles`, `driver_availability` | `driver_id`, vehicle, driver status, city, rating | `dim_driver` |
| Ride/Trip Service | `trips`, `trip_state_logs`, `qr_codes` | Trip data, status history, FELO-Now marker | `fact_trip`, `fact_trip_status_event` |
| Food Order Service | `orders`, `order_items`, `validation_logs` | Order, item detail, COD/geofence validation | `fact_food_order`, `fact_food_order_item`, `fact_geofence_event` |
| Merchant Service | `merchants`, `merchant_categories` | Merchant profile, category, city, status | `dim_merchant` |
| Send/Shipment Service | `send_orders`, `shipments`, `shipment_tracking` | Sender, receiver zone, package, lifecycle status | `fact_shipment` |
| Matching Service | `match_requests`, `match_attempts` | Matching attempt, accepted/rejected/timeout | `fact_matching_attempt` |
| Tracking Service | `tracking_sessions`, `gps_points` | Route, pickup/dropoff tracking, distance/duration enrichment | Fact enrichment, optional route mart |
| Pricing Service | `pricing_configs`, `pricing_logs` | Base fare, fee, surge, price calculation | Measures in facts, optional pricing audit |
| Payment Service | `payment_transactions` | Payment amount, status, method, idempotency key | `fact_payment` |
| Wallet Service | `wallets`, `transactions`, `idempotency_keys` | Credit/debit, settlement, balance | `fact_wallet_transaction` |
| Invoice Service | `invoices` | Invoice total, payer type, line amount | Reconciliation for fact measures |
| Feedback Service | `feedbacks`, `ratings` | Rating user-driver/merchant | `fact_feedback` |
| Location Service | `zones`, `geofence_events`, `location_logs` | Zone mapping, violation, risk score | `dim_location_zone`, `fact_geofence_event` |
| Promo/Voucher Service | `promos`, `promo_usages` | Promo campaign, discount usage | `dim_promo`, discount measures |

## Load Order

1. Load static and low-cardinality dimensions: `dim_date`, `dim_time`, `dim_service_type`, `dim_payment_method`, `dim_status`, `dim_cancel_reason`, `dim_vehicle`.
2. Load SCD Type 2 dimensions: `dim_user`, `dim_driver`, `dim_merchant`.
3. Load spatial dimensions: `dim_location_zone`.
4. Load campaign dimensions: `dim_promo`.
5. Load fact tables: `fact_trip`, `fact_trip_status_event`, `fact_food_order`, `fact_food_order_item`, `fact_shipment`, `fact_matching_attempt`, `fact_payment`, `fact_wallet_transaction`, `fact_feedback`, `fact_geofence_event`.
6. Load mart tables from warehouse facts and dimensions.

Jika foreign key dimensi tidak ditemukan, loader memakai row default seperti `unknown_user`, `unknown_driver`, `unknown_merchant`, `unknown_location`, `unknown_status`, dan `unknown_payment_method`.

## Standardisasi Staging

| Area | Aturan |
|---|---|
| Timestamp | Semua timestamp dikonversi ke `Asia/Jakarta` |
| Natural key | Semua source ID disimpan sebagai string agar aman lintas service |
| Amount | Semua amount memakai currency minor/major yang disepakati dan tipe decimal |
| Status | Status source dipetakan ke `dim_status.status_code` dan `status_group` |
| Location | Koordinat dipetakan ke `dim_location_zone`, bukan disimpan sebagai alamat lengkap |
| PII | Email, nomor HP, alamat lengkap, dan identitas pribadi tidak dimuat mentah |
| Flags | Boolean analitik disimpan sebagai `0`/`1` agar bisa dijumlahkan |

## Transformasi User

Sumber utama adalah `users` dari User Service. Auth Service dipakai sebagai enrichment untuk informasi verifikasi.

Aturan:

1. Ambil `user_id`, `created_at`, `status`, `city`, dan `verification_status` dari `users`.
2. Ambil `auth_created_at`, `phone_verified`, dan `email_verified` dari `auth_users` jika tersedia.
3. Jika `users.created_at` berbeda dari `auth_users.created_at`, gunakan `users.created_at` sebagai tanggal registrasi utama.
4. Jika `users.verification_status = verified` tetapi auth verification tidak lengkap, set `verification_status = partially_verified` atau tandai staging dengan `VERIFICATION_MISMATCH`.
5. Buat `user_hash = SHA256(user_id + salt)` dan jangan simpan PII mentah.
6. Jika `city`, `verification_status`, `user_type`, atau status aktif berubah, tutup versi lama dan buat versi baru di `dim_user`.

## Transformasi Driver

Sumber utama adalah `driver_profiles`. Availability realtime hanya dipakai untuk enrichment terbatas dan tidak menjadi atribut permanen jika terlalu sering berubah.

Aturan:

1. Ambil `driver_id`, vehicle, status operasional, city, dan join timestamp dari `driver_profiles`.
2. Hitung `rating_bucket` dari rata-rata rating driver di Feedback Service.
3. Gunakan bucket `<4.0`, `4.0-4.5`, dan `>4.5`.
4. Jika driver mengganti `vehicle_type`, `driver_status`, city, atau rating bucket, buat versi baru SCD Type 2.
5. Hash identitas driver dengan `driver_hash = SHA256(driver_id + salt)`.

## Transformasi Merchant

Sumber utama adalah `merchants` dan `merchant_categories`.

Aturan:

1. Ambil `merchant_id`, nama merchant, kategori, city, district, status aktif, dan trust score.
2. Masking nama merchant jika aturan privasi atau bisnis memerlukan.
3. Jika kategori kosong, isi `unknown`, bukan `NULL`.
4. Jika kategori, status aktif, city, district, atau trust score bucket berubah, buat versi baru SCD Type 2.

## Transformasi Trip FELO-City

Sumber: `trips`, `trip_state_logs`, `match_requests`, `match_attempts`, `pricing_logs`, `payment_transactions`, `invoices`, dan `qr_codes`.

Target: `fact_trip` dan `fact_trip_status_event`.

Aturan:

1. Satu baris `trips` menjadi satu baris `fact_trip`.
2. Satu baris `trip_state_logs` menjadi satu baris `fact_trip_status_event`.
3. Status akhir di `fact_trip.status_key` memakai status terakhir dari `trip_state_logs` berdasarkan event timestamp terbaru.
4. Jika `trips.status` berbeda dari status terakhir state log, pakai state log dan isi `data_quality_flag = STATUS_MISMATCH`.
5. `matching_duration_second` dihitung dari `matched_at - requested_at` atau event `requested` ke `matched`.
6. Jika trip terkait `qr_codes`, set `service_type_key = FELO_NOW`, `is_felo_now = 1`, dan `matching_duration_second` boleh `0` atau `NULL`.
7. Harga awal dapat diambil dari `pricing_logs`, tetapi nominal final mengacu ke `invoices` dan `payment_transactions`.
8. Jika `pricing_logs.final_fare` berbeda dari `invoices.total_amount`, pakai invoice sebagai `fare_amount` dan isi `AMOUNT_MISMATCH`.
9. `distance_km` diambil dari field final trip jika tersedia. Jika kosong, hitung dari pickup-dropoff atau GPS points.
10. Pickup dan dropoff dipetakan ke `dim_location_zone`.

## Transformasi Food Order

Sumber: `orders`, `order_items`, `merchants`, `validation_logs`, `payment_transactions`, `invoices`, dan `promo_usages`.

Target: `fact_food_order`, `fact_food_order_item`, dan `fact_geofence_event`.

Aturan:

1. Satu order menjadi satu baris `fact_food_order`.
2. Setiap item menjadi satu baris `fact_food_order_item`.
3. `gross_merchandise_value` = total subtotal item sebelum delivery fee, service fee, dan diskon.
4. `item_count` = jumlah `quantity` seluruh item.
5. Jika `orders.total_amount` berbeda dari item subtotal + fee - discount, pakai `invoices.total_amount` sebagai `total_paid_amount` dan isi `AMOUNT_MISMATCH`.
6. COD/geofencing berasal dari `validation_logs`.
7. Jika `validation_logs.rule_code = COD_DISTANCE_GT_1KM`, set `is_cod_locked = 1` dan buat baris `fact_geofence_event`.
8. Merchant zone berasal dari lokasi merchant. Dropoff zone berasal dari lokasi customer.

## Transformasi FELO-Send / Shipment

Sumber: `send_orders`, `shipments`, `shipment_tracking`, `invoices`, dan `payment_transactions`.

Target: `fact_shipment`.

Aturan:

1. Satu shipment menjadi satu baris `fact_shipment`.
2. `payer_type` distandarkan menjadi `sender` atau `receiver`.
3. Jika `send_orders.payer_type` berbeda dari `invoices.payer_type`, gunakan invoice sebagai sumber final.
4. `pickup_duration_minute` dihitung dari waktu driver matched sampai pickup.
5. `delivery_duration_minute` dihitung dari waktu pickup sampai delivered.
6. Jika proof of delivery tersedia di tracking atau media service, isi `has_proof_of_delivery = 1`.
7. Jika status shipment berbeda dari status terakhir tracking, gunakan status tracking dan isi `STATUS_MISMATCH`.

## Transformasi Matching

Sumber: `match_requests`, `match_attempts`, dan `driver_availability`.

Target: `fact_matching_attempt`.

Aturan:

1. Satu baris `match_attempts` menjadi satu baris `fact_matching_attempt`.
2. `attempt_number` dibuat berdasarkan urutan attempt dalam satu `match_request_id`.
3. `is_success = 1` jika status attempt `accepted`.
4. `is_rejected = 1` jika status attempt `rejected`.
5. `is_timeout = 1` jika status attempt `timeout`.
6. Jika lebih dari satu accepted driver dalam satu `match_request_id`, pilih accepted pertama berdasarkan timestamp dan tandai attempt lain dengan `DUPLICATE_ACCEPTED_ATTEMPT` di staging.
7. `matching_duration_second` dihitung dari `attempt_started_at` sampai `attempt_finished_at`.

## Transformasi Payment

Sumber: `payment_transactions`, `invoices`, `wallet_transactions`, dan `idempotency_keys`.

Target: `fact_payment`.

Aturan:

1. Satu payment transaction menjadi satu baris `fact_payment`.
2. Deduplicate payment memakai `idempotency_key`.
3. Jika ada dua transaksi dengan `idempotency_key` sama, pilih transaksi terbaru yang statusnya final.
4. Transaksi duplicate ditandai `DUPLICATE_IDEMPOTENCY_KEY` di staging atau audit table.
5. `amount` berasal dari `payment_transactions.amount`.
6. Jika payment sukses, isi `paid_amount = amount` dan `failed_amount = 0`.
7. Jika payment gagal, isi `failed_amount = amount` dan `paid_amount = 0`.
8. Jika payment refund, isi `refund_amount` sesuai nominal refund.
9. Jika payment amount berbeda dari invoice total, gunakan invoice sebagai nilai bisnis final dan isi `AMOUNT_MISMATCH`.

## Transformasi Wallet/Settlement

Sumber: `wallets`, `transactions`, dan `payment_transactions`.

Target: `fact_wallet_transaction`.

Aturan:

1. Satu transaksi wallet menjadi satu baris fact.
2. Jika `transaction_type = settlement`, isi `is_settlement = 1`.
3. Amount positif masuk ke `credit_amount`.
4. Amount negatif masuk ke `debit_amount` sebagai nilai absolut.
5. `balance_after_amount` diambil dari wallet transaction jika tersedia.
6. Jika transaksi tidak memiliki referensi payment/order, tetap simpan dengan `service_type_key = GENERAL`.

## Transformasi Lokasi

Sumber: pickup/dropoff lat-lng dari trip, order, send order, shipment, GPS points, dan `zones` dari Location Service.

Target: `dim_location_zone`.

Aturan:

1. Semua koordinat dipetakan ke zona analitik.
2. Zona dapat dibuat dari kecamatan/kota, geohash, H3 grid, atau custom business area.
3. Laporan manajemen memakai kota/kecamatan.
4. Analisis operasional seperti hotspot demand memakai H3/geohash.
5. Koordinat kosong atau tidak valid diarahkan ke zone `unknown`.
6. Alamat lengkap customer tidak disimpan di DWH.

## Conflict Handling

| Konflik | Sumber Final | Data Quality Flag |
|---|---|---|
| `trips.status` berbeda dari status terakhir `trip_state_logs` | `trip_state_logs` | `STATUS_MISMATCH` |
| `shipments.status` berbeda dari status terakhir `shipment_tracking` | `shipment_tracking` | `STATUS_MISMATCH` |
| Pricing final fare berbeda dari invoice total | `invoices` | `AMOUNT_MISMATCH` |
| Payment amount berbeda dari invoice total | `invoices` untuk nilai bisnis final | `AMOUNT_MISMATCH` |
| Dua payment memiliki idempotency key sama | Payment final terbaru | `DUPLICATE_IDEMPOTENCY_KEY` |
| Lebih dari satu accepted matching attempt | Accepted pertama berdasarkan timestamp | `DUPLICATE_ACCEPTED_ATTEMPT` |
| User verification tidak konsisten antara User dan Auth | User Service sebagai master profile | `VERIFICATION_MISMATCH` |
| Koordinat tidak valid | Zone `unknown` | `INVALID_LOCATION` |

## Data Quality Rules

| Check | Tujuan |
|---|---|
| Primary key tidak null | Memastikan data bisa dilacak |
| Natural key tidak duplikat pada dimensi current | Mencegah duplikasi user, driver, dan merchant aktif |
| Amount tidak negatif kecuali refund/debit | Mencegah nominal salah |
| Status harus valid di `dim_status` | Mencegah status liar |
| Payment sukses harus punya `paid_amount > 0` | Validasi transaksi berhasil |
| Completed trip/order/shipment harus punya timestamp selesai | Validasi lifecycle |
| Distance tidak ekstrem tanpa flag | Deteksi error GPS |
| Fact harus punya date key dan service type key | Memastikan data bisa dilaporkan |
| SCD Type 2 tidak boleh punya dua current row untuk natural key yang sama | Menjaga dimensi historis konsisten |
| PII tidak boleh masuk ke raw mart analitik | Menjaga privasi |

## Mart Layer Awal

| Mart | Grain | Sumber | Kegunaan |
|---|---|---|---|
| `mart_daily_revenue` | 1 baris = 1 tanggal + service group | `fact_trip`, `fact_food_order`, `fact_shipment` | Executive dashboard dan revenue trend |
| `mart_driver_performance` | 1 baris = 1 driver + periode | Trip, food, shipment, matching, feedback | Driver performance dashboard |
| `mart_food_merchant_performance` | 1 baris = 1 merchant + periode | `fact_food_order`, `fact_food_order_item`, `fact_feedback` | Merchant ranking dan food dashboard |
| `mart_payment_health` | 1 baris = 1 tanggal + payment method | `fact_payment` | Payment success rate dan failed amount |
| `mart_geofence_risk` | 1 baris = 1 tanggal + zone + rule | `fact_geofence_event` | Fraud/geofencing dashboard |

## Operasional ETL

- Jalankan dimensi sebelum fact.
- Gunakan idempotent load agar rerun tidak menggandakan data.
- Simpan metadata job seperti `job_name`, `started_at`, `finished_at`, `source_watermark`, `row_count`, dan `status`.
- Simpan rejected rows ke audit/staging error table, bukan langsung dibuang.
- Untuk finance, gunakan reconciliation harian antara invoice, payment, dan wallet.
- Untuk late arriving dimension, gunakan unknown key lalu lakukan backfill saat dimensi tersedia.
