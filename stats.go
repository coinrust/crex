package crex

import (
	"fmt"
	"time"
)

// Stats Back testing Statistics
type Stats struct {
	Start           time.Time     `json:"start"`
	End             time.Time     `json:"end"`
	Duration        time.Duration `json:"duration"`
	RunDuration     time.Duration `json:"run_duration"`
	EntryPrice      float64       `json:"entry_price"`
	ExitPrice       float64       `json:"exit_price"`
	EntryEquity     float64       `json:"entry_equity"`
	ExitEquity      float64       `json:"exit_equity"`
	BaHReturn       float64       `json:"bah_return"`     // Buy & Hold Return
	BaHReturnPnt    float64       `json:"bah_return_pnt"` // Buy & Hold Return
	EquityReturn    float64       `json:"equity_return"`
	EquityReturnPnt float64       `json:"equity_return_pnt"`
	AnnReturn       float64       `json:"ann_return"`    // 年化收益率
	MaxDrawDown     float64       `json:"max_draw_down"` // 最大回撤
}

func (s *Stats) PrintResult() {
	fmt.Printf("======================== RESULT ========================\n")
	fmt.Printf("Start: \t\t\t\t%v\n", s.Start)
	fmt.Printf("End: \t\t\t\t%v\n", s.End)
	fmt.Printf("Duration: \t\t\t%v\n", s.Duration.String())
	fmt.Printf("Run Duration: \t\t%v\n", s.RunDuration.String())
	fmt.Printf("Entry Price: \t\t%v\n", s.EntryPrice)
	fmt.Printf("Exit Price: \t\t%v\n", s.ExitPrice)
	fmt.Printf("Initial Equity: \t%v\n", s.EntryEquity)
	fmt.Printf("Exit Equity: \t\t%v\n", s.ExitEquity)
	fmt.Printf("Return: \t\t\t%.8f\n", s.EquityReturn)
	fmt.Printf("Return [%%]: \t\t%.4f%%\n", s.EquityReturnPnt*100)
	fmt.Printf("Buy & Hold Return: \t%.8f\n", s.BaHReturn)
	fmt.Printf("Buy & Hold Return [%%]: \t%.4f%%\n", s.BaHReturnPnt*100)
	fmt.Printf("Ann Return [%%]: \t\t%.4f%%\n", s.AnnReturn*100)
	fmt.Printf("Max Drawdown [%%]: \t\t%.4f%%\n", s.MaxDrawDown*100)
}
