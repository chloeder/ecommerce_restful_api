package models

type Checkout struct {
	Email    string            `json:"email"`
	Address  string            `json:"address"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
}

type Order struct {
	ID                string  `json:"id"`
	Email             string  `json:"email"`
	GrandTotal        int64   `json:"grand_total"`
	Passcode          *string `json:"passcode,omitempty"`
	PaidAt            *string `json:"paid_at,omitempty"`
	PaidBank          *string `json:"paid_bank,omitempty"`
	PaidAccountNumber *string `json:"paid_account_number,omitempty"`
}

type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}
