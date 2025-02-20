package model

type StockDecrease struct {
	ProductId uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
	OrderId   string `json:"order_id"`
}
