package okex_broker

import (
	"log"
	"testing"
	"time"

	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
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
	accountSummary, err := b.GetAccountSummary("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
}

func TestGetOrderBook(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-200327"
	//symbol := "BTC-USD-200626"
	for {
		ob, err := b.GetOrderBook(symbol, 5)
		if err != nil {
			t.Error(err)
			return
		}
		//t.Logf("%#v", ob)
		//t.Logf("Ask: %v Bid: %v", ob.AskPrice(), ob.BidPrice())
		log.Printf("Ask: %v Bid: %v", ob.AskPrice(), ob.BidPrice())
		time.Sleep(500 * time.Millisecond)
	}
	//for _, v := range ob.Asks {
	//	t.Logf("%v", v.Price)
	//}

	// for _, v := range ob.Bids {
	// 	t.Logf("%v", v.Price)
	// }
	//t.Logf("Time: %v", ob.Time)
}

func TestOKEXBroker_PlaceOrder(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-200327"
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
	symbol := "BTC-USD-200327"
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
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	order, err := b.GetOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOKEXBroker_CancelOrder(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	ret, err := b.CancelOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}

func TestOKEXBroker_GetPosition(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD-200327"
	position, err := b.GetPosition(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", position)
}
