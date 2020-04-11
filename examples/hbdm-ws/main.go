package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
	"log"
)

func main() {
	wsURL := "wss://api.hbdm.com/ws" // "wss://api.btcgateway.pro/ws"
	params := map[string]string{}
	params["wsURL"] = wsURL

	ws := brokers.NewWS(brokers.HBDM,
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
	ws.On(WSEventPosition, func(position *Position) {
		log.Printf("position: %#v", position)
	})

	// 订阅订单薄
	ws.SubscribeLevel2Snapshots(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅成交记录
	ws.SubscribeTrades(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅订单成交信息
	ws.SubscribeOrders(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅持仓信息
	ws.SubscribePositions(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})

	select {}
}
