package hbdmswap

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdmswap"
	"strconv"
	"strings"
	"time"
)

const StatusOK = "ok"

// HbdmSwap the Huobi DM Swap exchange
type HbdmSwap struct {
	client    *hbdmswap.Client
	ws        *SwapWebSocket
	params    *Parameters
	leverRate int // 杠杆倍数
}

func (b *HbdmSwap) GetName() (name string) {
	return "hbdmswap"
}

func (b *HbdmSwap) GetTime() (tm int64, err error) {
	var heartbeat hbdmswap.HeartbeatResult
	heartbeat, err = b.client.Heartbeat()
	if err != nil {
		return
	}
	if heartbeat.Status != StatusOK {
		err = fmt.Errorf("%v", heartbeat.Status)
		return
	}
	tm = heartbeat.Ts
	return
}

func (b *HbdmSwap) GetBalance(currency string) (result Balance, err error) {
	var account hbdmswap.AccountInfoResult
	account, err = b.client.GetAccountInfo(currency)
	if err != nil {
		return
	}

	if account.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			account.ErrCode,
			account.ErrMsg)
		return
	}

	for _, v := range account.Data {
		if v.Symbol == currency {
			result.Equity = v.MarginBalance
			result.Available = v.MarginBalance
			result.RealizedPnl = v.ProfitReal
			result.UnrealisedPnl = v.ProfitUnreal
			break
		}
	}

	return
}

func (b *HbdmSwap) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ret hbdmswap.MarketDepthResult

	var _type = "step0" // 使用step0时，不合并深度获取150档数据
	if depth <= 20 {
		_type = "step6" // 使用step6时，不合并深度获取20档数据
	}
	ret, err = b.client.GetMarketDepth(symbol, _type)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
		return
	}
	for _, v := range ret.Tick.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	for _, v := range ret.Tick.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  v[0],
			Amount: v[1],
		})
	}
	result.Time = time.Unix(0, ret.Ts*int64(time.Millisecond))
	return
}

func (b *HbdmSwap) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	var _period string
	if strings.HasSuffix(period, "m") {
		_period = period[:len(period)-1] + "min"
	} else if strings.HasSuffix(period, "h") {
		_period = period[:len(period)-1] + "hour"
	} else if strings.HasSuffix(period, "d") {
		_period = period[:len(period)-1] + "day"
	} else if strings.HasSuffix(period, "d") {
		_period = period[:len(period)-1] + "day"
	} else if strings.HasSuffix(period, "w") {
		//_period = interval[:len(period)-1]+"week"
	} else if strings.HasSuffix(period, "M") {
		_period = period[:len(period)-1] + "mon"
	} else {
		_period = period + "min"
	}
	// 1min, 5min, 15min, 30min, 60min, 4hour, 1day, 1mon
	var ret hbdmswap.KLineResult
	ret, err = b.client.GetKLine(symbol, _period, limit, from, end)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		//err = fmt.Errorf("error code=%v msg=%v",
		//	ret.ErrCode,
		//	ret.ErrMsg)
		return
	}
	for _, v := range ret.Data {
		records = append(records, Record{
			Symbol:    symbol,
			Timestamp: time.Unix(int64(v.ID), 0),
			Open:      v.Open,
			High:      v.High,
			Low:       v.Low,
			Close:     v.Close,
			Volume:    float64(v.Vol),
		})
	}
	return
}

// 设置合约类型
func (b *HbdmSwap) SetContractType(pair string, contractType string) (err error) {
	return
}

func (b *HbdmSwap) GetContractID() (symbol string, err error) {
	return "", fmt.Errorf("not found")
}

// 设置杠杆大小
func (b *HbdmSwap) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *HbdmSwap) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *HbdmSwap) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *HbdmSwap) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *HbdmSwap) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

// PlaceOrder 下单
// params:
// order_price_type: 订单报价类型
// 订单报价类型:
// "limit": 限价
// "opponent": 对手价
// "post_only":只做maker单 post only下单只受用户持仓数量限制
// optimal_5：最优5档
// optimal_10：最优10档
// optimal_20：最优20档
// "fok":FOK订单
// "ioc":IOC订单
// opponent_ioc"： 对手价-IOC下单
// "optimal_5_ioc"：最优5档-IOC下单
// "optimal_10_ioc"：最优10档-IOC下单
// "optimal_20_ioc"：最优20档-IOC下单
// "opponent_fok"： 对手价-FOK下单
// "optimal_5_fok"：最优5档-FOK下单
// "optimal_10_fok"：最优10档-FOK下单
// "optimal_20_fok"：最优20档-FOK下单
func (b *HbdmSwap) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...PlaceOrderOption) (result Order, err error) {
	params := ParsePlaceOrderParameter()
	var orderResult hbdmswap.OrderResult
	var _direction string
	var offset string
	var orderPriceType string
	if direction == Buy {
		_direction = "buy"
	} else if direction == Sell {
		_direction = "sell"
	}
	if params.ReduceOnly {
		offset = "close"
	} else {
		offset = "open"
	}
	if orderType == OrderTypeLimit {
		orderPriceType = "limit"
	} else if orderType == OrderTypeMarket {
		orderPriceType = "optimal_5"
	}
	if params.PostOnly {
		orderPriceType = "post_only"
	}
	if params.PriceType != "" {
		orderPriceType = params.PriceType
	}
	orderResult, err = b.client.Order(
		symbol,
		0,
		price,
		size,
		_direction,
		offset,
		b.leverRate,
		orderPriceType)
	if err != nil {
		return
	}
	if orderResult.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			orderResult.ErrCode,
			orderResult.ErrMsg)
		return
	}
	result.Symbol = symbol
	result.ID = fmt.Sprint(orderResult.Data.OrderID)
	result.Status = OrderStatusNew
	return
}

