package gotrader

import "log"

// Strategy interface
type Strategy interface {
	Setup(brokers ...Broker)
	OnInit()
	OnTick()
	OnDeinit()
}

// StrategyBase Strategy base class
type StrategyBase struct {
	Brokers []Broker
	Broker  Broker
}

// Setup Setup the brokers
func (s *StrategyBase) Setup(brokers ...Broker) {
	if len(brokers) == 0 {
		log.Fatal("empty brokers")
	}
	s.Brokers = append(s.Brokers, brokers...)
	s.Broker = brokers[0]
}
