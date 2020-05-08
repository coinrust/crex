package main

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/backtest"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/exchanges/deribitsim"
	"github.com/coinrust/crex/log"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() error {
	return nil
}

func (s *BasicStrategy) OnTick() error {
	ts, _ := s.Exchange.GetTime()
	log.Infof("OnTick %v", ts)
	currency := "BTC"
	symbol := "BTC-PERPETUAL"

	s.Exchange.GetBalance(currency)
	s.Exchange.GetBalance(currency)

	s.Exchange.GetOrderBook(symbol, 10)
	s.Exchange.GetOrderBook(symbol, 10)

	//s.Exchange.PlaceOrder(symbol, Buy, OrderTypeLimit, 1000.0, 10, OrderPostOnlyOption(true))

	s.Exchange.GetOpenOrders(symbol)
	s.Exchange.GetPositions(symbol)
	return nil
}

func (s *BasicStrategy) Run() error {
	return nil
}

func (s *BasicStrategy) OnExit() error {
	return nil
}

func main() {
	data := dataloader.NewCsvData("../../data-samples/deribit/deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv")
	var exchanges []ExchangeSim
	for i := 0; i < 2; i++ {
		ex := deribitsim.NewDeribitSim(data, 5.0, -0.00025, 0.00075)
		exchanges = append(exchanges, ex)
	}
	s := &BasicStrategy{}
	outputDir := "./output"
	bt := backtest.NewBacktest(data,
		"BTC",
		s,
		exchanges,
		outputDir)
	bt.Run()

	logs := bt.GetLogs()
	for _, v := range logs {
		fmt.Printf("Time: %v Price: %v Equity: %v\n", v.Time, v.Price(), v.TotalEquity())
	}

	bt.ComputeStats().PrintResult()
	//bt.Plot()
}
