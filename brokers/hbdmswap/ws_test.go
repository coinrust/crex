package hbdmswap

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newTestWS() *WS {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	params := &Parameters{}
	params.AccessKey = viper.GetString("access_key")
	params.SecretKey = viper.GetString("secret_key")
	params.Testnet = true
	ws := NewWS(params)
	return ws
}

func TestWS_AllInOne(t *testing.T) {
	ws := newTestWS()

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

func TestWS_SubscribeOrders(t *testing.T) {
	ws := newTestWS()

	ws.On(WSEventOrder, func(orders []Order) {
		log.Printf("orders: %#v", orders)
	})

	ws.SubscribeOrders(Market{
		ID:     "BTC-USD",
		Params: "",
	})

	select {}
}
