package hbdmswap

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
)

func testWebSocket() *SwapWebSocket {
	testConfig := configtest.LoadTestConfig("hbdmswap")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.ProxyURL = testConfig.ProxyURL
	params.Testnet = testConfig.Testnet
	ws := NewSwapWebSocket(params)
	return ws
}

func TestSwapWebSocket_AllInOne(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots(Market{
		Symbol: "BTC-USD",
	}, func(ob *OrderBook) {
		t.Logf("%#v", ob)
	})
	ws.SubscribeTrades(Market{
		Symbol: "BTC-USD",
	}, func(trades []Trade) {
		t.Logf("%#v", trades)
	})

	select {}
}

func TestSwapWebSocket_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeOrders(Market{
		Symbol: "BTC-USD",
	}, func(orders []Order) {
		log.Printf("%#v", orders)
	})

	select {}
}
