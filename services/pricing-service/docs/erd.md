# pricing-service ERD

## Ownership

This ERD belongs only to `pricing-service`.
References to ride or order subjects are logical references only.

## Entities

### pricing_rules

| Column | Type | Notes |
|---|---|---|
| `rule_id` | text | PK |
| `service_type` | text | ride, send, food |
| `base_fare` | bigint | base amount |
| `surge_multiplier` | numeric | surge setting |
| `active_from` | timestamptz | rule start |
| `active_to` | timestamptz | rule end |

### fare_audit

| Column | Type | Notes |
|---|---|---|
| `trip_id` | text | ref ke ride/trip |
| `distance_km` | numeric | actual distance |
| `duration_mins` | numeric | actual duration |
| `demand_level` | int | demand at calculation time |
| `supply_level` | int | supply at calculation time |
| `base_fare` | bigint | computed base amount |
| `surge_multiplier` | numeric | applied surge multiplier |
| `final_fare` | bigint | final billed amount |
| `currency` | text | currency code |
| `calculated_at` | timestamptz | calculation timestamp |
