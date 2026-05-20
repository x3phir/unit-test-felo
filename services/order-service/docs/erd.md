# order-service ERD

## Ownership

This ERD belongs only to `order-service`.
Customer, merchant, and shipment references are logical references only.

## Entities

### orders

| Column | Type | Notes |
|---|---|---|
| `order_id` | text | PK |
| `customer_id` | text | logical reference to user-service |
| `merchant_id` | text | logical reference to merchant-service |
| `shipment_ref` | text | logical reference to shipment-service |
| `status` | text | requested through completed |
| `created_at` | timestamptz | creation time |
