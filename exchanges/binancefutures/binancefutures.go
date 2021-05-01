package binancefutures

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/utils"
)

// BinanceFutures the Binance futures exchange
type BinanceFutures struct {
	client *futures.Client
	symbol string
}

func (b *BinanceFutures) GetName() (name string) {
	return "binancefutures"
}

func (b *BinanceFutures) GetTime() (tm int64, err error) {
	tm, err = b.client.NewServerTimeService().
		Do(context.Background())
	return
}

// SetProxy ...
// proxyURL: http://127.0.0.1:1080
func (b *BinanceFutures) SetProxy(proxyURL string) error {
	proxyURL_, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL_),
	}

	//adding the Transport object to the http Client
	b.client.HTTPClient.Transport = transport
	return nil
}

// currency: USDT
func (b *BinanceFutures) GetBalance(currency string) (result *Balance, err error) {
	var res []*futures.Balance
	res, err = b.client.NewGetBalanceService().
		Do(context.Background())
	if err != nil {
		return
	}
	result = &Balance{}
	for _, v := range res {
		if v.Asset == currency { // USDT
			value := utils.ParseFloat64(v.Balance)
			result.Equity = value
			result.Available = value
			break
		}
	}
	return
}

func (b *BinanceFutures) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
	result = &OrderBook{}
	if depth <= 5 {
		depth = 5
	} else if depth <= 10 {
		depth = 10
	} else if depth <= 20 {
		depth = 20
	} else if depth <= 50 {
		depth = 50
	} else if depth <= 100 {
		depth = 100
	} else if depth <= 500 {
		depth = 500
	} else {
		depth = 1000
	}
	var res *futures.DepthResponse
	res, err = b.client.NewDepthService().
		Symbol(symbol).
		Limit(depth).
		Do(context.Background())
	if err != nil {
		return
	}
	for _, v := range res.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  utils.ParseFloat64(v.Price),
			Amount: utils.ParseFloat64(v.Quantity),
		})
	}
	for _, v := range res.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  utils.ParseFloat64(v.Price),
			Amount: utils.ParseFloat64(v.Quantity),
		})
	}
	result.Time = time.Now()
	return
}

func (b *BinanceFutures) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error) {
	var res []*futures.Kline
	service := b.client.NewKlinesService().
		Symbol(symbol).
		Interval(b.IntervalKlinePeriod(period)).
		Limit(limit)
	if from > 0 {
		service = service.StartTime(from * 1000)
	}
	if end > 0 {
		service = service.EndTime(end * 1000)
	}
	res, err = service.Do(context.Background())
	if err != nil {
		return
	}
	for _, v := range res {
		records = append(records, &Record{
			Symbol:    symbol,
			Timestamp: time.Unix(0, v.OpenTime*int64(time.Millisecond)),
			Open:      utils.ParseFloat64(v.Open),
			High:      utils.ParseFloat64(v.High),
			Low:       utils.ParseFloat64(v.Low),
			Close:     utils.ParseFloat64(v.Close),
			Volume:    utils.ParseFloat64(v.Volume),
		})
	}
	return
}

func (b *BinanceFutures) IntervalKlinePeriod(period string) string {
	m := map[string]string{
		PERIOD_1WEEK: "7d",
	}
	if v, ok := m[period]; ok {
		return v
	}
	return period
}

func (b *BinanceFutures) SetContractType(currencyPair string, contractType string) (err error) {
	b.symbol = currencyPair
	return
}

func (b *BinanceFutures) GetContractID() (symbol string, err error) {
	return b.symbol, nil
}

func (b *BinanceFutures) SetLeverRate(value float64) (err error) {
	return
}

