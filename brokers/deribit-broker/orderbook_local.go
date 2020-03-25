package deribit_broker

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/deribit-api/models"
	"sort"
	"strconv"
	"sync"
	"time"
)

type OrderBookLocal struct {
	symbol string
	asks   map[string]*Item // key: price
	bids   map[string]*Item // key: price
	m      sync.RWMutex
}

func NewOrderBookLocal(symbol string) *OrderBookLocal {
	o := &OrderBookLocal{
		symbol: symbol,
		asks:   make(map[string]*Item),
		bids:   make(map[string]*Item),
	}
	return o
}

func (o *OrderBookLocal) GetOrderbook() (ob OrderBook) {
	o.m.RLock()
	defer o.m.RUnlock()

	for _, v := range o.asks {
		ob.Asks = append(ob.Asks, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	for _, v := range o.bids {
		ob.Bids = append(ob.Bids, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}

	sort.Slice(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	sort.Slice(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	ob.Time = time.Now()
	ob.Symbol = o.symbol

	return
}

func (o *OrderBookLocal) Key(price float64) string {
	return strconv.FormatFloat(price, 'f', 1, 64)
}

func (o *OrderBookLocal) Update(newOrderBook *models.OrderBookNotification) {
	o.m.Lock()
	defer o.m.Unlock()

	// [action, price, amount]
	// action: new, change or delete.

	if newOrderBook.Type == "snapshot" {
		o.asks = make(map[string]*Item)
		o.bids = make(map[string]*Item)

		for _, v := range newOrderBook.Asks {
			//action := v[0].(string)
			price := v[1].(float64)
			amount := v[2].(float64)
			key := o.Key(price)
			o.asks[key] = &Item{
				Price:  price,
				Amount: amount,
			}
		}
		for _, v := range newOrderBook.Bids {
			price := v[1].(float64)
			amount := v[2].(float64)
			key := o.Key(price)
			o.bids[key] = &Item{
				Price:  price,
				Amount: amount,
			}
		}
	} else if newOrderBook.Type == "change" || newOrderBook.Type == "" {
		for _, v := range newOrderBook.Asks {
			action := v[0].(string)
			price := v[1].(float64)
			amount := v[2].(float64)
			key := o.Key(price)
			if action == "new" {
				o.asks[key] = &Item{
					Price:  price,
					Amount: amount,
				}
			} else if action == "change" {
				o.asks[key] = &Item{
					Price:  price,
					Amount: amount,
				}
			} else if action == "delete" {
				delete(o.asks, key)
			}
		}

		for _, v := range newOrderBook.Bids {
			action := v[0].(string)
			price := v[1].(float64)
			amount := v[2].(float64)
			key := o.Key(price)
			if action == "new" {
				o.bids[key] = &Item{
					Price:  price,
					Amount: amount,
				}
			} else if action == "change" {
				o.bids[key] = &Item{
					Price:  price,
					Amount: amount,
				}
			} else if action == "delete" {
				delete(o.bids, key)
			}
		}
	}
}
