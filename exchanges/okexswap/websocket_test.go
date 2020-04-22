package okexswap

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
)

func testWebSocket() *SwapWebSocket {
	testConfig := configtest.LoadTestConfig("okexswap")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.Passphrase = testConfig.Passphrase
	params.ProxyURL = testConfig.ProxyURL
	params.Testnet = testConfig.Testnet
	ws := NewSwapWebSocket(params)
	return ws
}

func TestWS_AllInOne(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots(Market{
		Symbol: "BTC-USD-SWAP",
	}, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	ws.SubscribeTrades(Market{
		Symbol: "BTC-USD-SWAP",
	}, func(trades []Trade) {
		log.Printf("%#v", trades)
	})

	select {}
}

func TestWS_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeOrders(Market{
		Symbol: "BTC-USD-SWAP",
	}, func(orders []Order) {
		log.Printf("%#v", orders)
	})

	select {}
}
