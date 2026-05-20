# matching-service ERD

## Ownership

This ERD belongs only to `matching-service`.
It contains no cross-service foreign keys.
Ride and driver references are logical references.

## Entities

### drivers

| Column | Type | Notes |
|---|---|---|
| `id` | text | PK |
| `status` | text | availability state |
| `lat` | double | latest matching latitude snapshot |
| `lng` | double | latest matching longitude snapshot |

### assignments

| Column | Type | Notes |
|---|---|---|
| `ride_id` | text | PK, logical reference to ride-service |
| `driver_id` | text | logical reference to driver-service |
| `assigned_at` | timestamptz | assignment timestamp |
