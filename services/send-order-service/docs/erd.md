# send-order-service ERD

## Ownership

This ERD belongs only to `send-order-service`.
User and shipment identifiers are logical references only.

## Entities

### send_orders

| Column | Type | Notes |
|---|---|---|
| `send_order_id` | text | PK |
| `sender_id` | text | logical user reference |
| `receiver_ref` | text | logical receiver reference |
| `shipment_ref` | text | logical shipment reference |
| `status` | text | draft through delivered |
| `created_at` | timestamptz | creation time |
