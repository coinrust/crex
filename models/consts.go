package models

// Direction 委托/持仓方向
type Direction int

const (
	Buy  Direction = iota // 做多
	Sell                  // 做空
)

// OrderType 委托类型
type OrderType int

const (
	OrderTypeMarket     OrderType = iota // 市价单
	OrderTypeLimit                       // 限价单
	OrderTypeStopMarket                  // 市价止损单
	OrderTypeStopLimit                   // 限价止损单
)

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
