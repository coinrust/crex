package deribit_broker

import (
	. "github.com/coinrust/gotrader/models"
	"github.com/frankrap/deribit-api"
	"github.com/frankrap/deribit-api/models"
	"time"
)

// DiribitBroker the deribit broker
type DiribitBroker struct {
	client *deribit.Client
}

func (b *DiribitBroker) Subscribe(event string, param string, listener interface{}) {
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
	b.client.Subscribe([]string{event})
	b.client.On(event, listener)
}

func (b *DiribitBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
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

func (b *DiribitBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
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
	result.Time = time.Unix(0, ret.Timestamp*1e6) // 1581819533335
	return
}

func (b *DiribitBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	amount float64, postOnly bool, reduceOnly bool) (result Order, err error) {
	var _orderType string
	if orderType == OrderTypeLimit {
		_orderType = models.OrderTypeLimit
	} else if orderType == OrderTypeMarket {
		_orderType = models.OrderTypeMarket
	}
	if direction == Buy {
		var ret models.BuyResponse
		ret, err = b.client.Buy(&models.BuyParams{
			InstrumentName: symbol,
			Amount:         amount,
			Type:           _orderType,
			//Label:          "",
			Price: price,
			//TimeInForce:    "",
			//MaxShow:        nil,
			PostOnly:   postOnly,
			ReduceOnly: reduceOnly,
			//StopPrice:      0,
			//Trigger:        "",
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
			Amount:         amount,
			Type:           _orderType,
			//Label:          "",
			Price: price,
			//TimeInForce:    "",
			//MaxShow:        nil,
			PostOnly:   postOnly,
			ReduceOnly: reduceOnly,
			//StopPrice:      0,
			//Trigger:        "",
			//Advanced:       "",
		})
		if err != nil {
			return
		}
		result = b.convertOrder(&ret.Order)
	}
	return
}

func (b *DiribitBroker) GetOpenOrders(symbol string) (result []Order, err error) {
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

func (b *DiribitBroker) GetOrder(symbol string, id string) (result Order, err error) {
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

func (b *DiribitBroker) CancelOrder(symbol string, id string) (result Order, err error) {
	var order models.Order
	order, err = b.client.Cancel(&models.CancelParams{OrderID: id})
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *DiribitBroker) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllByInstrument(&models.CancelAllByInstrumentParams{
		InstrumentName: symbol,
	})
	return
}

func (b *DiribitBroker) GetPosition(symbol string) (result Position, err error) {
	var ret models.Position
	ret, err = b.client.GetPosition(&models.GetPositionParams{InstrumentName: symbol})
	if err != nil {
		return
	}
	result.Symbol = ret.InstrumentName
	result.Size = ret.Size
	result.AvgPrice = ret.AveragePrice
	return
}

func (b *DiribitBroker) convertOrder(order *models.Order) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.InstrumentName
	result.Price = order.Price.ToFloat64()
	result.Amount = order.Amount
	result.Direction = b.convertDirection(order.Direction)
	result.Type = b.convertOrderType(order.OrderType)
	result.AvgPrice = order.AveragePrice
	result.FilledAmount = order.FilledAmount
	result.PostOnly = order.PostOnly
	result.ReduceOnly = order.ReduceOnly
	result.Status = b.orderStatus(order)
	return
}

func (b *DiribitBroker) convertDirection(direction string) Direction {
	switch direction {
	case models.DirectionBuy:
		return Buy
	case models.DirectionSell:
		return Sell
	default:
		return Buy
	}
}

func (b *DiribitBroker) convertOrderType(orderType string) OrderType {
	switch orderType {
	case models.OrderTypeLimit:
		return OrderTypeLimit
	case models.OrderTypeMarket:
		return OrderTypeMarket
	case models.OrderTypeStopLimit:
		return OrderTypeLimit
	case models.OrderTypeStopMarket:
		return OrderTypeMarket
	default:
		return OrderTypeLimit
	}
}

func (b *DiribitBroker) orderStatus(order *models.Order) OrderStatus {
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
		return OrderStatusCreated
	default:
		return OrderStatusCreated
	}
}

func (b *DiribitBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, apiKey string, secretKey string) *DiribitBroker {
	cfg := &deribit.Configuration{
		Addr:          addr,
		ApiKey:        apiKey,
		SecretKey:     secretKey,
		AutoReconnect: true,
	}
	client := deribit.New(cfg)
	return &DiribitBroker{
		client: client,
	}
}
