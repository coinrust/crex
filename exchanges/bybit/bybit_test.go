package bybit

import (
	. "github.com/coinrust/crex"
	"testing"
	"time"
)

func testExchange() Exchange {
	params := &Parameters{
		AccessKey: "6IASD6KDBdunn5qLpT",
		SecretKey: "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF",
		Testnet:   true,
	}
	ex := NewBybit(params)
	return ex
}

func TestBybit_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestBybit_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	orders, err := ex.GetOpenOrders("BTCUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBybit_GetRecords(t *testing.T) {
	ex := testExchange()
	start := time.Now().Add(-time.Hour)
	records, err := ex.GetRecords("BTCUSD",
		"1", start.Unix(), 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
