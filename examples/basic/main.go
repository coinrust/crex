package main

import (
	. "github.com/coinrust/gotrader"
	"github.com/coinrust/gotrader/backtest"
	"github.com/coinrust/gotrader/brokers/deribit-sim-broker"
	"github.com/coinrust/gotrader/data"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() {

}

func (s *BasicStrategy) OnTick() {

}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	data := data.NewDeribitData("../../deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01_sample100000.csv")
	var brokers []Broker
	for i := 0; i < 2; i++ {
		broker := deribit_sim_broker.NewBroker(data, 5.0, -0.00025, 0.00075)
		brokers = append(brokers, broker)
	}
	s := &BasicStrategy{}
	bt := backtest.NewBacktest(data,
		s,
		brokers)
	bt.Run()
	bt.ComputeStats().PrintResult()
	//bt.Plot()
}
