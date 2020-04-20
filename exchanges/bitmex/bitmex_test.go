package bitmex

import (
	. "github.com/coinrust/crex"
	"testing"
	"time"
)

func testExchange() *BitMEX {
	params := &Parameters{
		AccessKey: "eEtTUdma5LgAmryFerX-DAdp",
		SecretKey: "kPjKmu-EIe1E73poRTnUraQWCMWbRq7PZ2-bzP8cnemniMXu",
		Testnet:   true,
	}
	ex := NewBitMEX(params)
	ex.client.SetProxy("127.0.0.1:1080")
	return ex
}

func TestBitMEX_GetOrderBook(t *testing.T) {
	ex := testExchange()
	ob, err := ex.GetOrderBook("XBTUSD", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestBitMEX_GetRecords(t *testing.T) {
	ex := testExchange()
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := ex.GetRecords("XBTUSD",
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
	ex := testExchange()
	order, err := ex.PlaceOrder("XBTUSD",
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
	ex := testExchange()
	orders, err := ex.GetOpenOrders("XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBitMEX_GetOrder(t *testing.T) {
	ex := testExchange()
	order, err := ex.GetOrder("XBTUSD", "c90d5194-0f6a-31db-7942-c3caf1f8a055")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}
