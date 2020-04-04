package bitmex_broker

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/bitmex-api"
	"testing"
)

func newBrokerForTest() *BitMEXBroker {
	apiKey := "eEtTUdma5LgAmryFerX-DAdp"
	secretKey := "kPjKmu-EIe1E73poRTnUraQWCMWbRq7PZ2-bzP8cnemniMXu"
	b := NewBroker(bitmex.HostTestnet, apiKey, secretKey)
	b.client.SetProxy("127.0.0.1:1080")
	return b
}

func TestBitMEXBroker_GetOrderBook(t *testing.T) {
	b := newBrokerForTest()
	ob, err := b.GetOrderBook("XBTUSD", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestBitMEXBroker_PlaceOrder(t *testing.T) {
	b := newBrokerForTest()
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

func TestBitMEXBroker_GetOpenOrders(t *testing.T) {
	b := newBrokerForTest()
	orders, err := b.GetOpenOrders("XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBitMEXBroker_GetOrder(t *testing.T) {
	b := newBrokerForTest()
	order, err := b.GetOrder("XBTUSD", "c90d5194-0f6a-31db-7942-c3caf1f8a055")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}
