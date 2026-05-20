# wallet-service ERD

## Ownership

This ERD belongs only to `wallet-service`.
Customer and driver IDs are logical references only.

## Entities

### wallets

| Column | Type | Notes |
|---|---|---|
| `owner_id` | text | PK, logical owner reference |
| `owner_type` | text | customer or driver |
| `balance` | bigint | current balance |
| `currency` | text | currency code |
| `updated_at` | timestamptz | last update |

### wallet_ledger

| Column | Type | Notes |
|---|---|---|
| `reference` | text | PK, idempotency key |
| `owner_id` | text | logical owner reference |
| `delta` | bigint | amount change |
| `reason` | text | settlement, payment, adjustment |
| `created_at` | timestamptz | entry time |
