# tracking-service ERD

## Ownership

This ERD belongs only to `tracking-service`.
It does not define foreign keys to other services.
References to shipment or driver are logical references only.

## Entities

### tracking_sessions

| Column | Type | Notes |
|---|---|---|
| `id` | text | PK |
| `shipment_id` | text | logical reference to shipment domain |
| `driver_id` | text | logical reference to driver domain |
| `status` | text | active / paused / completed |
| `started_at` | timestamptz | session start time |
| `updated_at` | timestamptz | last update time |
| `ended_at` | timestamptz | session end time (nullable) |

### tracking_records

| Column | Type | Notes |
|---|---|---|
| `id` | text | PK |
| `session_id` | text | FK to tracking_sessions |
| `lat` | double | latitude |
| `lng` | double | longitude |
| `speed` | double | speed in km/h |
| `heading` | double | heading in degrees |
| `recorded_at` | timestamptz | when the position was recorded |
