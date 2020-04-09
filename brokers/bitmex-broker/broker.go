package bitmex_broker

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/bitmex-api"
	"github.com/frankrap/bitmex-api/swagger"
	"strings"
)

// BitMEXBroker the BitMEX broker
type BitMEXBroker struct {
	client *bitmex.BitMEX
	symbol string
}

func (b *BitMEXBroker) GetName() (name string) {
	return "BitMEX"
}

func (b *BitMEXBroker) Subscribe(event string, param string, listener interface{}) {

}

func (b *BitMEXBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
	var margin swagger.Margin
	margin, err = b.client.GetMargin()
	if err != nil {
		return
	}
	result.Equity = float64(margin.MarginBalance)
	result.Balance = float64(margin.WalletBalance)
	result.Pnl = 0
	return
}

func (b *BitMEXBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ret bitmex.OrderBook
	ret, err = b.client.GetOrderBook(depth, symbol)
	if err != nil {
		return
	}
	for _, v := range ret.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	result.Time = ret.Timestamp
	return
}

func (b *BitMEXBroker) SetContractType(currencyPair string, contractType string) (err error) {
	b.symbol = currencyPair
	return
}

func (b *BitMEXBroker) GetContractID() (symbol string, err error) {
	return b.symbol, nil
}

func (b *BitMEXBroker) SetLeverRate(value float64) (err error) {
	return
}

func (b *BitMEXBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool, params map[string]interface{}) (result Order, err error) {
	var side string
	var _orderType string
	if direction == Buy {
		side = bitmex.SIDE_BUY
	} else if direction == Sell {
		side = bitmex.SIDE_SELL
	}
	if orderType == OrderTypeLimit {
		_orderType = bitmex.ORD_TYPE_LIMIT
	} else if orderType == OrderTypeMarket {
		_orderType = bitmex.ORD_TYPE_MARKET
	} else if orderType == OrderTypeStopLimit {
		_orderType = bitmex.ORD_TYPE_STOP_LIMIT
	} else if orderType == OrderTypeStopMarket {
		_orderType = bitmex.ORD_TYPE_STOP
	}
	var execInst string
	if postOnly {
		execInst = "ParticipateDoNotInitiate"
	}
	if reduceOnly {
		if execInst != "" {
			execInst += ","
		}
		execInst += "ReduceOnly"
	}
	var order swagger.Order
	order, err = b.client.PlaceOrder(side, _orderType, stopPx, price, int32(size), "", execInst, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BitMEXBroker) GetOpenOrders(symbol string) (result []Order, err error) {
	var ret []swagger.Order
	ret, err = b.client.GetOrders(symbol)
	if err != nil {
		return
	}
	for _, v := range ret {
		result = append(result, b.convertOrder(&v))
	}
	return
}

func (b *BitMEXBroker) GetOrder(symbol string, id string) (result Order, err error) {
	var ret swagger.Order
	ret, err = b.client.GetOrder(id, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&ret)
	return
}

func (b *BitMEXBroker) CancelOrder(symbol string, id string) (result Order, err error) {
	var order swagger.Order
	order, err = b.client.CancelOrder(id)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BitMEXBroker) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllOrders(symbol)
	return
}

func (b *BitMEXBroker) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	var resp swagger.Order
	resp, err = b.client.AmendOrder2(id, "", "", 0, float32(size), 0, 0, price, 0, 0, "")
	if err != nil {
		return
	}
	result = b.convertOrder(&resp)
	return
}

func (b *BitMEXBroker) GetPosition(symbol string) (result Position, err error) {
	var ret swagger.Position
	ret, err = b.client.GetPosition(symbol)
	if err != nil {
		return
	}
	result.Symbol = ret.Symbol
	result.Size = float64(ret.CurrentQty)
	result.AvgPrice = ret.AvgEntryPrice
	return
}

func (b *BitMEXBroker) convertOrder(order *swagger.Order) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.Symbol
	result.Price = order.Price
	result.StopPx = order.StopPx
	result.Size = float64(order.OrderQty)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.OrdType)
	result.AvgPrice = order.AvgPx
	result.FilledAmount = float64(order.CumQty)
	if strings.Contains(order.ExecInst, "ParticipateDoNotInitiate") {
		result.PostOnly = true
	}
	if strings.Contains(order.ExecInst, "ReduceOnly") {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *BitMEXBroker) convertDirection(side string) Direction {
	switch side {
	case bitmex.SIDE_BUY:
		return Buy
	case bitmex.SIDE_SELL:
		return Sell
	default:
		return Buy
	}
}

func (b *BitMEXBroker) convertOrderType(orderType string) OrderType {
	switch orderType {
	case bitmex.ORD_TYPE_LIMIT:
		return OrderTypeLimit
	case bitmex.ORD_TYPE_MARKET:
		return OrderTypeMarket
	case bitmex.ORD_TYPE_STOP_LIMIT:
		return OrderTypeStopLimit
	case bitmex.ORD_TYPE_STOP:
		return OrderTypeStopMarket
	default:
		return OrderTypeLimit
	}
}

func (b *BitMEXBroker) orderStatus(order *swagger.Order) OrderStatus {
	orderState := order.OrdStatus
	switch orderState {
	case bitmex.OS_NEW:
		return OrderStatusNew
	case bitmex.OS_PARTIALLY_FILLED:
		return OrderStatusPartiallyFilled
	case bitmex.OS_FILLED:
		return OrderStatusFilled
	case bitmex.OS_CANCELED:
		return OrderStatusCancelled
	case bitmex.OS_REJECTED:
		return OrderStatusRejected
	default:
		return OrderStatusCreated
	}
}

func (b *BitMEXBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, accessKey string, secretKey string) *BitMEXBroker {
	client := bitmex.New(addr, accessKey, secretKey)
	return &BitMEXBroker{
		client: client,
	}
}
