# merchant-service ERD

## Ownership

This ERD belongs only to `merchant-service`.
Owner and catalog references are logical references only.

## Entities

### merchants

| Column | Type | Notes |
|---|---|---|
| `merchant_id` | text | PK |
| `owner_user_id` | text | logical reference to user-service |
| `name` | text | merchant name |
| `status` | text | active or inactive |
| `updated_at` | timestamptz | last change |
