package spotsim

import (
	"errors"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/dataloader"
	"log"
	"time"
)

type SpotSim struct {
	data          *dataloader.Data
	makerFeeRate  float64 // -0.00025	// Maker fee rate
	takerFeeRate  float64 // 0.00075	// Taker fee rate
	initBalance   SpotBalance
	balance       SpotBalance
	backtest      IBacktest
	eLog          ExchangeLogger
	orders        map[string]*Order // All orders key: OrderID value: Order
	openOrders    map[string]*Order // Open orders
	historyOrders map[string]*Order // History orders

}

func New(data *dataloader.Data, initBalance SpotBalance, makerFeeRate float64, takerFeeRate float64) *SpotSim {
	return &SpotSim{
		data:          data,
		makerFeeRate:  makerFeeRate,
		takerFeeRate:  takerFeeRate,
		initBalance:   initBalance,
		balance:       initBalance,
		orders:        make(map[string]*Order),
		openOrders:    make(map[string]*Order),
		historyOrders: make(map[string]*Order),
	}
}

// 获取 Exchange 名称
func (s *SpotSim) GetName() (name string) {
	return "spot_sim"
}

// 获取交易所时间(ms)
func (s *SpotSim) GetTime() (tm int64, err error) {
	if s.data != nil && s.data.GetOrderBook() != nil {
		return s.data.GetOrderBook().Time.UnixNano() / int64(time.Millisecond), nil
	}
	return time.Now().UnixNano() / (int64(time.Millisecond)), nil
}

// 获取账号余额
func (s *SpotSim) GetBalance(currency string) (result *SpotBalance, err error) {
	result = &s.balance
	return
}

// 获取订单薄(OrderBook)
func (s *SpotSim) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
	result = s.data.GetOrderBook()
	return
}

// 获取K线数据
// period: 数据周期. 分钟或者关键字1m(minute) 1h 1d 1w 1M(month) 1y 枚举值：1 3 5 15 30 60 120 240 360 720 "5m" "4h" "1d" ...
func (s *SpotSim) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error) {
	return
}

// 买
func (s *SpotSim) Buy(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Buy, orderType, price, size)
}

// 卖
func (s *SpotSim) Sell(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Sell, orderType, price, size)
}

// 下单
func (s *SpotSim) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, size float64,
	opts ...PlaceOrderOption) (result *Order, err error) {
	if size == 0 {
		err = errors.New("size is zero")
		return
	}
	params := ParsePlaceOrderParameter(opts...)
	id := GenOrderId()
	ob := s.data.GetOrderBook()
	order := &Order{
		ID:           id,
		Symbol:       symbol,
		Time:         ob.Time,
		Price:        price,
		Amount:       size,
		AvgPrice:     0,
		FilledAmount: 0,
		Direction:    direction,
		Type:         orderType,
		PostOnly:     params.PostOnly,
		ReduceOnly:   params.ReduceOnly,
		UpdateTime:   ob.Time,
		Status:       OrderStatusNew,
	}
	s.eLog.Infow(
		"PlaceOrder",
		"symbol", symbol,
		"direction", direction,
		"orderType", orderType.String(),
		"price", price,
		"size", size,
		"params", params,
	)

	_, err = s.matchOrder(order, true)
	if err != nil {
		s.eLog.Error(err)
		return
	}

	if order.IsOpen() {
		s.openOrders[id] = order
	} else {
		s.historyOrders[id] = order
	}

	s.orders[id] = order
	result = order
	s.logOrderInfo("Place order", SimEventOrder, order)
	return

}

// 撮合成交
func (s *SpotSim) matchOrder(order *Order, immediate bool) (match bool, err error) {
	switch order.Type {
	case OrderTypeMarket:
		match, err = s.matchMarketOrder(order)
	case OrderTypeLimit:
		match, err = s.matchLimitOrder(order, immediate)
	}
	return
}

