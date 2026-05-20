# location-service ERD

## Ownership

This ERD belongs only to `location-service`.
Driver references are logical references, not foreign keys.

## Entities

### driver_locations

| Column | Type | Notes |
|---|---|---|
| `id` | bigint | PK |
| `driver_id` | text | logical reference to driver domain |
| `lat` | double | latitude |
| `lng` | double | longitude |
| `recorded_at` | timestamptz | sample time |
