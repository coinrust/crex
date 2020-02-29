package backtest

import (
	. "github.com/coinrust/gotrader"
	data "github.com/coinrust/gotrader/data"
	"log"
	"time"
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

	for {
		b.strategy.OnTick()

		b.runEventLoopOnce()

		b.addItemStats()

		if !b.data.Next() {
			break
		}
	}

	// Deinit
	b.strategy.OnDeinit()
}

func (b *Backtest) runEventLoopOnce() {
	for _, broker := range b.brokers {
		broker.RunEventLoopOnce()
	}
}

func (b *Backtest) addItemStats() {
	ob := b.data.GetOrderBook()
	tm := ob.Time
	update := false
	timestamp := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute()+1, 0, 0, time.UTC)
	var lastItem *LogItem

	if len(b.logs) > 0 {
		lastItem = b.logs[len(b.logs)-1]
		if timestamp.Unix() == lastItem.Time.Unix() {
			update = true
			return
		}
	}
	var item *LogItem
	if update {
		item = lastItem
		item.RawTime = ob.Time
		item.Ask = ob.AskPrice()
		item.Bid = ob.BidPrice()
		item.Stats = nil
		b.fetchItemStats(item)
	} else {
		item = &LogItem{
			Time:    timestamp,
			RawTime: ob.Time,
			Ask:     ob.AskPrice(),
			Bid:     ob.BidPrice(),
			Stats:   nil,
		}
		b.fetchItemStats(item)
		b.logs = append(b.logs, item)
		//log.Printf("%v / %v", tick.Timestamp, timestamp)
	}
}

func (b *Backtest) fetchItemStats(item *LogItem) {
	n := len(b.brokers)
	for i := 0; i < n; i++ {
		accountSummary, err := b.brokers[i].GetAccountSummary("BTC")
		if err != nil {
			log.Fatal(err)
		}
		item.Stats = append(item.Stats, LogStats{
			Balance: accountSummary.Balance,
			Equity:  accountSummary.Equity,
		})
	}
}

func (b *Backtest) GetLogs() LogItems {
	return b.logs
}

// ComputeStats Calculating Backtest Statistics
func (b *Backtest) ComputeStats() (result *Stats) {
	result = &Stats{}

	if len(b.logs) == 0 {
		return
	}

	logs := b.logs

	n := len(logs)

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