func (s *SpotSim) matchMarketOrder(order *Order) (match bool, err error) {
	if !order.IsOpen() {
		err = errors.New("order is closed")
		return
	}

	ob := s.data.GetOrderBook()

	// 市价成交
	if order.Direction == Buy {
		size := order.Amount
		price := ob.AskAvePrice(size)
		if price <= 0 {
			err = errors.New("size is bigger than orderbook")
			return
		}
		value := size * price
		fee := value * s.takerFeeRate
		if fee+value > s.balance.Quote.Available {
			err = errors.New("no more money")
			return
		}

		order.FilledAmount = size
		order.AvgPrice = price
		order.Commission += fee
		s.balance.Quote.Available = s.balance.Quote.Available - fee - value
		s.balance.Base.Available += size
		order.Status = OrderStatusFilled
	} else if order.Direction == Sell {

		size := order.Amount
		price := ob.BidAvePrice(size)
		if price <= 0 {
			err = errors.New("size is bigger than orderbook")
			return
		}

		fee := size * s.takerFeeRate
		if fee+size > s.balance.Quote.Available {
			err = errors.New("no more stock")
			return
		}

		order.FilledAmount = size
		order.AvgPrice = price

		// Update balance
		order.Commission += fee
		s.balance.Base.Available = s.balance.Base.Available - size - fee
		s.balance.Quote.Available = s.balance.Quote.Available - size*price
		order.Status = OrderStatusFilled
	}
	order.UpdateTime = ob.Time
	match = true
	return
}

func (s *SpotSim) matchLimitOrder(order *Order, immediate bool) (match bool, err error) {
	if !order.IsOpen() {
		return
	}

	ob := s.data.GetOrderBook()
	if order.Direction == Buy { // Bid order
		if order.Price >= ob.AskPrice() {
			if immediate && order.PostOnly {
				order.Status = OrderStatusRejected
				return
			}
			value := order.Price * order.Amount
			fee := value * s.takerFeeRate
			if fee+value > s.balance.Quote.Available {
				err = errors.New("no more money")
				return
			}

			size, price := s.matchAsk(order.Price, order.Amount, ob.Asks)
			if price <= 0 {
				err = errors.New("size is bigger than orderbook")
				return
			}
			value = size * price
			fee = value * s.takerFeeRate

			order.FilledAmount = size
			order.AvgPrice = price
			order.Commission += fee
			s.balance.Quote.Available = s.balance.Quote.Available - fee - value
			s.balance.Base.Available += size

			if size < order.Amount {
				order.Status = OrderStatusPartiallyFilled
				value = (order.Amount - size) * order.Price
				fee = value * s.takerFeeRate
				s.balance.Quote.Available = s.balance.Quote.Available - fee - value
				s.balance.Quote.Frozen = fee + value
			} else {
				order.Status = OrderStatusFilled
			}
			match = true
		}
	} else { // Ask order
		if order.Price <= ob.BidPrice() {
			if immediate && order.PostOnly {
				order.Status = OrderStatusRejected
				return
			}

			value := order.Amount
			fee := value * s.takerFeeRate
			if fee+value > s.balance.Base.Available {
				err = errors.New("no more stock")
				return
			}

			size, price := s.matchBid(order.Price, order.Amount, ob.Bids)
			if price <= 0 {
				err = errors.New("size is bigger than orderbook")
				return
			}
			value = size
			fee = value * s.takerFeeRate

			//s.balance.Base.Available = s.balance.Base.Available - size - fee
			//s.balance.Quote.Available = s.balance.Quote.Available - size * price

			order.FilledAmount = size
			order.AvgPrice = price
			order.Commission += fee
			s.balance.Base.Available = s.balance.Base.Available - fee - value
			s.balance.Quote.Available += size * price

			if size < order.Amount {
				order.Status = OrderStatusPartiallyFilled
				value = order.Amount - size
				fee = value * s.takerFeeRate
				s.balance.Base.Available = s.balance.Base.Available - fee - value
				s.balance.Base.Frozen = fee + value
			} else {
				order.Status = OrderStatusFilled
			}
			match = true
		}
	}
	return
}

// 获取活跃委托单列表
func (s *SpotSim) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	for _, v := range s.openOrders {
		if v.Symbol == symbol {
			result = append(result, v)
		}
	}
	return
}

// 获取历史委托列表
func (s *SpotSim) GetHistoryOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	for _, v := range s.historyOrders {
		if v.Symbol == symbol {
			result = append(result, v)
		}
	}
	return
}

