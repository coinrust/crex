package crex

import "time"

// Record 表示K线数据
type Record struct {
	Symbol    string    `json:"symbol"`    // 标
	Timestamp time.Time `json:"timestamp"` // 时间
	Open      float64   `json:"open"`      // 开盘价
	High      float64   `json:"high"`      // 最高价
	Low       float64   `json:"low"`       // 最低价
	Close     float64   `json:"close"`     // 收盘价
	Volume    float64   `json:"volume"`    // 量
}

// Trade 成交记录
type Trade struct {
	ID        string    `json:"id"`     // ID
	Direction Direction `json:"type"`   // 主动成交方向
	Price     float64   `json:"price"`  // 价格
	Amount    float64   `json:"amount"` // 成交量(张)，买卖双边成交量之和
	Ts        int64     `json:"ts"`     // 订单成交时间 unix time (ms)
	Symbol    string    `json:"omitempty"`
}

// Order 委托
type Order struct {
	ID           string      `json:"id"`            // ID
	Symbol       string      `json:"symbol"`        // 标
	Price        float64     `json:"price"`         // 价格
	StopPx       float64     `json:"stop_px"`       // 触发价
	Size         float64     `json:"size"`          // 委托数量
	AvgPrice     float64     `json:"avg_price"`     // 平均成交价
	FilledAmount float64     `json:"filled_amount"` // 成交数量
	Direction    Direction   `json:"direction"`     // 委托方向
	Type         OrderType   `json:"type"`          // 委托类型
	PostOnly     bool        `json:"post_only"`     // 只做Maker选项
	ReduceOnly   bool        `json:"reduce_only"`   // 只减仓选项
	Status       OrderStatus `json:"status"`        // 委托状态
}

// IsOpen 是否活跃委托
func (o *Order) IsOpen() bool {
	return o.Status == OrderStatusCreated || o.Status == OrderStatusNew || o.Status == OrderStatusPartiallyFilled
}
