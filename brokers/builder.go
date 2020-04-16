package brokers

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/brokers/binancefutures"
	"github.com/coinrust/crex/brokers/bitmex"
	"github.com/coinrust/crex/brokers/bybit"
	"github.com/coinrust/crex/brokers/deribit"
	"github.com/coinrust/crex/brokers/hbdm"
	"github.com/coinrust/crex/brokers/hbdmswap"
	"github.com/coinrust/crex/brokers/okexfutures"
	"github.com/coinrust/crex/brokers/okexswap"
)

func New(name string, accessKey string, secret string, testnet bool, params map[string]string) Broker {
	switch name {
	case BinanceFutures:
		return binancefutures.New(accessKey, secret)
	case BitMEX:
		return bitmex.New(accessKey, secret, testnet)
	case Deribit:
		return deribit.New(accessKey, secret, testnet)
	case Bybit:
		return bybit.New(accessKey, secret, testnet)
	case HBDM:
		return hbdm.New(accessKey, secret, testnet)
	case HBDMSwap:
		return hbdmswap.New(accessKey, secret, testnet)
	case OKEXFutures:
		passphrase := getParamsString(params, "passphrase")
		return okexfutures.New(accessKey, secret, passphrase, testnet)
	case OKEXSwap:
		passphrase := getParamsString(params, "passphrase")
		return okexswap.New(accessKey, secret, passphrase, testnet)
	default:
		panic(fmt.Sprintf("broker error [%v]", name))
	}
}

func NewWS(name string, accessKey string, secret string, testnet bool, params map[string]string) WebSocket {
	switch name {
	case Bybit:
		return bybit.NewWS(accessKey, secret, testnet)
	case HBDM:
		return hbdm.NewWS(accessKey, secret, testnet)
	case HBDMSwap:
		return hbdmswap.NewWS(accessKey, secret, testnet)
	case OKEXFutures:
		passphrase := getParamsString(params, "passphrase")
		return okexfutures.NewWS(accessKey, secret, passphrase, testnet)
	case OKEXSwap:
		passphrase := getParamsString(params, "passphrase")
		return okexswap.NewWS(accessKey, secret, passphrase, testnet)
	default:
		panic(fmt.Sprintf("broker error [%v]", name))
	}
}

func getParamsString(params map[string]string, key string) string {
	if params == nil {
		return ""
	}
	if v, ok := params[key]; ok {
		return v
	} else {
		return ""
	}
}
