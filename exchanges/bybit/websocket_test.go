package bybit

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
)

func testWebSocket() *BybitWebSocket {
	params := &Parameters{
		DebugMode:  true,
		HttpClient: nil,
		AccessKey:  "6IASD6KDBdunn5qLpT",
		SecretKey:  "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF",
		Passphrase: "",
		Testnet:    true,
	}
	ws := NewBybitWebSocket(params)
	return ws
}

func TestNewWebSocket(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots(Market{Symbol: "BTCUSD"}, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	ws.SubscribeTrades(Market{Symbol: "BTCUSD"}, func(trades []Trade) {
		log.Printf("%#v", trades)
	})

	select {}
}

func TestBybitWebSocket_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	market := Market{Symbol: "BTCUSD"}
	ws.SubscribeOrders(market, func(orders []Order) {
		log.Printf("Orders: %#v", orders)
	})

	ws.SubscribePositions(market, func(positions []Position) {
		log.Printf("Positions: %#v", positions)
	})

	select {}
}
