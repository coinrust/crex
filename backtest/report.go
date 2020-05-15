package backtest

import (
	. "github.com/coinrust/crex"
	"time"
)

// SOrder "event":"order"/"deal"
type SOrder struct {
	Ts        time.Time   // ts: 2019-10-02T07:03:53.584+0800
	Order     *Order      // order
	OrderBook *OrderBook  // orderbook
	Positions []*Position // positions
	Balance   float64     // balance
	Comment   string      // msg: Place order/Match order
}
