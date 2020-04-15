package bybit

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
	"time"
)

func TestNewWS(t *testing.T) {
	ws := NewWS("wss://stream-testnet.bybit.com/realtime",
		"6IASD6KDBdunn5qLpT", "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF")
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
