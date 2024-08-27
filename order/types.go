package order

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"order_id"`
	BookID   int `json:"book_id"`
	Quantity int `json:"quantity"`
}

type CreateOrderParams struct {
	UserID int              `json:"user_id"`
	Items  []OrderItemInput `json:"items"`
}

type OrderItemInput struct {
	BookID   int `json:"book_id"`
	Quantity int `json:"quantity"`
}

type CreateOrderResponse struct {
	OrderID int `json:"order_id"`
}

type GetOrderHistoryResponse struct {
	Orders []OrderWithItems `json:"orders"`
}

type OrderWithItems struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}
