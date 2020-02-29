package gotrader

// Order 委托
type Order struct {
	ID           string      // ID
	Symbol       string      // 标
	Price        float64     // 价格
	Size         float64     // 委托数量
	AvgPrice     float64     // 平均成交价
	FilledAmount float64     // 成交数量
	Direction    Direction   // 委托方向
	Type         OrderType   // 委托类型
	PostOnly     bool        // 只做Maker选项
	ReduceOnly   bool        // 只减仓选项
	Status       OrderStatus // 委托状态
}

// IsOpen 是否活跃委托
func (o *Order) IsOpen() bool {
	return o.Status == OrderStatusCreated || o.Status == OrderStatusNew || o.Status == OrderStatusPartiallyFilled
}
