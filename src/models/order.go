package models

type Order struct {
	ID            int         `json:"id"`
	TableID       int         `json:"table_id"`
	UserID        int         `json:"user_id"`
	Status        string      `json:"status"`
	Items         []OrderItem `json:"items"`
	Notes         string      `json:"notes"`
	PaymentMethod string      `json:"payment_method"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	TotalAmount   float64     `json:"total_amount"`
	Discount      float64     `json:"discount"`
	Tax           float64     `json:"tax"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	MenuID    int     `json:"menu_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

var Orders []Order

const (
	OrderStatusPending  = "pending"
	OrderStatusCooking  = "cooking"
	OrderStatusReady    = "ready"
	OrderStatusServed   = "served"
	OrderStatusPaid     = "paid"
	OrderStatusCanceled = "canceled"
)

const (
	PaymentMethodCash   = "cash"
	PaymentMethodDebit  = "debit"
	PaymentMethodCredit = "credit"
)
