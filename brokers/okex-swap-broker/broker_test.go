package okex_futures_broker

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newTestBroker() Broker {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	passphrase := viper.GetString("passphrase")
	baseURL := "https://www.okex.me" // https://www.okex.com
	return NewBroker(baseURL, accessKey, secretKey, passphrase)
}

func TestGetAccountSummary(t *testing.T) {
	b := newTestBroker()
	accountSummary, err := b.GetAccountSummary("BTC-USD-SWAP")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
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
		false)
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
	position, err := b.GetPosition(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", position)
}
