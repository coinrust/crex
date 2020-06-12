package backtest

import (
	"testing"
	"time"
)

func TestBacktest_HtmlReport(t *testing.T) {
	start, _ := time.Parse("2006-01-02 15:04:05", "2019-10-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2019-10-02 00:00:00")
	b := NewBacktest(nil,
		"BTC-USDT", start, end, nil, nil, "")
	path := `../testdata/trade_0.log`
	err := b.strategyTester.htmlReport(path)
	if err != nil {
		t.Error(err)
		return
	}
}
