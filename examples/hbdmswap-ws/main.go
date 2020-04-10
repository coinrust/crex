package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
	"log"
)

func main() {
	wsURL := "wss://api.hbdm.com/swap-ws" // "wss://api.btcgateway.pro/swap-ws"
	params := map[string]string{}
	params["wsURL"] = wsURL

	ws := brokers.NewWS(brokers.HBDMSwap,
		"", "", false, params)

	// 订单薄事件方法
	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	// 成交记录事件方法
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
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

	select {}
}
