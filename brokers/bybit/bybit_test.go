package bybit

import (
	. "github.com/coinrust/crex"
	"testing"
	"time"
)

func newForTest() Broker {
	params := &Parameters{
		AccessKey: "6IASD6KDBdunn5qLpT",
		SecretKey: "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF",
		Testnet:   true,
	}
	b := New(params)
	return b
}

func TestBybit_GetBalance(t *testing.T) {
	b := newForTest()
	balance, err := b.GetBalance("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestBybit_GetOpenOrders(t *testing.T) {
	b := newForTest()
	orders, err := b.GetOpenOrders("BTCUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBybit_GetRecords(t *testing.T) {
	b := newForTest()
	start := time.Now().Add(-time.Hour)
	records, err := b.GetRecords("BTCUSD",
		"1", start.Unix(), 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
