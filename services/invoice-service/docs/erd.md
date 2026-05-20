# invoice-service ERD

## Ownership

This ERD belongs only to `invoice-service`.
Payment, order, and customer references are logical references only.

## Entities

### invoices

| Column | Type | Notes |
|---|---|---|
| `invoice_id` | text | PK |
| `subject_ref` | text | logical ride/order/payment reference |
| `customer_id` | text | logical user reference |
| `amount` | bigint | billed amount |
| `currency` | text | currency code |
| `status` | text | issued, paid, cancelled |
