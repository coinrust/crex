package gotrader

import (
	"time"
)

type Item struct {
	Price  float64
	Amount float64
}

type OrderBook struct {
	Symbol string
	Time   time.Time
	Asks   []Item
	Bids   []Item
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

// AskPrice 卖一价
func (o *OrderBook) AskPrice() (result float64) {
	if len(o.Asks) > 0 {
		result = o.Asks[0].Price
	}
	return
}

// BidPrice 买一价
func (o *OrderBook) BidPrice() (result float64) {
	if len(o.Bids) > 0 {
		result = o.Bids[0].Price
	}
	return
}

// Price returns the middle of Bid and Ask.
func (o *OrderBook) Price() float64 {
	latest := (o.Bid().Price + o.Ask().Price) / float64(2)
	return latest
}
