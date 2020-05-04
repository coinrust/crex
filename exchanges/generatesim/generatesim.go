package generatesim

import (
	"errors"
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/data"
	"github.com/coinrust/crex/util"
	"log"
	"math"
	"time"
)

type MarginInfo struct {
	Leverage              float64
	MaintMargin           float64
	LiquidationPriceLong  float64
	LiquidationPriceShort float64
}

// GenerateSim the generate/common exchange for backtest
type GenerateSim struct {
	data               *data.Data
	makerFeeRate       float64 // -0.00025	// Maker fee rate
	takerFeeRate       float64 // 0.00075	// Taker fee rate
	balance            float64
	orders             map[string]*Order     // All orders key: OrderID value: Order
	openOrders         map[string]*Order     // Open orders
	historyOrders      map[string]*Order     // History orders
	positions          map[string][]Position // Position key: symbol, index:0, long; 1, short, when dual side mode, its size is 2, otherwise size is 1
	isDualSidePosition bool                  // dual side position
	isForwardContract  bool                  // forward contract, otherwise reverse contract, UNUSED, TBD
}

func NewGenerateSim(data *data.Data, cash float64, makerFeeRate float64, takerFeeRate float64, isForwardContract bool, posMode ...bool) *GenerateSim {
	isDualSidePosition := false
	if len(posMode) > 0 {
		isDualSidePosition = posMode[0]
	}
	return &GenerateSim{
		data:               data,
		balance:            cash,
		makerFeeRate:       makerFeeRate, // -0.00025 // Maker 费率
		takerFeeRate:       takerFeeRate, // 0.00075	// Taker 费率
		orders:             make(map[string]*Order),
		openOrders:         make(map[string]*Order),
		historyOrders:      make(map[string]*Order),
		positions:          make(map[string][]Position),
		isDualSidePosition: isDualSidePosition,
		isForwardContract:  isForwardContract,
	}
}

func (s *GenerateSim) GetName() (name string) {
	return "generate"
}

func (s *GenerateSim) GetTime() (tm int64, err error) {
	return time.Now().UnixNano() / (int64(time.Millisecond)), nil
}

func (s *GenerateSim) GetBalance(symbol string) (result Balance, err error) {
	result.Available = s.balance

	position := s.getPosition(symbol)
	var price float64
	var pnl float64
	ob := s.data.GetOrderBook()
	for _, pos := range position {
		side := pos.Side()
		if side == Buy {
			price = ob.AskPrice()
		} else if side == Sell {
			price = ob.BidPrice()
		}
		pnl2, _ := CalcPnl(side, math.Abs(pos.Size), pos.AvgPrice, price)
		pnl += pnl2
	}
	result.Equity = result.Available + pnl
	return
}

func (s *GenerateSim) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	result = *s.data.GetOrderBook()
	return
}

func (s *GenerateSim) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	return
}

func (s *GenerateSim) SetContractType(pair string, contractType string) (err error) {
	return
}

func (s *GenerateSim) GetContractID() (symbol string, err error) {
	return
}

func (s *GenerateSim) SetLeverRate(value float64) (err error) {
	return
}

func (s *GenerateSim) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (s *GenerateSim) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (s *GenerateSim) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (s *GenerateSim) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return s.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (s *GenerateSim) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...PlaceOrderOption) (result *Order, err error) {
	if size == 0 {
		err = errors.New("size is zero")
		return
	}
	params := ParsePlaceOrderParameter(opts...)
	_id, _ := util.NextID()
	id := fmt.Sprintf("%v", _id)
	order := &Order{
		ID:           id,
		Symbol:       symbol,
		Price:        price,
		Amount:       size,
		AvgPrice:     0,
		FilledAmount: 0,
		Direction:    direction,
		Type:         orderType,
		PostOnly:     params.PostOnly,
		ReduceOnly:   params.ReduceOnly,
		Status:       OrderStatusNew,
	}

	err = s.matchOrder(order, true)
	if err != nil {
		return
	}

	if order.IsOpen() {
		s.openOrders[id] = order
	} else {
		s.historyOrders[id] = order
	}

	s.orders[id] = order
	result = order
	return
}

