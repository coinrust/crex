package okexfutures

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
)

func testWebSocket() *FuturesWebSocket {
	testConfig := configtest.LoadTestConfig("okexfutures")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.Passphrase = testConfig.Passphrase
	params.ProxyURL = testConfig.ProxyURL
	params.Testnet = testConfig.Testnet
	ws := NewFuturesWebSocket(params)
	return ws
}

func TestFuturesWebSocket_AllInOne(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots(Market{
		Symbol: "BTC-USD-200626",
	}, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	ws.SubscribeTrades(Market{
		Symbol: "BTC-USD-200626",
	}, func(trades []Trade) {
		log.Printf("%#v", trades)
	})

	select {}
}

func TestFuturesWebSocket_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeOrders(Market{
		Symbol: "BTC-USD-200626",
	}, func(orders []Order) {
		log.Printf("%#v", orders)
	})

	select {}
}
