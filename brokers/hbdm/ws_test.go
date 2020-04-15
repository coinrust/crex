package hbdm

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

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	ws := NewWS(accessKey, secretKey, false)
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
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	ws.SubscribeTrades(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})

	select {}
}

func TestWS_SubscribeOrders(t *testing.T) {
	ws := newTestWS()

	ws.On(WSEventOrder, func(orders []Order) {
		log.Printf("orders: %#v", orders)
	})

	ws.SubscribeOrders(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})

	select {}
}
