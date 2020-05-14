package okexfutures

import (
	"fmt"
	"github.com/coinrust/crex/utils"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"

	. "github.com/coinrust/crex"
	"github.com/frankrap/okex-api"
)

// OkexFutures the Okex futures exchange
type OkexFutures struct {
	client *okex.Client
	ws     *FuturesWebSocket
	params *Parameters

	currencyPair  string // contract currencyPair 合约交易对
	contractType  string // contract type 合约类型
	contractAlias string // okex contract type 合约类型
	leverRate     int    // lever rate 杠杆倍数
}

func (b *OkexFutures) GetName() (name string) {
	return "okexfutures"
}

func (b *OkexFutures) GetTime() (tm int64, err error) {
	var serverTime okex.ServerTime
	serverTime, err = b.client.GetServerTime()
	if err != nil {
		return
	}
	// 1420674445.201
	tm = cast.ToInt64(cast.ToFloat64(serverTime.Epoch) * 1000)
	return
}

func (b *OkexFutures) GetBalance(currency string) (result *Balance, err error) {
	var account okex.FuturesCurrencyAccount
	account, err = b.client.GetFuturesAccountsByCurrency(currency)
	if err != nil {
		return
	}

	result = &Balance{}
	result.Equity = account.Equity
	result.Available = account.TotalAvailBalance
	result.RealizedPnl = account.RealizedPnl
	result.UnrealisedPnl = account.UnRealizedPnl

	return
}

func (b *OkexFutures) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
	params := map[string]string{}
	params["size"] = fmt.Sprintf("%v", depth) // "10"
	//params["depth"] = fmt.Sprintf("%v", 0.01) // BTC: "0.1"

	var ret okex.FuturesInstrumentBookResult
	ret, err = b.client.GetFuturesInstrumentBook(symbol, params)
	if err != nil {
		return
	}

	result = &OrderBook{}

	for _, v := range ret.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  utils.ParseFloat64(v[0]),
			Amount: utils.ParseFloat64(v[1]),
		})
	}

	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  utils.ParseFloat64(v[0]),
			Amount: utils.ParseFloat64(v[1]),
		})
	}

	// 2019-07-04T09:35:07.752Z
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", ret.Timestamp)
	result.Time = timestamp.Local()
	return
}

func (b *OkexFutures) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error) {
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
	var ret [][]string
	ret, err = b.client.GetFuturesInstrumentCandles(symbol, optional)
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
	for _, v := range ret {
		var timestamp time.Time
		timestamp, err = time.Parse(time.RFC3339, v[0]) // 2020-04-09T09:16:00.000Z
		if err != nil {
			return
		}
		records = append(records, &Record{
			Symbol:    symbol,
			Timestamp: timestamp.Local(),
			Open:      utils.ParseFloat64(v[1]),
			High:      utils.ParseFloat64(v[2]),
			Low:       utils.ParseFloat64(v[3]),
			Close:     utils.ParseFloat64(v[4]),
			Volume:    utils.ParseFloat64(v[5]),
		})
	}
	return
}

// 设置合约类型
// currencyPair: BTC-USD
// contractType: W1,W2,Q1,Q2,...
func (b *OkexFutures) SetContractType(currencyPair string, contractType string) (err error) {
	b.currencyPair = currencyPair
	b.contractType = contractType
	var contractAlias string
	switch contractType {
	case ContractTypeNone:
	case ContractTypeW1:
		contractAlias = "this_week"
	case ContractTypeW2:
		contractAlias = "next_week"
	case ContractTypeQ1:
		contractAlias = "quarter"
	case ContractTypeQ2:
		contractAlias = "bi_quarter"
	}
	b.contractAlias = contractAlias
	return
}

func (b *OkexFutures) GetContractID() (symbol string, err error) {
	var ret []okex.FuturesInstrumentsResult
	ret, err = b.client.GetFuturesInstruments()
	if err != nil {
		return
	}
	for _, v := range ret {
		//log.Printf("%v %v %#v", v.Alias, v.InstrumentId, v)
		if v.Underlying == b.currencyPair &&
			v.Alias == b.contractAlias {
			symbol = v.InstrumentId
			return
		}
	}
	return "", fmt.Errorf("not found")
}

// 设置杠杆大小
func (b *OkexFutures) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *OkexFutures) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *OkexFutures) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *OkexFutures) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *OkexFutures) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *OkexFutures) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
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
	var matchPrice int
	if params.PostOnly {
		_orderType = 1
	}
	if orderType == OrderTypeMarket {
		price = 0
		_orderType = 4
	}
	var newOrderParams okex.FuturesNewOrderParams
	newOrderParams.InstrumentId = symbol                      // "BTC-USD-190705"
	newOrderParams.Leverage = fmt.Sprintf("%v", b.leverRate)  // "10"
	newOrderParams.Type = fmt.Sprintf("%v", pType)            // "1"       // 1:开多2:开空3:平多4:平空
	newOrderParams.OrderType = fmt.Sprintf("%v", _orderType)  // "0"  // 参数填数字，0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC） 4: 市价委托
	newOrderParams.Price = fmt.Sprintf("%v", price)           // "3000.0" // 每张合约的价格
	newOrderParams.Size = fmt.Sprintf("%v", size)             // "1"       // 买入或卖出合约的数量（以张计数）
	newOrderParams.MatchPrice = fmt.Sprintf("%v", matchPrice) // "0" // 是否以对手价下单(0:不是 1:是)，默认为0，当取值为1时。price字段无效，当以对手价下单，order_type只能选择0:普通委托
	var ret okex.FuturesNewOrderResult
	var resp []byte
	resp, ret, err = b.client.FuturesOrder(newOrderParams)
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
	result = &Order{}
	result.Symbol = symbol
	result.ID = ret.OrderId
	result.Status = OrderStatusNew
	////log.Printf("%v", string(resp))
	//result, err = b.GetOrder(symbol, ret.OrderId)
	return
}

