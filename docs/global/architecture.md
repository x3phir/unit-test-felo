# FELO Global Architecture

## Principles

1. Services own their own schema and ERD.
2. No cross-service foreign keys are allowed.
3. Cross-service joins are replaced by logical references and contracts.
4. Service-local tests stay inside each service.
5. Multi-service E2E tests live only in `tests/e2e`.

## Current Service Boundaries

- `ride-service`
- `matching-service`
- `wallet-service`
- `payment-service`
- `location-service`

## Deployment Rule

Each service is expected to be buildable, testable, and deployable independently.
Shared runtime helpers must not weaken service ownership.
