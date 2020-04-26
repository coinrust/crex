package okexfutures

import (
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
	"time"

	. "github.com/coinrust/crex"
)

func testExchange() Exchange {
	testConfig := configtest.LoadTestConfig("okexfutures")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.Passphrase = testConfig.Passphrase
	params.Testnet = testConfig.Testnet
	params.ProxyURL = testConfig.ProxyURL
	params.ApiURL = "https://www.okex.me"
	return NewOkexFutures(params)
}

func TestOkexFutures_GetTime(t *testing.T) {
	ex := testExchange()
	tm, err := ex.GetTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", tm)
}

func TestOkexFutures_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestOkexFutures_GetOrderBook(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	//symbol := "BTC-USD-200626"
	for {
		ob, err := ex.GetOrderBook(symbol, 5)
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

func TestOkexFutures_GetRecords(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200410"
	start := time.Now().Add(-20 * time.Hour)
	end := time.Now()
	records, err := ex.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%v: %#v", v.Timestamp.String(), v)
	}
}

func TestOkexFutures_GetContractID(t *testing.T) {
	ex := testExchange()
	ex.SetContractType("BTC-USD", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
}

func TestOkexFutures_PlaceOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	order, err := ex.PlaceOrder(
		symbol,
		Buy,
		OrderTypeLimit,
		3000,
		1,
		OrderPostOnlyOption(true))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOkexFutures_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	orders, err := ex.GetOpenOrders(symbol)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestOkexFutures_GetOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	order, err := ex.GetOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOkexFutures_CancelOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	id := "4605829824487425"
	ret, err := ex.CancelOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}

func TestOkexFutures_GetPosition(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-200327"
	position, err := ex.GetPositions(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", position)
}
