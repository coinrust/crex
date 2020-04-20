package deribit

import (
	. "github.com/coinrust/crex"
	"testing"
	"time"
)

func newForTest() Exchange {
	params := &Parameters{
		AccessKey: "AsJTU16U",
		SecretKey: "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		Testnet:   true,
	}
	b := NewDeribit(params)
	return b
}

func TestDiribit_GetOrderBook(t *testing.T) {
	b := newForTest()
	b.GetOrderBook("BTC-PERPETUAL", 10)
}

func TestDiribit_GetRecords(t *testing.T) {
	b := newForTest()
	start := time.Now().Add(-time.Hour)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	records, err := b.GetRecords("BTC-PERPETUAL",
		"1", start.Unix(), end, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}

func TestDiribit_PlaceStopOrder(t *testing.T) {
	b := newForTest()
	order, err := b.PlaceOrder("BTC-PERPETUAL",
		Buy,
		OrderTypeStopMarket,
		0,
		10,
		OrderStopPxOption(8900))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
	t.Logf("Status: %v", order.Status.String())
}

func TestDiribit_GetOpenOrders(t *testing.T) {
	b := newForTest()
	orders, err := b.GetOpenOrders("BTC-PERPETUAL")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v Type: %v Status: %v",
			v,
			v.Type.String(),
			v.Status.String(),
		)
	}
}
