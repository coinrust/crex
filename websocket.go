package crex

// WS 事件
type WSEvent int

// WS 事件枚举
const (
	WSEventTrade WSEvent = iota + 1
	WSEventOrder
	WSEventPosition
	WSEventL2Snapshot
	WSEventError
	WSEventDisconnected
	WSEventReconnected
)

// Market 市场信息
type Market struct {
	ID     string // BTCUSDT(OKEX)/XBTUSD(BitMEx)/...
	Params string
}

// WebSocket 代表WS连接
type WebSocket interface {
	On(event WSEvent, listener interface{})
	SubscribeTrades(market Market)
	SubscribeLevel2Snapshots(market Market)
	SubscribeOrders(market Market)
	SubscribePositions(market Market)
}
