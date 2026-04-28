# FELO Functional Test Guide

## Overview

Functional tests are intended for end-to-end and cross-service validation using:

- real `PostgreSQL`
- real `Redis`
- real `RabbitMQ`
- real `gRPC` service endpoints

These tests are separate from unit tests and use the `functional` build tag.

## Current Status

The infrastructure and test harness are scaffolded now.
The actual FELO service implementations do not exist yet in this repository, so the functional tests currently skip until real gRPC adapters are connected.

That is expected.

## Working Directory

```powershell
cd 'C:\Users\Harri Supriadi\Documents\unit-test-felo'
```

## 1. Start Infrastructure

```powershell
docker compose -f .\docker-compose.functional.yml up -d
```

This starts:

- `postgres-ride`
- `postgres-matching`
- `postgres-wallet`
- `postgres-payment`
- `postgres-location`
- `redis`
- `rabbitmq`

## 1A. Seed Database

Set the shared auth token first:

```powershell
$env:FELO_AUTH_JWT='demo-functional-token'
```

Then seed all FELO demo databases:

```powershell
go run .\cmd\felo-seed
```

Seeder source:

- `functional/testdata/seeds/customers.json`
- `functional/testdata/seeds/drivers.json`
- `functional/testdata/seeds/locations.json`

## 1B. Start Real FELO gRPC Demo Services

Run this in a separate terminal:

```powershell
$env:FELO_AUTH_JWT='demo-functional-token'
go run .\cmd\felo-demo
```

This starts real gRPC endpoints for:

- ride-service on `127.0.0.1:50051`
- matching-service on `127.0.0.1:50052`
- wallet-service on `127.0.0.1:50053`
- payment-service on `127.0.0.1:50054`
- location-service on `127.0.0.1:50055`

## 2. Verify Containers

```powershell
docker compose -f .\docker-compose.functional.yml ps
```

## 3. Enable Functional Tests

Functional tests are guarded by environment variables so they do not run accidentally.

```powershell
$env:FELO_FUNCTIONAL_ENABLED='1'
$env:FELO_TEST_SUITE='smoke'
$env:FELO_AUTH_JWT='demo-functional-token'
```

Available suite values:

- `smoke`
- `critical-flow`
- `full-regression`

## 4. Run Functional Tests

```powershell
go test -tags=functional ./functional/...
```

Verbose mode:

```powershell
go test -v -tags=functional ./functional/...
```

## 5. Run One Suite Category

Smoke:

```powershell
$env:FELO_TEST_SUITE='smoke'
go test -v -tags=functional ./functional/...
```

Critical flow:

```powershell
$env:FELO_TEST_SUITE='critical-flow'
go test -v -tags=functional ./functional/...
```

Full regression:

```powershell
$env:FELO_TEST_SUITE='full-regression'
go test -v -tags=functional ./functional/...
```

## 6. Run with JUnit Output

```powershell
go test -json -tags=functional ./functional/... | Tee-Object -FilePath 'functional-gotest.json'
Get-Content -LiteralPath 'functional-gotest.json' | go run ./tools/gotest2junit | Set-Content 'functional-junit.xml'
```

## 7. Stop Infrastructure

```powershell
docker compose -f .\docker-compose.functional.yml down
```

To remove volumes too:

```powershell
docker compose -f .\docker-compose.functional.yml down -v
```

## 8. Required Environment Variables for Real Adapter Wiring

These are the endpoint variables the harness expects:

```text
FELO_RIDE_GRPC_ADDR
FELO_MATCHING_GRPC_ADDR
FELO_WALLET_GRPC_ADDR
FELO_PAYMENT_GRPC_ADDR
FELO_LOCATION_GRPC_ADDR
FELO_AUTH_JWT
```

Example:

```powershell
$env:FELO_RIDE_GRPC_ADDR='127.0.0.1:50051'
$env:FELO_MATCHING_GRPC_ADDR='127.0.0.1:50052'
$env:FELO_WALLET_GRPC_ADDR='127.0.0.1:50053'
$env:FELO_PAYMENT_GRPC_ADDR='127.0.0.1:50054'
$env:FELO_LOCATION_GRPC_ADDR='127.0.0.1:50055'
$env:FELO_AUTH_JWT='replace-with-valid-test-jwt'
```

## 9. Seed Data

Fixtures live under `functional/testdata/`.
They define deterministic IDs for:

- customers
- drivers
- wallets
- locations
- QR and trip flows

## 10. Expected Behavior Right Now

At the current runnable demo stage:

1. functional tests compile
2. infrastructure can be started with Docker
3. databases can be seeded with `go run .\cmd\felo-seed`
4. gRPC demo services can be started with `go run .\cmd\felo-demo`
5. functional tests run against real local endpoints
