package crex

import "net/http"

type Parameters struct {
	HttpClient *http.Client
	DebugMode  bool
	AccessKey  string
	SecretKey  string
	Passphrase string
	Testnet    bool
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

// Broker 交易所接口
type Broker interface {
	// 获取当前Broker名称
	GetName() (name string)

	// 获取账号余额
	GetBalance(currency string) (result Balance, err error)

	// 获取订单薄(OrderBook)
	GetOrderBook(symbol string, depth int) (result OrderBook, err error)

	// 获取K线数据
	// period: 数据周期. 分钟或者关键字1m(minute) 1h 1d 1w 1M(month) 1y 枚举值：1 3 5 15 30 60 120 240 360 720 "5m" "4h" "1d" ...
	GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error)

	// 设置合约类型
	// currencyPair: 交易对，如: BTC-USD(OKEX) BTC(HBDM)
	// contractType: W1,W2,Q1,Q2
	SetContractType(currencyPair string, contractType string) (err error)

	// 获取当前设置的合约ID
	GetContractID() (symbol string, err error)

	// 设置杠杆大小
	SetLeverRate(value float64) (err error)

	// 下单
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, stopPx float64, size float64,
		postOnly bool, reduceOnly bool, params map[string]interface{}) (result Order, err error)

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
	GetPositions(symbol string) (result []Position, err error)

	// 返回WebSocket对象
	WS() (ws WebSocket, err error)

	// 运行一次(回测系统调用)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}
