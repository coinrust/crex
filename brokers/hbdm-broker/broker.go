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

// HuobiBroker the Huobi broker
type HuobiBroker struct {
	client       *hbdm.Client
	contractType string // 合约类型
	leverRate    int    // 杠杆倍数
}

func (b *HuobiBroker) Subscribe(event string, param string, listener interface{}) {

}

func (b *HuobiBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
	var account hbdm.AccountInfoResult
	account, err = b.client.GetAccountInfo(currency)
	if err != nil {
		return
	}

	if account.Status != StatusOK {
		err = fmt.Errorf("error")
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

func (b *HuobiBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ret hbdm.MarketDepthResult

	ret, err = b.client.GetMarketDepth(b.contractType, "step0")
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error")
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
// 如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
func (b *HuobiBroker) SetContractType(contractType string) (err error) {
	b.contractType = contractType
	return
}

func (b *HuobiBroker) GetContractType() (symbol string, err error) {
	return
}

// 设置杠杆大小
func (b *HuobiBroker) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *HuobiBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool) (result Order, err error) {
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
	}
	if postOnly {
		orderPriceType += ",post_only"
	}
	orderResult, err = b.client.Order(
		symbol,
		b.contractType,
		"",
		0,
		price,
		size,
		_direction,
		offset,
		b.leverRate,
		orderPriceType,
	)
	if err != nil {
		return
	}
	if orderResult.Status != StatusOK {
		err = fmt.Errorf("error")
		return
	}
	var order hbdm.OrderInfoResult
	order, err = b.client.OrderInfo(
		symbol,
		orderResult.Data.OrderID,
		0,
	)
	if err != nil {
		return
	}
	if order.Status != StatusOK {
		err = fmt.Errorf("error")
		return
	}
	if len(order.Data) != 1 {
		err = fmt.Errorf("error")
		return
	}
	result = b.convertOrder(symbol, &order.Data[0])
	return
}

func (b *HuobiBroker) GetOpenOrders(symbol string) (result []Order, err error) {
	var ret hbdm.OpenOrdersResult
	ret, err = b.client.GetOpenOrders(
		symbol,
		0,
		0,
	)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error")
		return
	}
	for _, v := range ret.Data.Orders {
		result = append(result, b.convertOrder(symbol, &v))
	}
	return
}

func (b *HuobiBroker) GetOrder(symbol string, id string) (result Order, err error) {
	var ret hbdm.OrderInfoResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.OrderInfo(symbol, _id, 0)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("error")
		return
	}
	if len(ret.Data) != 1 {
		err = fmt.Errorf("not found")
		return
	}
	result = b.convertOrder(symbol, &ret.Data[0])
	return
}

func (b *HuobiBroker) CancelOrder(symbol string, id string) (result Order, err error) {
	var ret hbdm.CancelResult
	var _id, _ = strconv.ParseInt(id, 10, 64)
	ret, err = b.client.Cancel(symbol, _id, 0)
	if err != nil {
		return
	}
	if ret.Status != StatusOK {
		err = fmt.Errorf("err")
		return
	}
	orderID := ret.Data.Successes
	result.ID = orderID
	return
}

func (b *HuobiBroker) CancelAllOrders(symbol string) (err error) {
	return
}

func (b *HuobiBroker) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	return
}

func (b *HuobiBroker) GetPosition(symbol string) (result Position, err error) {
	result.Symbol = symbol

	var ret hbdm.PositionInfoResult
	ret, err = b.client.GetPositionInfo(symbol)
	if err != nil {
		return
	}

	if ret.Status != StatusOK {
		err = fmt.Errorf("error")
		return
	}

	if len(ret.Data) != 1 {
		return
	}

	//position := ret.Data[0]

	//if position > 0 {
	//	result.Size = position.BuyAmount
	//	result.AvgPrice = position.BuyPriceAvg
	//} else if position.SellAmount > 0 {
	//	result.Size = -position.SellAmount
	//	result.AvgPrice = position.SellPriceAvg
	//}
	return
}

func (b *HuobiBroker) convertOrder(symbol string, order *hbdm.Order) (result Order) {
	result.ID = order.OrderIDStr
	result.Symbol = symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Size = order.Volume
	result.Direction = b.orderDirection(order)
	result.Type = b.orderType(order)
	result.AvgPrice = order.TradeAvgPrice
	result.FilledAmount = order.TradeVolume
	if strings.Contains(order.OrderPriceType.String(), "post_only") {
		result.PostOnly = true
	}
	if order.Offset == "close" {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *HuobiBroker) orderDirection(order *hbdm.Order) Direction {
	if order.Direction == "buy" {
		return Buy
	} else if order.Direction == "sell" {
		return Sell
	}
	return Buy
}

func (b *HuobiBroker) orderType(order *hbdm.Order) OrderType {
	/*
		order_price_type 订单报价类型	订单报价类型 订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制,optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档，ioc:IOC订单，fok：FOK订单
	*/

	opt := order.OrderPriceType.String()
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

func (b *HuobiBroker) orderStatus(order *hbdm.Order) OrderStatus {
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

func (b *HuobiBroker) RunEventLoopOnce() (err error) {
	return
}

func NewBroker(addr string, accessKey string, secretKey string) *HuobiBroker {
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
	return &HuobiBroker{
		client: client,
	}
}
