package orderbook

type Order struct {
	Amount uint64    `json:"amount"`
	Price  uint64    `json:"price"`
	ID     string    `json:"id"`
	Side   OrderType `json:"side"`
}
