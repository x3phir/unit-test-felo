# driver-service ERD

## Ownership

This ERD belongs only to `driver-service`.
Identifiers for users or vehicles are logical references only.

## Entities

### drivers

| Column | Type | Notes |
|---|---|---|
| `driver_id` | text | PK |
| `user_id` | text | logical reference to user-service |
| `status` | text | active, suspended, offline |
| `availability` | text | online or offline |
| `updated_at` | timestamptz | last change |
