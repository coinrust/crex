<div align=center><img src="https://raw.githubusercontent.com/coinrust/crex/master/images/logo.png" /></div>

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

**CREX** A real-time quantitative trading library in Golang.

### Backtesting
See [@backtest](https://github.com/coinrust/crex/blob/master/examples/backtest/main.go)

### Live trading
See [@live trading](https://github.com/coinrust/crex/blob/master/examples/live/main.go)

### Main Features
* Ease of use.
* Tick-level backtesting.
* WebSocket Supported

### Supported Exchanges
The CREX library currently supports the following 8 cryptocurrency exchange markets and trading APIs:

| logo                                                                                                                                             | id             | name                                                                      | ver | ws  | doc                                                               |
| ------------------------------------------------------------------------------------------------------------------------------------------------ | -------------- | ------------------------------------------------------------------------- | --- | --- | ----------------------------------------------------------------- |
| [![binance](https://raw.githubusercontent.com/coinrust/crex/master/images/binance.jpg)](https://www.binance.com/en/register?ref=10916733)        | binancefutures | [Binance Futures](https://www.binance.com/en/register?ref=10916733)       | 1   | N   | [API](https://binance-docs.github.io/apidocs/futures/en/)         |
| [![bitmex](https://raw.githubusercontent.com/coinrust/crex/master/images/bitmex.jpg)](https://www.bitmex.com/register/o0Duru)                    | bitmex         | [BitMEX](https://www.bitmex.com/register/o0Duru)                          | 1   | N   | [API](https://www.bitmex.com/app/apiOverview)                     |
| [![deribit](https://raw.githubusercontent.com/coinrust/crex/master/images/deribit.jpg)](https://www.deribit.com/reg-7357.93)                     | deribit        | [Deribit](https://www.deribit.com/reg-7357.93)                            | 2   | Y   | [API](https://docs.deribit.com/)                                  |
| [![bybit](https://raw.githubusercontent.com/coinrust/crex/master/images/bybit.jpg)](https://www.bybit.com/app/register?ref=qQggy)                | bybit          | [Bybit](https://www.bybit.com/app/register?ref=qQggy)                     | 2   | Y   | [API](https://bybit-exchange.github.io/docs/inverse/)             |
| [![huobi](https://raw.githubusercontent.com/coinrust/crex/master/images/huobi.jpg)](https://www.huobi.io/en-us/topic/invited/?invite_code=7hzc5) | hbdm           | [Huobi DM](https://www.huobi.io/en-us/topic/invited/?invite_code=7hzc5)   | 1   | Y   | [API](https://docs.huobigroup.com/docs/dm/v1/en/)                 |
| [![huobi](https://raw.githubusercontent.com/coinrust/crex/master/images/huobi.jpg)](https://www.huobi.io/en-us/topic/invited/?invite_code=7hzc5) | hbdmswap       | [Huobi Swap](https://www.huobi.io/en-us/topic/invited/?invite_code=7hzc5) | 1   | Y   | [API](https://docs.huobigroup.com/docs/coin_margined_swap/v1/en/) |
| [![okex](https://raw.githubusercontent.com/coinrust/crex/master/images/okex.jpg)](https://www.okex.com/join/1890951)                             | okexfutures    | [OKEX Futures](https://www.okex.com/join/1890951)                         | 3   | Y   | [API](https://www.okex.me/docs/en/#futures-README)                |
| [![okex](https://raw.githubusercontent.com/coinrust/crex/master/images/okex.jpg)](https://www.okex.com/join/1890951)                             | okexswap       | [OKEX Swap](https://www.okex.com/join/1890951)                            | 3   | Y   | [API](https://www.okex.me/docs/en/#swap-README)                   |

### Example
```golang
package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
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

	balance, err := s.Exchange.GetBalance(currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("balance: %#v", balance)

	s.Exchange.GetOrderBook(symbol, 10)

	s.Exchange.OpenLong(symbol, OrderTypeLimit, 5000, 10)
	s.Exchange.CloseLong(symbol, OrderTypeLimit, 6000, 10)

	s.Exchange.PlaceOrder(symbol,
		Buy, OrderTypeLimit, 1000.0, 10, OrderPostOnlyOption(true))

	s.Exchange.GetOpenOrders(symbol)
	s.Exchange.GetPositions(symbol)
}

func (s *BasicStrategy) OnDeinit() {

}

func main() {
	exchange := exchanges.NewExchange(exchanges.Deribit,
		ApiProxyURLOption("socks5://127.0.0.1:1080"), // 使用代理
		//ApiAccessKeyOption("[accessKey]"),
		//ApiSecretKeyOption("[secretKey]"),
		ApiTestnetOption(true))

	s := &BasicStrategy{}
	s.Setup(TradeModeLiveTrading, exchange)

	// run loop
	for {
		s.OnTick()
		time.Sleep(1 * time.Second)
	}
}
```

## WebSocket example
```golang
package main

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
	"log"
)

func main() {
	ws := exchanges.NewExchange(exchanges.OkexFutures,
		ApiProxyURLOption("socks5://127.0.0.1:1080"), // 使用代理
		//ApiAccessKeyOption("[accessKey]"),
		//ApiSecretKeyOption("[secretKey]"),
		//ApiPassPhraseOption("[passphrase]"),
		ApiWebSocketOption(true)) // 开启 WebSocket

	market := Market{
		Symbol: "BTC-USD-200626",
	}
	// 订阅订单薄
	ws.SubscribeLevel2Snapshots(market, func(ob *OrderBook) {
		log.Printf("%#v", ob)
	})
	// 订阅成交记录
	ws.SubscribeTrades(market, func(trades []Trade) {
		log.Printf("%#v", trades)
	})
	// 订阅订单成交信息
	ws.SubscribeOrders(market, func(orders []Order) {
		log.Printf("%#v", orders)
	})
	// 订阅持仓信息
	ws.SubscribePositions(market, func(positions []Position) {
		log.Printf("%#v", positions)
	})

	select {}
}
```

### Standard CSV data types formats
* columns delimiter: , (comma)
* new line marker: \n (LF)
* date time format: Unix time (ms)

### Data format
| column name      | description                      |
| ---------------- |--------------------------------- |
| t                | Unix time (ms)                   |
| asks[0-X].price  | asks prices in ascending order   |
| asks[0-X].amount | asks amounts in ascending order  |
| bids[0-X].price  | bids prices in descending order  |
| bids[0-X].amount | bids amounts in descending order |

### Sample data rows preview
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

### TODO
* Paper trading.

### QQ group
QQ group: [932289088](https://jq.qq.com/?_wv=1027&k=5rg0FEK)

### Donate
| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| BTC     | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |
| ETH     | 0xa74eade7ea08a8c48d7de4d582fac145afc86e3d  |

### LICENSE
MIT [©coinrust](https://github.com/coinrust)