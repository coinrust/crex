package bybit

import (
	"errors"
	. "github.com/coinrust/crex"
	"github.com/frankrap/bybit-api/rest"
	"log"
	"strings"
	"time"
)

// Bybit the Bybit exchange
type Bybit struct {
	client *rest.ByBit
	ws     *BybitWebSocket
	params *Parameters
	symbol string
}

func (b *Bybit) GetName() (name string) {
	return "bybit"
}

func (b *Bybit) GetBalance(currency string) (result Balance, err error) {
	var balance rest.Balance
	balance, err = b.client.GetWalletBalance(currency)
	if err != nil {
		return
	}

	result.Equity = balance.Equity
	result.Available = balance.AvailableBalance
	result.RealizedPnl = balance.RealisedPnl
	result.UnrealisedPnl = balance.UnrealisedPnl
	return
}

func (b *Bybit) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
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

func (b *Bybit) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	var values []rest.OHLC
	values, err = b.client.GetKLine(symbol, period, from, limit)
	if err != nil {
		return
	}
	for _, v := range values {
		records = append(records, Record{
			Symbol:    v.Symbol,
			Timestamp: time.Unix(v.OpenTime, 0),
			Open:      v.Open,
			High:      v.High,
			Low:       v.Low,
			Close:     v.Close,
			Volume:    v.Volume,
		})
	}
	return
}

func (b *Bybit) SetContractType(currencyPair string, contractType string) (err error) {
	b.symbol = currencyPair
	return
}

func (b *Bybit) GetContractID() (symbol string, err error) {
	return b.symbol, nil
}

func (b *Bybit) SetLeverRate(value float64) (err error) {
	return
}

func (b *Bybit) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *Bybit) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *Bybit) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *Bybit) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *Bybit) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...OrderOption) (result Order, err error) {
	params := ParseOrderParameter(opts...)
	if orderType == OrderTypeLimit || orderType == OrderTypeMarket {
		return b.placeOrder(symbol, direction, orderType, price, size, params.PostOnly, params.ReduceOnly)
	} else if orderType == OrderTypeStopLimit || orderType == OrderTypeStopMarket {
		return b.placeStopOrder(symbol, direction, orderType, price, params.StopPx, size, params.PostOnly, params.ReduceOnly)
	} else {
		err = errors.New("error")
		return
	}
}

func (b *Bybit) placeOrder(symbol string, direction Direction, orderType OrderType, price float64,
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

func (b *Bybit) placeStopOrder(symbol string, direction Direction, orderType OrderType, price float64,
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

func (b *Bybit) GetOpenOrders(symbol string) (result []Order, err error) {
	limit := 10
	orderStatus := "Created,NewBybit,PartiallyFilled,PendingCancel"
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

func (b *Bybit) GetOrder(symbol string, id string) (result Order, err error) {
	var ret rest.OrderV2
	ret, err = b.client.GetOrderByID(id, "", symbol)
	if err != nil {
		return
	}
	result = b.convertOrderV2(&ret)
	return
}

func (b *Bybit) CancelOrder(symbol string, id string) (result Order, err error) {
	var order rest.OrderV2
	order, err = b.client.CancelOrderV2(id, "", symbol)
	if err != nil {
		return
	}
	result = b.convertOrderV2(&order)
	return
}

func (b *Bybit) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllOrder(symbol)
	return
}

func (b *Bybit) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	var order rest.Order
	order, err = b.client.ReplaceOrder(symbol, id, int(size), price)
	if err != nil {
		return
	}

	result = b.convertOrder(&order)
	return
}

func (b *Bybit) GetPositions(symbol string) (result []Position, err error) {
	var ret rest.Position
	ret, err = b.client.GetPosition(symbol)
	if err != nil {
		return
	}
	result = []Position{
		{
			Symbol:    symbol,
			OpenTime:  time.Time{},
			OpenPrice: ret.EntryPrice,
			Size:      ret.Size,
			AvgPrice:  ret.EntryPrice,
		},
	}
	return
}

func (b *Bybit) convertOrder(order *rest.Order) (result Order) {
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
	if order.ExtFields != nil {
		result.ReduceOnly = order.ExtFields.ReduceOnly
	}
	result.Status = b.orderStatus(order.OrderStatus)
	return
}

func (b *Bybit) convertOrderV2(order *rest.OrderV2) (result Order) {
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
	result.Status = b.orderStatus(order.OrderStatus)
	return
}

func (b *Bybit) convertDirection(side string) Direction {
	switch side {
	case "Buy":
		return Buy
	case "Sell":
		return Sell
	default:
		return Buy
	}
}

func (b *Bybit) convertOrderType(orderType string) OrderType {
	switch orderType {
	case "Limit":
		return OrderTypeLimit
	case "Market":
		return OrderTypeMarket
	default:
		return OrderTypeLimit
	}
}

func (b *Bybit) orderStatus(orderStatus string) OrderStatus {
	switch orderStatus {
	case "Created":
		return OrderStatusCreated
	case "NewBybit":
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

func (b *Bybit) SubscribeTrades(market Market, callback func(trades []Trade)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeTrades(market, callback)
}

func (b *Bybit) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeLevel2Snapshots(market, callback)
}

func (b *Bybit) SubscribeOrders(market Market, callback func(orders []Order)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeOrders(market, callback)
}

func (b *Bybit) SubscribePositions(market Market, callback func(positions []Position)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribePositions(market, callback)
}

func (b *Bybit) RunEventLoopOnce() (err error) {
	return
}

func NewBybit(params *Parameters) *Bybit {
	baseUri := "https://api.bybit.com/"
	if params.Testnet {
		baseUri = "https://api-testnet.bybit.com/"
	}
	client := rest.New(params.HttpClient, baseUri, params.AccessKey, params.SecretKey)
	for i := 0; i < 3; i++ {
		err := client.SetCorrectServerTime()
		if err != nil {
			log.Printf("%v", err)
			continue
		}
	}
	var ws *BybitWebSocket
	if params.WebSocket {
		ws = NewBybitWebSocket(params)
	}
	return &Bybit{
		client: client,
		ws:     ws,
		params: params,
	}
}