func (b *HbdmSwap) GetOpenOrders(symbol string, opts ...OrderOption) (result []Order, err error) {
	var ret hbdmswap.OpenOrdersResult
	ret, err = b.client.GetOpenOrders(
		symbol,
		1,
		50,
	)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
		return
	}
	for _, v := range ret.Data.Orders {
		result = append(result, b.convertOrder(symbol, &v))
	}
	return
}

func (b *HbdmSwap) GetOrder(symbol string, id string, opts ...OrderOption) (result Order, err error) {
	var ret hbdmswap.OrderInfoResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.OrderInfo(symbol, _id, 0)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
		return
	}
	if len(ret.Data) != 1 {
		err = fmt.Errorf("not found")
		return
	}
	result = b.convertOrder(symbol, &ret.Data[0])
	return
}

func (b *HbdmSwap) CancelOrder(symbol string, id string, opts ...OrderOption) (result Order, err error) {
	var ret hbdmswap.CancelResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.Cancel(symbol, _id, 0)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
		return
	}
	orderID := ret.Data.Successes
	result.ID = orderID
	return
}

func (b *HbdmSwap) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	return
}

func (b *HbdmSwap) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result Order, err error) {
	return
}

func (b *HbdmSwap) GetPositions(symbol string) (result []Position, err error) {
	var ret hbdmswap.PositionInfoResult
	ret, err = b.client.GetPositionInfo(symbol)
	if err != nil {
		return
	}

	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
		return
	}

	if len(ret.Data) < 1 {
		return
	}

	for _, v := range ret.Data {
		position := Position{}
		position.Symbol = v.Symbol
		if v.Direction == "buy" {
			position.Size = v.Volume
		} else if v.Direction == "sell" {
			position.Size = -v.Volume
		}
		position.AvgPrice = v.CostHold
		position.OpenPrice = v.CostOpen
		result = append(result, position)
	}
	return
}

func (b *HbdmSwap) convertOrder(symbol string, order *hbdmswap.Order) (result Order) {
	result.ID = order.OrderIDStr
	result.Symbol = symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Amount = order.Volume
	result.Direction = b.orderDirection(order)
	result.Type = b.orderType(order)
	result.AvgPrice = order.TradeAvgPrice
	result.FilledAmount = order.TradeVolume
	if strings.Contains(order.OrderPriceType(), "post_only") {
		result.PostOnly = true
	}
	if order.Offset == "close" {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *HbdmSwap) orderDirection(order *hbdmswap.Order) Direction {
	if order.Direction == "buy" {
		return Buy
	} else if order.Direction == "sell" {
		return Sell
	}
	return Buy
}

func (b *HbdmSwap) orderType(order *hbdmswap.Order) OrderType {
	/*
		order_price_type 订单报价类型	订单报价类型 订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制,optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档，ioc:IOC订单，fok：FOK订单
	*/

	opt := order.OrderPriceType()
	if strings.Contains(opt, "limit") {
		return OrderTypeLimit
	} else if strings.Contains(opt, "opponent") ||
		strings.Contains(opt, "optimal_5") ||
		strings.Contains(opt, "optimal_10") ||
		strings.Contains(opt, "optimal_20") {
		return OrderTypeMarket
	}
	return OrderTypeLimit
}

func (b *HbdmSwap) orderStatus(order *hbdmswap.Order) OrderStatus {
	/*
		订单状态	(1准备提交 2准备提交 3已提交 4部分成交 5部分成交已撤单
		6全部成交 7已撤单 11撤单中)
	*/
	switch order.Status {
	case 1, 2, 3:
		return OrderStatusNew
	case 4:
		return OrderStatusPartiallyFilled
	case 5:
		return OrderStatusCancelled
	case 6:
		return OrderStatusFilled
	case 7:
		return OrderStatusCancelled
	case 11:
		return OrderStatusCancelPending
	default:
		return OrderStatusCreated
	}
}

func (b *HbdmSwap) SubscribeTrades(market Market, callback func(trades []Trade)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeTrades(market, callback)
}

func (b *HbdmSwap) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeLevel2Snapshots(market, callback)
}

func (b *HbdmSwap) SubscribeOrders(market Market, callback func(orders []Order)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribeOrders(market, callback)
}

func (b *HbdmSwap) SubscribePositions(market Market, callback func(positions []Position)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	return b.ws.SubscribePositions(market, callback)
}

func (b *HbdmSwap) RunEventLoopOnce() (err error) {
	return
}

func NewHbdmSwap(params *Parameters) *HbdmSwap {
	baseUri := "https://api.hbdm.com"
	if params.ApiURL != "" {
		baseUri = params.ApiURL
	}
	apiParams := &hbdmswap.ApiParameter{
		Debug:              false,
		AccessKey:          params.AccessKey,
		SecretKey:          params.SecretKey,
		EnablePrivateSign:  false,
		Url:                baseUri,
		PrivateKeyPrime256: "",
		HttpClient:         params.HttpClient,
	}
	client := hbdmswap.NewClient(apiParams)
	var ws *SwapWebSocket
	if params.WebSocket {
		ws = NewSwapWebSocket(params)
	}
	return &HbdmSwap{
		client: client,
		ws:     ws,
		params: params,
	}
}
