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

func New(brokerName string, accessKey string, secret string, testnet bool, params map[string]string) Broker {
	var baseURL string
	switch brokerName {
	case BitMEX:
		if testnet {
			baseURL = bitmexapi.HostTestnet
		} else {
			baseURL = bitmexapi.HostReal
		}
		return bitmex.New(baseURL, accessKey, secret)
	case Deribit:
		if testnet {
			baseURL = deribitapi.TestBaseURL
		} else {
			baseURL = deribitapi.RealBaseURL
		}
		return deribit.New(baseURL, accessKey, secret)
	case Bybit:
		if testnet {
			baseURL = "https://api-testnet.bybit.com/"
		} else {
			baseURL = "https://api.bybit.com/"
		}
		return bybit.New(baseURL, accessKey, secret)
	case HBDM:
		if testnet {
			baseURL = "https://api.btcgateway.pro"
		} else {
			baseURL = "https://api.hbdm.com"
		}
		return hbdm.New(baseURL, accessKey, secret)
	case HBDMSwap:
		if testnet {
			baseURL = "https://api.btcgateway.pro"
		} else {
			baseURL = "https://api.hbdm.com"
		}
		return hbdm_swap.New(baseURL, accessKey, secret)
	case OKEXFutures:
		if testnet {
			baseURL = "https://testnet.okex.me"
		} else {
			baseURL = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseURL"]; ok {
			baseURL = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okex_futures.New(baseURL, accessKey, secret, passphrase)
	case OKEXSwap:
		if testnet {
			baseURL = "https://testnet.okex.me"
		} else {
			baseURL = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseURL"]; ok {
			baseURL = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okex_swap.New(baseURL, accessKey, secret, passphrase)
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
