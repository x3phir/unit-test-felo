# FELO Data Warehouse Dashboard Requirements

Dokumen ini mendefinisikan dashboard dan laporan utama yang memakai FELO Data Warehouse. Dashboard dibangun di atas star schema `dim_*` dan `fact_*`, serta dapat memakai mart agregasi untuk performa query.

## Prinsip Dashboard

- Semua metric harus punya formula yang eksplisit.
- Dashboard tidak langsung membaca database OLTP microservices.
- Filter utama yang konsisten: date range, service group, service type, city, district/zone, payment method, status group.
- Amount finansial memakai nilai yang sudah direkonsiliasi dengan invoice/payment.
- Status lifecycle memakai status final dari log/tracking jika terjadi konflik.
- Flag seperti `is_completed`, `is_cancelled`, dan `is_success` disimpan sebagai `0`/`1` agar bisa dijumlahkan.

## Shared Filters

| Filter | Dimension |
|---|---|
| Date range | `dim_date` |
| Time bucket atau peak hour | `dim_time` |
| Service group/type | `dim_service_type` |
| City/district/zone | `dim_location_zone`, `dim_user`, `dim_driver`, `dim_merchant` |
| Payment method | `dim_payment_method` |
| Status/status group | `dim_status` |
| Promo/campaign | `dim_promo` |
| Vehicle type | `dim_vehicle`, `dim_driver` |

## Executive Dashboard

Tujuan: melihat kondisi bisnis FELO secara keseluruhan.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_trip`, `fact_food_order`, `fact_shipment`, `fact_payment`, `fact_wallet_transaction` |
| KPI | Total transaksi, total GMV, total revenue/platform fee, active user, active driver, completion rate, cancellation rate, payment success rate |
| Chart | Line chart GMV per hari, bar chart transaksi per layanan, KPI card revenue/order/user/driver, donut chart kontribusi revenue, area chart pertumbuhan bulanan |
| Drill-down | Service group, city, payment method, status group |

Metric utama:

| Metric | Formula |
|---|---|
| Total Transactions | `SUM(completed/counted rows dari fact_trip, fact_food_order, fact_shipment)` |
| Total GMV | `SUM(fact_trip.fare_amount) + SUM(fact_food_order.gross_merchandise_value) + SUM(fact_shipment.delivery_fee_amount)` untuk transaksi completed |
| Total Revenue | `SUM(platform_fee_amount)` dari fact layanan |
| Completion Rate | `SUM(is_completed) * 1.0 / COUNT(*)` |
| Cancellation Rate | `SUM(is_cancelled) * 1.0 / COUNT(*)` |
| Payment Success Rate | `SUM(fact_payment.is_success) * 1.0 / COUNT(fact_payment.payment_id)` |

## FELO-City Operations Dashboard

Tujuan: memantau performa layanan ride-hailing.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_trip`, `fact_trip_status_event`, `fact_matching_attempt` |
| KPI | Trip requested, matched, completed, cancelled, average matching time, average trip duration, average distance, average fare |
| Chart | Funnel requested to completed, line chart completed trip per jam/hari, heatmap pickup zone, bar chart cancellation reason, bar chart reguler vs FELO-Now |
| Drill-down | City, pickup zone, dropoff zone, time bucket, vehicle type, service type |

Metric utama:

| Metric | Formula |
|---|---|
| Total Trip | `COUNT(fact_trip.trip_id)` |
| Completed Trip | `SUM(fact_trip.is_completed)` |
| Cancelled Trip | `SUM(fact_trip.is_cancelled)` |
| Average Matching Time | `AVG(fact_trip.matching_duration_second)` untuk non FELO-Now |
| Average Trip Duration | `AVG(fact_trip.duration_minute)` |
| FELO-Now Share | `SUM(is_felo_now) * 1.0 / COUNT(trip_id)` |

## FELO-Food Dashboard

Tujuan: melihat performa order makanan dan merchant.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_food_order`, `fact_food_order_item`, `fact_feedback`, `fact_geofence_event` |
| KPI | Total food order, GMV FELO-Food, completed order, average preparation time, average delivery time, COD locked order |
| Chart | Top 10 merchant by GMV, top 10 menu/item, food order trend, heatmap merchant zone, cancellation by actor/reason |
| Drill-down | Merchant category, merchant city, dropoff zone, payment method, promo, status group |

Metric utama:

| Metric | Formula |
|---|---|
| Food GMV | `SUM(fact_food_order.gross_merchandise_value)` |
| Completed Food Order | `SUM(fact_food_order.is_completed)` |
| Average Preparation Time | `AVG(preparation_duration_minute)` |
| Average Delivery Time | `AVG(delivery_duration_minute)` |
| COD Locked Order | `SUM(is_cod_locked)` |
| Top Item Quantity | `SUM(fact_food_order_item.quantity)` by item |

## FELO-Send Dashboard

Tujuan: memantau performa layanan pengiriman barang.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_shipment`, `fact_payment`, `fact_feedback` |
| KPI | Total shipment, shipment completed, shipment cancelled, average delivery duration, average distance, proof of delivery count |
| Chart | Shipment per hari, payer type sender vs receiver, top pickup zone, top dropoff zone, completion rate per city |
| Drill-down | Pickup zone, dropoff zone, payer type, vehicle type, payment method, status group |

Metric utama:

| Metric | Formula |
|---|---|
| Total Shipment | `COUNT(fact_shipment.shipment_id)` |
| Completed Shipment | `SUM(fact_shipment.is_completed)` |
| Cancelled Shipment | `SUM(fact_shipment.is_cancelled)` |
| Average Delivery Duration | `AVG(delivery_duration_minute)` |
| POD Rate | `SUM(has_proof_of_delivery) * 1.0 / COUNT(shipment_id)` |
| Sender Payer Share | `COUNT(*) FILTER (WHERE payer_type = 'sender') * 1.0 / COUNT(*)` |

