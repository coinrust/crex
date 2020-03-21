package bybit_broker

import (
	"errors"
	. "github.com/coinrust/gotrader"
	"github.com/frankrap/bybit-api/rest"
	"log"
	"strings"
)

// BybitBroker the Bybit broker
type BybitBroker struct {
	client *rest.ByBit
}

func (b *BybitBroker) Subscribe(event string, param string, listener interface{}) {

}

func (b *BybitBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
	var balance rest.Balance
	balance, err = b.client.GetWalletBalance(currency)
	if err != nil {
		return
	}

	result.Equity = balance.Equity
	result.Balance = balance.WalletBalance
	result.Pnl = balance.UnrealisedPnl
	return
}

func (b *BybitBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ob rest.OrderBook
	ob, err = b.client.GetOrderBook(symbol)
	if err != nil {
		return
	}

	for _, v := range ob.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  v.Price,
			Amount: v.Size,
		})
	}

	for _, v := range ob.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  v.Price,
			Amount: v.Size,
		})
	}

	result.Time = ob.Time
	return
}

func (b *BybitBroker) SetContractType(contractType string) (err error) {
	return
}

func (b *BybitBroker) SetLeverRate(value float64) (err error) {
	return
}

func (b *BybitBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool) (result Order, err error) {
	if orderType == OrderTypeLimit || orderType == OrderTypeMarket {
		return b.placeOrder(symbol, direction, orderType, price, size, postOnly, reduceOnly)
	} else if orderType == OrderTypeStopLimit || orderType == OrderTypeStopMarket {
		return b.placeStopOrder(symbol, direction, orderType, price, stopPx, size, postOnly, reduceOnly)
	} else {
		err = errors.New("error")
		return
	}
}

func (b *BybitBroker) placeOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, postOnly bool, reduceOnly bool) (result Order, err error) {
	var side string
	var _orderType string
	var timeInForce string

	if direction == Buy {
		side = "Buy"
	} else if direction == Sell {
		side = "Sell"
	}
	if orderType == OrderTypeLimit {
		_orderType = "Limit"
	} else if orderType == OrderTypeMarket {
		_orderType = "Market"
	}

	if postOnly {
		timeInForce = "PostOnly"
	} else {
		timeInForce = "GoodTillCancel"
	}
	var order rest.Order
	order, err = b.client.CreateOrder(
		side,
		_orderType,
		price,
		int(size),
		timeInForce,
		reduceOnly,
		symbol,
	)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BybitBroker) placeStopOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool) (result Order, err error) {
	var side string
	var _orderType string
	var timeInForce string

	if direction == Buy {
		side = "Buy"
	} else if direction == Sell {
		side = "Sell"
	}
	if orderType == OrderTypeStopLimit {
		_orderType = "Limit"
	} else if orderType == OrderTypeStopMarket {
		_orderType = "Market"
	}

	if postOnly {
		timeInForce = "PostOnly"
	} else {
		timeInForce = "GoodTillCancel"
	}
	basePrice := stopPx // 触发价
	var order rest.Order
	order, err = b.client.CreateStopOrder(
		side,
		_orderType,
		price,
		basePrice,
		0,
		int(size),
		"",
		timeInForce,
		reduceOnly,
		symbol,
	)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BybitBroker) GetOpenOrders(symbol string) (result []Order, err error) {
	limit := 10
	orderStatus := "Created,New,PartiallyFilled,PendingCancel"
	for page := 1; page <= 5; page++ {
		var orders []rest.Order
		orders, err = b.client.GetOrders("", "", page, limit, orderStatus, symbol)
		log.Printf("page=%v %#v", page, orders)
		if err != nil {
			return
		}
		for _, v := range orders {
			//log.Printf("%#v", v)
			result = append(result, b.convertOrder(&v))
		}
		if len(orders) < limit {
			break
		}
	}
	return
}

func (b *BybitBroker) GetOrder(symbol string, id string) (result Order, err error) {
	var ret rest.OrderV2
	ret, err = b.client.GetOrderByID(id, "", symbol)
	if err != nil {
		return
	}
	result = b.convertOrderV2(&ret)
	return
}

func (b *BybitBroker) CancelOrder(symbol string, id string) (result Order, err error) {
	var order rest.OrderV2
	order, err = b.client.CancelOrderV2(id, "", symbol)
	if err != nil {
		return
	}
	result = b.convertOrderV2(&order)
	return
}

func (b *BybitBroker) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllOrder(symbol)
	return
}

func (b *BybitBroker) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	var order rest.Order
	order, err = b.client.ReplaceOrder(symbol, id, int(size), price)
	if err != nil {
		return
	}

	result = b.convertOrder(&order)
	return
}

func (b *BybitBroker) GetPosition(symbol string) (result Position, err error) {
	var ret rest.Position
	ret, err = b.client.GetPosition(symbol)
	if err != nil {
		return
	}
	result.Symbol = ret.Symbol
	result.Size = ret.Size
	result.AvgPrice = ret.EntryPrice
	return
}

func (b *BybitBroker) convertOrder(order *rest.Order) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.Symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Size = order.Qty
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.OrderType)
	if order.CumExecQty > 0 && order.CumExecValue > 0 {
		result.AvgPrice = order.CumExecQty / order.CumExecValue
	}
	result.FilledAmount = order.CumExecQty
	if strings.Contains(order.TimeInForce, "PostOnly") {
		result.PostOnly = true
	}
	result.ReduceOnly = false
	result.Status = b.orderStatus(order.OrderStatus)
	return
}

func (b *BybitBroker) convertOrderV2(order *rest.OrderV2) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.Symbol
	result.Price, _ = order.Price.Float64()
	result.StopPx = 0
	result.Size = order.Qty
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.OrderType)
	cumExecValue, err := order.CumExecValue.Float64()
	if err == nil && order.CumExecQty > 0 && cumExecValue > 0 {
		result.AvgPrice = float64(order.CumExecQty) / cumExecValue
	}
	result.FilledAmount = float64(order.CumExecQty)
	if strings.Contains(order.TimeInForce, "PostOnly") {
		result.PostOnly = true
	}
	result.ReduceOnly = false
	result.Status = b.orderStatus(order.OrderStatus)
	return
}

func (b *BybitBroker) convertDirection(side string) Direction {
	switch side {
	case "Buy":
		return Buy
	case "Sell":
		return Sell
	default:
		return Buy
	}
}

func (b *BybitBroker) convertOrderType(orderType string) OrderType {
	switch orderType {
	case "Limit":
		return OrderTypeLimit
	case "Market":
		return OrderTypeMarket
	default:
		return OrderTypeLimit
	}
}

func (b *BybitBroker) orderStatus(orderStatus string) OrderStatus {
	switch orderStatus {
	case "Created":
		return OrderStatusCreated
	case "New":
		return OrderStatusNew
	case "PartiallyFilled":
		return OrderStatusPartiallyFilled
	case "Filled":
		return OrderStatusFilled
	case "PendingCancel":
		return OrderStatusCancelPending
	case "Cancelled":
		return OrderStatusCancelled
	case "Rejected":
		return OrderStatusRejected
	case "Untriggered":
		return OrderStatusUntriggered
	case "Triggered":
		return OrderStatusTriggered
	default:
		return OrderStatusCreated
	}
}

func (b *BybitBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, apiKey string, secretKey string) *BybitBroker {
	client := rest.New(addr, apiKey, secretKey)
	for i := 0; i < 3; i++ {
		err := client.SetCorrectServerTime()
		if err != nil {
			log.Printf("%v", err)
			continue
		}
	}
	return &BybitBroker{
		client: client,
	}
}
