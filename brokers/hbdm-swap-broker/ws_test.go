package hbdm_swap_broker

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
)

func TestWS_AllInOne(t *testing.T) {
	ws := NewWS("wss://api.btcgateway.pro/swap-ws",
		"", "")
	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	ws.SubscribeLevel2Snapshots(Market{
		ID:     "BTC-USD",
		Params: "",
	})
	ws.SubscribeTrades(Market{
		ID:     "BTC-USD",
		Params: "",
	})

	select {}
}
