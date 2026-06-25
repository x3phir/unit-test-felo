# FELO Backend Monorepo

Tugas besar **Cloud Computing** Kelompok 24. Repo ini berisi unit test dan functional test untuk 18 microservice FELO.

**Stack:** Go 1.25, PostgreSQL, Redis, RabbitMQ, gRPC JSON codec.

## Cara Menjalankan Test

### 1. Unit Test

Unit test tidak mengakses database atau service eksternal

```bash
go test ./services/...
```

### 2. Functional Test

Functional test mengakses PostgreSQL asli. Jalankan database test terlebih dahulu:

```bash
docker compose -f docker-compose.functional.yml up -d
go test -tags=functional ./services/...
docker compose -f docker-compose.functional.yml down --remove-orphans
```

Jika command functional test dijalankan tanpa database, test akan gagal karena koneksi PostgreSQL tidak tersedia.

### 3. Test Per Service

```bash
go test ./services/auth-service/...
go test -tags=functional ./services/auth-service/...
```

Ganti `auth-service` dengan service lain sesuai kebutuhan.

## Jenkins CI/CD

Pipeline Jenkins menjalankan:

```text
checkout -> vet -> build image -> infrastructure -> per-service tests -> push image
```

Stage `Infrastructure` otomatis menjalankan `docker-compose.functional.yml`, menghubungkan container Jenkins ke network Compose, lalu membersihkan container setelah test selesai.

Setelah `Infrastructure`, Jenkins menampilkan stage terpisah untuk setiap service, misalnya `Auth Service`, `User Service`, `Ride Service`, dan seterusnya. Setiap service stage menjalankan unit test dan functional test untuk service tersebut.

Parameter Jenkins:

| Parameter | Default | Keterangan |
|---|---|---|
| `DOCKERHUB_NAMESPACE` | `piipapoy` | Username/organization Docker Hub tujuan push image |
| `DOCKER_IMAGE_NAME` | `felo-backend` | Nama repository image Docker |
| `DOCKER_CREDENTIAL_ID` | `dockerhub-login` | ID credential Jenkins untuk Docker Hub |
| `RUN_PUSH_IMAGE` | `true` | Push image setelah test selesai |
| `RUN_DEPLOY` | `false` | Jalankan deploy Kubernetes opsional |

Untuk penilaian test saja, `RUN_DEPLOY` dapat dibiarkan `false`.

## Struktur Test

```text
services/{service}/
├── internal/                 # domain, port, service logic
├── tests/unit/               # unit test, tanpa DB
├── tests/functional/         # functional test, dengan PostgreSQL
└── docs/erd.md               # ERD service
```

## Service dan Database

Repo ini memiliki 18 service dan 18 PostgreSQL instance untuk functional test.

| Service | Port DB |
|---|---:|
| `ride-service` | 54321 |
| `matching-service` | 54322 |
| `location-service` | 54325 |
| `invoice-service` | 54326 |
| `merchant-service` | 54327 |
| `notification-service` | 54328 |
| `auth-service` | 54329 |
| `user-service` | 54330 |
| `driver-service` | 54331 |
| `feedback-service` | 54332 |
| `order-service` | 54333 |
| `cart-service` | 54334 |
| `send-order-service` | 54335 |
| `shipment-service` | 54336 |
| `pricing-service` | 54337 |
| `payment-service` | 54338 |
| `wallet-service` | 54339 |
| `tracking-service` | 54340 |

ERD tiap service tersedia di `services/{service}/docs/erd.md`.
