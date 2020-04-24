package deribit

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
	"time"
)

func newForTest() Exchange {
	return testExchange(false)
}

func testExchange(websocket bool) Exchange {
	testConfig := configtest.LoadTestConfig("deribit")
	params := &Parameters{
		DebugMode: true,
		AccessKey: testConfig.AccessKey,
		SecretKey: testConfig.SecretKey,
		Testnet:   testConfig.Testnet,
		WebSocket: websocket,
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

func TestDeribit_PlaceOrder(t *testing.T) {
	b := newForTest()
	order, err := b.PlaceOrder("BTC-PERPETUAL",
		Buy,
		OrderTypeMarket,
		0,
		1,
		OrderReduceOnlyOption(false))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
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

func TestDeribit_SubscribeTrades(t *testing.T) {
	b := newForTest()
	b.SubscribeTrades(Market{Symbol: "BTC-PERPETUAL"}, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	select {}
}

func TestDeribit_SubscribeLevel2Snapshots(t *testing.T) {
	b := testExchange(true)
	b.SubscribeLevel2Snapshots(Market{Symbol: "BTC-PERPETUAL"}, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})

	select {}
}
