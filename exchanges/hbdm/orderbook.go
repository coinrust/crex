package hbdm

import (
	"fmt"
	"github.com/MauriceGit/skiplist"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdm"
)

type DobItem struct {
	Price  float64
	Amount float64
}

func (e DobItem) ExtractKey() float64 {
	return e.Price
}

func (e DobItem) String() string {
	return fmt.Sprintf("%.2f", e.Price)
}

type DepthOrderBook struct {
	symbol string
	asks   skiplist.SkipList
	bids   skiplist.SkipList
}

func (d *DepthOrderBook) GetSymbol() string {
	return d.symbol
}

func (d *DepthOrderBook) Update(data *hbdm.WSDepthHF) {
	if data.Tick.Event == "snapshot" {
		d.asks = skiplist.New()
		d.bids = skiplist.New()
		for _, ask := range data.Tick.Asks {
			d.asks.Insert(DobItem{
				Price:  ask[0],
				Amount: ask[1],
			})
		}
		for _, bid := range data.Tick.Bids {
			d.bids.Insert(DobItem{
				Price:  bid[0],
				Amount: bid[1],
			})
		}
		return
	}

	if data.Tick.Event == "update" {
		for _, ask := range data.Tick.Asks {
			price := ask[0]
			amount := ask[1]
			if amount == 0 {
				d.asks.Delete(DobItem{
					Price:  price,
					Amount: amount,
				})
			} else {
				item := DobItem{
					Price:  price,
					Amount: amount,
				}
				elem, ok := d.asks.Find(item)
				if ok {
					d.asks.ChangeValue(elem, item)
				} else {
					d.asks.Insert(item)
				}
			}
		}
		for _, bid := range data.Tick.Bids {
			price := bid[0]
			amount := bid[1]
			if amount == 0 {
				d.bids.Delete(DobItem{
					Price:  price,
					Amount: amount,
				})
			} else {
				item := DobItem{
					Price:  price,
					Amount: amount,
				}
				elem, ok := d.bids.Find(item)
				if ok {
					d.bids.ChangeValue(elem, item)
				} else {
					d.bids.Insert(item)
				}
			}
		}
	}
}

func (d *DepthOrderBook) GetOrderBook(depth int) (result OrderBook) {
	result.Symbol = d.symbol
	smallest := d.asks.GetSmallestNode()
	if smallest != nil {
		item := smallest.GetValue().(DobItem)
		result.Asks = append(result.Asks, Item{
			Price:  item.Price,
			Amount: item.Amount,
		})
		count := 1
		node := smallest
		for count < depth {
			node = d.asks.Next(node)
			if node == nil {
				break
			}
			item := node.GetValue().(DobItem)
			result.Asks = append(result.Asks, Item{
				Price:  item.Price,
				Amount: item.Amount,
			})
			count++
		}
	}

	largest := d.bids.GetLargestNode()
	if largest != nil {
		item := largest.GetValue().(DobItem)
		result.Bids = append(result.Bids, Item{
			Price:  item.Price,
			Amount: item.Amount,
		})
		count := 1
		node := largest
		for count < depth {
			node = d.bids.Prev(node)
			if node == nil {
				break
			}
			item := node.GetValue().(DobItem)
			result.Bids = append(result.Bids, Item{
				Price:  item.Price,
				Amount: item.Amount,
			})
			count++
		}
	}
	return
}

func NewDepthOrderBook(symbol string) *DepthOrderBook {
	return &DepthOrderBook{
		symbol: symbol,
		asks:   skiplist.New(),
		bids:   skiplist.New(),
	}
}
