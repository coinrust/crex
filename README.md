# GoTrader

[README](README.md) | [English README](README_en.md)

### GoTrader
GoTrader 是一个用Golang语言开发的量化平台。支持tick级别数字币期货平台的回测和实盘。

### 标准 CSV 数据格式
* 列定界符: , (逗号)
* 换行标记: \n (LF)
* 日期时间格式: Unix 时间戳 (ms)

### 时间格式
| 列名              | 描述                             |
| ---------------- |--------------------------------- |
| t                | Unix 时间戳 (ms)                  |
| asks[0-X].price  | 卖单价(升序)                      |
| asks[0-X].amount | 卖单量                            |
| bids[0-X].price  | 买单价(降序)                      |
| bids[0-X].amount | 买单量                            |

### 样本数据示例
```csv
t,asks[0].price,asks[0].amount,asks[1].price,asks[1].amount,asks[2].price,asks[2].amount,asks[3].price,asks[3].amount,asks[4].price,asks[4].amount,asks[5].price,asks[5].amount,asks[6].price,asks[6].amount,asks[7].price,asks[7].amount,asks[8].price,asks[8].amount,asks[9].price,asks[9].amount,bids[0].price,bids[0].amount,bids[1].price,bids[1].amount,bids[2].price,bids[2].amount,bids[3].price,bids[3].amount,bids[4].price,bids[4].amount,bids[5].price,bids[5].amount,bids[6].price,bids[6].amount,bids[7].price,bids[7].amount,bids[8].price,bids[8].amount,bids[9].price,bids[9].amount
1569888000143,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,200,8307.5,1650,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,90,8300,71320,8299.5,310
1569888000285,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,200,8307.5,1650,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000307,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,200,8307.5,11010,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000309,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,200,8307.5,20370,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000406,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,8960,8307.5,11010,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000500,8304.5,7010,8305,60,8305.5,1220,8306,80,8307,200,8307.5,20370,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000522,8304.5,10270,8305,60,8305.5,1220,8306,80,8307,200,8307.5,20370,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185750,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
1569888000527,8304.5,10270,8305,60,8305.5,1220,8306,80,8307,200,8307.5,20370,8308,68260,8308.5,120000,8309,38400,8309.5,8400,8304,185010,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,8301,18000,8300.5,5090,8300,71320,8299.5,310
```

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
	s.Setup(TradeModeLiveTrading, broker)

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