package okexswap

import (
	"fmt"
	"github.com/coinrust/crex/util"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"

	. "github.com/coinrust/crex"
	"github.com/frankrap/okex-api"
)

// OkexSwap the Okex swap exchange
type OkexSwap struct {
	client *okex.Client
	ws     *SwapWebSocket
	params *Parameters

	currencyPair  string // contract currencyPair 合约交易对
	contractType  string // contract type 合约类型
	contractAlias string // okex contract type 合约类型
	leverRate     int    // lever rate 杠杆倍数
}

func (b *OkexSwap) GetName() (name string) {
	return "okexswap"
}

func (b *OkexSwap) GetTime() (tm int64, err error) {
	var serverTime okex.ServerTime
	serverTime, err = b.client.GetServerTime()
	if err != nil {
		return
	}
	// 1420674445.201
	tm = cast.ToInt64(cast.ToFloat64(serverTime.Epoch) * 1000)
	return
}

func (b *OkexSwap) GetBalance(currency string) (result Balance, err error) {
	var account okex.SwapAccount
	account, err = b.client.GetSwapAccount(currency)
	if err != nil {
		return
	}

	result.Equity, _ = strconv.ParseFloat(account.Info.Equity, 64)
	result.Available, _ = strconv.ParseFloat(account.Info.TotalAvailBalance, 64)
	result.RealizedPnl, _ = strconv.ParseFloat(account.Info.RealizedPnl, 64)
	result.UnrealisedPnl, _ = strconv.ParseFloat(account.Info.UnrealizedPnl, 64)

	return
}

func (b *OkexSwap) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	params := map[string]string{}
	params["size"] = fmt.Sprintf("%v", depth) // "10"
	//params["depth"] = fmt.Sprintf("%v", 0.01) // BTC: "0.1"

	var ret okex.SwapInstrumentDepth
	ret, err = b.client.GetSwapDepthByInstrumentId(symbol, params)
	if err != nil {
		return
	}

	for _, v := range ret.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  util.ParseFloat64(v[0]),
			Amount: util.ParseFloat64(v[1]),
		})
	}

	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  util.ParseFloat64(v[0]),
			Amount: util.ParseFloat64(v[1]),
		})
	}

	// 2019-07-04T09:35:07.752Z
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", ret.Time)
	result.Time = timestamp.Local()
	return
}

func (b *OkexSwap) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	var granularity int64
	var intervalValue string
	var intervalF int64
	if strings.HasSuffix(period, "m") {
		intervalValue = period[:len(period)-1]
		intervalF = 60
	} else if strings.HasSuffix(period, "h") {
		intervalValue = period[:len(period)-1]
		intervalF = 60 * 60
	} else if strings.HasSuffix(period, "d") {
		intervalValue = period[:len(period)-1]
		intervalF = 60 * 60 * 24
	} else if strings.HasSuffix(period, "w") {
		intervalValue = period[:len(period)-1]
		intervalF = 60 * 60 * 24 * 7
	} else if strings.HasSuffix(period, "M") {
		intervalValue = period[:len(period)-1]
		intervalF = 60 * 60 * 24 * 30
	} else if strings.HasSuffix(period, "y") {
		intervalValue = period[:len(period)-1]
		intervalF = 60 * 60 * 24 * 365
	} else {
		var i int64
		i, err = strconv.ParseInt(period, 10, 64)
		if err != nil {
			return
		}
		granularity = i * 60
	}
	if intervalValue != "" {
		var i int64
		i, err = strconv.ParseInt(intervalValue, 10, 64)
		if err != nil {
			return
		}
		granularity = i * intervalF
	}
	optional := map[string]string{}
	// 2018-06-20T02:31:00Z
	if from != 0 {
		optional["start"] = time.Unix(from, 0).UTC().Format("2006-01-02T15:04:05Z")
	}
	if end != 0 {
		optional["end"] = time.Unix(end, 0).UTC().Format("2006-01-02T15:04:05Z")
	}
	optional["granularity"] = fmt.Sprint(granularity)
	//log.Printf("%#v", optional)
	var ret *okex.SwapCandleList
	ret, err = b.client.GetSwapCandlesByInstrument(symbol, optional)
	if err != nil {
		return
	}
	/*
		timestamp	String	开始时间
		open	String	开盘价格
		high	String	最高价格
		low	String	最低价格
		close	String	收盘价格
		volume	String	交易量（张）
		currency_volume	String	按币种折算的交易量
	*/
	for _, v := range *ret {
		var timestamp time.Time
		timestamp, err = time.Parse(time.RFC3339, v[0]) // 2020-04-09T09:16:00.000Z
		if err != nil {
			return
		}
		records = append(records, Record{
			Symbol:    symbol,
			Timestamp: timestamp.Local(),
			Open:      util.ParseFloat64(v[1]),
			High:      util.ParseFloat64(v[2]),
			Low:       util.ParseFloat64(v[3]),
			Close:     util.ParseFloat64(v[4]),
			Volume:    util.ParseFloat64(v[5]),
		})
	}
	return
}

