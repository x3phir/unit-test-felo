# FELO Backend Monorepo

Test-first backend scaffold for FELO microservices in Go.

## Why This Repo Exists

This repository is designed to let FELO backend services be built from tests first.
The initial scaffold focuses on the five highest-risk services from the PRD:

- `ride-service`
- `matching-service`
- `wallet-service`
- `payment-service`
- `pricing-service`
- `location-service`

The remaining services are tracked in [docs/service-roadmap.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\docs\service-roadmap.md).
How to run the tests is documented in [docs/run-testing-guide.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\docs\run-testing-guide.md).
Functional test design and execution are documented in [docs/functional-testing-strategy.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\docs\functional-testing-strategy.md) and [docs/run-functional-testing-guide.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\docs\run-functional-testing-guide.md).

## Demo Runtime

To try the FELO functional flow locally with real gRPC endpoints and seeded databases:

```powershell
docker compose -f .\docker-compose.functional.yml up -d
$env:FELO_AUTH_JWT='demo-functional-token'
go run .\cmd\felo-seed
go run .\cmd\felo-demo
```

Then, in another terminal:

```powershell
$env:FELO_FUNCTIONAL_ENABLED='1'
$env:FELO_TEST_SUITE='critical-flow'
$env:FELO_AUTH_JWT='demo-functional-token'
go test -v -tags=functional ./functional/...
```

## Principles

- Keep business rules in service packages and unit-test them first.
- Mock all external dependencies for unit tests.
- Split integration tests from unit tests.
- Keep CI deterministic with only standard Go tooling at the core.

## Layout

```text
docs/                       Strategy, naming, rollout notes
services/                   Service-owned code and tests
tools/                      CI helper tools written in Go
Jenkinsfile                 Jenkins pipeline with coverage and JUnit output
go.mod                      Root Go module
```

## Commands

```powershell
go test -json -covermode=atomic -coverprofile='coverage.out' ./services/... | Tee-Object -FilePath 'gotest.json'
go test ./tools/...
go tool cover -func 'coverage.out'
go tool cover -html='coverage.out' -o 'coverage.html'
Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```

## Coverage Defaults

- Overall service package target: `>= 70%`
- Core business logic target: `>= 80%`
- Mandatory focus areas:
  - state transitions
  - idempotency
  - timeout and context cancellation
  - concurrency and race-safety for shared state

`-race` should be enabled in CI whenever the Jenkins agent supports it. The current local environment uses `windows/386`, which does not support Go race detection.
