package hbdmswap

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func testWebSocket() *SwapWebSocket {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	params := &Parameters{}
	params.AccessKey = viper.GetString("access_key")
	params.SecretKey = viper.GetString("secret_key")
	params.ProxyURL = viper.GetString("proxy_url")
	params.Testnet = true
	ws := NewSwapWebSocket(params)
	return ws
}

func TestSwapWebSocket_AllInOne(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots(Market{
		Symbol:       "BTC-USD",
		ContractType: "",
	}, func(ob *OrderBook) {
		t.Logf("%#v", ob)
	})
	ws.SubscribeTrades(Market{
		Symbol:       "BTC-USD",
		ContractType: "",
	}, func(trades []Trade) {
		t.Logf("%#v", trades)
	})

	select {}
}

func TestSwapWebSocket_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeOrders(Market{
		Symbol:       "BTC-USD",
		ContractType: "",
	}, func(orders []Order) {
		log.Printf("%#v", orders)
	})

	select {}
}
