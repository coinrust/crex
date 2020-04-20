package bybit

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
)

func TestNewWebSocket(t *testing.T) {
	params := &Parameters{
		DebugMode:  true,
		HttpClient: nil,
		AccessKey:  "6IASD6KDBdunn5qLpT",
		SecretKey:  "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF",
		Passphrase: "",
		Testnet:    true,
	}
	ws := NewBybitWebSocket(params)

	ws.SubscribeLevel2Snapshots(Market{Symbol: "BTCUSD"}, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	ws.SubscribeTrades(Market{Symbol: "BTCUSD"}, func(trades []Trade) {
		log.Printf("%#v", trades)
	})

	select {}
}
