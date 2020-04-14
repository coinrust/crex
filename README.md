<div align=center><img src="https://github.com/coinrust/crex/raw/master/images/logo.png" /></div>

<p align="center">
  <a href="https://github.com/golang/go">
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/coinrust/crex">
  </a>

  <a href="https://github.com/coinrust/crex/master/LICENSE">
    <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
  <a href="https://www.travis-ci.com/coinrust/crex">
    <img src="https://www.travis-ci.com/coinrust/crex.svg?branch=master" alt="build status">
  </a>
</p>

# CREX

[中文](README.md) | [English](README_en.md)

**CREX** 是一个用Golang语言开发的量化交易库。支持`tick`级别数字币期货平台的回测和实盘。实盘与回测无缝切换，无需更改代码。

## 回测
示例 [@backtest](https://github.com/coinrust/crex/blob/master/examples/backtest/main.go)

## 实盘
示例 [@live trading](https://github.com/coinrust/crex/blob/master/examples/live/main.go)

## 主要特性
* 使用简单
* Tick级别回测
* 支持 WebSocket

## 支持交易所
| name                                                                      | id             | ver | doc                                                               |
| ------------------------------------------------------------------------- | -------------- | --- | ----------------------------------------------------------------- |
| [Binance Futures](https://www.binance.com/cn/register?ref=10916733)       | binancefutures | 1   | [API](https://binance-docs.github.io/apidocs/futures/cn/)         |
| [BitMEX](https://www.bitmex.com/register/o0Duru)                          | bitmex         | 1   | [API](https://www.bitmex.com/app/apiOverview)                     |
| [Deribit](https://www.deribit.com/reg-7357.93)                            | deribit        | 2   | [API](https://docs.deribit.com/)                                  |
| [Bybit](https://www.bybit.com/app/register?ref=qQggy)                     | bybit          | 2   | [API](https://bybit-exchange.github.io/docs/inverse/)             |
| [Huobi DM](https://www.huobi.io/zh-cn/topic/invited/?invite_code=7hzc5)   | hbdm           | 1   | [API](https://docs.huobigroup.com/docs/dm/v1/cn/)                 |
| [Huobi Swap](https://www.huobi.io/zh-cn/topic/invited/?invite_code=7hzc5) | hbdmswap       | 1   | [API](https://docs.huobigroup.com/docs/coin_margined_swap/v1/cn/) |
| [OKEX Futures](https://www.okex.com/join/1890951)                         | okexfutures    | 3   | [API](https://www.okex.me/docs/zh/#futures-README)                |
| [OKEX Swap](https://www.okex.com/join/1890951)                            | okexswap       | 3   | [API](https://www.okex.me/docs/zh/#swap-README)                   |

## 示例
```golang
package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
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

	accountSummary, err := s.Broker.GetAccountSummary(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("accountSummary: %#v", accountSummary)

	s.Broker.GetOrderBook(symbol, 10)

	s.Broker.PlaceOrder(symbol,
		Buy, OrderTypeLimit, 1000.0, 10, 1, true, false, nil)

	s.Broker.GetOpenOrders(symbol)
	s.Broker.GetPosition(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	accessKey := "[accessKey]"
	secretKey := "[secretKey]"
	broker := brokers.New(brokers.Deribit,
		accessKey, secretKey, true, map[string]string{})

	s := &BasicStrategy{}
	s.Setup(TradeModeLiveTrading, broker)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
```

## WebSocket 示例
```golang
package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers"
	"log"
)

func main() {
	wsURL := "wss://api.hbdm.com/ws" // "wss://api.btcgateway.pro/ws"
	params := map[string]string{}
	params["wsURL"] = wsURL

	ws := brokers.NewWS(brokers.HBDM,
		"[accessKey]", "[secretKey]", false, params)

	// 订单薄事件方法
	ws.On(WSEventL2Snapshot, func(ob *OrderBook) {
		log.Printf("ob: %#v", ob)
	})
	// 成交记录事件方法
	ws.On(WSEventTrade, func(trades []Trade) {
		log.Printf("trades: %#v", trades)
	})

	// 订单事件方法
	ws.On(WSEventOrder, func(order *Order) {
		log.Printf("order: %#v", order)
	})
	// 持仓事件方法
	ws.On(WSEventPosition, func(position *Position) {
		log.Printf("position: %#v", position)
	})

	// 订阅订单薄
	ws.SubscribeLevel2Snapshots(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅成交记录
	ws.SubscribeTrades(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅订单成交信息
	ws.SubscribeOrders(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})
	// 订阅持仓信息
	ws.SubscribePositions(Market{
		ID:     "BTC",
		Params: ContractTypeW1,
	})

	select {}
}
```

## 回测数据
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

## TODO
* Paper trading.

## QQ群
QQ群: [932289088](https://jq.qq.com/?_wv=1027&k=5rg0FEK)

## 捐赠

欢迎支持项目，金额随意:

| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| BTC     | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |
| ETH     | 0xa74eade7ea08a8c48d7de4d582fac145afc86e3d  |

## LICENSE
MIT [©coinrust](https://github.com/coinrust)