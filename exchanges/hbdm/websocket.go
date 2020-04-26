package hbdm

import (
	"fmt"
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdm"
	"strings"
	"time"
)

type HbdmWebSocket struct {
	ws      *hbdm.WS
	nws     *hbdm.NWS
	dobMap  map[string]*DepthOrderBook
	emitter *emission.Emitter
}

func (s *HbdmWebSocket) SubscribeTrades(symbol string, contractType string, callback func(trades []Trade)) error {
	s.emitter.On(WSEventTrade, callback)
	s.ws.SubscribeTrade("trade_1",
		s.convertToSymbol(symbol, contractType))
	return nil
}

func (s *HbdmWebSocket) SubscribeLevel2Snapshots(symbol string, contractType string, callback func(ob *OrderBook)) error {
	s.emitter.On(WSEventL2Snapshot, callback)
	//s.ws.SubscribeDepth("depth_1",
	//	s.convertToSymbol(symbol, contractType))
	s.ws.SubscribeDepthHF("depth_1", s.convertToSymbol(symbol, contractType), 20, "incremental")
	return nil
}

func (s *HbdmWebSocket) SubscribeOrders(symbol string, contractType string, callback func(orders []Order)) error {
	if s.nws == nil {
		return ErrApiKeysRequired
	}
	s.emitter.On(WSEventOrder, callback)
	s.nws.SubscribeOrders("order_1", symbol)
	return nil
}

func (s *HbdmWebSocket) SubscribePositions(symbol string, contractType string, callback func(positions []Position)) error {
	if s.nws == nil {
		return ErrApiKeysRequired
	}
	s.emitter.On(WSEventPosition, callback)
	s.nws.SubscribePositions("position_1", symbol)
	return nil
}

func (s *HbdmWebSocket) convertToSymbol(currencyPair string, contractType string) string {
	var symbol string
	switch contractType {
	case ContractTypeW1, "this_week":
		symbol = currencyPair + "_CW"
	case ContractTypeW2, "next_week":
		symbol = currencyPair + "_NW"
	case ContractTypeQ1, "quarter":
		symbol = currencyPair + "_CQ"
	}
	return symbol
}

func (s *HbdmWebSocket) depthCallback(depth *hbdm.WSDepth) {
	// log.Printf("depthCallback %#v", *depth)
	// ch: market.BTC_CQ.depth.step0
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

func (s *HbdmWebSocket) depthHFCallback(depth *hbdm.WSDepthHF) {
	// ch: market.BTC_CQ.depth.size_20.high_freq
	symbol := depth.Ch
	if v, ok := s.dobMap[symbol]; ok {
		v.Update(depth)
		ob := v.GetOrderBook(20)
		s.emitter.Emit(WSEventL2Snapshot, &ob)
	} else {
		dob := NewDepthOrderBook(symbol)
		dob.Update(depth)
		s.dobMap[symbol] = dob
		ob := dob.GetOrderBook(20)
		s.emitter.Emit(WSEventL2Snapshot, &ob)
	}
}

func (s *HbdmWebSocket) tradeCallback(trade *hbdm.WSTrade) {
	// log.Printf("tradeCallback")
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

func (s *HbdmWebSocket) ordersCallback(order *hbdm.WSOrder) {
	//log.Printf("ordersCallback")
	var o Order
	o.ID = fmt.Sprint(order.OrderID)
	o.Symbol = order.Symbol
	o.Price = order.Price
	o.AvgPrice = order.TradeAvgPrice
	// o.StopPx = 0
	o.Amount = order.Volume
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

func (s *HbdmWebSocket) positionsCallback(positions *hbdm.WSPositions) {
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

func NewHbdmWebSocket(params *Parameters) *HbdmWebSocket {
	wsURL := "wss://api.hbdm.com/ws"
	if params.WsURL != "" {
		wsURL = params.WsURL
	}
	s := &HbdmWebSocket{
		dobMap:  make(map[string]*DepthOrderBook),
		emitter: emission.NewEmitter(),
	}
	ws := hbdm.NewWS(wsURL, "", "", params.DebugMode)
	if params.ProxyURL != "" {
		ws.SetProxy(params.ProxyURL)
	}
	ws.SetDepthCallback(s.depthCallback)
	ws.SetDepthHFCallback(s.depthHFCallback)
	ws.SetTradeCallback(s.tradeCallback)
	ws.Start()
	s.ws = ws
	if params.AccessKey != "" && params.SecretKey != "" {
		nwsURL := strings.Replace(wsURL,
			"/ws", "/notification", -1)
		nws := hbdm.NewNWS(nwsURL, params.AccessKey, params.SecretKey)
		if params.ProxyURL != "" {
			nws.SetProxy(params.ProxyURL)
		}
		nws.SetOrdersCallback(s.ordersCallback)
		nws.SetPositionsCallback(s.positionsCallback)
		nws.Start()
		s.nws = nws
	}
	return s
}
