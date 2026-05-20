# FELO E2E Tests

Root-level end-to-end tests for flows that span multiple services.

These tests intentionally live outside any single service because they validate:

- gRPC calls across service boundaries
- RabbitMQ event propagation
- multi-database side effects
- end-to-end business journeys

Use the `e2e` build tag when running them.
