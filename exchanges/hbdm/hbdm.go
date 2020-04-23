package hbdm

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdm"
	"strconv"
	"strings"
	"time"
)

const StatusOK = "ok"

// Hbdm the Huobi DM exchange
type Hbdm struct {
	client        *hbdm.Client
	ws            *HbdmWebSocket
	params        *Parameters
	pair          string // 交易对 BTC/ETH/...
	_contractType string // 合约类型
	contractType  string // 合约类型(Hbdm)
	symbol        string // 合约Symbol(Hbdm) BTC_CQ
	leverRate     int    // 杠杆倍数
}

func (b *Hbdm) GetName() (name string) {
	return "hbdm"
}

func (b *Hbdm) GetBalance(currency string) (result Balance, err error) {
	var account hbdm.AccountInfoResult
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

func (b *Hbdm) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ret hbdm.MarketDepthResult

	var _type = "step0" // 使用step0时，不合并深度获取150档数据
	if depth <= 20 {
		_type = "step6" // 使用step6时，不合并深度获取20档数据
	}
	ret, err = b.client.GetMarketDepth(b.symbol, _type)
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

func (b *Hbdm) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
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
	var ret hbdm.KLineResult
	ret, err = b.client.GetKLine(b.symbol, _period, limit, from, end)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			ret.ErrCode,
			ret.ErrMsg)
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
// currencyPair: BTC/ETH/...
func (b *Hbdm) SetContractType(currencyPair string, contractType string) (err error) {
	// // 如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
	b.pair = currencyPair
	b._contractType = contractType
	var contractAlias string
	var symbol string
	switch contractType {
	case ContractTypeNone:
	case ContractTypeW1:
		contractAlias = "this_week"
		symbol = currencyPair + "_CW"
	case ContractTypeW2:
		contractAlias = "next_week"
		symbol = currencyPair + "_NW"
	case ContractTypeQ1:
		contractAlias = "quarter"
		symbol = currencyPair + "_CQ"
	}
	b.contractType = contractAlias
	b.symbol = symbol
	return
}

func (b *Hbdm) GetContractID() (symbol string, err error) {
	var ret hbdm.ContractInfoResult
	ret, err = b.client.GetContractInfo(b.pair, b.contractType, "")
	if err != nil {
		return
	}
	for _, v := range ret.Data {
		// log.Printf("%#v", v)
		if v.Symbol == b.pair &&
			v.ContractType == b.contractType {
			symbol = v.ContractCode
			return
		}
	}
	return "", fmt.Errorf("not found")
}

// 设置杠杆大小
func (b *Hbdm) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *Hbdm) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *Hbdm) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *Hbdm) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *Hbdm) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

// PlaceOrder 下单
// params:
// order_price_type: 订单报价类型
// "limit": 限价
// "opponent": 对手价
// "post_only": 只做maker单,post only下单只受用户持仓数量限制
// optimal_5：最优5档
// optimal_10：最优10档
// optimal_20：最优20档
// ioc: IOC订单
// fok：FOK订单
// "opponent_ioc"： 对手价-IOC下单
// "optimal_5_ioc"：最优5档-IOC下单
// "optimal_10_ioc"：最优10档-IOC下单
// "optimal_20_ioc"：最优20档-IOC下单
// "opponent_fok"： 对手价-FOK下单
// "optimal_5_fok"：最优5档-FOK下单
// "optimal_10_fok"：最优10档-FOK下单
// "optimal_20_fok"：最优20档-FOK下单
// -----------------------------------------------------
// 对手价下单price价格参数不用传，对手价下单价格是买一和卖一价
// optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档下单price价格参数不用传
// "limit":限价，"post_only":只做maker单 需要传价格
// "fok"：全部成交或立即取消，"ioc":立即成交并取消剩余。
func (b *Hbdm) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...OrderOption) (result Order, err error) {
	params := ParseOrderParameter(opts...)
	var orderResult hbdm.OrderResult
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
		price = 0
	}
	if params.PostOnly {
		orderPriceType = "post_only"
	}
	if params.PriceType != "" {
		orderPriceType = params.PriceType
	}
	orderResult, err = b.client.Order(
		"",
		"",
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
	//var order hbdm.OrderInfoResult
	//order, err = b.client.OrderInfo(
	//	b.pair,
	//	orderResult.Data.OrderID,
	//	0,
	//)
	//if err != nil {
	//	return
	//}
	//if order.Status != StatusOK {
	//	err = fmt.Errorf("error code=%v msg=%v",
	//		orderResult.ErrCode,
	//		orderResult.ErrMsg)
	//	return
	//}
	//if len(order.Data) != 1 {
	//	err = fmt.Errorf("missing data")
	//	return
	//}
	//result = b.convertOrder(symbol, &order.Data[0])
	return
}