func (b *BinanceFutures) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *BinanceFutures) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *BinanceFutures) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *BinanceFutures) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *BinanceFutures) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...PlaceOrderOption) (result *Order, err error) {
	params := ParsePlaceOrderParameter(opts...)
	service := b.client.NewCreateOrderService().
		Symbol(symbol).
		Quantity(fmt.Sprint(size)).
		ReduceOnly(params.ReduceOnly)
	var side futures.SideType
	if direction == Buy {
		side = futures.SideTypeBuy
	} else if direction == Sell {
		side = futures.SideTypeSell
	}
	var _orderType futures.OrderType
	switch orderType {
	case OrderTypeLimit:
		_orderType = futures.OrderTypeLimit
	case OrderTypeMarket:
		_orderType = futures.OrderTypeMarket
	case OrderTypeStopMarket:
		_orderType = futures.OrderTypeStopMarket
		service = service.StopPrice(fmt.Sprint(params.StopPx))
	case OrderTypeStopLimit:
		_orderType = futures.OrderTypeStop
		service = service.StopPrice(fmt.Sprint(params.StopPx))
	}

	if orderType != OrderTypeMarket {
		service = service.TimeInForce(resolveTimeInForce(params.TimeInForce))
		if price > 0 {
			service = service.Price(fmt.Sprint(price))
		}
	}

	if params.PostOnly {
		service = service.TimeInForce(futures.TimeInForceTypeGTX)
	}

	service = service.Side(side).Type(_orderType)
	var res *futures.CreateOrderResponse
	res, err = service.Do(context.Background())
	if err != nil {
		return
	}
	result = b.convertOrder1(res)
	return
}

func resolveTimeInForce(timeInForce string) futures.TimeInForceType {
	var futuresTimeInForce futures.TimeInForceType

	switch timeInForce {
	case string(futures.TimeInForceTypeGTC):
		futuresTimeInForce = futures.TimeInForceTypeGTC
	case string(futures.TimeInForceTypeFOK):
		futuresTimeInForce = futures.TimeInForceTypeFOK
	case string(futures.TimeInForceTypeIOC):
		futuresTimeInForce = futures.TimeInForceTypeIOC
	case string(futures.TimeInForceTypeGTX):
		futuresTimeInForce = futures.TimeInForceTypeGTX
	default:
		futuresTimeInForce = futures.TimeInForceTypeGTC
	}

	return futuresTimeInForce
}

func (b *BinanceFutures) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	service := b.client.NewListOpenOrdersService().
		Symbol(symbol)
	var res []*futures.Order
	res, err = service.Do(context.Background())
	if err != nil {
		return
	}
	for _, v := range res {
		result = append(result, b.convertOrder(v))
	}
	return
}

func (b *BinanceFutures) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var orderID int64
	orderID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}
	var res *futures.Order
	res, err = b.client.NewGetOrderService().
		Symbol(symbol).
		OrderID(orderID).
		Do(context.Background())
	if err != nil {
		return
	}
	result = b.convertOrder(res)
	return
}

func (b *BinanceFutures) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var orderID int64
	orderID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}
	var res *futures.CancelOrderResponse
	res, err = b.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderID).
		Do(context.Background())
	if err != nil {
		return
	}
	result = b.convertOrder2(res)
	return
}

func (b *BinanceFutures) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	err = b.client.NewCancelAllOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	return
}

func (b *BinanceFutures) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (b *BinanceFutures) GetPositions(symbol string) (result []*Position, err error) {
	var res []*futures.PositionRisk
	res, err = b.client.NewGetPositionRiskService().
		Do(context.Background())
	if err != nil {
		return
	}

	useFilter := symbol != ""

	for _, v := range res {
		if useFilter && v.Symbol != symbol {
			continue
		}
		position := &Position{}
		position.Symbol = v.Symbol
		size := utils.ParseFloat64(v.PositionAmt)
		if size != 0 {
			position.Size = size
			position.OpenPrice = utils.ParseFloat64(v.EntryPrice)
			position.AvgPrice = position.OpenPrice
		}
		result = append(result, position)
	}
	return
}

func (b *BinanceFutures) convertOrder(order *futures.Order) (result *Order) {
	result = &Order{}
	result.ID = fmt.Sprint(order.OrderID)
	result.ClientOId = order.ClientOrderID
	result.Symbol = order.Symbol
	result.Price = utils.ParseFloat64(order.Price)
	result.StopPx = utils.ParseFloat64(order.StopPrice)
	result.Amount = utils.ParseFloat64(order.OrigQuantity)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.Type)
	result.AvgPrice = utils.ParseFloat64(order.AvgPrice)
	result.FilledAmount = utils.ParseFloat64(order.ExecutedQuantity)
	if order.TimeInForce == futures.TimeInForceTypeGTX {
		result.PostOnly = true
	}
	result.ReduceOnly = order.ReduceOnly
	result.Status = b.orderStatus(order.Status)
	result.Time = time.Unix(order.Time/int64(1e3), 0)
	result.UpdateTime = time.Unix(order.UpdateTime/int64(1e3), 0)
	return
}

