package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
	"log"
)

func main() {
	wsURL := "wss://real.okex.com:8443/ws/v3"
	params := map[string]string{}
	params["wsURL"] = wsURL
	params["passphrase"] = "[passphrase]"

	ws := brokers.NewWS(brokers.OKEXSwap,
		"[accessKey]", "[secretKey]", false, params)

	// 订单薄事件方法
	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	// 成交记录事件方法
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	// 订单事件方法
	ws.On(WSEventOrder, func(order *Order) {
		log.Printf("order: %#v", order)
	})
	// 持仓事件方法
	ws.On(WSEventPosition, func(positions []Position) {
		log.Printf("positions: %#v", positions)
	})

	// 订阅订单薄
	ws.SubscribeLevel2Snapshots(Market{
		ID:     "BTC-USD-SWAP",
		Params: "",
	})
	// 订阅成交记录
	ws.SubscribeTrades(Market{
		ID:     "BTC-USD-SWAP",
		Params: "",
	})
	// 订阅订单成交信息
	ws.SubscribeOrders(Market{
		ID:     "BTC-USD-SWAP",
		Params: "",
	})
	// 订阅持仓信息
	ws.SubscribePositions(Market{
		ID:     "BTC-USD-SWAP",
		Params: "",
	})

	select {}
}
