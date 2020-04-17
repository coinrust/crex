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

func New(name string, opts ...ApiOption) Broker {
	params := &Parameters{}

	for _, opt := range opts {
		opt(params)
	}

	return NewFromParameters(name, params)
}

func NewWS(name string, opts ...ApiOption) WebSocket {
	params := &Parameters{}

	for _, opt := range opts {
		opt(params)
	}

	return NewWSFromParameters(name, params)
}

func NewFromParameters(name string, params *Parameters) Broker {
	switch name {
	case BinanceFutures:
		return binancefutures.New(params)
	case BitMEX:
		return bitmex.New(params)
	case Deribit:
		return deribit.New(params)
	case Bybit:
		return bybit.New(params)
	case HBDM:
		return hbdm.New(params)
	case HBDMSwap:
		return hbdmswap.New(params)
	case OKEXFutures:
		return okexfutures.New(params)
	case OKEXSwap:
		return okexswap.New(params)
	default:
		panic(fmt.Sprintf("broker error [%v]", name))
	}
}

func NewWSFromParameters(name string, params *Parameters) WebSocket {
	switch name {
	case Bybit:
		return bybit.NewWS(params)
	case HBDM:
		return hbdm.NewWS(params)
	case HBDMSwap:
		return hbdmswap.NewWS(params)
	case OKEXFutures:
		return okexfutures.NewWS(params)
	case OKEXSwap:
		return okexswap.NewWS(params)
	default:
		panic(fmt.Sprintf("broker error [%v]", name))
	}
}
