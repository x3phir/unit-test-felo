# notification-service ERD

## Ownership

This ERD belongs only to `notification-service`.
Recipient references are logical references only.

## Entities

### notifications

| Column | Type | Notes |
|---|---|---|
| `notification_id` | text | PK |
| `recipient_ref` | text | logical user/driver/merchant reference |
| `channel` | text | push, sms, email, wa |
| `template_code` | text | message template |
| `status` | text | queued, sent, failed |
| `created_at` | timestamptz | enqueue time |
