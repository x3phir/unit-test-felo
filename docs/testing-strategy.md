# FELO Unit Test Strategy

## Scope

This repo uses a test-first approach for FELO microservices implemented in Go.
Unit tests target business logic first and mock all external dependencies:

- PostgreSQL repositories
- Redis cache and lock behavior
- RabbitMQ publishers and consumers
- gRPC clients between services
- JWT validation helpers

Integration tests are intentionally separated from unit tests and should run in a different Jenkins stage later.

## Priorities

Phase 1 covers the highest-risk services from the PRD:

1. `ride-service`
2. `matching-service`
3. `wallet-service`
4. `payment-service`
5. `location-service`

Phase 2 extends the same structure to:

1. `tracking-service`
2. `user-service`
3. `driver-service`
4. `auth-service`
5. `invoice-service`
6. `feedback-service`
7. `notification-service`

## Conventions

- Use table-driven tests for behavior matrices.
- Name tests with the format `Test<Subject>_<Condition>_<ExpectedBehavior>`.
- Create explicit test cases for:
  - happy path
  - validation failure
  - invalid state transition
  - dependency failure
  - idempotency
  - context timeout or cancellation
  - concurrency or race-sensitive behavior
- Keep business rules in service packages and do not hide them inside handlers.
- Keep mocks narrow and owned by the test package unless shared behavior justifies extraction.

## Coverage Rules

- Core service logic must target `>= 80%`.
- Overall package coverage must not drop below `70%`.
- Branch-heavy flows must include at least one invalid-path assertion.
- State machines must be tested for legal and illegal transitions.

## CI Rules

- Run `go test -json -covermode=atomic -coverprofile=coverage.out ./services/...`.
- Run `go test ./tools/...`.
- Generate `coverage.out` and `coverage.html`.
- Convert `go test -json` output into `junit.xml`.
- Fail the build if the configured threshold is not reached.
- Enable `-race` on Jenkins agents that support it. Skip it only on unsupported combinations such as `windows/386`.

## Package Shape

Each service follows the same broad layout even if its internal details differ:

```text
services/<service-name>/
  internal/domain/
  internal/ports/
  internal/service/
```

This keeps service ownership clear while allowing variations in implementation.
