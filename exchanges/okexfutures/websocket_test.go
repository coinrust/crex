package okexfutures

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func testWebSocket() *FuturesWebSocket {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	params := &Parameters{}
	params.AccessKey = viper.GetString("access_key")
	params.SecretKey = viper.GetString("secret_key")
	params.Passphrase = viper.GetString("passphrase")
	params.ProxyURL = viper.GetString("proxy_url")
	params.Testnet = true
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
