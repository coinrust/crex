package brokers

import (
	"fmt"
	. "github.com/coinrust/gotrader"
	bitmex_broker "github.com/coinrust/gotrader/brokers/bitmex-broker"
	bybit_broker "github.com/coinrust/gotrader/brokers/bybit-broker"
	deribit_broker "github.com/coinrust/gotrader/brokers/deribit-broker"
	"github.com/frankrap/bitmex-api"
	"github.com/frankrap/deribit-api"
)

func NewBroker(brokerName string, apiKey string, secret string, testnet bool) Broker {
	var addr string
	switch brokerName {
	case BitMEX:
		if testnet {
			addr = bitmex.HostTestnet
		} else {
			addr = bitmex.HostReal
		}
		return bitmex_broker.NewBroker(addr, apiKey, secret)
	case Deribit:
		if testnet {
			addr = deribit.TestBaseURL
		} else {
			addr = deribit.RealBaseURL
		}
		return deribit_broker.NewBroker(addr, apiKey, secret)
	case Bybit:
		if testnet {
			addr = "https://api-testnet.bybit.com/"
		} else {
			addr = "https://api.bybit.com/"
		}
		return bybit_broker.NewBroker(addr, apiKey, secret)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}
