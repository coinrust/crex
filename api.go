package crex

import "net/http"

type Parameters struct {
	DebugMode  bool
	HttpClient *http.Client
	ProxyURL   string // example: socks5://127.0.0.1:1080 | http://127.0.0.1:1080
	ApiURL     string
	WsURL      string
	Testnet    bool
	AccessKey  string
	SecretKey  string
	Passphrase string
	WebSocket  bool // Enable websocket option
}

type ApiOption func(p *Parameters)

func ApiDebugModeOption(debugMode bool) ApiOption {
	return func(p *Parameters) {
		p.DebugMode = debugMode
	}
}

func ApiHttpClientOption(httpClient *http.Client) ApiOption {
	return func(p *Parameters) {
		p.HttpClient = httpClient
	}
}

func ApiProxyURLOption(proxyURL string) ApiOption {
	return func(p *Parameters) {
		p.ProxyURL = proxyURL
	}
}

func ApiApiURLOption(apiURL string) ApiOption {
	return func(p *Parameters) {
		p.ApiURL = apiURL
	}
}

func ApiWsURLOption(wsURL string) ApiOption {
	return func(p *Parameters) {
		p.WsURL = wsURL
	}
}

func ApiAccessKeyOption(accessKey string) ApiOption {
	return func(p *Parameters) {
		p.AccessKey = accessKey
	}
}

func ApiSecretKeyOption(secretKey string) ApiOption {
	return func(p *Parameters) {
		p.SecretKey = secretKey
	}
}

func ApiPassPhraseOption(passPhrase string) ApiOption {
	return func(p *Parameters) {
		p.Passphrase = passPhrase
	}
}

func ApiTestnetOption(testnet bool) ApiOption {
	return func(p *Parameters) {
		p.Testnet = testnet
	}
}

func ApiWebSocketOption(enabled bool) ApiOption {
	return func(p *Parameters) {
		p.WebSocket = enabled
	}
}

type OrderParameter struct {
	Stop bool // 是否是触发委托
}

// 订单选项
type OrderOption func(p *OrderParameter)

// 触发委托选项
func OrderStopOption(stop bool) OrderOption {
	return func(p *OrderParameter) {
		p.Stop = stop
	}
}

func ParseOrderParameter(opts ...OrderOption) *OrderParameter {
	p := &OrderParameter{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type PlaceOrderParameter struct {
	BasePrice  float64
	StopPx     float64
	PostOnly   bool
	ReduceOnly bool
	PriceType  string
}

// 订单选项
type PlaceOrderOption func(p *PlaceOrderParameter)

// 基础价格选项(如: bybit 需要提供此参数)
func OrderBasePriceOption(basePrice float64) PlaceOrderOption {
	return func(p *PlaceOrderParameter) {
		p.BasePrice = basePrice
	}
}

// 触发价格选项
func OrderStopPxOption(stopPx float64) PlaceOrderOption {
	return func(p *PlaceOrderParameter) {
		p.StopPx = stopPx
	}
}

// 被动委托选项
func OrderPostOnlyOption(postOnly bool) PlaceOrderOption {
	return func(p *PlaceOrderParameter) {
		p.PostOnly = postOnly
	}
}

// 只减仓选项
func OrderReduceOnlyOption(reduceOnly bool) PlaceOrderOption {
	return func(p *PlaceOrderParameter) {
		p.ReduceOnly = reduceOnly
	}
}

// OrderPriceType 选项
func OrderPriceTypeOption(priceType string) PlaceOrderOption {
	return func(p *PlaceOrderParameter) {
		p.PriceType = priceType
	}
}

func ParsePlaceOrderParameter(opts ...PlaceOrderOption) *PlaceOrderParameter {
	p := &PlaceOrderParameter{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Exchange 交易所接口
type Exchange interface {

	// 获取 Exchange 名称
	GetName() (name string)

	// 获取交易所时间(ms)
	GetTime() (tm int64, err error)

	// 获取账号余额
	GetBalance(currency string) (result *Balance, err error)

	// 获取订单薄(OrderBook)
	GetOrderBook(symbol string, depth int) (result *OrderBook, err error)

	// 获取K线数据
	// period: 数据周期. 分钟或者关键字1m(minute) 1h 1d 1w 1M(month) 1y 枚举值：1 3 5 15 30 60 120 240 360 720 "5m" "4h" "1d" ...
	GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error)

	// 设置合约类型
	// currencyPair: 交易对，如: BTC-USD(OKEX) BTC(HBDM)
	// contractType: W1,W2,Q1,Q2
	SetContractType(currencyPair string, contractType string) (err error)

	// 获取当前设置的合约ID
	GetContractID() (symbol string, err error)

	// 设置杠杆大小
	SetLeverRate(value float64) (err error)

	// 开多
	OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 开空
	OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 平多
	CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 平空
	CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error)

	// 下单
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, size float64,
		opts ...PlaceOrderOption) (result *Order, err error)

	// 获取活跃委托单列表
	GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error)

	// 获取委托信息
	GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error)

	// 撤销全部委托单
	CancelAllOrders(symbol string, opts ...OrderOption) (err error)

	// 撤销单个委托单
	CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error)

	// 修改委托
	AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error)

	// 获取持仓
	GetPositions(symbol string) (result []*Position, err error)

	// 订阅成交记录
	SubscribeTrades(market Market, callback func(trades []*Trade)) error

	// 订阅L2 OrderBook
	SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error

	// 订阅Balance
	//SubscribeBalances(market Market, callback func(balance *Balance)) error

	// 订阅委托
	SubscribeOrders(market Market, callback func(orders []*Order)) error

	// 订阅持仓
	SubscribePositions(market Market, callback func(positions []*Position)) error

	// 运行一次(回测系统调用)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}
