# feedback-service ERD

## Ownership

This ERD belongs only to `feedback-service`.
Ride, order, customer, and driver IDs are logical references only.

## Entities

### feedback_entries

| Column | Type | Notes |
|---|---|---|
| `feedback_id` | text | PK |
| `subject_ref` | text | logical ride/order reference |
| `author_id` | text | logical customer reference |
| `target_id` | text | logical driver or merchant reference |
| `rating` | integer | star rating |
| `comment` | text | optional review |
