# FELO Functional Testing Strategy

## Purpose

Functional tests validate multi-service behavior across real infrastructure:

- `gRPC` communication
- `PostgreSQL` per service
- `Redis`
- `RabbitMQ`

The initial rollout covers only the five priority services:

1. `ride-service`
2. `matching-service`
3. `wallet-service`
4. `payment-service`
5. `location-service`

## Scope

Phase 1 functional coverage focuses on:

1. Regular ride flow
2. FELO-Now QR flow
3. Critical negative flows
4. RabbitMQ event verification
5. Seeded data with deterministic scenarios

Deferred from this phase:

- order fiktif mitigation
- notification fallback
- feedback and rating

## Test Categories

Three categories are used in CI and local execution:

1. `smoke`
   Checks service reachability, gRPC connectivity, and basic health.
2. `critical-flow`
   Runs the core business journeys that must pass before merge.
3. `full-regression`
   Runs every functional flow, including negative-path and async checks.

## Foundational Rules

- Tests run serially to reduce flakiness in async flows.
- Every test case uses seeded data with explicit IDs.
- RabbitMQ assertions are first-class pass criteria.
- Time-bounded eventual assertions are required for async workflows.
- E2E tests are separate from unit tests and must use the `e2e` build tag.

## Core Scenarios

### Regular Ride Critical Flow

1. customer requests ride
2. ride enters matching
3. driver is matched
4. driver location is reported
5. ride starts
6. ride completes
7. payment completes
8. wallet settlement completes
9. location history is queryable
10. required events are observed

### FELO-Now Critical Flow

1. customer generates QR
2. driver scans QR
3. driver accepts ride
4. ride starts from QR flow
5. completion triggers payment and settlement
6. required events are observed

### Negative Flows

1. no driver available after retries
2. expired QR cannot be scanned
3. payment failure publishes `payment.failed.v1`
4. invalid JWT request is rejected
5. duplicate settlement event does not double-credit wallet

## Seed Data

Functional tests depend on stable seeded fixtures:

- active customer with sufficient wallet balance
- active customer with insufficient wallet balance
- three available drivers
- one offline driver
- one driver with existing location history
- fixed trip, QR, and event IDs for deterministic assertions

## Execution Model

E2E tests run from the root folder `tests/e2e` and target real local service endpoints. Service-local functional coverage stays under each service directory.

The repository provides:

- environment configuration
- scenario definitions
- readiness hooks
- eventual assertion helpers
- real gRPC demo runtime for local validation

## CI Strategy

- `smoke` runs on every branch build
- `critical-flow` runs on merge validation
- `full-regression` runs nightly or before release

Functional test results should be exported as:

- `JUnit XML`
- structured logs
- optional event capture artifacts
