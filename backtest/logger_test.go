package backtest

import (
	"github.com/coinrust/crex/log"
	"testing"
)

func TestBtLogger(t *testing.T) {
	//bt := NewBacktest(nil, "", [strategy], nil)
	logger := NewBtLogger(nil,
		"../testdata/btlog.log", log.DebugLevel, false)
	defer logger.Sync()

	logger.Debug("hello", "world")
	logger.Info("hello")
	logger.Info("hello", "world")
}
