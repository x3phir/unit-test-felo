# FELO Backend Monorepo

Go microservices monorepo untuk FELO (ride-hailing, food delivery, shipping platform).

## Services

- `auth-service` - Authentication & OTP
- `cart-service` - Shopping cart
- `driver-service` - Driver management
- `feedback-service` - Feedback & ratings
- `invoice-service` - Invoice management
- `location-service` - Location tracking
- `matching-service` - Driver-passenger matching
- `merchant-service` - Merchant management
- `notification-service` - Push notifications
- `order-service` - Food order management
- `payment-service` - Payment processing
- `pricing-service` - Pricing calculation
- `ride-service` - Ride/trip management
- `send-order-service` - Package delivery
- `shipment-service` - Shipment tracking
- `user-service` - User management
- `wallet-service` - Wallet & settlements

## Struktur Project

```
services/           # 17 microservices
├── {service-name}/
│   ├── internal/   # domain, ports, service
│   └── tests/
│       ├── unit/   # Unit test (mock, tidak akses DB)
│       └── functional/  # Functional test (akses DB)
tests/e2e/          # End-to-end cross-service tests
tools/              # CI tools (coveragecheck, gotest2junit)
cmd/                # Demo entrypoints
```

## Quick Start

```powershell
# Run semua unit test
go test ./...

# Run unit test per service
go test ./services/ride-service/...

# Run specific test
go test -v ./services/... -run TestName

# Coverage
go test -covermode=atomic -coverprofile='coverage.out' ./services/...
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```

Lihat [guide.md](guide.md) untuk panduan testing lengkap.

## Requirements

- Go 1.25+
- Docker (untuk functional/E2E tests)
- PostgreSQL, Redis, RabbitMQ (untuk functional tests)