## Driver Performance Dashboard

Tujuan: mengevaluasi performa driver lintas FELO-City, FELO-Food, dan FELO-Send.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_trip`, `fact_food_order`, `fact_shipment`, `fact_matching_attempt`, `fact_feedback`, `fact_wallet_transaction` |
| KPI | Completed job per driver, driver earning, average rating, acceptance rate, cancellation rate by driver, average pickup time |
| Chart | Leaderboard top driver, scatter rating vs completed order, bar chart earning per periode, acceptance rate trend |
| Drill-down | Driver city, vehicle type, rating bucket, service group, time period |

Metric utama:

| Metric | Formula |
|---|---|
| Completed Jobs | `SUM(is_completed)` dari trip, food order, dan shipment per driver |
| Driver Earning | `SUM(driver_earning_amount)` dari fact layanan atau settlement wallet |
| Average Rating | `AVG(fact_feedback.rating_value)` |
| Acceptance Rate | `SUM(fact_matching_attempt.is_success) * 1.0 / COUNT(fact_matching_attempt.matching_attempt_key)` |
| Timeout Rate | `SUM(is_timeout) * 1.0 / COUNT(matching_attempt_key)` |

## Payment & Wallet Dashboard

Tujuan: memantau kesehatan transaksi finansial.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_payment`, `fact_wallet_transaction` |
| KPI | Total payment success, total payment failed, payment success rate, total refund, wallet topup, wallet settlement, payment processing time |
| Chart | Payment success rate by method, total paid amount per hari, failed amount trend, settlement trend, table audit amount mismatch |
| Drill-down | Payment method, service group, status group, date, city |

Metric utama:

| Metric | Formula |
|---|---|
| Paid Amount | `SUM(fact_payment.paid_amount)` |
| Failed Amount | `SUM(fact_payment.failed_amount)` |
| Refund Amount | `SUM(fact_payment.refund_amount)` |
| Payment Success Rate | `SUM(is_success) * 1.0 / COUNT(payment_id)` |
| Wallet Topup | `SUM(credit_amount)` where `transaction_type = 'topup'` |
| Driver Settlement | `SUM(credit_amount)` where `transaction_type = 'settlement'` and `driver_key` is not null |

## Fraud & Geofencing Dashboard

Tujuan: memantau event risiko seperti COD lock, geo mismatch, device risk, dan potensi fake order.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `fact_geofence_event`, `fact_food_order`, `fact_payment` |
| KPI | Total fraud/geofence event, COD locked order, OTP required transaction, blocked transaction, average risk score |
| Chart | Top fraud/geofence rule, heatmap risk event per zone, blocked transaction trend, table high-risk event |
| Drill-down | Rule code, service group, city/zone, merchant, user hash, driver hash, status |

Metric utama:

| Metric | Formula |
|---|---|
| Fraud Event Count | `SUM(fact_geofence_event.event_count)` |
| Blocked Transaction | `SUM(is_blocked)` |
| OTP Required | `SUM(is_otp_required)` |
| Average Risk Score | `AVG(risk_score)` |
| Average Distance Violation | `AVG(distance_violation_km)` |

## Customer Growth & Retention Dashboard

Tujuan: memahami pertumbuhan, aktivitas, dan retensi user.

| Area | Kebutuhan |
|---|---|
| Fakta utama | `dim_user`, `fact_trip`, `fact_food_order`, `fact_shipment`, optional login mart |
| KPI | New registered user, DAU, WAU, MAU, repeat order rate, multi-service user, average order per user |
| Chart | New user growth, cohort retention table, multi-service behavior, average order per user trend |
| Drill-down | Registration month, city, service group, user type, verification status |

Metric utama:

| Metric | Formula |
|---|---|
| New Registered User | `COUNT(dim_user.user_id)` by `registered_date_key` current first version |
| Active User | Distinct user with completed trip/order/shipment in period |
| Repeat Order Rate | Users with more than 1 completed transaction / active users |
| Multi-Service User | Users with completed transaction in more than 1 service group |
| Average Order per User | Completed transactions / active users |

## Recommended Mart Tables

| Mart | Dashboard | Grain | Metric Utama |
|---|---|---|---|
| `mart_daily_revenue` | Executive | Date + service group | GMV, revenue, transaction count |
| `mart_city_operations` | FELO-City | Date + pickup zone + service type | Trip count, completion rate, matching time |
| `mart_food_merchant_performance` | FELO-Food | Date + merchant | GMV, order count, prep time, rating |
| `mart_send_operations` | FELO-Send | Date + pickup/dropoff zone | Shipment count, delivery duration, POD rate |
| `mart_driver_performance` | Driver Performance | Date + driver | Completed job, earning, rating, acceptance rate |
| `mart_payment_health` | Payment & Wallet | Date + payment method | Paid amount, failed amount, success rate, refund |
| `mart_geofence_risk` | Fraud & Geofencing | Date + zone + rule | Event count, blocked count, risk score |
| `mart_customer_retention` | Customer Growth | Cohort month + activity month | Registered user, retained user, retention rate |

## Acceptance Rules for Reports

- Dashboard metric harus bisa ditelusuri ke fact table dan formula.
- Semua chart finansial harus memakai fact yang sudah melewati rekonsiliasi invoice/payment.
- Semua status harus join ke `dim_status`, bukan hardcode string dari source service.
- Semua lokasi customer/merchant/driver ditampilkan sebagai zona, bukan alamat lengkap.
- Semua user dan driver identifier di dashboard memakai hash atau surrogate key, bukan PII.
- Semua dashboard wajib memiliki filter periode tanggal.
