package okexfutures

import (
	"log"
	"testing"
	"time"

	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
)

func newForTest() Broker {
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
	return New(params)
}

func TestGetAccountSummary(t *testing.T) {
	b := newForTest()
	accountSummary, err := b.GetAccountSummary("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
}

func TestGetOrderBook(t *testing.T) {
	b := newForTest()
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

func TestOKEXFutures_GetRecords(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD-200410"
	start := time.Now().Add(-20 * time.Hour)
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

func TestOKEXFutures_GetContractID(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC-USD", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
}

func TestOKEXFutures_PlaceOrder(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD-200327"
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

func TestOKEXFutures_GetOpenOrders(t *testing.T) {
	b := newForTest()
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

func TestOKEXFutures_GetOrder(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	order, err := b.GetOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOKEXFutures_CancelOrder(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	ret, err := b.CancelOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}

func TestOKEXFutures_GetPosition(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD-200327"
	position, err := b.GetPosition(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", position)
}
