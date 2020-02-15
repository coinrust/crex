package gotrader

import (
	"fmt"
	"time"
)

// Stats 回测的统计信息
type Stats struct {
	Start           time.Time     `json:"start"`
	End             time.Time     `json:"end"`
	Duration        time.Duration `json:"duration"`
	EntryPrice      float64       `json:"entry_price"`
	ExitPrice       float64       `json:"exit_price"`
	EntryEquity     float64       `json:"entry_equity"`
	ExitEquity      float64       `json:"exit_equity"`
	BaHReturn       float64       `json:"bah_return"`     // Buy & Hold Return
	BaHReturnPnt    float64       `json:"bah_return_pnt"` // Buy & Hold Return
	EquityReturn    float64       `json:"equity_return"`
	EquityReturnPnt float64       `json:"equity_return_pnt"`
}

func (s *Stats) PrintResult() {
	fmt.Printf("Start: %v\n", s.Start)
	fmt.Printf("End: %v\n", s.End)
	fmt.Printf("Duration: %v(min)\n", s.Duration.Minutes())
	fmt.Printf("EntryPrice: %v\n", s.EntryPrice)
	fmt.Printf("ExitPrice: %v\n", s.ExitPrice)
	fmt.Printf("EntryEquity: %v\n", s.EntryEquity)
	fmt.Printf("ExitEquity: %v\n", s.ExitEquity)
	fmt.Printf("Buy & Hold Return: %v\n", s.BaHReturn)
	fmt.Printf("Buy & Hold Return Pnt: %v%%\n", s.BaHReturnPnt*100)
	fmt.Printf("EquityReturn: %v\n", s.EquityReturn)
	fmt.Printf("EquityReturnPnt: %v%%\n", s.EquityReturnPnt*100)
}
