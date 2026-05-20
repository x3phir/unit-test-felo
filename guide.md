# Testing Guide

##Jenis Test

### 1. Unit Test
- Lokasi: `services/{service}/tests/unit/`
- Menggunakan mock (go.uber.org/mock)
- Tidak mengakses database atau web service lain
- Build tag: default (tanpa tag)

### 2. Functional Test
- Lokasi: `services/{service}/tests/functional/`
- Dapat mengakses database PostgreSQL
- Menggunakan build tag: `//go:build functional`

### 3. E2E Test
- Lokasi: `tests/e2e/`
- Cross-service integration test
- Menggunakan build tag: `//go:build e2e`

## Menjalankan Test

### Unit Test

```powershell
# Semua unit test
go test ./...

# Per service
go test ./services/ride-service/...
go test ./services/wallet-service/...

# Specific test
go test -v ./services/... -run TestTripService_Create

# Dengan coverage
go test -covermode=atomic -coverprofile='coverage.out' ./services/...
go tool cover -func coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Functional Test

```powershell
# Mulai infrastruktur dulu
docker compose -f docker-compose.functional.yml up -d

# Run functional test
go test -tags=functional ./services/ride-service/...

# Run specific functional test
go test -v -tags=functional ./services/... -run TestRideFunctional

# Stop infrastructure
docker compose -f docker-compose.functional.yml down
```

### E2E Test

```powershell
# Seed database
$env:FELO_AUTH_JWT='demo-functional-token'
go run .\cmd\felo-seed

# Start services (terminal lain)
$env:FELO_AUTH_JWT='demo-functional-token'
go run .\cmd\felo-demo

# Run E2E
$env:FELO_E2E_ENABLED='1'
$env:FELO_E2E_SUITE='smoke'
$env:FELO_AUTH_JWT='demo-functional-token'
go test -tags=e2e ./tests/e2e/...
```

## Pipeline (Jenkins)

1. **Checkout** - Clone repo
2. **Unit Tests** - `go test ./...`
3. **Vet/Lint** - `go vet ./...`
4. **Build Image** - Build Docker image
5. **Functional Tests** - `go test -tags=functional ./services/...`
6. **Push Image** - Push ke registry
7. **Deploy** - Deploy ke Kubernetes
8. **Verify** - Verify deployment

## Coverage Threshold

- Overall: 70%
- Business logic: 80%

```powershell
go run ./tools/coveragecheck -file coverage.out -threshold 70
```

## Environment Variables

### Functional Tests
```
FELO_RIDE_PG_DSN=postgres://felo:felo@127.0.0.1:54321/ride_db
FELO_WALLET_PG_DSN=postgres://felo:felo@127.0.0.1:54323/wallet_db
```

### E2E Tests
```
FELO_AUTH_JWT=<token>
FELO_RIDE_GRPC_ADDR=127.0.0.1:50051
FELO_MATCHING_GRPC_ADDR=127.0.0.1:50052
FELO_WALLET_GRPC_ADDR=127.0.0.1:50053
FELO_PAYMENT_GRPC_ADDR=127.0.0.1:50054
FELO_LOCATION_GRPC_ADDR=127.0.0.1:50055
```

## Troubleshooting

### Race Detection
```powershell
# Windows/386 tidak support -race
go test ./...

# Linux/Darwin amd64
go test -race ./...
```

### Database Connection
```powershell
# Cek PostgreSQL
docker compose -f docker-compose.functional.yml ps
pg_isready -h 127.0.0.1 -p 54321 -U felo
```