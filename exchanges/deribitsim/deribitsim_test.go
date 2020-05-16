package deribitsim

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/math"
	"github.com/coinrust/crex/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func testExchange() Exchange {
	start, _ := time.Parse("2006-01-02 15:04:05", "2019-10-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2019-10-02 00:00:00")
	SetIdGenerate(utils.NewIdGenerate(start))
	data := dataloader.NewCsvData("../../data-samples/deribit/deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv")
	data.Reset(start, end)
	ex := NewDeribitSim(data, 10000, -0.00025, 0.00075)
	ex.SetExchangeLogger(&EmptyExchangeLogger{})
	return ex
}

func TestReduceOrder(t *testing.T) {
	ex := testExchange()
	_, err := ex.PlaceOrder("BTC-PERPETUAL",
		Buy, OrderTypeMarket, 3000, 10,
		OrderReduceOnlyOption(true))
	assert.Equal(t, err, ErrInvalidAmount)
}

// 计算公式:
// https://www.reddit.com/r/DeribitExchange/
// 选择菜单: Calculation Sheets-Liquidation Price

// 一个关于计算的参考:
// https://blog.deribit.com/education/introduction-to-leverage-and-margin/

// http://www.qmwxb.com/article/2060187.html

// https://www.reddit.com/r/BitMEX/comments/7tes5v/does_anyone_knows_how_does_bitmex_calculate_the/
// total_quantity / ((quantity_1 / price_1) + (quantity_2 / price_2)) = entry_price

// Deribit:
// Rejected, maximum size of future position is $1,000,000
// 开仓总量不能大于 1000000
// Invalid size - not multiple of contract size ($10)
// 数量必须是10的整数倍
func TestUpdateBalance(t *testing.T) {
	// 5.07045231

	/*
		Order Type: Market
		BUY $10 of BTC-PERPETUAL
		Margin =  0.00000980
		Total =  0.0010
		Estimated Liquidation Price = $1.98 ↓
	*/

	/*
		Cross Liquidation Price =1/(1/Entry Price+Minimum Margin Balance/Position Amount)

		Minimum Margin Balance = Account Balance *(1- (Maintenance Margin + Taker Fees + Funding Rate))
	*/
	feeRate := 0.00075 //0.075%
	price := 10216.00
	amount := 10.0
	fee := amount * (1.0 / price) * feeRate
	// 0.00000073
	t.Logf("fee: %.10f", fee)
	fee2 := math.ToFixed(fee, 8)
	t.Logf("fee2: %.8f", fee2)
	if fee2 != 0.00000073 {
		t.Error("fee2 error")
		return
	}

	price = 10210.06
	amount = 1010.0
	leverage := 100.0

	// margin := positionSize(BTC)/leverage
	positionSizeBTC := amount * (1.0 / price)
	margin := positionSizeBTC / leverage
	t.Logf("margin: %v", margin)

	//availableBalance := equity - margin

	/*
		Order Type: Limit
		BUY $1000 of BTC-PERPETUAL
		Price = $10210.50
		Margin =  0.00097973
		Total =  0.0979
		Estimated Liquidation Price = $196.51 ↓
	*/
}

func TestCalcLiquidationPrice(t *testing.T) {
	/*
		DERIBIT BTC PERPETUAL/FUTURES
		Starting Balance BTC	0.05000000		<-- Edit the values in yellow
		Entry price $	6500
		Position size $	6500
		Position size BTC	1.00000000
		Leverage	20.00
		Maintenance Margin	0.00580000
		Balance - Maint. Margin	0.04420000		This is how much equity you need to lose before maintenance margin > equity

		Liquidation price $ (long)	6224.86		This is the price at which maintenance margin > equity and therefore liquidation would begin to trigger for a long

		Liquidation price $ (short)	6800.59		This is the price at which maintenance margin > equity and therefore liquidation would begin to trigger for a short
	*/
	balance := 0.05 //0.05
	entryPrice := 6500.0
	positionSize := 6500.0
	positionSizeC := positionSize / entryPrice //1.0
	leverage := positionSizeC / balance        //20.0

	// BTC: (0.00575+(C6*0.00005))*C6
	// ETH: (0.01+(C6*0.000001))*C6

	// 维持保证金maintenanceMargin
	maintenanceMargin := (0.00575 + (positionSizeC * 0.00005)) * positionSizeC // 0.0058
	// 保证金余额
	lose := balance - maintenanceMargin // maintenanceMargin

	// https://medium.com/deribitofficial/introduction-to-leverage-and-margin-f0045676c07

	// https://www.deribit.com/pages/docs/portfoliomargin
	// Initial margin is Maintenance Margin + 30%. Example: If Maintenance Margin is 10 BTC, Initial margin will be 10 BTC+30% = 13 BTC.

	t.Logf("positionSizeC: %v", positionSizeC)
	t.Logf("leverage: %v", leverage)
	t.Logf("maintenanceMargin: %v", maintenanceMargin)
	t.Logf("lose: %v", lose)

	liquidationPriceLong := entryPrice / (((lose * entryPrice) / positionSize) + 1)
	t.Logf("liquidationPriceLong: %v", liquidationPriceLong)

	// =IF((C5/C4)<C3,"No Liq Price",C4/(1-((C9*C4)/C5)))
	// C4/(1-((C9*C4)/C5))
	if positionSize/entryPrice < balance {
		// 没有用到杠杆
		t.Logf("No Liq Price")
	} else {
		liquidationPriceShort := entryPrice / (1 - ((lose * entryPrice) / positionSize))
		t.Logf("liquidationPriceShort: %v", liquidationPriceShort)
	}
}

func TestDiribitSim_MarginInfo(t *testing.T) {
	info := CalcMarginInfo(10, 6500, 10)
	t.Logf("%#v", info)
}

func TestOrderAmountNotMultiple(t *testing.T) {
	amount := 110.0
	iAmount := int(amount)
	a := iAmount % 10
	t.Logf("a=%v", a)
}
