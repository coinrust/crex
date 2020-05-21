package crex

// 现货资产
type SpotAsset struct {
	Name      string  // BTC
	Available float64 // 可用
	Frozen    float64 // 冻结
}

// SpotBalance 现货账号资产
type SpotBalance struct {
	Base  SpotAsset // 基础货币
	Quote SpotAsset // 交易的资产
}

// SpotExchange 现货交易所接口
type SpotExchange interface {

	// 获取 Exchange 名称
	GetName() (name string)

	// 获取交易所时间(ms)
	GetTime() (tm int64, err error)

	// 获取账号余额
	GetBalance(currency string) (result *SpotBalance, err error)

	// 获取订单薄(OrderBook)
	GetOrderBook(symbol string, depth int) (result *OrderBook, err error)

	// 获取K线数据
	// period: 数据周期. 分钟或者关键字1m(minute) 1h 1d 1w 1M(month) 1y 枚举值：1 3 5 15 30 60 120 240 360 720 "5m" "4h" "1d" ...
	GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error)

	// 买
	Buy(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 卖
	Sell(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 下单
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, size float64,
		opts ...PlaceOrderOption) (result *Order, err error)

	// 获取活跃委托单列表
	GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error)

	// 获取历史委托列表
	GetHistoryOrders(symbol string, opts ...OrderOption) (result []*Order, err error)

	// 获取委托信息
	GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error)

	// 撤销全部委托单
	CancelAllOrders(symbol string, opts ...OrderOption) (err error)

	// 撤销单个委托单
	CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error)
}

// SpotExchangeSim 模拟交易所接口
type SpotExchangeSim interface {
	SpotExchange

	// 设置交易撮合日志组件
	SetExchangeLogger(l ExchangeLogger)

	// 运行一次(回测系统调用)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}
