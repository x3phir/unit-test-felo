# auth-service ERD

## Ownership

This ERD belongs only to `auth-service`.
Any user reference is a logical reference only and must not be modeled as a cross-service foreign key.

## Entities

### auth_sessions

| Column | Type | Notes |
|---|---|---|
| `session_id` | text | PK |
| `user_id` | text | logical reference to user-service |
| `access_token` | text | issued token |
| `refresh_token` | text | refresh token |
| `expires_at` | timestamptz | token expiry |
