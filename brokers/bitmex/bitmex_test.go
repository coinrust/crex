package bitmex

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/bitmex-api"
	"testing"
	"time"
)

func newForTest() *BitMEX {
	apiKey := "eEtTUdma5LgAmryFerX-DAdp"
	secretKey := "kPjKmu-EIe1E73poRTnUraQWCMWbRq7PZ2-bzP8cnemniMXu"
	b := New(bitmex.HostTestnet, apiKey, secretKey)
	b.client.SetProxy("127.0.0.1:1080")
	return b
}

func TestBitMEX_GetOrderBook(t *testing.T) {
	b := newForTest()
	ob, err := b.GetOrderBook("XBTUSD", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestBitMEX_GetRecords(t *testing.T) {
	b := newForTest()
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := b.GetRecords("XBTUSD",
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}

func TestBitMEX_PlaceOrder(t *testing.T) {
	b := newForTest()
	order, err := b.PlaceOrder("XBTUSD",
		Buy,
		OrderTypeLimit,
		8000,
		0,
		10,
		true,
		false,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestBitMEX_GetOpenOrders(t *testing.T) {
	b := newForTest()
	orders, err := b.GetOpenOrders("XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBitMEX_GetOrder(t *testing.T) {
	b := newForTest()
	order, err := b.GetOrder("XBTUSD", "c90d5194-0f6a-31db-7942-c3caf1f8a055")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}
