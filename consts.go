package gotrader

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
	Buy  Direction = iota // 做多
	Sell                  // 做空
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
