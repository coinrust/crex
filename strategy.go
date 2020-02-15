package gotrader

import "log"

// 策略接口
type Strategy interface {
	Setup(brokers ...Broker)
	OnInit()
	OnTick()
	OnDeinit()
}

// 策略基类
type StrategyBase struct {
	Brokers []Broker
	Broker  Broker
}

// 设置 Brokers
func (s *StrategyBase) Setup(brokers ...Broker) {
	if len(brokers) == 0 {
		log.Fatal("empty brokers")
	}
	s.Brokers = append(s.Brokers, brokers...)
	s.Broker = brokers[0]
}
