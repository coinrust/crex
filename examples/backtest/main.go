package main

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/backtest"
	"github.com/coinrust/crex/brokers/deribitsim"
	"github.com/coinrust/crex/data"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() {

}

func (s *BasicStrategy) OnTick() {
	currency := "BTC"
	symbol := "BTC-PERPETUAL"

	s.Brokers[0].GetAccountSummary(currency)
	s.Brokers[1].GetAccountSummary(currency)

	s.Brokers[0].GetOrderBook(symbol, 10)
	s.Brokers[1].GetOrderBook(symbol, 10)

	//s.Brokers[0].PlaceOrder(symbol, Buy, OrderTypeLimit, 1000.0, 10, true, false, nil)

	s.Brokers[0].GetOpenOrders(symbol)
	s.Brokers[0].GetPositions(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	data := data.NewCsvData("../../data-samples/deribit/deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv")
	var brokers []Broker
	for i := 0; i < 2; i++ {
		broker := deribitsim.New(data, 5.0, -0.00025, 0.00075)
		brokers = append(brokers, broker)
	}
	s := &BasicStrategy{}
	bt := backtest.NewBacktest(data,
		s,
		brokers)
	bt.Run()

	logs := bt.GetLogs()
	for _, v := range logs {
		fmt.Printf("Time: %v Price: %v Equity: %v\n", v.Time, v.Price(), v.TotalEquity())
	}

	bt.ComputeStats().PrintResult()
	//bt.Plot()
}
