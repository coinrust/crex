package hbdm_broker

import (
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/frankrap/huobi-api/hbdm"
	"strconv"
	"strings"
	"time"
)

const StatusOK = "ok"

// HBDMBroker the Huobi DM broker
type HBDMBroker struct {
	client        *hbdm.Client
	pair          string // 交易对 BTC/ETH/...
	_contractType string // 合约类型
	contractType  string // 合约类型(HBDM)
	symbol        string // 合约Symbol(HBDM) BTC_CQ
	leverRate     int    // 杠杆倍数
}

func (b *HBDMBroker) Subscribe(event string, param string, listener interface{}) {

}

func (b *HBDMBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
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
			result.Balance = v.MarginAvailable
			result.Pnl = v.ProfitReal
			break
		}
	}

	return
}

func (b *HBDMBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
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
	result.Time = time.Unix(0, ret.Ts*1e6)
	return
}

// 设置合约类型
// pair: BTC/ETH/...
func (b *HBDMBroker) SetContractType(pair string, contractType string) (err error) {
	// // 如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
	b.pair = pair
	b._contractType = contractType
	var contractAlias string
	var symbol string
	switch contractType {
	case ContractTypeNone:
	case ContractTypeW1:
		contractAlias = "this_week"
		symbol = pair + "_CW"
	case ContractTypeW2:
		contractAlias = "next_week"
		symbol = pair + "_NW"
	case ContractTypeQ1:
		contractAlias = "quarter"
		symbol = pair + "_CQ"
	}
	b.contractType = contractAlias
	b.symbol = symbol
	return
}

func (b *HBDMBroker) GetContractID() (symbol string, err error) {
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
func (b *HBDMBroker) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
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
func (b *HBDMBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool, params map[string]interface{}) (result Order, err error) {
	var orderResult hbdm.OrderResult
	var _direction string
	var offset string
	var orderPriceType string
	if direction == Buy {
		_direction = "buy"
	} else if direction == Sell {
		_direction = "sell"
	}
	if reduceOnly {
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
	if postOnly {
		orderPriceType = "post_only"
	}
	if params != nil {
		if v, ok := params["order_price_type"]; ok {
			orderPriceType = v.(string)
		}
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
	var order hbdm.OrderInfoResult
	order, err = b.client.OrderInfo(
		b.pair,
		orderResult.Data.OrderID,
		0,
	)
	if err != nil {
		return
	}
	if order.Status != StatusOK {
		err = fmt.Errorf("error code=%v msg=%v",
			orderResult.ErrCode,
			orderResult.ErrMsg)
		return
	}
	if len(order.Data) != 1 {
		err = fmt.Errorf("missing data")
		return
	}
	result = b.convertOrder(symbol, &order.Data[0])
	return
}

func (b *HBDMBroker) GetOpenOrders(symbol string) (result []Order, err error) {
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

func (b *HBDMBroker) GetOrder(symbol string, id string) (result Order, err error) {
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

func (b *HBDMBroker) CancelOrder(symbol string, id string) (result Order, err error) {
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

func (b *HBDMBroker) CancelAllOrders(symbol string) (err error) {
	return
}

func (b *HBDMBroker) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	return
}

func (b *HBDMBroker) GetPosition(symbol string) (result Position, err error) {
	result.Symbol = symbol

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

	if len(ret.Data) != 1 {
		return
	}

	position := ret.Data[0]
	if position.Direction == "buy" {
		result.Size = position.Volume
	} else if position.Direction == "sell" {
		result.Size = -position.Volume
	}
	result.AvgPrice = position.CostHold
	result.OpenPrice = position.CostOpen
	return
}

func (b *HBDMBroker) convertOrder(symbol string, order *hbdm.Order) (result Order) {
	result.ID = order.OrderIDStr
	result.Symbol = symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Size = order.Volume
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

func (b *HBDMBroker) orderDirection(order *hbdm.Order) Direction {
	if order.Direction == "buy" {
		return Buy
	} else if order.Direction == "sell" {
		return Sell
	}
	return Buy
}

func (b *HBDMBroker) orderType(order *hbdm.Order) OrderType {
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

func (b *HBDMBroker) orderStatus(order *hbdm.Order) OrderStatus {
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

func (b *HBDMBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, accessKey string, secretKey string) *HBDMBroker {
	//baseURL := "https://api.hbdm.com"
	apiParams := &hbdm.ApiParameter{
		Debug:              false,
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		EnablePrivateSign:  false,
		Url:                addr,
		PrivateKeyPrime256: "",
	}
	client := hbdm.NewClient(apiParams)
	return &HBDMBroker{
		client: client,
	}
}
