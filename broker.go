package gotrader

import (
	. "github.com/coinrust/gotrader/models"
)

type Broker interface {
	Subscribe(event string, param string, listener interface{})
	GetAccountSummary(currency string) (result AccountSummary, err error)
	GetOrderBook(symbol string, depth int) (result OrderBook, err error)
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, amount float64,
		postOnly bool, reduceOnly bool) (result Order, err error)
	GetOpenOrders(symbol string) (result []Order, err error)
	GetOrder(symbol string, id string) (result Order, err error)
	CancelAllOrders(symbol string) (err error)
	CancelOrder(symbol string, id string) (result Order, err error)
	GetPosition(symbol string) (result Position, err error)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}
