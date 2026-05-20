# shipment-service ERD

## Ownership

This ERD belongs only to `shipment-service`.
Order and driver references are logical references only.

## Entities

### shipments

| Column | Type | Notes |
|---|---|---|
| `shipment_id` | text | PK |
| `order_ref` | text | logical order reference |
| `driver_id` | text | logical driver reference |
| `status` | text | packed through delivered |
| `eta_minutes` | integer | estimated delivery |
| `updated_at` | timestamptz | last change |
