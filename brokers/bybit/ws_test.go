package bybit

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
	"time"
)

func TestNewWS(t *testing.T) {
	params := &Parameters{
		HttpClient: nil,
		AccessKey:  "6IASD6KDBdunn5qLpT",
		SecretKey:  "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF",
		Passphrase: "",
		Testnet:    true,
	}
	ws := NewWS(params)
	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("%#v", trades)
	})

	time.Sleep(3 * time.Second)

	ws.SubscribeLevel2Snapshots(Market{ID: "BTCUSD"})
	ws.SubscribeTrades(Market{ID: "BTCUSD"})

	time.Sleep(5 * time.Second)
	ws.SubscribeLevel2Snapshots(Market{ID: "BTCUSD"})

	select {}
}
