package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
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

	accountSummary, err := s.Broker.GetAccountSummary(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("accountSummary: %#v", accountSummary)

	s.Broker.GetOrderBook(symbol, 10)

	s.Broker.PlaceOrder(symbol,
		Buy, OrderTypeLimit, 1000.0, 10, 1, true, false, nil)

	s.Broker.GetOpenOrders(symbol)
	s.Broker.GetPositions(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	accessKey := "[accessKey]"
	secretKey := "[secretKey]"
	broker := brokers.New(brokers.Deribit,
		accessKey, secretKey, true, nil)

	s := &BasicStrategy{}
	s.Setup(TradeModeLiveTrading, broker)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
