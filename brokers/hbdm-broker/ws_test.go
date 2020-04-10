package hbdm_broker

import (
	. "github.com/coinrust/crex"
	"log"
	"testing"
)

func TestWS_AllInOne(t *testing.T) {
	ws := NewWS("wss://api.btcgateway.pro/ws",
		"", "")

	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	ws.SubscribeLevel2Snapshots(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	ws.SubscribeTrades(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})

	select {}
}
