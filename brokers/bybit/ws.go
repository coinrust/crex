package bybit

import (
	"fmt"
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	bws "github.com/frankrap/bybit-api/ws"
	"time"
)

type WS struct {
	ws      *bws.ByBitWS
	emitter *emission.Emitter
}

func (s *WS) On(event WSEvent, listener interface{}) {
	s.emitter.On(event, listener)
}

func (s *WS) SubscribeTrades(market Market) {
	var arg = bws.WSTrade
	if market.ID != "" {
		arg += "." + market.ID
	}
	s.ws.Subscribe(arg)
}

func (s *WS) SubscribeLevel2Snapshots(market Market) {
	arg := bws.WSOrderBook25L1 + "." + market.ID
	s.ws.Subscribe(arg)
}

func (s *WS) SubscribeOrders(market Market) {
	s.ws.Subscribe(bws.WSOrder)
}

func (s *WS) SubscribePositions(market Market) {
	s.ws.Subscribe(bws.WSPosition)
}

func (s *WS) handleOrderBook(symbol string, data bws.OrderBook) {
	//log.Printf("handleOrderBook symbol: %v", symbol)
	ob := &OrderBook{
		Symbol: symbol,
	}
	for _, v := range data.Asks {
		ob.Asks = append(ob.Asks, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	for _, v := range data.Bids {
		ob.Bids = append(ob.Bids, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	ob.Time = data.Timestamp
	s.emitter.Emit(WSEventL2Snapshot, ob)
}

func (s *WS) handleTrade(symbol string, data []*bws.Trade) {
	var trades []Trade
	for _, v := range data {
		var direction Direction
		if v.Side == "Buy" {
			direction = Buy
		} else if v.Side == "Sell" {
			direction = Sell
		}
		trades = append(trades, Trade{
			ID:        v.TradeID,
			Direction: direction,
			Price:     v.Price,
			Amount:    float64(v.Size),
			Ts:        v.Timestamp.UnixNano() / int64(time.Millisecond),
			Symbol:    v.Symbol,
		})
	}
	s.emitter.Emit(WSEventTrade, trades)
}

func (s *WS) handlePosition(data []*bws.Position) {
	var eventData []Position
	now := time.Now()
	for _, v := range data {
		var o Position
		o.Symbol = v.Symbol
		o.OpenTime = now
		o.OpenPrice = v.EntryPrice
		switch v.Side {
		case "Buy":
			o.Size = v.Size
		case "Sell":
			o.Size = -v.Size
		}
		o.AvgPrice = v.EntryPrice
		eventData = append(eventData, o)
	}
	s.emitter.Emit(WSEventPosition, eventData)
}

func (s *WS) handleOrder(data []*bws.Order) {
	var orders []Order
	for _, v := range data {
		var o Order
		o.ID = fmt.Sprint(v.OrderID)
		o.Symbol = v.Symbol
		o.Price = v.Price
		//o.AvgPrice = 0
		// o.StopPx = 0
		o.Size = v.Qty
		o.FilledAmount = v.CumExecQty
		if v.Side == "Buy" {
			o.Direction = Buy
		} else if v.Side == "Sell" {
			o.Direction = Sell
		}
		switch v.OrderType {
		case "Limit":
			o.Type = OrderTypeLimit
		case "Market":
			o.Type = OrderTypeMarket
		}
		if v.TimeInForce == "PostOnly" {
			o.PostOnly = true
		}
		o.Status = s.orderStatus(v.OrderStatus)
	}
	s.emitter.Emit(WSEventOrder, orders)
}

func (s *WS) orderStatus(orderStatus string) OrderStatus {
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

func NewWS(params *Parameters) *WS {
	wsURL := "wss://stream.bybit.com/realtime"
	if params.Testnet {
		wsURL = "wss://stream-testnet.bybit.com/realtime"
	}
	s := &WS{
		emitter: emission.NewEmitter(),
	}
	cfg := &bws.Configuration{
		Addr:          wsURL,
		ApiKey:        params.AccessKey,
		SecretKey:     params.SecretKey,
		AutoReconnect: true,
	}
	ws := bws.New(cfg)
	s.ws = ws
	ws.On(bws.WSOrderBook25L1, s.handleOrderBook)
	ws.On(bws.WSTrade, s.handleTrade)
	ws.On(bws.WSOrder, s.handleOrder)
	ws.On(bws.WSPosition, s.handlePosition)
	ws.Start()
	return s
}
