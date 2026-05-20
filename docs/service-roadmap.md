# FELO Service Test Roadmap

| Service | Priority | Status | Notes |
|---|---|---|---|
| ride-service | P1 | Scaffolded | State machine and FELO-Now QR flow |
| matching-service | P1 | Scaffolded | Retry radius and nearest-driver matching |
| wallet-service | P1 | Scaffolded | Settlement and idempotency |
| payment-service | P1 | Scaffolded | Wallet charge and invoice orchestration |
| pricing-service | P1 | Scaffolded | Dynamic pricing, surge multiplier, fare audit |
| location-service | P1 | Scaffolded | History, latest cache, route estimation |
| tracking-service | P2 | Planned | GPS stream ownership and active trip broadcast |
| user-service | P2 | Planned | Profile read/write and validation |
| driver-service | P2 | Planned | Availability and KYC status surfaces |
| auth-service | P2 | Planned | JWT issue/validate and session edges |
| invoice-service | P2 | Planned | Invoice generation and retrieval |
| feedback-service | P2 | Planned | Rating rules and duplicate submission guard |
| notification-service | P2 | Planned | Push, WA, SMS fallback behavior |

## Functional Test Status

| Area | Status | Notes |
|---|---|---|
| Service-level functional tests | Scaffolded | Live under `services/<service>/tests/functional` |
| Docker infra | Scaffolded | PostgreSQL, Redis, RabbitMQ |
| E2E harness | Scaffolded | Root-level `tests/e2e` with build tag `e2e` |
| Smoke suite | Scaffolded | Cross-service reachability |
| Critical-flow suite | Scaffolded | Regular ride and FELO-Now |
| Full regression suite | Scaffolded | Includes negative flow placeholders |
