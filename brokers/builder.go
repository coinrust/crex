package brokers

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers/bitmex"
	"github.com/coinrust/crex/brokers/bybit"
	"github.com/coinrust/crex/brokers/deribit"
	"github.com/coinrust/crex/brokers/hbdm"
	"github.com/coinrust/crex/brokers/hbdmswap"
	"github.com/coinrust/crex/brokers/okexfutures"
	"github.com/coinrust/crex/brokers/okexswap"
	"log"
)

func New(brokerName string, accessKey string, secret string, testnet bool, params map[string]string) Broker {
	var baseUri string
	switch brokerName {
	case BitMEX:
		if testnet {
			baseUri = "testnet.bitmex.com"
		} else {
			baseUri = "www.bitmex.com"
		}
		return bitmex.New(baseUri, accessKey, secret)
	case Deribit:
		if testnet {
			baseUri = "wss://test.deribit.com/ws/api/v2/"
		} else {
			baseUri = "wss://www.deribit.com/ws/api/v2/"
		}
		return deribit.New(baseUri, accessKey, secret)
	case Bybit:
		if testnet {
			baseUri = "https://api-testnet.bybit.com/"
		} else {
			baseUri = "https://api.bybit.com/"
		}
		return bybit.New(baseUri, accessKey, secret)
	case HBDM:
		if testnet {
			baseUri = "https://api.btcgateway.pro"
		} else {
			baseUri = "https://api.hbdm.com"
		}
		return hbdm.New(baseUri, accessKey, secret)
	case HBDMSwap:
		if testnet {
			baseUri = "https://api.btcgateway.pro"
		} else {
			baseUri = "https://api.hbdm.com"
		}
		return hbdmswap.New(baseUri, accessKey, secret)
	case OKEXFutures:
		if testnet {
			baseUri = "https://testnet.okex.me"
		} else {
			baseUri = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseUri"]; ok {
			baseUri = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexfutures.New(baseUri, accessKey, secret, passphrase)
	case OKEXSwap:
		if testnet {
			baseUri = "https://testnet.okex.me"
		} else {
			baseUri = "https://www.okex.com"
		}
		if params == nil {
			log.Fatalf("missing params")
		}
		if v, ok := params["baseUri"]; ok {
			baseUri = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexswap.New(baseUri, accessKey, secret, passphrase)
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
		return hbdmswap.NewWS(wsURL, accessKey, secret)
	case OKEXFutures:
		wsURL := "wss://real.okex.com:8443/ws/v3"
		if v, ok := params["wsURL"]; ok {
			wsURL = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexfutures.NewWS(wsURL, accessKey, secret, passphrase)
	case OKEXSwap:
		wsURL := "wss://real.okex.com:8443/ws/v3"
		if v, ok := params["wsURL"]; ok {
			wsURL = v
		}
		var passphrase string
		if v, ok := params["passphrase"]; ok {
			passphrase = v
		} else {
			log.Fatalf("passphrase missing")
		}
		return okexswap.NewWS(wsURL, accessKey, secret, passphrase)
	default:
		panic(fmt.Sprintf("broker error [%v]", brokerName))
	}
}
