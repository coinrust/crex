package crex

// WS 事件
type WSEvent int

// WS 事件枚举
const (
	WSEventTrade WSEvent = iota + 1
	WSEventL2Snapshot
	WSEventBalance
	WSEventOrder
	WSEventPosition
	WSEventError
	WSEventDisconnected
	WSEventReconnected
)

// Market 市场信息
type Market struct {
	Symbol string // BTCUSDT(OKEX)/XBTUSD(BitMEX)/...
}

// WebSocket 代表WS连接
type WebSocket interface {
	SubscribeTrades(market Market, callback func(trades []Trade)) error
	SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error
	//SubscribeBalances(market Market, callback func(balance *Balance)) error
	SubscribeOrders(market Market, callback func(orders []Order)) error
	SubscribePositions(market Market, callback func(positions []Position)) error
}
