package generatesim

import (
	. "github.com/coinrust/crex"
)

func CalcPnl(side Direction, positionSize float64, entryPrice float64, exitPrice float64, isForwardContract bool) (pnl float64) {
	if positionSize == 0 {
		return
	}
	if isForwardContract {
		if side == Buy {
			pnl = positionSize * (exitPrice - entryPrice)
		} else if side == Sell {
			pnl = positionSize * (entryPrice - exitPrice)
		}
	} else {
		if side == Buy {
			pnl = positionSize * (1/entryPrice - 1/exitPrice)
		} else if side == Sell {
			pnl = positionSize * (1/exitPrice - 1/entryPrice)
		}
	}

	return
}
