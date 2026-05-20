# FELO Global Data Contracts

This folder contains only cross-service contracts and architecture-level documentation.

## Rules

- service ERDs stay under each service directory
- event payloads may reference foreign IDs, but never as database foreign keys
- global contracts are transport- and integration-oriented

## Current Event Contracts

- `ride.created.v1`
- `driver.matched.v1`
- `ride.completed.v1`
- `payment.completed.v1`
- `payment.failed.v1`
