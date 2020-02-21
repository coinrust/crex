# GoTrader

[README](README.md) | [中文文档](README_zh.md)

### GoTrader
GoTrader 是一个用Golang语言开发的量化平台。支持tick级别数字币期货平台的回测和实盘。

### 回测
示例 [@backtest](https://github.com/coinrust/gotrader/blob/master/examples/backtest/main.go)

### 实盘
示例 [@live trading](https://github.com/coinrust/gotrader/blob/master/examples/live/main.go)

### 主要特性
* 使用简单
* Tick级别回测
* 支持实盘

### 支持交易所
| 交易所                                                 | 回测               | 实盘              | Broker            |
| ----------------------------------------------------- |------------------ | ----------------- | ----------------- |
| [BitMEX](https://www.bitmex.com/register/o0Duru)      | Yes               | Yes               | [Live](https://github.com/coinrust/gotrader/tree/master/brokers/bitmex-broker) |
| [Deribit](https://www.deribit.com/reg-7357.93)        | Yes               | Yes               | [Sim](https://github.com/coinrust/gotrader/tree/master/brokers/deribit-sim-broker) / [Live](https://github.com/coinrust/gotrader/tree/master/brokers/deribit-broker) |
| [Bybit](https://www.bybit.com/app/register?ref=qQggy) | No                | Yes               | [Live](https://github.com/coinrust/gotrader/tree/master/brokers/bybit-broker) |

### 示例
```golang
package main

import (
	. "github.com/coinrust/gotrader"
	"github.com/coinrust/gotrader/brokers"
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

	accountSummary, err := s.Brokers[0].GetAccountSummary(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("accountSummary: %#v", accountSummary)

	s.Brokers[0].GetOrderBook(symbol, 10)

	//s.Brokers[0].PlaceOrder(symbol, Buy, OrderTypeLimit, 1000.0, 10, true, false)

	s.Brokers[0].GetOpenOrders(symbol)
	s.Brokers[0].GetPosition(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	broker := brokers.NewBroker(brokers.Deribit, apiKey, secretKey, true)
	s := &BasicStrategy{}
	s.Setup(broker)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
```

### TODO
* 支持 Bybit 平台的回测.
* Paper trading.

### QQ群
Coinrust QQ群: 932289088

### 捐赠

欢迎支持项目，金额随意:

| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| BTC     | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |
| ETH     | 0xa74eade7ea08a8c48d7de4d582fac145afc86e3d  |
| USDT    | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |

### LICENSE
MIT [©coinrust](https://github.com/coinrust)