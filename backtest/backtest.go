package backtest

import (
	. "github.com/coinrust/gotrader"
	data "github.com/coinrust/gotrader/data"
	"log"
)

type Backtest struct {
	data     *data.Data
	strategy Strategy
	brokers  []Broker
	logs     LogItems
}

// NewBacktest Create backtest
// data: The data
func NewBacktest(data *data.Data, strategy Strategy, brokers []Broker) *Backtest {
	b := &Backtest{
		data:     data,
		strategy: strategy,
	}
	b.brokers = brokers
	strategy.Setup(b.brokers...)
	b.logs = LogItems{}
	return b
}

// Run Run backtest
func (b *Backtest) Run() {
	b.data.Reset()

	// Init
	b.strategy.OnInit()

	nBrokers := len(b.brokers)

	for {
		b.strategy.OnTick()

		b.addItemStats(nBrokers)

		if !b.data.Next() {
			break
		}
	}

	// Deinit
	b.strategy.OnDeinit()
}

func (b *Backtest) addItemStats(nBrokers int) {
	tick := b.data.GetTick()
	item := &LogItem{
		Time:  tick.Timestamp,
		Ask:   tick.Ask,
		Bid:   tick.Bid,
		Stats: nil,
	}
	for i := 0; i < nBrokers; i++ {
		accountSummary, err := b.brokers[i].GetAccountSummary("BTC")
		if err != nil {
			log.Fatal(err)
		}
		item.Stats = append(item.Stats, LogStats{
			Balance: accountSummary.Balance,
			Equity:  accountSummary.Equity,
		})
	}
	b.logs = append(b.logs, item)
}

// ComputeStats Calculating Backtest Statistics
func (b *Backtest) ComputeStats() (result *Stats) {
	if len(b.logs) == 0 {
		return
	}

	logs := b.logs

	n := len(logs)

	result = &Stats{}
	result.Start = logs[0].Time
	result.End = logs[n-1].Time
	result.Duration = result.End.Sub(result.Start)
	result.EntryPrice = logs[0].Price()
	result.ExitPrice = logs[n-1].Price()
	result.EntryEquity = logs[0].TotalEquity()
	result.ExitEquity = logs[n-1].TotalEquity()
	result.BaHReturn = (result.ExitPrice - result.EntryPrice) / result.EntryPrice * result.EntryEquity
	result.BaHReturnPnt = (result.ExitPrice - result.EntryPrice) / result.EntryPrice
	result.EquityReturn = result.ExitEquity - result.EntryEquity
	result.EquityReturnPnt = result.EquityReturn / result.EntryEquity

	return
}

// Plot Output backtest results
func (b *Backtest) Plot() {

}
