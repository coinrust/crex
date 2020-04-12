package brokers

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers/bitmex"
	"github.com/coinrust/crex/brokers/bybit"
	"github.com/coinrust/crex/brokers/deribit"
	"github.com/coinrust/crex/brokers/hbdm"
	"github.com/coinrust/crex/brokers/hbdm-swap"
	"github.com/coinrust/crex/brokers/okex-futures"
	"github.com/coinrust/crex/brokers/okex-swap"
	bitmexapi "github.com/frankrap/bitmex-api"
	deribitapi "github.com/frankrap/deribit-api"
	"log"
)

func NewBroker(brokerName string, accessKey string, secret string, testnet bool, params map[string]string) Broker {
	var addr string
	switch brokerName {
	case BitMEX:
		if testnet {
			addr = bitmexapi.HostTestnet
		} else {
			addr = bitmexapi.HostReal
		}
		return bitmex.New(addr, accessKey, secret)
	case Deribit:
		if testnet {
			addr = deribitapi.TestBaseURL
		} else {
			addr = deribitapi.RealBaseURL
		}
		return deribit.New(addr, accessKey, secret)
	case Bybit:
		if testnet {
			addr = "https://api-testnet.bybit.com/"
		} else {
			addr = "https://api.bybit.com/"
		}
		return bybit.New(addr, accessKey, secret)
	case HBDM:
		if testnet {
			addr = "https://api.btcgateway.pro"
		} else {
			addr = "https://api.hbdm.com"
		}
		return hbdm.New(addr, accessKey, secret)
	case HBDMSwap:
		if testnet {
			addr = "https://api.btcgateway.pro"
		} else {
			addr = "https://api.hbdm.com"
		}
		return hbdm_swap.New(addr, accessKey, secret)
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
		return okex_futures.New(addr, accessKey, secret, passphrase)
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
		return okex_swap.New(addr, accessKey, secret, passphrase)
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
		return hbdm.NewWS(wsURL, accessKey, secret)
	case HBDMSwap:
		wsURL := "wss://api.hbdm.com/swap-ws"
		if v, ok := params["wsURL"]; ok {
			wsURL = v
		}
		return hbdm_swap.NewWS(wsURL, accessKey, secret)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}
