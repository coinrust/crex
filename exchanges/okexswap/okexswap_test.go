package okexswap

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"testing"
	"time"
)

func testExchange() Exchange {
	testConfig := configtest.LoadTestConfig("okexswap")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.Passphrase = testConfig.Passphrase
	params.Testnet = testConfig.Testnet
	params.ProxyURL = testConfig.ProxyURL
	return NewOkexSwap(params)
}

func TestOKEXSwap_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("BTC-USD-SWAP")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestOKEXSwap_GetOrderBook(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"

	ob, err := ex.GetOrderBook(symbol, 5)
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

func TestOKEXSwap_GetRecords(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
	start := time.Now().Add(-20 * time.Minute)
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

func TestOKEXSwap_PlaceOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
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

func TestOKEXSwap_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
	orders, err := ex.GetOpenOrders(symbol)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestOKEXSwap_GetOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
	id := "469142537568198656"
	order, err := ex.GetOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestOKEXSwap_CancelOrder(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
	id := "469142537568198656"
	ret, err := ex.CancelOrder(symbol, id)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}

func TestOKEXSwap_GetPosition(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD-SWAP"
	positions, err := ex.GetPositions(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", positions)
}
