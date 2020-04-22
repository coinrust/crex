package deribit

import (
	"fmt"
	"github.com/MauriceGit/skiplist"
	. "github.com/coinrust/crex"
	"github.com/frankrap/deribit-api/models"
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

func (d *DepthOrderBook) Update(data *models.OrderBookRawNotification) {
	if data.PrevChangeID == 0 {
		d.asks = skiplist.New()
		d.bids = skiplist.New()
		// 举例: ["411.8", "10", "1", "4"]
		// 411.8为深度价格，10为此价格的合约张数，1为此价格的强平单个数，4为此价格的订单个数。
		for _, ask := range data.Asks {
			d.asks.Insert(DobItem{
				Price:  ask.Price,
				Amount: ask.Amount,
			})
		}
		for _, bid := range data.Bids {
			d.bids.Insert(DobItem{
				Price:  bid.Price,
				Amount: bid.Amount,
			})
		}
		return
	}

	if data.PrevChangeID > 0 {
		for _, ask := range data.Asks {
			if ask.Action == "delete" {
				d.asks.Delete(DobItem{
					Price:  ask.Price,
					Amount: ask.Amount,
				})
			} else if ask.Action == "new" || ask.Action == "change" {
				item := DobItem{
					Price:  ask.Price,
					Amount: ask.Amount,
				}
				elem, ok := d.asks.Find(item)
				if ok {
					d.asks.ChangeValue(elem, item)
				} else {
					d.asks.Insert(item)
				}
			}
		}
		for _, bid := range data.Bids {
			if bid.Action == "delete" {
				d.bids.Delete(DobItem{
					Price:  bid.Price,
					Amount: bid.Amount,
				})
			} else if bid.Action == "new" || bid.Action == "change" {
				item := DobItem{
					Price:  bid.Price,
					Amount: bid.Amount,
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
