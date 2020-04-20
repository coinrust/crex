package binancefutures

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"testing"
	"time"
)

func testExchange() Exchange {
	testConfig := configtest.LoadTestConfig("binancefutures")
	params := &Parameters{
		AccessKey: testConfig.AccessKey,
		SecretKey: testConfig.SecretKey,
		Testnet:   testConfig.Testnet,
		ProxyURL:  testConfig.ProxyURL,
	}
	ex := NewBinanceFutures(params)
	return ex
}

func TestBinanceFutures_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("USDT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestBinanceFutures_GetOrderBook(t *testing.T) {
	ex := testExchange()
	ob, err := ex.GetOrderBook("BTCUSDT", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestBinanceFutures_GetRecords(t *testing.T) {
	ex := testExchange()
	now := time.Now()
	start := now.Add(-300 * time.Minute)
	end := now
	records, err := ex.GetRecords("BTCUSDT",
		PERIOD_1MIN, start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("Timestamp: %v %#v", v.Timestamp, v)
	}
}

func TestBinanceFutures_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	orders, err := ex.GetOpenOrders("BTCUSDT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orders)
}
