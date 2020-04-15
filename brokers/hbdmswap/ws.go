package hbdmswap

import (
	"fmt"
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdmswap"
	"strings"
	"time"
)

type WS struct {
	ws      *hbdmswap.WS
	nws     *hbdmswap.NWS
	emitter *emission.Emitter
}

func (s *WS) On(event WSEvent, listener interface{}) {
	s.emitter.On(event, listener)
}

func (s *WS) SubscribeTrades(market Market) {
	s.ws.SubscribeTrade("trade_1", market.ID)
}

func (s *WS) SubscribeLevel2Snapshots(market Market) {
	s.ws.SubscribeDepth("depth_1", market.ID)
}

func (s *WS) SubscribeOrders(market Market) {
	if s.nws == nil {
		return
	}
	s.nws.SubscribeOrders("order_1", market.ID)
}

func (s *WS) SubscribePositions(market Market) {
	if s.nws == nil {
		return
	}
	s.nws.SubscribePositions("position_1", market.ID)
}

func (s *WS) depthCallback(depth *hbdmswap.WSDepth) {
	// log.Printf("depthCallback %#v", *depth)
	// ch: market.BTC-USD.depth.step0
	ob := &OrderBook{
		Symbol: depth.Ch,
		Time:   time.Unix(0, depth.Ts*int64(time.Millisecond)),
		Asks:   nil,
		Bids:   nil,
	}
	for _, v := range depth.Tick.Asks {
		ob.Asks = append(ob.Asks, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	for _, v := range depth.Tick.Bids {
		ob.Bids = append(ob.Bids, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	s.emitter.Emit(WSEventL2Snapshot, ob)
}

func (s *WS) tradeCallback(trade *hbdmswap.WSTrade) {
	var trades []Trade
	for _, v := range trade.Tick.Data {
		var direction Direction
		if v.Direction == "buy" {
			direction = Buy
		} else if v.Direction == "sell" {
			direction = Sell
		}
		t := Trade{
			ID:        fmt.Sprint(v.ID),
			Direction: direction,
			Price:     v.Price,
			Amount:    float64(v.Amount),
			Ts:        v.Ts,
			Symbol:    "",
		}
		trades = append(trades, t)
	}
	s.emitter.Emit(WSEventTrade, trades)
}

func (s *WS) ordersCallback(order *hbdmswap.WSOrder) {
	//log.Printf("ordersCallback")
	var o Order
	o.ID = fmt.Sprint(order.OrderID)
	o.Symbol = order.Symbol
	o.Price = order.Price
	o.AvgPrice = order.TradeAvgPrice
	// o.StopPx = 0
	o.Size = order.Volume
	o.FilledAmount = order.TradeVolume
	if order.Direction == "buy" {
		o.Direction = Buy
	} else if order.Direction == "sell" {
		o.Direction = Sell
	}
	// 订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制
	switch order.OrderPriceType {
	case "limit":
		o.Type = OrderTypeLimit
	case "opponent":
		o.Type = OrderTypeMarket
	case "post_only":
		o.Type = OrderTypeLimit
		o.PostOnly = true
	}
	// "open":开 "close":平
	switch order.Offset {
	case "open":
	case "close":
		o.ReduceOnly = true
	}
	// 订单状态(1准备提交 2准备提交 3已提交 4部分成交 5部分成交已撤单 6全部成交 7已撤单)
	switch order.Status {
	case 1:
		o.Status = OrderStatusNew
	case 2:
		o.Status = OrderStatusNew
	case 3:
		o.Status = OrderStatusNew
	case 4:
		o.Status = OrderStatusPartiallyFilled
	case 5:
		o.Status = OrderStatusCancelled
	case 6:
		o.Status = OrderStatusFilled
	case 7:
		o.Status = OrderStatusCancelled
	case 11:
		o.Status = OrderStatusCancelPending
	default:
		o.Status = OrderStatusCreated
	}
	s.emitter.Emit(WSEventOrder, []Order{o})
}

func (s *WS) positionsCallback(positions *hbdmswap.WSPositions) {
	//log.Printf("positionsCallback")
	var eventData []Position
	for _, v := range positions.Data {
		var o Position
		o.Symbol = v.Symbol
		o.OpenTime = time.Unix(0, positions.Ts*int64(time.Millisecond))
		o.OpenPrice = v.CostOpen
		switch v.Direction {
		case "buy":
			o.Size = v.Volume
		case "sell":
			o.Size = -v.Volume
		}
		o.AvgPrice = v.CostHold
		eventData = append(eventData, o)
	}
	s.emitter.Emit(WSEventPosition, eventData)
}

func NewWS(accessKey string, secretKey string, testnet bool) *WS {
	wsURL := "wss://api.hbdm.com/swap-ws"
	s := &WS{
		emitter: emission.NewEmitter(),
	}
	ws := hbdmswap.NewWS(wsURL, accessKey, secretKey)
	ws.SetDepthCallback(s.depthCallback)
	ws.SetTradeCallback(s.tradeCallback)
	ws.Start()
	s.ws = ws
	if accessKey != "" && secretKey != "" {
		nwsURL := strings.Replace(wsURL,
			"/swap-ws", "/swap-notification", -1)
		nws := hbdmswap.NewNWS(nwsURL, accessKey, secretKey)
		nws.SetOrdersCallback(s.ordersCallback)
		nws.SetPositionsCallback(s.positionsCallback)
		nws.Start()
		s.nws = nws
	}
	return s
}
