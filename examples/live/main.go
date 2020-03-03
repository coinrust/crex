package main

import (
	. "github.com/coinrust/gotrader"
	"github.com/coinrust/gotrader/brokers"
	"log"
	"time"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() {

}

func (s *BasicStrategy) OnTick() {
	currency := "BTC"
	symbol := "BTC-PERPETUAL"

	accountSummary, err := s.Brokers[0].GetAccountSummary(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("accountSummary: %#v", accountSummary)

	s.Brokers[0].GetOrderBook(symbol, 10)

	//s.Brokers[0].PlaceOrder(symbol, Buy, OrderTypeLimit, 1000.0, 10, true, false)

	s.Brokers[0].GetOpenOrders(symbol)
	s.Brokers[0].GetPosition(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	broker := brokers.NewBroker(brokers.Deribit, apiKey, secretKey, true)
	s := &BasicStrategy{}
	s.Setup(TradeModeLiveTrading, broker)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
