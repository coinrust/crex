package crex

import "log"

// Strategy interface
type Strategy interface {
	Setup(mode TradeMode, brokers ...Broker)
	GetTradeMode() TradeMode
	OnInit()
	OnTick()
	OnDeinit()
}

// StrategyBase Strategy base class
type StrategyBase struct {
	tradeMode TradeMode
	Brokers   []Broker
	Broker    Broker
}

// Setup Setup the brokers
func (s *StrategyBase) Setup(mode TradeMode, brokers ...Broker) {
	if len(brokers) == 0 {
		log.Fatal("empty brokers")
	}
	s.tradeMode = mode
	s.Brokers = append(s.Brokers, brokers...)
	s.Broker = brokers[0]
}

func (s *StrategyBase) GetTradeMode() TradeMode {
	return s.tradeMode
}
