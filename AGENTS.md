# AGENTS.md

## Project Overview
Go microservices monorepo for FELO (ride-hailing, food delivery, shipping platform). Currently scaffolded for 5 services: `ride-service`, `matching-service`, `wallet-service`, `payment-service`, `location-service`.

## Key Commands

### Unit Testing
```powershell
go test ./...                                    # All tests
go test ./services/ride-service/...             # Single service
go test -v ./services/... -run TestName         # Specific test
go test -race ./...                             # Race detection (linux/amd64, darwin only)
```

### Coverage
```powershell
go test -covermode=atomic -coverprofile='coverage.out' ./services/...
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```
- Threshold: 70% overall, 80% for business logic
- Race detection disabled on windows/386 - CI skips `-race` automatically

### Functional/E2E Tests
Requires Docker desktop with `docker-compose.functional.yml`:
```powershell
docker compose -f .\docker-compose.functional.yml up -d  # Start infra
$env:FELO_AUTH_JWT='demo-functional-token'
go run .\cmd\felo-seed                                    # Seed databases
go run .\cmd\felo-demo                                    # Start gRPC services (separate terminal)
$env:FELO_E2E_ENABLED='1'
$env:FELO_E2E_SUITE='smoke'  # smoke|critical-flow|full-regression
go test -tags=e2e ./tests/e2e/...
```

### Local CI Pipeline (Jenkins-equivalent)
```powershell
go test -json -covermode=atomic -coverprofile='coverage.out' ./services/... | Tee-Object -FilePath 'gotest.json'
go test ./tools/...
go tool cover -html='coverage.out' -o 'coverage.html'
Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```

## Architecture

- **services/**: Each microservice owns its code and tests (internal/, tests/unit/)
- **tests/e2e/**: Cross-service E2E scenarios using `e2e` build tag
- **tools/**: Internal CI tools (coveragecheck, gotest2junit)
- **cmd/**: Demo entrypoints (felo-seed, felo-demo)
- **Services use RabbitMQ for async communication** (see prd.md)
- **No database sharing between services** (each has its own PostgreSQL)

## Required Environment Variables

### Functional tests
```
FELO_AUTH_JWT=<token>
FELO_RIDE_GRPC_ADDR=127.0.0.1:50051
FELO_MATCHING_GRPC_ADDR=127.0.0.1:50052
FELO_WALLET_GRPC_ADDR=127.0.0.1:50053
FELO_PAYMENT_GRPC_ADDR=127.0.0.1:50054
FELO_LOCATION_GRPC_ADDR=127.0.0.1:50055
```

## Documentation References

- `docs/run-testing-guide.md` - Detailed unit testing commands
- `docs/run-functional-testing-guide.md` - E2E test setup
- `docs/global/architecture.md` - System architecture
- `prd.md` - Product requirements (18 microservices planned)