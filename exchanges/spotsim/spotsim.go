package spotsim

import (
	"errors"
	"github.com/chuckpreslar/emission"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/dataloader"
	"log"
	"sync"
	"time"
)

type SpotSim struct {
	name          string
	data          *dataloader.Data
	makerFeeRate  float64 // -0.00025	// Maker fee rate
	takerFeeRate  float64 // 0.00075	// Taker fee rate
	initBalance   SpotBalance
	balance       SpotBalance
	backtest      IBacktest
	eLog          ExchangeLogger
	orders        sync.Map // All orders key: OrderID value: Order
	openOrders    sync.Map // Open orders
	historyOrders sync.Map // History orders
	emitter       *emission.Emitter
}

func New(name string, data *dataloader.Data, initBalance SpotBalance, makerFeeRate float64, takerFeeRate float64) *SpotSim {
	return &SpotSim{
		name:         name,
		data:         data,
		makerFeeRate: makerFeeRate,
		takerFeeRate: takerFeeRate,
		initBalance:  initBalance,
		balance:      initBalance,
		emitter:      emission.NewEmitter(),
	}
}

// 获取 Exchange 名称
func (s *SpotSim) GetName() (name string) {
	return s.name + "_spot_sim"
}

// 获取交易所时间(ms)
func (s *SpotSim) GetTime() (tm int64, err error) {
	tm = s.backtest.GetTime().UnixNano() / int64(time.Millisecond)
	return
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
	//ob := s.data.GetOrderBook()
	order := &Order{
		ID:           id,
		Symbol:       symbol,
		Time:         s.backtest.GetTime(),
		Price:        price,
		Amount:       size,
		AvgPrice:     0,
		FilledAmount: 0,
		Direction:    direction,
		Type:         orderType,
		PostOnly:     params.PostOnly,
		ReduceOnly:   params.ReduceOnly,
		UpdateTime:   s.backtest.GetTime(),
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
		s.openOrders.Store(id, order)
	} else {
		s.historyOrders.Store(id, order)
	}

	s.orders.Store(id, order)
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
		fee := size * s.takerFeeRate // 买扣币，卖扣钱
		if fee+value > s.balance.Quote.Available {
			err = errors.New("no more money")
			return
		}

		order.FilledAmount = size
		order.AvgPrice = price
		order.Commission += fee
		s.balance.Quote.Available -= value
		s.balance.Base.Available += size - fee
		order.Status = OrderStatusFilled
	} else if order.Direction == Sell {

		size := order.Amount
		price := ob.BidAvePrice(size)
		if price <= 0 {
			err = errors.New("size is bigger than orderbook")
			return
		}

		value := size * price
		fee := value * s.takerFeeRate // 买扣币，卖扣钱
		if fee+size > s.balance.Quote.Available {
			err = errors.New("no more stock")
			return
		}

		order.FilledAmount = size
		order.AvgPrice = price

		// Update balance
		order.Commission += fee
		s.balance.Base.Available = s.balance.Base.Available - size
		s.balance.Quote.Available = s.balance.Quote.Available + value - fee
		order.Status = OrderStatusFilled
	}
	order.UpdateTime = s.backtest.GetTime()
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

			if value > s.balance.Quote.Available {
				err = errors.New("no more money")
				return
			}

			size, price := s.matchAsk(order.Price, order.Amount, ob.Asks, immediate)
			if price <= 0 {
				err = errors.New("size is bigger than orderbook")
				return
			}
			value = size * price
			fee := 0.0
			if immediate {
				fee = size * s.takerFeeRate
			} else {
				fee = size * s.makerFeeRate
			}

			order.FilledAmount = size
			order.AvgPrice = price
			order.Commission += fee
			s.balance.Quote.Available = s.balance.Quote.Available - value
			s.balance.Base.Available += size - fee

			if size < order.Amount {
				order.Status = OrderStatusPartiallyFilled
				value = (order.Amount - size) * order.Price
				s.balance.Quote.Frozen = value
			} else {
				order.Status = OrderStatusFilled
			}
			match = true
			order.UpdateTime = s.backtest.GetTime()
		}
	} else { // Ask order
		if order.Price <= ob.BidPrice() {
			if immediate && order.PostOnly {
				order.Status = OrderStatusRejected
				return
			}

			value := order.Amount
			if order.Amount > s.balance.Base.Available {
				err = errors.New("no more stock")
				return
			}

			size, price := s.matchBid(order.Price, order.Amount, ob.Bids, immediate)
			if price <= 0 {
				err = errors.New("size is bigger than orderbook")
				return
			}
			value = size * price
			fee := 0.0
			if immediate {
				fee = value * s.takerFeeRate
			} else {
				fee = value * s.makerFeeRate
			}

			order.FilledAmount = size
			order.AvgPrice = price
			order.Commission += fee
			s.balance.Base.Available = s.balance.Base.Available - size
			s.balance.Quote.Available += value - fee

			if size < order.Amount {
				order.Status = OrderStatusPartiallyFilled
				value = order.Amount - size
				s.balance.Base.Frozen = value
			} else {
				order.Status = OrderStatusFilled
			}
			order.UpdateTime = s.backtest.GetTime()
			match = true
		}
	}
	return
}

