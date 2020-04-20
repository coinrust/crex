package crex

import "log"

// Strategy interface
type Strategy interface {
	Setup(mode TradeMode, exchanges ...Exchange)
	GetTradeMode() TradeMode
	OnInit()
	OnTick()
	OnDeinit()
}

// StrategyBase Strategy base class
type StrategyBase struct {
	tradeMode TradeMode
	Exchanges []Exchange
	Exchange  Exchange
}

// Setup Setups the exchanges
func (s *StrategyBase) Setup(mode TradeMode, exchanges ...Exchange) {
	if len(exchanges) == 0 {
		log.Fatal("empty exchanges")
	}
	s.tradeMode = mode
	s.Exchanges = append(s.Exchanges, exchanges...)
	s.Exchange = exchanges[0]
}

func (s *StrategyBase) GetTradeMode() TradeMode {
	return s.tradeMode
}
