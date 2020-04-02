package brokers

import (
	"fmt"
	. "github.com/coinrust/crex"
	bitmexbroker "github.com/coinrust/crex/brokers/bitmex-broker"
	bybitbroker "github.com/coinrust/crex/brokers/bybit-broker"
	deribitbroker "github.com/coinrust/crex/brokers/deribit-broker"
	hbdmbroker "github.com/coinrust/crex/brokers/hbdm-broker"
	okexfuturesbroker "github.com/coinrust/crex/brokers/okex-futures-broker"
	"github.com/frankrap/bitmex-api"
	"github.com/frankrap/deribit-api"
	"log"
)

func NewBroker(brokerName string, accessKey string, secret string, testnet bool, params map[string]string) Broker {
	var addr string
	switch brokerName {
	case BitMEX:
		if testnet {
			addr = bitmex.HostTestnet
		} else {
			addr = bitmex.HostReal
		}
		return bitmexbroker.NewBroker(addr, accessKey, secret)
	case Deribit:
		if testnet {
			addr = deribit.TestBaseURL
		} else {
			addr = deribit.RealBaseURL
		}
		return deribitbroker.NewBroker(addr, accessKey, secret)
	case Bybit:
		if testnet {
			addr = "https://api-testnet.bybit.com/"
		} else {
			addr = "https://api.bybit.com/"
		}
		return bybitbroker.NewBroker(addr, accessKey, secret)
	case HBDM:
		if testnet {
			addr = "https://api.btcgateway.pro"
		} else {
			addr = "https://api.hbdm.com"
		}
		return hbdmbroker.NewBroker(addr, accessKey, secret)
	case OKEXFutures:
		if testnet {
			addr = "https://www.okex.me"
		} else {
			addr = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("passphrase missing")
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexfuturesbroker.NewBroker(addr, accessKey, secret, passphrase)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}