func (b *Hbdm) GetOpenOrders(symbol string) (result []Order, err error) {
	var ret hbdm.OpenOrdersResult
	ret, err = b.client.GetOpenOrders(
		b.pair,
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

func (b *Hbdm) GetOrder(symbol string, id string) (result Order, err error) {
	var ret hbdm.OrderInfoResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.OrderInfo(b.pair, _id, 0)
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

func (b *Hbdm) CancelOrder(symbol string, id string) (result Order, err error) {
	var ret hbdm.CancelResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.Cancel(b.pair, _id, 0)
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

func (b *Hbdm) CancelAllOrders(symbol string) (err error) {
	return
}

func (b *Hbdm) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	return
}

func (b *Hbdm) GetPositions(symbol string) (result []Position, err error) {
	var ret hbdm.PositionInfoResult
	ret, err = b.client.GetPositionInfo(b.pair)
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

func (b *Hbdm) convertOrder(symbol string, order *hbdm.Order) (result Order) {
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

func (b *Hbdm) orderDirection(order *hbdm.Order) Direction {
	if order.Direction == "buy" {
		return Buy
	} else if order.Direction == "sell" {
		return Sell
	}
	return Buy
}

func (b *Hbdm) orderType(order *hbdm.Order) OrderType {
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

func (b *Hbdm) orderStatus(order *hbdm.Order) OrderStatus {
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

func (b *Hbdm) GetContractInfo(symbol string) (rawSymbol string, contractType string, err error) {
	var info hbdm.ContractInfoResult
	info, err = b.client.GetContractInfo("", "", symbol)
	if err != nil {
		return
	}
	if info.ErrCode != 0 {
		err = fmt.Errorf("error code=%v msg=%v",
			info.ErrCode, info.ErrMsg)
		return
	}
	for _, v := range info.Data {
		if symbol == v.ContractCode {
			rawSymbol = v.Symbol
			contractType = v.ContractType
			break
		}
	}
	if rawSymbol == "" || contractType == "" {
		err = fmt.Errorf("error")
	}
	return
}

func (b *Hbdm) SubscribeTrades(market Market, callback func(trades []Trade)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	rawSymbol, contractType, err := b.GetContractInfo(market.Symbol)
	if err != nil {
		return err
	}
	return b.ws.SubscribeTrades(rawSymbol, contractType, callback)
}

func (b *Hbdm) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	rawSymbol, contractType, err := b.GetContractInfo(market.Symbol)
	if err != nil {
		return err
	}
	return b.ws.SubscribeLevel2Snapshots(rawSymbol, contractType, callback)
}

func (b *Hbdm) SubscribeOrders(market Market, callback func(orders []Order)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	rawSymbol, contractType, err := b.GetContractInfo(market.Symbol)
	if err != nil {
		return err
	}
	return b.ws.SubscribeOrders(rawSymbol, contractType, callback)
}

func (b *Hbdm) SubscribePositions(market Market, callback func(positions []Position)) error {
	if b.ws == nil {
		return ErrWebSocketDisabled
	}
	rawSymbol, contractType, err := b.GetContractInfo(market.Symbol)
	if err != nil {
		return err
	}
	return b.ws.SubscribePositions(rawSymbol, contractType, callback)
}

func (b *Hbdm) RunEventLoopOnce() (err error) {
	return
}

func NewHbdm(params *Parameters) *Hbdm {
	baseUri := "https://api.hbdm.com"
	apiParams := &hbdm.ApiParameter{
		Debug:              false,
		AccessKey:          params.AccessKey,
		SecretKey:          params.SecretKey,
		EnablePrivateSign:  false,
		Url:                baseUri,
		PrivateKeyPrime256: "",
		HttpClient:         params.HttpClient,
		ProxyURL:           params.ProxyURL,
	}
	client := hbdm.NewClient(apiParams)
	var ws *HbdmWebSocket
	if params.WebSocket {
		ws = NewHbdmWebSocket(params)
	}
	return &Hbdm{
		client: client,
		ws:     ws,
		params: params,
	}
}
