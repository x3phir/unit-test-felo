# ride-service

Owns the trip lifecycle and FELO-Now QR flow.

## Build

```powershell
go test ./services/ride-service/...
```

## Test Surfaces

- unit tests in `internal/service`
- isolated functional tests in `tests/functional`
- cross-service E2E coverage in `tests/e2e`

## Docs

- `docs/erd.md`