// 撮合成交
func (s *GenerateSim) matchOrder(order *Order, immediate bool) (err error) {
	switch order.Type {
	case OrderTypeMarket:
		err = s.matchMarketOrder(order)
	case OrderTypeLimit:
		err = s.matchLimitOrder(order, immediate)
	}
	return
}

func (s *GenerateSim) matchMarketOrder(order *Order) (err error) {
	if !order.IsOpen() {
		return
	}

	ob := s.data.GetOrderBook()

	// 市价成交
	if order.Direction == Buy {
		size := order.Amount
		price := ob.AskAvePrice(size)
		if price <= 0 {
			return
		}

		// Update position

		size, err = s.updatePosition(order.Symbol, size, price, order.ReduceOnly)
		if err != nil {
			order.Status = OrderStatusRejected
			err = errors.New("order rejected")
			return
		}
		// trade fee
		fee := size / price * s.takerFeeRate
		order.FilledAmount = size
		order.AvgPrice = price
		// Update balance
		s.addBalance(-fee)

	} else if order.Direction == Sell {

		size := order.Amount
		price := ob.BidAvePrice(size)
		// Update position
		size, err = s.updatePosition(order.Symbol, -size, price, order.ReduceOnly)
		if err != nil {
			order.Status = OrderStatusRejected
			err = errors.New("order rejected")
			return
		}

		// trade fee
		fee := size / price * s.takerFeeRate
		order.FilledAmount = size
		order.AvgPrice = price

		// Update balance
		s.addBalance(-fee)

	}
	order.Status = OrderStatusFilled
	return
}

func (s *GenerateSim) matchLimitOrder(order *Order, immediate bool) (err error) {
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

			// match trade
			size := order.Amount
			var fee float64

			// trade fee
			if immediate {
				fee = size / order.Price * s.takerFeeRate
			} else {
				fee = size / order.Price * s.makerFeeRate
			}

			// Update balance
			s.addBalance(-fee)

			// Update position
			s.updatePosition(order.Symbol, size, order.Price, order.ReduceOnly)
		}
	} else { // Ask order
		if order.Price <= ob.BidPrice() {
			if immediate && order.PostOnly {
				order.Status = OrderStatusRejected
				return
			}

			// match trade
			size := order.Amount
			var fee float64

			// trade fee
			if immediate {
				fee = size / order.Price * s.takerFeeRate
			} else {
				fee = size / order.Price * s.makerFeeRate
			}

			// Update balance
			s.addBalance(-fee)

			// Update position
			s.updatePosition(order.Symbol, -size, order.Price, order.ReduceOnly)
		}
	}
	return
}

// 更新持仓
func (s *GenerateSim) updatePosition(symbol string, size float64, price float64, isReduce bool) (amount float64, err error) {
	position := s.getPosition(symbol)
	if position == nil {
		err = fmt.Errorf("position error symbol=%v", symbol)
		return
	}

	if !s.isDualSidePosition {
		if position[0].Size > 0 && size < 0 || position[0].Size < 0 && size > 0 {
			return s.closePosition(&position[0], size, price, isReduce)
		} else {
			return s.addPosition(&position[0], size, price)
		}
	} else {
		if size < 0 {
			if isReduce {
				return s.closePosition(&position[1], size, price, isReduce)
			} else {
				return s.addPosition(&position[1], size, price)
			}
		} else if size > 0 {
			if isReduce {
				return s.closePosition(&position[0], size, price, isReduce)
			} else {
				return s.addPosition(&position[0], size, price)
			}
		}
	}
	err = errors.New("error")
	return
}