// 设置合约类型
func (b *OkexSwap) SetContractType(currencyPair string, contractType string) (err error) {
	return
}

func (b *OkexSwap) GetContractID() (symbol string, err error) {
	return
}

// 设置杠杆大小
func (b *OkexSwap) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *OkexSwap) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *OkexSwap) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *OkexSwap) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *OkexSwap) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *OkexSwap) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...PlaceOrderOption) (result *Order, err error) {
	params := ParsePlaceOrderParameter(opts...)
	var pType int
	if direction == Buy {
		if params.ReduceOnly {
			pType = 4
		} else {
			pType = 1
		}
	} else if direction == Sell {
		if params.ReduceOnly {
			pType = 3
		} else {
			pType = 2
		}
	}
	var _orderType int
	if params.PostOnly {
		_orderType = 1
	}
	var matchPrice int
	if orderType == OrderTypeMarket {
		_orderType = 4
	}
	var newOrderParams okex.BasePlaceOrderInfo
	newOrderParams.Type = fmt.Sprintf("%v", pType)            // "1"       // 1:开多2:开空3:平多4:平空
	newOrderParams.OrderType = fmt.Sprintf("%v", _orderType)  // "0"  // 参数填数字，0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
	newOrderParams.Price = fmt.Sprintf("%v", price)           // "3000.0" // 每张合约的价格
	newOrderParams.Size = fmt.Sprintf("%v", size)             // "1"       // 买入或卖出合约的数量（以张计数）
	newOrderParams.MatchPrice = fmt.Sprintf("%v", matchPrice) // "0" // 是否以对手价下单(0:不是 1:是)，默认为0，当取值为1时。price字段无效，当以对手价下单，order_type只能选择0:普通委托
	var ret okex.SwapOrderResult
	var resp []byte
	resp, ret, err = b.client.PostSwapOrder(symbol, newOrderParams)
	if err != nil {
		err = fmt.Errorf("%v [%v]", err, string(resp))
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("code: %v message: %v [%v]",
			ret.Code,
			ret.Message,
			string(resp))
		return
	}
	//log.Printf("%v", string(resp))
	//result, err = b.GetOrder(symbol, ret.OrderId)
	result = &Order{}
	result.Symbol = symbol
	result.ID = ret.OrderId
	result.Status = OrderStatusNew
	return
}

func (b *OkexSwap) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	// 6: 未完成（等待成交+部分成交）
	// 7: 已完成（撤单成功+完全成交）
	var ret *okex.SwapOrdersInfo
	paramMap := map[string]string{}
	paramMap["instrument_id"] = symbol
	paramMap["status"] = "6"
	ret, err = b.client.GetSwapOrderByInstrumentId(symbol, paramMap)
	if err != nil {
		return
	}
	for _, v := range ret.OrderInfo {
		result = append(result, b.convertOrder(symbol, &v))
	}
	return
}

func (b *OkexSwap) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var ret okex.BaseOrderInfo
	ret, err = b.client.GetSwapOrderById(symbol, id)
	if err != nil {
		return
	}
	result = b.convertOrder(symbol, &ret)
	return
}

func (b *OkexSwap) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var ret okex.SwapCancelOrderResult
	var resp []byte
	resp, ret, err = b.client.PostSwapCancelOrder(symbol, id)
	if err != nil {
		err = fmt.Errorf("%v [%v]", err, string(resp))
		return
	}
	if ret.ErrorCode != "0" {
		err = fmt.Errorf("code: %v message: %v [%v]",
			ret.ErrorCode,
			ret.ErrorMessage,
			string(resp))
		return
	}
	result = &Order{}
	result.ID = ret.OrderId
	return
}

func (b *OkexSwap) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	return
}