func (b *BinanceFutures) convertOrder1(order *futures.CreateOrderResponse) (result *Order) {
	result = &Order{}
	result.ID = fmt.Sprint(order.OrderID)
	result.ClientOId = order.ClientOrderID
	result.Symbol = order.Symbol
	result.Price = utils.ParseFloat64(order.Price)
	result.StopPx = utils.ParseFloat64(order.StopPrice)
	result.Amount = utils.ParseFloat64(order.OrigQuantity)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.Type)
	result.AvgPrice = utils.ParseFloat64(order.AvgPrice)
	result.FilledAmount = utils.ParseFloat64(order.ExecutedQuantity)
	result.Time = time.Unix(order.UpdateTime/int64(1e3), 0)
	if order.TimeInForce == futures.TimeInForceTypeGTX {
		result.PostOnly = true
	}
	result.ReduceOnly = order.ReduceOnly
	result.Status = b.orderStatus(order.Status)
	return
}

func (b *BinanceFutures) convertOrder2(order *futures.CancelOrderResponse) (result *Order) {
	result = &Order{}
	result.ID = fmt.Sprint(order.OrderID)
	result.ClientOId = order.ClientOrderID
	result.Symbol = order.Symbol
	result.Price = utils.ParseFloat64(order.Price)
	result.StopPx = utils.ParseFloat64(order.StopPrice)
	result.Amount = utils.ParseFloat64(order.OrigQuantity)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.Type)
	result.AvgPrice = 0
	result.FilledAmount = utils.ParseFloat64(order.ExecutedQuantity)
	if order.TimeInForce == futures.TimeInForceTypeGTX {
		result.PostOnly = true
	}
	result.ReduceOnly = order.ReduceOnly
	result.Status = b.orderStatus(order.Status)
	return
}

func (b *BinanceFutures) convertDirection(side futures.SideType) Direction {
	switch side {
	case futures.SideTypeBuy:
		return Buy
	case futures.SideTypeSell:
		return Sell
	default:
		return Buy
	}
}

func (b *BinanceFutures) convertOrderType(orderType futures.OrderType) OrderType {
	/*
		OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"
		OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET"
	*/
	switch orderType {
	case futures.OrderTypeLimit:
		return OrderTypeLimit
	case futures.OrderTypeMarket:
		return OrderTypeMarket
	case futures.OrderTypeStop:
		return OrderTypeStopLimit
	case futures.OrderTypeStopMarket:
		return OrderTypeStopMarket
	default:
		return OrderTypeLimit
	}
}

func (b *BinanceFutures) orderStatus(status futures.OrderStatusType) OrderStatus {
	switch status {
	case futures.OrderStatusTypeNew:
		return OrderStatusNew
	case futures.OrderStatusTypePartiallyFilled:
		return OrderStatusPartiallyFilled
	case futures.OrderStatusTypeFilled:
		return OrderStatusFilled
	case futures.OrderStatusTypeCanceled:
		return OrderStatusCancelled
	case futures.OrderStatusTypeRejected:
		return OrderStatusRejected
	case futures.OrderStatusTypeExpired:
		return OrderStatusCancelled
	default:
		return OrderStatusCreated
	}
}

func (b *BinanceFutures) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	return ErrNotImplemented
}

func (b *BinanceFutures) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	return ErrNotImplemented
}

func (b *BinanceFutures) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	return ErrNotImplemented
}

func (b *BinanceFutures) SubscribePositions(market Market, callback func(positions []*Position)) error {
	return ErrNotImplemented
}

func (b *BinanceFutures) IO(name string, params string) (string, error) {
	return "", nil
}

func NewBinanceFutures(params *Parameters) *BinanceFutures {
	futures.UseTestnet = params.Testnet
	client := futures.NewClient(params.AccessKey, params.SecretKey)
	b := &BinanceFutures{
		client: client,
	}
	if params.ProxyURL != "" {
		b.SetProxy(params.ProxyURL)
	}
	return b
}
