# FELO Backend Monorepo

Proyek tugas besar **Cloud Computing** Kelompok 24 — platform transportasi online (ride-hailing, food delivery, package delivery).

**Stack:** Go 1.25, gRPC + JSON codec, PostgreSQL, Redis, RabbitMQ

---

## Team & Service Ownership

| Anggota | Services | Port Range |
|---|---|---|
| **Anas Miftakhul Falah** | auth, user, driver, feedback | 54329–54332 |
| **Harri Supriadi** | ride, matching, location, tracking | 54321, 54322, 54325, 54340 |
| **M. Raffa Mizanul Insan** | order (food), cart, send-order, shipment | 54333–54336 |
| **Muhammad Adwar Salman** | pricing, payment, wallet | 54337–54339 |
| **Rafi Ahmad Al Farisi** | invoice, notification, merchant | 54326–54328 |

## Architecture

Setiap microservice mengikuti **Clean Architecture**:

```
services/{name}/
├── internal/
│   ├── domain/{name}.go        # Entity/struct definitions
│   ├── ports/interfaces.go     # Repository & client interfaces
│   └── service/{name}_service.go  # Business logic
├── tests/
│   ├── unit/{name}_test.go     # Unit test (mock-based, no DB)
│   └── functional/{name}_test.go  # Functional test (real PostgreSQL)
└── docs/
    └── erd.md                  # ERD diagram
```

### Service List

| Service | Description | PJ |
|---|---|---|
| `auth-service` | Autentikasi, OTP, manajemen sesi | Anas |
| `user-service` | Profil & alamat pengguna | Anas |
| `driver-service` | Manajemen driver & KYC | Anas |
| `feedback-service` | Rating & ulasan | Anas |
| `pricing-service` | Kalkulasi tarif & surge pricing | Adwar |
| `payment-service` | Proses pembayaran ride | Adwar |
| `wallet-service` | Dompet digital & settlement | Adwar |
| `invoice-service` | Nota digital & penanggung biaya | Rafi |
| `notification-service` | Push/WhatsApp/SMS notification | Rafi |
| `merchant-service` | Restoran & menu | Rafi |
| `matching-service` | Pencarian driver terdekat | Harri |
| `location-service` | Data lokasi & geofence | Harri |
| `ride-service` | Manajemen perjalanan (City) | Harri |
| `tracking-service` | Live tracking perjalanan & pengiriman | Harri |
| `order-service` | Pesanan makanan (Food) | Raffa |
| `cart-service` | Keranjang belanja | Raffa |
| `send-order-service` | Pesanan pengiriman (Send) | Raffa |
| `shipment-service` | Tracking pengiriman | Raffa |

## Testing

Tiga level testing:

### Unit Test
- Mock-based (go.uber.org/mock), no DB access
- Covers business logic & edge cases

```bash
go test ./...
go test -v ./services/ride-service/...
```

### Functional Test
- Real PostgreSQL via `pgxpool`
- Build tag: `//go:build functional`
- `CREATE TABLE IF NOT EXISTS` di setup
- Pola: `initXxxTables` → seed → call service → assert DB state

```bash
# Start containers
docker compose -f docker-compose.functional.yml up -d

# Run semua functional test
go test -tags=functional ./services/...

# Per service
go test -tags=functional ./services/order-service/tests/functional/...
```

### E2E Test
- Cross-service integration
- Build tag: `//go:build e2e`

```bash
go test -tags=e2e ./tests/e2e/...
```

### Coverage

```bash
go test -covermode=atomic -coverprofile=coverage.out ./services/...
go run ./tools/coveragecheck -file coverage.out -threshold 70
```

Target: **70% overall**, **80% business logic**.

## Quick Start

```bash
go test ./...                          # Unit test semua service
go vet ./...                           # Lint
docker compose -f docker-compose.functional.yml up -d  # DB untuk functional test
go test -tags=functional ./services/...  # Functional test
```

## Infrastruktur

- **PostgreSQL** — 18 instance (satu per service)
- **Redis** — Caching & geospatial
- **RabbitMQ** — Event bus (async komunikasi antar service)
- **gRPC** — JSON codec (custom encoding, no protobuf)

CI pipeline (Jenkins): `checkout → unit test → vet → build → functional test → deploy`.

Lihat [guide.md](./guide.md) untuk panduan testing detail dan `services/{name}/docs/erd.md` untuk ERD tiap service.
