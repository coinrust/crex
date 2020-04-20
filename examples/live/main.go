package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
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

	balance, err := s.Exchange.GetBalance(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("balance: %#v", balance)

	s.Exchange.GetOrderBook(symbol, 10)

	s.Exchange.OpenLong(symbol, OrderTypeLimit, 5000, 10)
	s.Exchange.CloseLong(symbol, OrderTypeLimit, 6000, 10)

	s.Exchange.PlaceOrder(symbol,
		Buy, OrderTypeLimit, 1000.0, 10, OrderPostOnlyOption(true))

	s.Exchange.GetOpenOrders(symbol)
	s.Exchange.GetPositions(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	exchange := exchanges.NewExchange(exchanges.Deribit,
		ApiProxyURLOption("socks5://127.0.0.1:1080"), // 使用代理
		//ApiAccessKeyOption("[accessKey]"),
		//ApiSecretKeyOption("[secretKey]"),
		ApiTestnetOption(true))

	s := &BasicStrategy{}
	s.Setup(TradeModeLiveTrading, exchange)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
