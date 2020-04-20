package hbdm

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func testWebSocket() *HbdmWebSocket {
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
	ws := NewHbdmWebSocket(params)
	return ws
}

func TestHbdmWebSocket_AllInOne(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeLevel2Snapshots("BTC", ContractTypeW1, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	ws.SubscribeTrades("BTC", ContractTypeW1, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	select {}
}

func TestHbdmWebSocket_SubscribeOrders(t *testing.T) {
	ws := testWebSocket()

	ws.SubscribeOrders("BTC", ContractTypeW1, func(orders []Order) {
		log.Printf("orders: %#v", orders)
	})

	select {}
}
