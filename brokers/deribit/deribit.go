package deribit

import (
	"errors"
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	"github.com/frankrap/deribit-api"
	"github.com/frankrap/deribit-api/models"
	"time"
)

// Diribit the deribit broker
type Diribit struct {
	client           *deribit.Client
	orderBookManager *OrderBookManager
	emitter          *emission.Emitter
}

func (b *Diribit) GetName() (name string) {
	return "deribit"
}

func (b *Diribit) Subscribe(event string, param string, listener interface{}) {
	//b.client.Subscribe([]string{
	//	"announcements",
	//	"book.BTC-PERPETUAL.100.1.100ms",
	//	"book.BTC-PERPETUAL.100ms",
	//	"deribit_price_index.btc_usd",
	//	"deribit_price_ranking.btc_usd",
	//	"estimated_expiration_price.btc_usd",
	//	"markprice.options.btc_usd",
	//	"perpetual.BTC-PERPETUAL.raw",
	//	"quote.BTC-PERPETUAL",
	//	"ticker.BTC-PERPETUAL.raw",
	//	"user.changes.BTC-PERPETUAL.raw",
	//	"user.changes.future.BTC.raw",
	//	"user.orders.BTC-PERPETUAL.raw",
	//	"user.orders.future.BTC.100ms",
	//	"user.portfolio.btc",
	//	"user.trades.BTC-PERPETUAL.raw",
	//	"user.trades.future.BTC.100ms",
	//})
	if event == "orderbook" {
		b.emitter.On(event, listener)
		b.client.On(param, b.handleOrderBook)
		b.client.Subscribe([]string{param})
	}
}

func (b *Diribit) handleOrderBook(m *models.OrderBookNotification) {
	b.orderBookManager.Update(m)
	ob, ok := b.orderBookManager.GetOrderBook(m.InstrumentName)
	if !ok {
		return
	}
	b.emitter.Emit("orderbook", &ob)
}

func (b *Diribit) GetAccountSummary(currency string) (result AccountSummary, err error) {
	params := &models.GetAccountSummaryParams{
		Currency: currency,
		Extended: false,
	}
	var ret models.AccountSummary
	ret, err = b.client.GetAccountSummary(params)
	if err != nil {
		return
	}
	result.Equity = ret.Equity
	result.Balance = ret.Balance
	result.Pnl = ret.TotalPl
	return
}

func (b *Diribit) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	params := &models.GetOrderBookParams{
		InstrumentName: symbol,
		Depth:          depth,
	}
	var ret models.GetOrderBookResponse
	ret, err = b.client.GetOrderBook(params)
	if err != nil {
		return
	}
	for _, v := range ret.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	result.Time = time.Unix(0, ret.Timestamp*int64(time.Millisecond)) // 1581819533335
	return
}

func (b *Diribit) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	if end == 0 {
		end = time.Now().Unix()
	}
	params := &models.GetTradingviewChartDataParams{
		InstrumentName: symbol,
		StartTimestamp: from * 1000,
		EndTimestamp:   end * 1000,
		Resolution:     period,
	}
	var resp models.GetTradingviewChartDataResponse
	resp, err = b.client.GetTradingviewChartData(params)
	if err != nil {
		return
	}
	n := len(resp.Ticks)
	for i := 0; i < n; i++ {
		records = append(records, Record{
			Symbol:    symbol,
			Timestamp: time.Unix(0, resp.Ticks[i]*int64(time.Millisecond)),
			Open:      resp.Open[i],
			High:      resp.High[i],
			Low:       resp.Low[i],
			Close:     resp.Close[i],
			Volume:    resp.Volume[i],
		})
	}
	return
}

func (b *Diribit) SetContractType(currencyPair string, contractType string) (err error) {
	return
}

func (b *Diribit) GetContractID() (symbol string, err error) {
	return
}

func (b *Diribit) SetLeverRate(value float64) (err error) {
	return
}

func (b *Diribit) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool, params map[string]interface{}) (result Order, err error) {
	var _orderType string
	var trigger string
	if orderType == OrderTypeLimit {
		_orderType = models.OrderTypeLimit
		stopPx = 0
	} else if orderType == OrderTypeMarket {
		_orderType = models.OrderTypeMarket
		stopPx = 0
	} else if orderType == OrderTypeStopLimit {
		_orderType = models.OrderTypeStopLimit
		trigger = models.TriggerTypeLastPrice
	} else if orderType == OrderTypeStopMarket {
		_orderType = models.OrderTypeStopMarket
		trigger = models.TriggerTypeLastPrice
	}
	if direction == Buy {
		var ret models.BuyResponse
		ret, err = b.client.Buy(&models.BuyParams{
			InstrumentName: symbol,
			Amount:         size,
			Type:           _orderType,
			//Label:          "",
			Price: price,
			//TimeInForce:    "",
			//MaxShow:        nil,
			PostOnly:   postOnly,
			ReduceOnly: reduceOnly,
			StopPrice:  stopPx,
			Trigger:    trigger,
			//Advanced:       "",
		})
		if err != nil {
			return
		}
		result = b.convertOrder(&ret.Order)
	} else if direction == Sell {
		var ret models.SellResponse
		ret, err = b.client.Sell(&models.SellParams{
			InstrumentName: symbol,
			Amount:         size,
			Type:           _orderType,
			//Label:          "",
			Price: price,
			//TimeInForce:    "",
			//MaxShow:        nil,
			PostOnly:   postOnly,
			ReduceOnly: reduceOnly,
			StopPrice:  stopPx,
			Trigger:    trigger,
			//Advanced:       "",
		})
		if err != nil {
			return
		}
		result = b.convertOrder(&ret.Order)
	}
	return
}

