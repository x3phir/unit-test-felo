package domain

type CartItem struct {
	MenuItemID string
	Quantity   int
	Notes      string
	Price      int64
}

type Cart struct {
	UserID     string
	MerchantID string
	Items      []CartItem
	TotalPrice int64
}

type CartSummary struct {
	UserID     string
	MerchantID string
	ItemCount  int
	TotalPrice int64
}
