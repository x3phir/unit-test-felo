# cart-service ERD

## Ownership

This ERD belongs only to `cart-service`.
Customer and merchant identifiers are logical references only.

## Entities

### carts

| Column | Type | Notes |
|---|---|---|
| `cart_id` | text | PK |
| `customer_id` | text | logical reference to user-service |
| `merchant_id` | text | logical reference to merchant-service |
| `status` | text | active or checked_out |
| `updated_at` | timestamptz | last change |

### cart_items

| Column | Type | Notes |
|---|---|---|
| `item_id` | text | PK |
| `cart_id` | text | service-local relation |
| `product_ref` | text | logical product reference |
| `quantity` | integer | quantity in cart |