// 获取委托信息
func (s *SpotSim) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	order, ok := s.orders[id]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = order
	return
}

// 撤销全部委托单
func (s *SpotSim) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	var idsToBeRemoved []string

	for _, order := range s.openOrders {
		if !order.IsOpen() {
			log.Printf("Order error: %#v", order)
			continue
		}
		switch order.Status {
		case OrderStatusCreated, OrderStatusNew, OrderStatusPartiallyFilled:
			order.Status = OrderStatusCancelled
			idsToBeRemoved = append(idsToBeRemoved, order.ID)
		default:
			err = errors.New("error")
		}
	}

	for _, id := range idsToBeRemoved {
		delete(s.openOrders, id)
	}
	return
}

// 撤销单个委托单
func (s *SpotSim) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	if order, ok := s.orders[id]; ok {
		if !order.IsOpen() {
			err = errors.New("status error")
			return
		}
		switch order.Status {
		case OrderStatusCreated, OrderStatusNew, OrderStatusPartiallyFilled:
			order.Status = OrderStatusCancelled
			result = order
			delete(s.openOrders, id)
		default:
			err = errors.New("error")
		}
	} else {
		err = errors.New("not found")
	}
	return
}

func (s *SpotSim) SetBacktest(backtest IBacktest) {
	s.backtest = backtest
}

func (s *SpotSim) SetExchangeLogger(l ExchangeLogger) {
	s.eLog = l
}

func (s *SpotSim) RunEventLoopOnce() (err error) {
	var match bool
	for _, order := range s.openOrders {
		match, err = s.matchOrder(order, false)
		if match {
			s.logOrderInfo("Match order", SimEventDeal, order)
			//var orders = []*Order{order}
			//s.emitter.Emit(WSEventOrder, orders)
		}
	}
	return
}

func (s *SpotSim) logOrderInfo(msg string, event string, order *Order) {
	//ob := s.data.GetOrderBook()
	//s.eLog.Infow(
	//	msg,
	//	SimEventKey,
	//	event,
	//	"order", order,
	//	"orderbook", ob,
	//	"balance", s.balance,
	//	"positions", position,
	//)
}

func (s *SpotSim) matchAsk(price, size float64, asks []Item) (filledSize float64, avgPrice float64) {
	type item = struct {
		Amount float64
		Price  float64
	}

	var items []item
	lSize := size
	for i := 0; i < len(asks); i++ {
		if price < asks[i].Price {
			break
		}
		if lSize >= asks[i].Amount {
			items = append(items, item{
				Amount: asks[i].Amount,
				Price:  asks[i].Price,
			})
			lSize -= asks[i].Amount
		} else {
			items = append(items, item{
				Amount: lSize,
				Price:  asks[i].Price,
			})
			lSize = 0
		}
		if lSize <= 0 {
			break
		}
	}

	if lSize != 0 {
		return
	}

	// 计算平均价
	amount := 0.0
	for _, v := range items {
		amount += v.Price * v.Amount
		filledSize += v.Amount
	}
	if filledSize == 0 {
		return
	}
	avgPrice = amount / filledSize
	return
}

func (s *SpotSim) matchBid(price, size float64, bids []Item) (filledSize float64, avgPrice float64) {
	type item = struct {
		Amount float64
		Price  float64
	}

	var items []item
	lSize := size
	for i := 0; i < len(bids); i++ {
		if price > bids[i].Price {
			break
		}
		if lSize >= bids[i].Amount {
			items = append(items, item{
				Amount: bids[i].Amount,
				Price:  bids[i].Price,
			})
			lSize -= bids[i].Amount
		} else {
			items = append(items, item{
				Amount: lSize,
				Price:  bids[i].Price,
			})
			lSize = 0
		}
		if lSize <= 0 {
			break
		}
	}

	if lSize != 0 {
		return
	}

	// 计算平均价
	amount := 0.0
	for _, v := range items {
		amount += v.Price * v.Amount
		filledSize += v.Amount
	}
	if filledSize == 0 {
		return
	}
	avgPrice = amount / filledSize
	return
}
