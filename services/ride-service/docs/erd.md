# ride-service ERD

## Ownership

This ERD belongs only to `ride-service`.
It does not define foreign keys to other services.
References to customer or driver are logical references only.

## Entities

### rides

| Column | Type | Notes |
|---|---|---|
| `id` | text | PK |
| `customer_id` | text | logical reference to identity domain |
| `driver_id` | text | logical reference to driver domain |
| `pickup_lat` | double | pickup latitude |
| `pickup_lng` | double | pickup longitude |
| `dest_lat` | double | destination latitude |
| `dest_lng` | double | destination longitude |
| `fare` | bigint | estimated or agreed fare |
| `state` | text | trip state machine |
| `qr_code` | text | FELO-Now QR token |
| `qr_expires_at` | timestamptz | QR expiry |
| `qr_locked_driver` | text | logical driver lock |
| `created_at` | timestamptz | creation time |
| `updated_at` | timestamptz | update time |
