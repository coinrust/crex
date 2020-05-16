package main

import (
	"testing"
	"time"
)

func TestStore_GetCollection(t *testing.T) {
	s := NewStore("mongodb://localhost:27017",
		"tick_db")
	exchange := "deribit"
	symbol = "BTC-PERPETUAL"
	s.GetCollection(exchange, symbol, true)
}

func TestStore_Insert(t *testing.T) {
	s := NewStore("mongodb://localhost:27017", "test1")
	n := 10000
	for i := 0; i < n; i++ {
		var asks []Item
		var bids []Item
		asks = append(asks, Item{
			Price:  7000,
			Amount: 1000,
		})
		bids = append(bids, Item{
			Price:  6999,
			Amount: 500,
		})
		ob := &OrderBook{
			Timestamp: time.Now(),
			Asks:      asks,
			Bids:      bids,
		}
		s.Insert(ob)
	}
	s.SyncBuffer(true)
}