func (b *OkexSwap) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (b *OkexSwap) GetPositions(symbol string) (result []*Position, err error) {
	var ret okex.SwapPosition
	ret, err = b.client.GetSwapPositionByInstrument(symbol)
	if err != nil {
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("%v [%v]", err, ret)
		return
	}

	if ret.MarginMode == "crossed" || ret.MarginMode == "fixed" { // 全仓
		for _, v := range ret.Holding {
			if v.InstrumentId != symbol {
				continue
			}
			position := Position{}
			position.Symbol = symbol
			// 2019-10-08T11:56:07.922Z
			timestamp, _ := time.ParseInLocation(v.Timestamp,
				"2006-01-02T15:04:05.000Z",
				time.Local)
			if v.Side == "long" {
				position.Size, _ = strconv.ParseFloat(v.Position, 64)
				position.AvgPrice, _ = strconv.ParseFloat(v.AvgCost, 64)
				position.OpenTime = timestamp
			} else if v.Side == "short" {
				position.Size, _ = strconv.ParseFloat(v.Position, 64)
				position.AvgPrice, _ = strconv.ParseFloat(v.AvgCost, 64)
				position.OpenTime = timestamp
			}
			result = append(result, &position)
		}
	}
	return
}

func (b *OkexSwap) convertOrder(symbol string, order *okex.BaseOrderInfo) (result *Order) {
	result = &Order{}
	result.ID = order.OrderId
	result.Symbol = symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Amount = float64(order.Size)
	result.Direction = b.orderDirection(order)
	result.Type = b.orderType(order)
	result.AvgPrice = order.PriceAvg
	result.FilledAmount = order.FilledQty
	if order.OrderType == "1" {
		result.PostOnly = true
	}
	if order.Type == 2 || order.Type == 3 {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *OkexSwap) orderDirection(order *okex.BaseOrderInfo) Direction {
	// 订单类型
	//1:开多
	//2:开空
	//3:平多
	//4:平空
	if order.Type == 1 || order.Type == 4 {
		return Buy
	} else if order.Type == 2 || order.Type == 3 {
		return Sell
	}
	return Buy
}

func (b *OkexSwap) orderType(order *okex.BaseOrderInfo) OrderType {
	if order.OrderType == "4" {
		return OrderTypeMarket
	}
	return OrderTypeLimit
}

func (b *OkexSwap) orderStatus(order *okex.BaseOrderInfo) OrderStatus {
	/*
		订单状态
		-2：失败
		-1：撤单成功
		0：等待成交
		1：部分成交
		2：完全成交
		3：下单中
		4：撤单中
	*/
	switch order.State {
	case -2:
		return OrderStatusRejected
	case -1:
		return OrderStatusCancelled
	case 0:
		return OrderStatusNew
	case 1:
		return OrderStatusPartiallyFilled
	case 2:
		return OrderStatusFilled
	case 3:
		return OrderStatusNew
	case 4:
		return OrderStatusCancelPending
	default:
		return OrderStatusCreated
	}
}

func (b *OkexSwap) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeTrades(market, callback)
}

func (b *OkexSwap) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeLevel2Snapshots(market, callback)
}

func (b *OkexSwap) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeOrders(market, callback)
}

func (b *OkexSwap) SubscribePositions(market Market, callback func(positions []*Position)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribePositions(market, callback)
}

func (b *OkexSwap) RunEventLoopOnce() (err error) {
	return
}

func NewOkexSwap(params *Parameters) *OkexSwap {
	var baseUri string
	if params.ApiURL != "" {
		baseUri = params.ApiURL
	} else {
		baseUri = "https://www.okex.com"
		if params.Testnet {
			baseUri = "https://testnet.okex.com"
		}
	}
	config := okex.Config{
		Endpoint:      baseUri,
		WSEndpoint:    "",
		ApiKey:        params.AccessKey,
		SecretKey:     params.SecretKey,
		Passphrase:    params.Passphrase,
		TimeoutSecond: 45,
		IsPrint:       false,
		I18n:          okex.ENGLISH,
		ProxyURL:      params.ProxyURL,
		HTTPClient:    params.HttpClient,
	}
	client := okex.NewClient(config)
	var ws *SwapWebSocket
	if params.WebSocket {
		ws = NewSwapWebSocket(params)
	}
	return &OkexSwap{
		client: client,
		ws:     ws,
		params: params,
	}
}
