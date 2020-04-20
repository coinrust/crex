package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
	"log"
)

func main() {
	ws := exchanges.NewExchange(exchanges.HbdmSwap,
		ApiProxyURLOption("socks5://127.0.0.1:1080"), // 使用代理
		//ApiAccessKeyOption("[accessKey]"),
		//ApiSecretKeyOption("[secretKey]"),
		ApiTestnetOption(false),
		ApiWebSocketOption(true)) // 开启 WebSocket

	market := Market{Symbol: "BTC-USD"}
	// 订阅订单薄
	ws.SubscribeLevel2Snapshots(market, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	// 订阅成交记录
	ws.SubscribeTrades(market, func(trades []Trade) {
		log.Printf("%#v", trades)
	})
	// 订阅订单成交信息
	ws.SubscribeOrders(market, func(orders []Order) {
		log.Printf("%#v", orders)
	})
	// 订阅持仓信息
	ws.SubscribePositions(market, func(positions []Position) {
		log.Printf("%#v", positions)
	})

	select {}
}
