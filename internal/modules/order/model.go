package order

type Order struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Item   string `json:"item"`
	Amount int64  `json:"amount"`
}

type CreateOrderRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Item   string `json:"item" binding:"required"`
	Amount int64  `json:"amount" binding:"required"`
}
