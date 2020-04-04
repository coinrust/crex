package deribit_broker

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/deribit-api"
	"log"
	"testing"
	"time"
)

func newBroker() Broker {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	b := NewBroker(deribit.TestBaseURL, apiKey, secretKey)
	return b
}

func TestDiribitBroker_GetOrderBook(t *testing.T) {
	b := newBroker()
	b.GetOrderBook("BTC-PERPETUAL", 10)
}

func TestDiribitBroker_Subscribe(t *testing.T) {
	b := newBroker()
	//event := "book.ETH-PERPETUAL.100.1.100ms"
	param := "book.BTC-PERPETUAL.100ms"
	b.Subscribe("orderbook", param, func(e *OrderBook) {
		log.Printf("OrderBook: %#v", *e)
	})

	for {
		time.Sleep(1 * time.Second)
	}
}

func TestDiribitBroker_PlaceStopOrder(t *testing.T) {
	b := newBroker()
	order, err := b.PlaceOrder("BTC-PERPETUAL",
		Buy,
		OrderTypeStopMarket,
		0,
		8900,
		10,
		false,
		false,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
	t.Logf("Status: %v", order.Status.String())
}

func TestDiribitBroker_GetOpenOrders(t *testing.T) {
	b := newBroker()
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
