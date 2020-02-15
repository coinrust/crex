package backtest

import (
	. "github.com/coinrust/gotrader"
	data "github.com/coinrust/gotrader/data"
	"log"
)

type Backtest struct {
	data     *data.Data
	strategy Strategy
	//cash     []float64
	brokers []Broker
	logs    LogItems
}

// NewBacktest 创建回测
// data: 数据
// cash: 初始资金,多个账号
//func NewBacktest(data *data.Data, strategy Strategy, cash []float64, makerFeeRate float64, takerFeeRate float64) *Backtest {
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

// Run 运行回测
func (b *Backtest) Run() {
	b.data.Reset()

	// 初始化
	b.strategy.OnInit()

	nBrokers := len(b.brokers)

	for {
		b.strategy.OnTick()

		b.addItemStats(nBrokers)

		if !b.data.Next() {
			break
		}
	}

	// 完成
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

// ComputeStats 计算回测的统计信息
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

// Plot 输出回测结果
func (b *Backtest) Plot() {

}
