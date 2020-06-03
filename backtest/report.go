package backtest

import (
	"fmt"
	. "github.com/coinrust/crex"
	"strings"
	"time"
)

// SOrder "event":"order"/"deal"
type SOrder struct {
	Ts        time.Time   // ts: 2019-10-02T07:03:53.584+0800
	Order     *Order      // order
	OrderBook *OrderBook  // orderbook
	Positions []*Position // positions
	Balances  []float64   // balances
	Comment   string      // msg: Place order/Match order
}

func (o *SOrder) BalancesString() string {
	n := len(o.Balances)
	if n == 0 {
		return ""
	}
	if n == 1 {
		return fmt.Sprintf("%v", o.Balances[0])
	}
	var list []string
	for _, v := range o.Balances {
		list = append(list, fmt.Sprintf("%v", v))
	}
	return strings.Join(list, " | ")
}
