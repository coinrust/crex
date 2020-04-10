package brokers

import (
	"fmt"
	. "github.com/coinrust/crex"
	bitmexbroker "github.com/coinrust/crex/brokers/bitmex-broker"
	bybitbroker "github.com/coinrust/crex/brokers/bybit-broker"
	deribitbroker "github.com/coinrust/crex/brokers/deribit-broker"
	hbdmbroker "github.com/coinrust/crex/brokers/hbdm-broker"
	hbdmswapbroker "github.com/coinrust/crex/brokers/hbdm-swap-broker"
	okexfuturesbroker "github.com/coinrust/crex/brokers/okex-futures-broker"
	okexswapbroker "github.com/coinrust/crex/brokers/okex-swap-broker"
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
	case HBDMSwap:
		if testnet {
			addr = "https://api.btcgateway.pro"
		} else {
			addr = "https://api.hbdm.com"
		}
		return hbdmswapbroker.NewBroker(addr, accessKey, secret)
	case OKEXFutures:
		if testnet {
			addr = "https://testnet.okex.me"
		} else {
			addr = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseURL"]; ok {
			addr = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexfuturesbroker.NewBroker(addr, accessKey, secret, passphrase)
	case OKEXSwap:
		if testnet {
			addr = "https://testnet.okex.me"
		} else {
			addr = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseURL"]; ok {
			addr = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexswapbroker.NewBroker(addr, accessKey, secret, passphrase)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}

func NewWS(brokerName string, accessKey string, secret string, testnet bool, params map[string]string) WebSocket {
	switch brokerName {
	case HBDM:
		wsURL := "wss://api.hbdm.com/ws"
		if v, ok := params["wsURL"]; ok {
			wsURL = v
		}
		return hbdmbroker.NewWS(wsURL, accessKey, secret)
	case HBDMSwap:
		wsURL := "wss://api.hbdm.com/swap-ws"
		if v, ok := params["wsURL"]; ok {
			wsURL = v
		}
		return hbdmswapbroker.NewWS(wsURL, accessKey, secret)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}