// 增加持仓
func (s *GenerateSim) addPosition(position *Position, size float64, price float64) (amount float64, err error) {
	if position.Size < 0 && size > 0 || position.Size > 0 && size < 0 {
		err = errors.New("方向错误")
		return
	}
	// 平均成交价
	// total_quantity / ((quantity_1 / price_1) + (quantity_2 / price_2)) = entry_price
	// 增加持仓
	var positionCost float64
	if position.Size != 0 && position.AvgPrice != 0 {
		positionCost = math.Abs(position.Size) / position.AvgPrice
	}

	newPositionCost := math.Abs(size) / price
	totalCost := positionCost + newPositionCost

	totalSize := math.Abs(position.Size + size)
	avgPrice := totalSize / totalCost

	position.AvgPrice = avgPrice
	position.Size += size
	amount = math.Abs(size)
	return
}

// 平仓，超过数量，则开立新仓
func (s *GenerateSim) closePosition(position *Position, size float64, price float64, isReduce bool) (amount float64, err error) {
	if position.Size == 0 {
		err = errors.New("当前无持仓")
		return
	}
	//if position.Size > 0 && size > 0 || position.Size < 0 && size < 0 {
	//	err = errors.New("方向错误")
	//	return
	//}
	remaining := math.Abs(size) - math.Abs(position.Size)
	if isReduce {
		remaining = 0
		amount = math.Abs(position.Size)
	} else {
		amount = math.Abs(size)
	}

	if remaining > 0 {
		// 先平掉原有持仓
		// 计算盈利
		pnl, _ := CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		s.addPnl(pnl)
		position.AvgPrice = price
		position.Size = position.Size + size
	} else if remaining == 0 {
		// 完全平仓
		pnl, _ := CalcPnl(position.Side(), math.Abs(size), position.AvgPrice, price)
		s.addPnl(pnl)
		position.AvgPrice = 0
		position.Size = 0
	} else {
		// 部分平仓
		pnl, _ := CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		s.addPnl(pnl)
		//position.AvgPrice = position.AvgPrice
		position.Size = position.Size + size
	}
	return
}

// 增加Balance
func (s *GenerateSim) addBalance(value float64) {
	s.balance += value
}

// 增加P/L
func (s *GenerateSim) addPnl(pnl float64) {
	s.balance += pnl
}

// 获取持仓
func (s *GenerateSim) getPosition(symbol string) []Position {
	if position, ok := s.positions[symbol]; ok {
		return position
	} else {
		position = append(position, Position{
			Symbol:    symbol,
			OpenTime:  time.Time{},
			OpenPrice: 0,
			Size:      0,
			AvgPrice:  0,
		})
		if s.isDualSidePosition {
			position = append(position, Position{
				Symbol:    symbol,
				OpenTime:  time.Time{},
				OpenPrice: 0,
				Size:      0,
				AvgPrice:  0,
			})
		}
		s.positions[symbol] = position
		return s.positions[symbol]
	}
}

func (s *GenerateSim) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	for _, v := range s.openOrders {
		if v.Symbol == symbol {
			result = append(result, v)
		}
	}
	return
}

func (s *GenerateSim) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	order, ok := s.orders[id]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = order
	return
}

func (s *GenerateSim) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
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

func (s *GenerateSim) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
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

func (s *GenerateSim) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (s *GenerateSim) GetPositions(symbol string) (result []*Position, err error) {
	ret := s.getPosition(symbol)
	for _, v := range ret {
		o := v
		result = append(result, &o)
	}
	return
}

func (s *GenerateSim) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	return nil
}

func (s *GenerateSim) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	return nil
}

func (s *GenerateSim) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	return nil
}

func (s *GenerateSim) SubscribePositions(market Market, callback func(positions []*Position)) error {
	return nil
}

func (s *GenerateSim) RunEventLoopOnce() (err error) {
	for _, order := range s.openOrders {
		s.matchOrder(order, false)
	}
	return
}
