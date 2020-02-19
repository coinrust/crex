### GoTrader
A real-time quantitative trading/backtesting platform in Golang.

### Backtesting
See [@backtest](https://github.com/coinrust/gotrader/blob/master/examples/backtest/main.go)

### Live trading
See [@live trading](https://github.com/coinrust/gotrader/blob/master/examples/live/main.go)

### Main Features
* Ease of use.
* Tick-level backtesting.
* Backtesting and live-trading functionality.

### Supported Exchanges
| Exchange Name                                         | Backtesting       | Live trading      | Broker            |
| ----------------------------------------------------- |------------------ | ----------------- | ----------------- |
| [BitMEX](https://www.bitmex.com/register/o0Duru)      | No                | Yes               | [Live](https://github.com/coinrust/gotrader/tree/master/brokers/bitmex-broker) |
| [Deribit](https://www.deribit.com/reg-7357.93)        | Yes               | Yes               | [Sim](https://github.com/coinrust/gotrader/tree/master/brokers/deribit-sim-broker) / [Live](https://github.com/coinrust/gotrader/tree/master/brokers/deribit-broker) |
| [Bybit](https://www.bybit.com/app/register?ref=qQggy) | No                | Yes               | [Live](https://github.com/coinrust/gotrader/tree/master/brokers/bybit-broker) |

### Example
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
	broker := brokers.NewBroker(brokers.BrokerDeribit, apiKey, secretKey, true)
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
* Support backtesting for BitMEX.
* Support backtesting for Bybit.
* Paper trading.

### Donate

Feel free to donate:

| METHOD  | ADDRESS                                     |
|-------- |-------------------------------------------- |
| BTC     | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |
| ETH     | 0xa74eade7ea08a8c48d7de4d582fac145afc86e3d  |
| USDT    | 1Nk4AsGj5HEJ5csRenTUPab1sjUySCZ3Pq          |

### LICENSE
MIT [©coinrust](https://github.com/coinrust)