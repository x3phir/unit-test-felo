# payment-service ERD

## Ownership

This ERD belongs only to `payment-service`.
Ride and customer identifiers are logical references, not foreign keys.

## Entities

### payments

| Column | Type | Notes |
|---|---|---|
| `event_id` | text | PK, source event reference |
| `ride_id` | text | logical reference to ride-service |
| `customer_id` | text | logical reference to identity domain |
| `amount` | bigint | payment amount |
| `currency` | text | currency code |
| `status` | text | completed or failed |
| `created_at` | timestamptz | processing timestamp |
