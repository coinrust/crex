package generatesim

import (
	. "github.com/coinrust/crex"
)

//im := (0.01 + sizeCurrency*0.00005) * sizeCurrency
//t.Logf("IM: %v/%v", 0, im) // 保留9位小数，四舍五入
//
//mm := (0.0055 + sizeCurrency*0.00005) * sizeCurrency

// https://www.deribit.com/pages/docs/perpetual

// 计算收益
// pnl: 收益(BTC/ETH)
// pnlUsd: 收益(USD)
func CalcPnl(side Direction, positionSize float64, entryPrice float64, exitPrice float64) (pnl float64, pnlUsd float64) {
	//side := "Short" // "Short"
	//positionSize := 3850.0
	//entryPrice := 3850.0
	//exitPrice := 3600.0
	//pnl := 0.0
	//pnlUsd := 0.0
	if positionSize == 0 {
		return
	}
	if side == Buy {
		pnl = (((entryPrice - exitPrice) / exitPrice) * (positionSize / entryPrice)) * -1
		pnlUsd = ((entryPrice - exitPrice) * (positionSize / entryPrice)) * -1
	} else if side == Sell {
		pnl = ((entryPrice - exitPrice) / exitPrice) * (positionSize / entryPrice)
		pnlUsd = (entryPrice - exitPrice) * (positionSize / entryPrice)
	}
	return
}

// 平均成交价
// total_quantity / ((quantity_1 / price_1) + (quantity_2 / price_2)) = entry_price

// 少于此保证金，则强制挂平仓单
func CalcMaintMargin(sizeCurrency float64) float64 {
	return (0.0055 + (sizeCurrency * 0.00005)) * sizeCurrency
}

// 少于此保证金，则爆仓
func CalcInitialMargin(sizeCurrency float64) float64 {
	return (0.01 + (sizeCurrency * 0.00005)) * sizeCurrency
}

// 计算保证金参数
func CalcMarginInfo(balance float64, entryPrice float64, positionSize float64) (result MarginInfo) {
	//balance := 0.05//0.05
	//entryPrice := 6500.0
	//positionSize := 6500.0
	var positionSizeC float64
	var leverage float64
	if positionSize > 0 {
		positionSizeC = positionSize / entryPrice //1.0
	}
	if balance > 0 {
		leverage = positionSizeC / balance //20.0
	}

	// BTC: (0.00575+(C6*0.00005))*C6
	// ETH: (0.01+(C6*0.000001))*C6

	// 维持保证金 Maintenance Margin
	maintMargin := (0.00575 + (positionSizeC * 0.00005)) * positionSizeC // 0.0058
	// Balance - Maint. Margin
	// This is how much equity you need to lose before maintenance margin > equity
	lose := balance - maintMargin

	//t.Logf("positionSizeC: %v", positionSizeC)
	//t.Logf("leverage: %v", leverage)
	//t.Logf("maintenanceMargin: %v", maintenanceMargin)
	//t.Logf("marginBalance: %v", marginBalance)

	//t.Logf("liquidationPriceLong: %v", liquidationPriceLong)

	// =IF((C5/C4)<C3,"No Liq Price",C4/(1-((C9*C4)/C5)))
	// C4/(1-((C9*C4)/C5))

	var liquidationPriceLong float64
	var liquidationPriceShort float64

	if positionSize > 0 {
		liquidationPriceLong = entryPrice / (((lose * entryPrice) / positionSize) + 1)

		if positionSize/entryPrice < balance {
			// 没有用到杠杆
			//t.Logf("No Liq Price")
		} else {
			liquidationPriceShort = entryPrice / (1 - ((lose * entryPrice) / positionSize))
			//t.Logf("liquidationPriceShort: %v", liquidationPriceShort)
		}
	}

	result.Leverage = leverage
	result.MaintMargin = maintMargin
	result.LiquidationPriceLong = liquidationPriceLong
	result.LiquidationPriceShort = liquidationPriceShort
	return
}
