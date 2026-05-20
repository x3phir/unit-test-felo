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
