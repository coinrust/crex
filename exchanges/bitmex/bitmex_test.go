package bitmex

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"log"
	"testing"
	"time"
)

func testExchange() *BitMEX {
	return testExchangeWS(false)
}

func testExchangeWS(websocket bool) *BitMEX {
	testConfig := configtest.LoadTestConfig("bitmex")
	params := &Parameters{
		DebugMode: true,
		AccessKey: testConfig.AccessKey,
		SecretKey: testConfig.SecretKey,
		Testnet:   testConfig.Testnet,
		ProxyURL:  testConfig.ProxyURL,
		WebSocket: websocket,
	}
	ex := NewBitMEX(params)
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
		10,
		OrderPostOnlyOption(true))
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

func TestBitMEX_SubscribeTrades(t *testing.T) {
	ex := testExchangeWS(true)
	ex.SubscribeTrades(Market{Symbol: "XBTUSD"}, func(trades []Trade) {
		log.Printf("trades=%#v", trades)
	})

	select {}
}

func TestBitMEX_SubscribeLevel2Snapshots(t *testing.T) {
	ex := testExchangeWS(true)
	ex.SubscribeLevel2Snapshots(Market{Symbol: "XBTUSD"}, func(ob *OrderBook) {
		log.Printf("ob=%#v", ob)
	})

	select {}
}

func TestBitMEX_SubscribeOrders(t *testing.T) {
	ex := testExchangeWS(true)
	ex.SubscribeOrders(Market{Symbol: "XBTUSD"}, func(orders []*Order) {
		log.Printf("orders=%#v", orders)
	})

	select {}
}
