package hbdm_swap_broker

import (
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdmswap"
	"time"
)

type WS struct {
	ws      *hbdmswap.WS
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

}

func (s *WS) SubscribePositions(market Market) {

}

func (s *WS) depthCallback(depth *hbdmswap.WSDepth) {
	// log.Printf("depthCallback %#v", *depth)
	// ch: market.BTC-USD.depth.step0
	ob := &OrderBook{
		Symbol: depth.Ch,
		Time:   time.Unix(0, depth.Ts*1e6),
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
			ID:        v.ID,
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

func NewWS(wsURL string, accessKey string, secretKey string) *WS {
	s := &WS{
		emitter: emission.NewEmitter(),
	}
	ws := hbdmswap.NewWS(wsURL, accessKey, secretKey)
	ws.SetDepthCallback(s.depthCallback)
	ws.SetTradeCallback(s.tradeCallback)
	ws.Start()
	s.ws = ws
	return s
}
