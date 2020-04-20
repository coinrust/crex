package main

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/backtest"
	"github.com/coinrust/crex/data"
	"github.com/coinrust/crex/exchanges/deribitsim"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() {

}

func (s *BasicStrategy) OnTick() {
	currency := "BTC"
	symbol := "BTC-PERPETUAL"

	s.Exchanges[0].GetBalance(currency)
	s.Exchanges[1].GetBalance(currency)

	s.Exchanges[0].GetOrderBook(symbol, 10)
	s.Exchanges[1].GetOrderBook(symbol, 10)

	//s.Exchanges[0].PlaceOrder(symbol, Buy, OrderTypeLimit, 1000.0, 10, true, false, nil)

	s.Exchanges[0].GetOpenOrders(symbol)
	s.Exchanges[0].GetPositions(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	data := data.NewCsvData("../../data-samples/deribit/deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv")
	var exchanges []Exchange
	for i := 0; i < 2; i++ {
		broker := deribitsim.NewDeribitSim(data, 5.0, -0.00025, 0.00075)
		exchanges = append(exchanges, broker)
	}
	s := &BasicStrategy{}
	bt := backtest.NewBacktest(data,
		s,
		exchanges)
	bt.Run()

	logs := bt.GetLogs()
	for _, v := range logs {
		fmt.Printf("Time: %v Price: %v Equity: %v\n", v.Time, v.Price(), v.TotalEquity())
	}

	bt.ComputeStats().PrintResult()
	//bt.Plot()
}