func (b *OkexFutures) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	// 6: 未完成（等待成交+部分成交）
	// 7: 已完成（撤单成功+完全成交）
	var ret okex.FuturesGetOrdersResult
	ret, err = b.client.GetFuturesOrders(symbol, 6, "", "", 100)
	if err != nil {
		return
	}
	for _, v := range ret.Orders {
		result = append(result, b.convertOrder(symbol, &v))
	}
	return
}

func (b *OkexFutures) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var ret okex.FuturesGetOrderResult
	ret, err = b.client.GetFuturesOrder(symbol, id)
	if err != nil {
		result.Symbol = symbol
		result.ID = id
		return
	}
	result = b.convertOrder(symbol, &ret)
	return
}

func (b *OkexFutures) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	var ret okex.FuturesCancelInstrumentOrderResult
	var resp []byte
	resp, ret, err = b.client.CancelFuturesInstrumentOrder(symbol, id)
	if err != nil {
		err = fmt.Errorf("%v [%v]", err, string(resp))
		return
	}
	if ret.ErrorCode != 0 {
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

func (b *OkexFutures) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	return
}

func (b *OkexFutures) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (b *OkexFutures) GetPositions(symbol string) (result []*Position, err error) {
	var ret okex.FuturesPosition
	ret, err = b.client.GetFuturesInstrumentPosition(symbol)
	if err != nil {
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("%v [%v]", err, ret)
		return
	}

	if ret.MarginMode == "crossed" { // 全仓
		for _, v := range ret.CrossPosition {
			if v.InstrumentId != symbol {
				continue
			}
			position := Position{}
			position.Symbol = symbol
			// 2019-10-08T11:56:07.922Z
			createAt, _ := time.ParseInLocation(v.CreatedAt,
				"2006-01-02T15:04:05.000Z",
				time.Local)
			if v.LongQty > 0 {
				position.Size = v.LongQty
				position.AvgPrice = v.LongAvgCost
				position.OpenTime = createAt
			} else if v.ShortQty > 0 {
				position.Size = -v.ShortQty
				position.AvgPrice = v.ShortAvgCost
				position.OpenTime = createAt
			}
			result = append(result, &position)
		}
	} else {
		for _, v := range ret.FixedPosition {
			if v.InstrumentId != symbol {
				continue
			}
			position := Position{}
			position.Symbol = symbol
			// 2019-10-08T11:56:07.922Z
			createAt, _ := time.ParseInLocation(v.CreatedAt,
				"2006-01-02T15:04:05.000Z",
				time.Local)
			if v.LongQty > 0 {
				position.Size = v.LongQty
				position.AvgPrice = v.LongAvgCost
				position.OpenTime = createAt
			} else if v.ShortQty > 0 {
				position.Size = -v.ShortQty
				position.AvgPrice = v.ShortAvgCost
				position.OpenTime = createAt
			}
			result = append(result, &position)
		}
	}
	return
}

func (b *OkexFutures) convertOrder(symbol string, order *okex.FuturesGetOrderResult) (result *Order) {
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
	if order.OrderType == 1 {
		result.PostOnly = true
	}
	if order.Type == 2 || order.Type == 3 {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *OkexFutures) orderDirection(order *okex.FuturesGetOrderResult) Direction {
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

func (b *OkexFutures) orderType(order *okex.FuturesGetOrderResult) OrderType {
	if order.OrderType == 4 {
		return OrderTypeMarket
	}
	return OrderTypeLimit
}

func (b *OkexFutures) orderStatus(order *okex.FuturesGetOrderResult) OrderStatus {
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

func (b *OkexFutures) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeTrades(market, callback)
}

func (b *OkexFutures) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeLevel2Snapshots(market, callback)
}

func (b *OkexFutures) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeOrders(market, callback)
}

func (b *OkexFutures) SubscribePositions(market Market, callback func(positions []*Position)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribePositions(market, callback)
}

func NewOkexFutures(params *Parameters) *OkexFutures {
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
	var ws *FuturesWebSocket
	if params.WebSocket {
		ws = NewFuturesWebSocket(params)
	}
	client := okex.NewClient(config)
	return &OkexFutures{
		client: client,
		ws:     ws,
		params: params,
	}
}