// 获取活跃委托单列表
func (s *SpotSim) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	s.openOrders.Range(func(key, value interface{}) bool {
		v := value.(*Order)
		if v.Symbol == symbol {
			result = append(result, v)
		}
		return true
	})
	return
}

// 获取历史委托列表
func (s *SpotSim) GetHistoryOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	s.historyOrders.Range(func(key, value interface{}) bool {
		v := value.(*Order)
		if v.Symbol == symbol {
			result = append(result, v)
		}
		return true
	})
	return
}

// 获取委托信息
func (s *SpotSim) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	order, ok := s.orders.Load(id)
	if !ok {
		err = errors.New("not found")
		return
	}
	result = order.(*Order)
	return
}

// 撤销全部委托单
func (s *SpotSim) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	var idsToBeRemoved []string

	s.openOrders.Range(func(key, value interface{}) bool {
		order := value.(*Order)
		if !order.IsOpen() {
			log.Printf("Order error: %#v", order)
			return true
		}
		switch order.Status {
		case OrderStatusCreated, OrderStatusNew, OrderStatusPartiallyFilled:
			order.Status = OrderStatusCancelled
			idsToBeRemoved = append(idsToBeRemoved, order.ID)
		default:
			err = errors.New("error")
			return false
		}
		return true
	})

	for _, id := range idsToBeRemoved {
		s.openOrders.Delete(id)
		orderValue, ok := s.orders.Load(id)
		if !ok {
			continue
		}
		order, ok := orderValue.(*Order)
		if !ok {
			continue
		}
		s.logOrderInfo("Cancel order", SimEventOrder, order)
	}
	return
}

// 撤销单个委托单
func (s *SpotSim) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	if value, ok := s.orders.Load(id); ok {
		order := value.(*Order)
		if !order.IsOpen() {
			err = errors.New("status error")
			return
		}
		switch order.Status {
		case OrderStatusCreated, OrderStatusNew, OrderStatusPartiallyFilled:
			order.Status = OrderStatusCancelled
			result = order
			s.openOrders.Delete(id)
			s.logOrderInfo("Cancel order", SimEventOrder, order)
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
	s.openOrders.Range(func(key, value interface{}) bool {
		order := value.(*Order)
		match, err = s.matchOrder(order, false)
		if match {
			s.logOrderInfo("Match order", SimEventDeal, order)
			var orders = []*Order{order}
			s.emitter.Emit(WSEventOrder, orders)
		}
		return true
	})
	return
}

func (s *SpotSim) logOrderInfo(msg string, event string, order *Order) {
	ob := s.data.GetOrderBook()
	baseBalance := s.balance.Base.Available + s.balance.Base.Frozen
	quoteBalance := s.balance.Quote.Available + s.balance.Quote.Frozen
	s.eLog.Infow(
		msg,
		SimEventKey,
		event,
		"order", order,
		"orderbook", ob,
		"balance", baseBalance*ob.Price()+quoteBalance,
		"balances", []float64{baseBalance, quoteBalance},
	)
}

func (s *SpotSim) matchAsk(price, size float64, asks []Item, immediate bool) (filledSize float64, avgPrice float64) {
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
	if immediate {
		value := 0.0
		for _, v := range items {
			value += v.Price * v.Amount
			filledSize += v.Amount
		}
		if filledSize == 0 {
			return
		}
		avgPrice = value / filledSize
	} else {
		filledSize, avgPrice = size, price
	}
	return
}

func (s *SpotSim) matchBid(price, size float64, bids []Item, immediate bool) (filledSize float64, avgPrice float64) {
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
	if immediate {
		amount := 0.0
		for _, v := range items {
			amount += v.Price * v.Amount
			filledSize += v.Amount
		}
		if filledSize == 0 {
			return
		}
		avgPrice = amount / filledSize
	} else {
		filledSize, avgPrice = size, price
	}
	return
}
