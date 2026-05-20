# user-service ERD

## Ownership

This ERD belongs only to `user-service`.
Other services may reference `user_id` logically but never through foreign keys.

## Entities

### users

| Column | Type | Notes |
|---|---|---|
| `user_id` | text | PK |
| `phone` | text | login phone |
| `name` | text | display name |
| `status` | text | active, blocked, deleted |
| `created_at` | timestamptz | creation time |
