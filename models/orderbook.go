package models

import (
	"time"
)

type Item struct {
	Price  float64
	Amount float64
}

type OrderBook struct {
	Time time.Time
	Asks []Item
	Bids []Item
}

// Ask 卖一
func (o *OrderBook) Ask() (result Item) {
	if len(o.Asks) > 0 {
		result = o.Asks[0]
	}
	return
}

// Bid 买一
func (o *OrderBook) Bid() (result Item) {
	if len(o.Bids) > 0 {
		result = o.Bids[0]
	}
	return
}
