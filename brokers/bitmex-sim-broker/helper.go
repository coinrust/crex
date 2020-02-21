package bitmex_sim_broker

import "github.com/coinrust/gotrader/models"

// 计算收益
// pnl: 收益(BTC/ETH)
// pnlUsd: 收益(USD)
func CalcPnl(side models.Direction, positionSize float64, entryPrice float64, exitPrice float64) (pnl float64, pnlUsd float64) {
	//side := "Short" // "Short"
	//positionSize := 3850.0
	//entryPrice := 3850.0
	//exitPrice := 3600.0
	//pnl := 0.0
	//pnlUsd := 0.0
	if positionSize == 0 {
		return
	}
	if side == models.Buy {
		pnl = (((entryPrice - exitPrice) / exitPrice) * (positionSize / entryPrice)) * -1
		pnlUsd = ((entryPrice - exitPrice) * (positionSize / entryPrice)) * -1
	} else if side == models.Sell {
		pnl = ((entryPrice - exitPrice) / exitPrice) * (positionSize / entryPrice)
		pnlUsd = (entryPrice - exitPrice) * (positionSize / entryPrice)
	}
	return
}
