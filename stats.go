package crex

import (
	"fmt"
	"time"
)

// Stats Backtesting Statistics
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
	fmt.Printf("======================== RESULT ========================\n")
	fmt.Printf("Start: \t\t\t%v\n", s.Start)
	fmt.Printf("End: \t\t\t%v\n", s.End)
	fmt.Printf("Duration: \t\t%v(min)\n", s.Duration.Minutes())
	fmt.Printf("EntryPrice: \t\t%v\n", s.EntryPrice)
	fmt.Printf("ExitPrice: \t\t%v\n", s.ExitPrice)
	fmt.Printf("EntryEquity: \t\t%v\n", s.EntryEquity)
	fmt.Printf("ExitEquity: \t\t%v\n", s.ExitEquity)
	fmt.Printf("Buy & Hold Return: \t%.8f\n", s.BaHReturn)
	fmt.Printf("Buy & Hold Return Pnt: \t%.4f%%\n", s.BaHReturnPnt*100)
	fmt.Printf("EquityReturn: \t\t%.8f\n", s.EquityReturn)
	fmt.Printf("EquityReturnPnt: \t%.4f%%\n", s.EquityReturnPnt*100)
}
