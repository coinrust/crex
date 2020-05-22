package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/backtest"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/exchanges/exsim"
	"github.com/coinrust/crex/log"
	"time"
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
	loader := dataloader.NewMongoDBDataLoader("mongodb://localhost:27017",
		"tick_db", "deribit", "BTC-PERPETUAL")
	data := dataloader.NewData(loader)

	start, _ := time.Parse("2006-01-02 15:04:05", "2019-10-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2019-10-02 00:00:00")

	var datas []*dataloader.Data
	var exchanges []ExchangeSim
	for i := 0; i < 2; i++ {
		datas = append(datas, data)
		ex := exsim.NewExSim(data,
			5.0, -0.00025, 0.00075, false, false)
		exchanges = append(exchanges, ex)
	}

	s := &BasicStrategy{}
	outputDir := "./output"
	bt := backtest.NewBacktest(datas,
		"BTC",
		start,
		end,
		s,
		exchanges,
		outputDir)
	bt.Run()

	//logs := bt.GetLogs()
	//for _, v := range logs {
	//	fmt.Printf("Time: %v Price: %v Equity: %v\n", v.Time, v.Price(), v.TotalEquity())
	//}

	bt.ComputeStats().PrintResult()
	bt.Plot()
	bt.HtmlReport()
}
