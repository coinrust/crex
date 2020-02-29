package gotrader

type Broker interface {
	// 订阅事件
	Subscribe(event string, param string, listener interface{})

	// 获取账号信息
	GetAccountSummary(currency string) (result AccountSummary, err error)

	// 获取订单薄(OrderBook)
	GetOrderBook(symbol string, depth int) (result OrderBook, err error)

	// 下单
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, stopPx float64, size float64,
		postOnly bool, reduceOnly bool) (result Order, err error)

	// 获取活跃委托单列表
	GetOpenOrders(symbol string) (result []Order, err error)

	// 获取委托信息
	GetOrder(symbol string, id string) (result Order, err error)

	// 撤销全部委托单
	CancelAllOrders(symbol string) (err error)

	// 撤销单个委托单
	CancelOrder(symbol string, id string) (result Order, err error)

	// 修改委托
	AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error)

	// 获取持仓
	GetPosition(symbol string) (result Position, err error)

	// 运行一次(回测系统调用)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}
