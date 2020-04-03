package okex_futures_broker

import (
	"fmt"
	"log"
	"strconv"
	"time"

	. "github.com/coinrust/crex"
	"github.com/frankrap/okex-api"
)

// OKEXSwapBroker the OKEX swap broker
type OKEXSwapBroker struct {
	client        *okex.Client
	pair          string // contract pair 合约交易对
	contractType  string // contract type 合约类型
	contractAlias string // okex contract type 合约类型
	leverRate     int    // lever rate 杠杆倍数
}

func (b *OKEXSwapBroker) Subscribe(event string, param string, listener interface{}) {

}

func (b *OKEXSwapBroker) GetAccountSummary(currency string) (result AccountSummary, err error) {
	var account okex.SwapAccount
	account, err = b.client.GetSwapAccount(currency)
	if err != nil {
		return
	}

	result.Equity, _ = strconv.ParseFloat(account.Info.Equity, 64)
	result.Balance, _ = strconv.ParseFloat(account.Info.TotalAvailBalance, 64)
	result.Pnl, _ = strconv.ParseFloat(account.Info.RealizedPnl, 64)

	return
}

func (b *OKEXSwapBroker) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
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
			Price:  ParseFloat(v[0]),
			Amount: ParseFloat(v[1]),
		})
	}

	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  ParseFloat(v[0]),
			Amount: ParseFloat(v[1]),
		})
	}

	// 2019-07-04T09:35:07.752Z
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", ret.Time)
	result.Time = timestamp.Local()
	return
}

// 设置合约类型
func (b *OKEXSwapBroker) SetContractType(pair string, contractType string) (err error) {
	return
}

func (b *OKEXSwapBroker) GetContractID() (symbol string, err error) {
	return
}

// 设置杠杆大小
func (b *OKEXSwapBroker) SetLeverRate(value float64) (err error) {
	b.leverRate = int(value)
	return
}

func (b *OKEXSwapBroker) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	stopPx float64, size float64, postOnly bool, reduceOnly bool) (result Order, err error) {
	var pType int
	if direction == Buy {
		if reduceOnly {
			pType = 4
		} else {
			pType = 1
		}
	} else if direction == Sell {
		if reduceOnly {
			pType = 3
		} else {
			pType = 2
		}
	}
	var _orderType int
	if postOnly {
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
	log.Printf("%v", string(resp))
	result, err = b.GetOrder(symbol, ret.OrderId)
	return
}

func (b *OKEXSwapBroker) GetOpenOrders(symbol string) (result []Order, err error) {
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

func (b *OKEXSwapBroker) GetOrder(symbol string, id string) (result Order, err error) {
	var ret okex.BaseOrderInfo
	ret, err = b.client.GetSwapOrderById(symbol, id)
	if err != nil {
		return
	}
	result = b.convertOrder(symbol, &ret)
	return
}

func (b *OKEXSwapBroker) CancelOrder(symbol string, id string) (result Order, err error) {
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
	result.ID = ret.OrderId
	return
}

func (b *OKEXSwapBroker) CancelAllOrders(symbol string) (err error) {
	return
}

func (b *OKEXSwapBroker) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	return
}

func (b *OKEXSwapBroker) GetPosition(symbol string) (result Position, err error) {
	var ret okex.SwapPosition
	ret, err = b.client.GetSwapPositionByInstrument(symbol)
	if err != nil {
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("%v [%v]", err, ret)
		return
	}

	result.Symbol = symbol

	if ret.MarginMode == "crossed" || ret.MarginMode == "fixed" { // 全仓
		for _, v := range ret.Holding {
			if v.InstrumentId == symbol {
				// 2019-10-08T11:56:07.922Z
				timestamp, _ := time.ParseInLocation(v.Timestamp,
					"2006-01-02T15:04:05.000Z",
					time.Now().Location())
				if v.Side == "long" {
					result.Size, _ = strconv.ParseFloat(v.Position, 64)
					result.AvgPrice, _ = strconv.ParseFloat(v.AvgCost, 64)
					result.OpenTime = timestamp
				} else if v.Side == "short" {
					result.Size, _ = strconv.ParseFloat(v.Position, 64)
					result.AvgPrice, _ = strconv.ParseFloat(v.AvgCost, 64)
					result.OpenTime = timestamp
				}
				break
			}
		}
	}
	return
}

func (b *OKEXSwapBroker) convertOrder(symbol string, order *okex.BaseOrderInfo) (result Order) {
	result.ID = order.OrderId
	result.Symbol = symbol
	result.Price = order.Price
	result.StopPx = 0
	result.Size = float64(order.Size)
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

func (b *OKEXSwapBroker) orderDirection(order *okex.BaseOrderInfo) Direction {
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

func (b *OKEXSwapBroker) orderType(order *okex.BaseOrderInfo) OrderType {
	if order.OrderType == "4" {
		return OrderTypeMarket
	}
	return OrderTypeLimit
}

func (b *OKEXSwapBroker) orderStatus(order *okex.BaseOrderInfo) OrderStatus {
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

func (b *OKEXSwapBroker) RunEventLoopOnce() (err error) {
	return
}

// addr: https://www.okex.com/
func NewBroker(addr string, accessKey string, secretKey string, passphrase string) *OKEXSwapBroker {
	config := okex.Config{
		Endpoint:      addr,
		WSEndpoint:    "",
		ApiKey:        accessKey,
		SecretKey:     secretKey,
		Passphrase:    passphrase,
		TimeoutSecond: 45,
		IsPrint:       false,
		I18n:          okex.ENGLISH,
		ProxyURL:      "",
	}
	client := okex.NewClient(config)
	return &OKEXSwapBroker{
		client: client,
	}
}

func ParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ParseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
