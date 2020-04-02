# CREX

[README](README.md) | [English README](README_en.md)

### CREX
A real-time quantitative trading library in Golang.

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

### Backtesting
See [@backtest](https://github.com/coinrust/crex/blob/master/examples/backtest/main.go)

### Live trading
See [@live trading](https://github.com/coinrust/crex/blob/master/examples/live/main.go)

### Main Features
* Ease of use.
* Tick-level backtesting.
* Backtesting and live-trading functionality.

### Supported Exchanges
| Exchange Name                                         | Backtesting       | Live trading      | Broker            |
| ----------------------------------------------------- |------------------ | ----------------- | ----------------- |
| ----------------------------------------------------- |------------------ | ----------------- | ----------------- |
| [BitMEX](https://www.bitmex.com/register/o0Duru)      | Yes               | Yes               | [Sim](https://github.com/coinrust/crex/tree/master/brokers/bitmex-sim-broker) / [Live](https://github.com/coinrust/crex/tree/master/brokers/bitmex-broker) |
| [Deribit](https://www.deribit.com/reg-7357.93)        | Yes               | Yes               | [Sim](https://github.com/coinrust/crex/tree/master/brokers/deribit-sim-broker) / [Live](https://github.com/coinrust/crex/tree/master/brokers/deribit-broker) |
| [Bybit](https://www.bybit.com/app/register?ref=qQggy) | No                | Yes               | [Live](https://github.com/coinrust/crex/tree/master/brokers/bybit-broker) |
| [Huobi DM](https://www.huobi.vc/zh-cn/topic/invited/?invite_code=7hzc5) | No                | Yes               | [Live](https://github.com/coinrust/crex/tree/master/brokers/huobi-broker) |
| [OKEX Futures](https://www.okex.me/join/1890951) | No                | Yes               | [Live](https://github.com/coinrust/crex/tree/master/brokers/okex-futures-broker) |

### Example
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
	accessKey := "[AccessKey]"
	secretKey := "[SecretKey]"
	broker := brokers.NewBroker(brokers.Deribit, accessKey, secretKey, true, map[string]string{})
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
* Paper trading.

### QQ group
QQ group: 932289088

### Donate

Feel free to donate:

| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| BTC     | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |
| ETH     | 0xa74eade7ea08a8c48d7de4d582fac145afc86e3d  |

### LICENSE
MIT [©coinrust](https://github.com/coinrust)