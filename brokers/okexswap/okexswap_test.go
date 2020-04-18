package okexswap

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func newTestBroker() Broker {
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
	params.Testnet = true
	return New(params)
}

func TestGetBalance(t *testing.T) {
	b := newTestBroker()
	balance, err := b.GetBalance("BTC-USD-SWAP")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestGetOrderBook(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"

	ob, err := b.GetOrderBook(symbol, 5)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range ob.Asks {
		t.Logf("Ask: %v", v.Price)
	}

	for _, v := range ob.Bids {
		t.Logf("Bid: %v", v.Price)
	}
	t.Logf("Time: %v", ob.Time)
}

func TestOKEXSwapBroker_GetRecords(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	start := time.Now().Add(-20 * time.Minute)
	end := time.Now()
	records, err := b.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%v: %#v", v.Timestamp.String(), v)
	}
}

func TestOKEXBroker_PlaceOrder(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	order, err := b.PlaceOrder(
		symbol,
		Buy,
		OrderTypeLimit,
		3000,
		0,
		1,
		true,
		false,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOKEXBroker_GetOpenOrders(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	orders, err := b.GetOpenOrders(symbol)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestOKEXBroker_GetOrder(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	id := "469142537568198656"
	order, err := b.GetOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOKEXBroker_CancelOrder(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	id := "469142537568198656"
	ret, err := b.CancelOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}

func TestOKEXBroker_GetPosition(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-SWAP"
	positions, err := b.GetPositions(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", positions)
}