func (b *Diribit) GetOpenOrders(symbol string) (result []Order, err error) {
	var ret []models.Order
	ret, err = b.client.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{
		InstrumentName: symbol,
		//Type:           "",
	})
	if err != nil {
		return
	}
	for _, v := range ret {
		result = append(result, b.convertOrder(&v))
	}
	return
}

func (b *Diribit) GetOrder(symbol string, id string) (result Order, err error) {
	var ret models.Order
	ret, err = b.client.GetOrderState(&models.GetOrderStateParams{
		OrderID: id,
	})
	if err != nil {
		return
	}
	result = b.convertOrder(&ret)
	return
}

func (b *Diribit) CancelOrder(symbol string, id string) (result Order, err error) {
	var order models.Order
	order, err = b.client.Cancel(&models.CancelParams{OrderID: id})
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *Diribit) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllByInstrument(&models.CancelAllByInstrumentParams{
		InstrumentName: symbol,
	})
	return
}

func (b *Diribit) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	params := &models.EditParams{
		OrderID:   id,
		Amount:    0,
		Price:     0,
		PostOnly:  false,
		Advanced:  "",
		StopPrice: 0,
	}
	if price <= 0 {
		err = errors.New("price is required")
		return
	}
	if size <= 0 {
		err = errors.New("size is required")
		return
	}
	params.Price = price
	params.Amount = size
	var resp models.EditResponse
	resp, err = b.client.Edit(params)
	if err != nil {
		return
	}
	result = b.convertOrder(&resp.Order)
	return
}

func (b *Diribit) GetPositions(symbol string) (result []Position, err error) {
	var ret models.Position
	ret, err = b.client.GetPosition(&models.GetPositionParams{InstrumentName: symbol})
	if err != nil {
		return
	}
	result = []Position{
		{
			Symbol:    symbol,
			OpenTime:  time.Time{},
			OpenPrice: ret.AveragePrice,
			Size:      ret.Size,
			AvgPrice:  ret.AveragePrice,
		},
	}
	return
}

func (b *Diribit) convertOrder(order *models.Order) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.InstrumentName
	result.Price = order.Price.ToFloat64()
	result.StopPx = order.StopPrice
	result.Size = order.Amount
	result.Direction = b.convertDirection(order.Direction)
	result.Type = b.convertOrderType(order.OrderType)
	result.AvgPrice = order.AveragePrice
	result.FilledAmount = order.FilledAmount
	result.PostOnly = order.PostOnly
	result.ReduceOnly = order.ReduceOnly
	result.Status = b.orderStatus(order)
	return
}

func (b *Diribit) convertDirection(direction string) Direction {
	switch direction {
	case models.DirectionBuy:
		return Buy
	case models.DirectionSell:
		return Sell
	default:
		return Buy
	}
}

func (b *Diribit) convertOrderType(orderType string) OrderType {
	switch orderType {
	case models.OrderTypeLimit:
		return OrderTypeLimit
	case models.OrderTypeMarket:
		return OrderTypeMarket
	case models.OrderTypeStopLimit:
		return OrderTypeStopLimit
	case models.OrderTypeStopMarket:
		return OrderTypeStopMarket
	default:
		return OrderTypeLimit
	}
}

func (b *Diribit) orderStatus(order *models.Order) OrderStatus {
	orderState := order.OrderState
	switch orderState {
	case models.OrderStateOpen:
		if order.FilledAmount > 0 {
			return OrderStatusPartiallyFilled
		}
		return OrderStatusNew
	case models.OrderStateFilled:
		return OrderStatusFilled
	case models.OrderStateRejected:
		return OrderStatusRejected
	case models.OrderStateCancelled:
		return OrderStatusCancelled
	case models.OrderStateUntriggered:
		return OrderStatusUntriggered
	default:
		return OrderStatusCreated
	}
}

func (b *Diribit) RunEventLoopOnce() (err error) {
	return
}

func New(addr string, accessKey string, secretKey string) *Diribit {
	cfg := &deribit.Configuration{
		Addr:          addr,
		ApiKey:        accessKey,
		SecretKey:     secretKey,
		AutoReconnect: true,
	}
	client := deribit.New(cfg)
	return &Diribit{
		client:           client,
		orderBookManager: NewOrderBookManager(),
		emitter:          emission.NewEmitter(),
	}
}
