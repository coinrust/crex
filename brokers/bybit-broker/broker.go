package bybit_broker

import (
	. "github.com/coinrust/gotrader/models"
	"github.com/frankrap/bybit-api/rest"
	"strings"
)

// BybitBroker the Bybit broker
type BybitBroker struct {
	client *rest.ByBit
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

func (b *BybitBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	amount float64, postOnly bool, reduceOnly bool) (result Order, err error) {
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
	order, err = b.client.CreateOrder(side, _orderType, price, int(amount), timeInForce, reduceOnly, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BybitBroker) GetOpenOrders(symbol string) (result []Order, err error) {
	limit := 10
	for page := 1; page <= 5; page++ {
		var orders []rest.Order
		orders, err = b.client.GetOrders("", "", page, limit, "", symbol)
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
	var ret rest.Order
	ret, err = b.client.GetOrderByID(id, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&ret)
	return
}

func (b *BybitBroker) CancelOrder(symbol string, id string) (result Order, err error) {
	var order rest.Order
	order, err = b.client.CancelOrder(id, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BybitBroker) CancelAllOrders(symbol string) (err error) {
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
	result.Amount = order.Qty
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
	result.Status = b.orderStatus(order)
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

func (b *BybitBroker) orderStatus(order *rest.Order) OrderStatus {
	switch order.OrderStatus {
	case "Created":
		return OrderStatusCreated
	case "New":
		return OrderStatusNew
	case "PartiallyFilled":
		return OrderStatusPartiallyFilled
	case "Filled":
		return OrderStatusFilled
	case "Cancelled":
		return OrderStatusCancelled
	case "Rejected":
		return OrderStatusRejected
	default:
		return OrderStatusCreated
	}
}

func (b *BybitBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, apiKey string, secretKey string) *BybitBroker {
	client := rest.New(addr, apiKey, secretKey)
	return &BybitBroker{
		client: client,
	}
}
