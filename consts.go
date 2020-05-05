package crex

// TradeMode 策略模式
type TradeMode int

const (
	TradeModeBacktest TradeMode = iota
	TradeModePaperTrading
	TradeModeLiveTrading
)

func (t TradeMode) String() string {
	switch t {
	case TradeModeBacktest:
		return "Backtest"
	case TradeModePaperTrading:
		return "PaperTrading"
	case TradeModeLiveTrading:
		return "LiveTrading"
	default:
		return "None"
	}
}

// Direction 委托/持仓方向
type Direction int

const (
	Buy       Direction = iota // 做多
	Sell                       // 做空
	CloseBuy                   // 平多
	CloseSell                  // 平空
)

func (d Direction) String() string {
	switch d {
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return "None"
	}
}

// OrderType 委托类型
type OrderType int

const (
	OrderTypeMarket     OrderType = iota // 市价单
	OrderTypeLimit                       // 限价单
	OrderTypeStopMarket                  // 市价止损单
	OrderTypeStopLimit                   // 限价止损单
)

func (t OrderType) String() string {
	switch t {
	case OrderTypeMarket:
		return "Market"
	case OrderTypeLimit:
		return "Limit"
	case OrderTypeStopMarket:
		return "StopMarket"
	case OrderTypeStopLimit:
		return "StopLimit"
	default:
		return "None"
	}
}

// OrderStatus 委托状态
type OrderStatus int

const (
	OrderStatusCreated         OrderStatus = iota // 创建委托
	OrderStatusRejected                           // 委托被拒绝
	OrderStatusNew                                // 委托待成交
	OrderStatusPartiallyFilled                    // 委托部分成交
	OrderStatusFilled                             // 委托完全成交
	OrderStatusCancelPending                      // 委托取消
	OrderStatusCancelled                          // 委托被取消
	OrderStatusUntriggered                        // 等待触发条件委托单
	OrderStatusTriggered                          // 已触发条件单
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusCreated:
		return "Created"
	case OrderStatusRejected:
		return "Rejected"
	case OrderStatusNew:
		return "New"
	case OrderStatusPartiallyFilled:
		return "PartiallyFilled"
	case OrderStatusFilled:
		return "Filled"
	case OrderStatusCancelPending:
		return "CancelPending"
	case OrderStatusCancelled:
		return "Cancelled"
	case OrderStatusUntriggered:
		return "Untriggered"
	case OrderStatusTriggered:
		return "Triggered"
	default:
		return "None"
	}
}

// ContractType 合约类型
const (
	ContractTypeNone = ""   // Non-delivery contract 非交割合约
	ContractTypeW1   = "W1" // week 当周合约
	ContractTypeW2   = "W2" // two week 次周合约
	ContractTypeM1   = "M1" // month 月合约
	ContractTypeQ1   = "Q1" // quarter 季度合约
	ContractTypeQ2   = "Q2" // two quarter 次季度合约
)

// K线周期
const (
	PERIOD_1MIN   = "1m"
	PERIOD_3MIN   = "3m"
	PERIOD_5MIN   = "5m"
	PERIOD_15MIN  = "15m"
	PERIOD_30MIN  = "30m"
	PERIOD_60MIN  = "60m"
	PERIOD_1H     = "1h"
	PERIOD_2H     = "2h"
	PERIOD_3H     = "3h"
	PERIOD_4H     = "4h"
	PERIOD_6H     = "6h"
	PERIOD_8H     = "8h"
	PERIOD_12H    = "12h"
	PERIOD_1DAY   = "1d"
	PERIOD_3DAY   = "3d"
	PERIOD_1WEEK  = "1w"
	PERIOD_1MONTH = "1M"
	PERIOD_1YEAR  = "1y"
)
