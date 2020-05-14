package deribitsim

import (
	"errors"
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/dataloader"
	"log"
	"math"
	"time"
)

const (
	PositionSizeLimit = 1000000 // Position size limit
)

type MarginInfo struct {
	Leverage              float64
	MaintMargin           float64
	LiquidationPriceLong  float64
	LiquidationPriceShort float64
}

// DeribitSim the deribit exchange for backtest
type DeribitSim struct {
	data          *dataloader.Data
	makerFeeRate  float64 // -0.00025	// Maker fee rate
	takerFeeRate  float64 // 0.00075	// Taker fee rate
	balance       float64
	orders        map[string]*Order    // All orders key: OrderID value: Order
	openOrders    map[string]*Order    // Open orders
	historyOrders map[string]*Order    // History orders
	positions     map[string]*Position // Position key: symbol

	eLog ExchangeLogger
}

func (b *DeribitSim) GetName() (name string) {
	return "deribit"
}

func (b *DeribitSim) GetTime() (tm int64, err error) {
	err = ErrNotImplemented
	return
}

func (b *DeribitSim) GetBalance(symbol string) (result *Balance, err error) {
	result = &Balance{}
	result.Available = b.balance
	position := b.getPosition(symbol)
	var price float64
	ob := b.data.GetOrderBook()
	side := position.Side()
	if side == Buy {
		price = ob.AskPrice()
	} else if side == Sell {
		price = ob.BidPrice()
	}
	pnl, _ := CalcPnl(side, math.Abs(position.Size), position.AvgPrice, price)
	result.Equity = result.Available + pnl
	return
}

func (b *DeribitSim) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
	result = b.data.GetOrderBook()
	return
}

func (b *DeribitSim) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error) {
	return
}

func (b *DeribitSim) SetContractType(pair string, contractType string) (err error) {
	return
}

func (b *DeribitSim) GetContractID() (symbol string, err error) {
	return
}

func (b *DeribitSim) SetLeverRate(value float64) (err error) {
	return
}

func (b *DeribitSim) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *DeribitSim) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *DeribitSim) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *DeribitSim) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *DeribitSim) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...PlaceOrderOption) (result *Order, err error) {
	params := ParsePlaceOrderParameter(opts...)
	id := GenOrderId()
	ob := b.data.GetOrderBook()
	order := &Order{
		ID:           id,
		Time:         ob.Time,
		Symbol:       symbol,
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

	b.eLog.Infow("PlaceOrder",
		"symbol", symbol,
		"direction", direction,
		"orderType", orderType.String(),
		"price", price,
		"size", size,
		"params", params)
	_, err = b.matchOrder(order, true)
	if err != nil {
		b.eLog.Error(err)
		return
	}

	if order.IsOpen() {
		b.openOrders[id] = order
	} else {
		b.historyOrders[id] = order
	}

	b.orders[id] = order
	result = order
	b.eLog.Infow("Place order",
		SimEventKey, SimEventOrder,
		"order", order,
		"orderbook", ob,
		"balance", b.balance,
		"position", b.positions)
	return
}

// 撮合成交
func (b *DeribitSim) matchOrder(order *Order, immediate bool) (changed bool, err error) {
	switch order.Type {
	case OrderTypeMarket:
		changed, err = b.matchMarketOrder(order)
	case OrderTypeLimit:
		changed, err = b.matchLimitOrder(order, immediate)
	}
	return
}

func (b *DeribitSim) matchMarketOrder(order *Order) (changed bool, err error) {
	if !order.IsOpen() {
		return
	}

	// 检查委托:
	// Rejected, maximum size of future position is $1,000,000
	// 开仓总量不能大于 1000000
	// Invalid size - not multiple of contract size ($10)
	// 数量必须是10的整数倍

	if int(order.Amount)%10 != 0 {
		err = errors.New("invalid size - not multiple of contract size ($10)")
		return
	}

	position := b.getPosition(order.Symbol)

	if int(position.Size+order.Amount) > PositionSizeLimit ||
		int(position.Size-order.Amount) < -PositionSizeLimit {
		err = errors.New("rejected, maximum size of future position is $1,000,000")
		return
	}

	ob := b.data.GetOrderBook()

	// 判断开仓数量
	margin := b.balance
	// sizeCurrency := order.Amount / price(ask/bid)
	// leverage := sizeCurrency / margin
	// 需要满足: sizeCurrency <= margin * 100
	// 可开仓数量: <= margin * 100 * price(ask/bid)
	var maxSize float64

	// 市价成交
	if order.Direction == Buy {
		maxSize = margin * 100 * ob.AskPrice()
		if order.Amount > maxSize {
			err = errors.New(fmt.Sprintf("rejected, maximum size of future position is %v", maxSize))
			return
		}

		price := ob.AskPrice()
		size := order.Amount

		// trade fee
		fee := size / price * b.takerFeeRate

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, size, price)

		order.AvgPrice = price
	} else if order.Direction == Sell {
		maxSize = margin * 100 * ob.BidPrice()
		if order.Amount > maxSize {
			err = errors.New(fmt.Sprintf("Rejected, maximum size of future position is %v", maxSize))
			return
		}

		price := ob.BidPrice()
		size := order.Amount

		// trade fee
		fee := size / price * b.takerFeeRate

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, -size, price)

		order.AvgPrice = price
	}
	order.FilledAmount = order.Amount
	order.UpdateTime = ob.Time
	order.Status = OrderStatusFilled
	changed = true
	return
}

func (b *DeribitSim) matchLimitOrder(order *Order, immediate bool) (changed bool, err error) {
	if !order.IsOpen() {
		return
	}

	ob := b.data.GetOrderBook()
	if order.Direction == Buy { // Bid order
		if order.Price < ob.AskPrice() {
			return
		}

		if immediate && order.PostOnly {
			order.UpdateTime = ob.Time
			order.Status = OrderStatusRejected
			changed = true
			return
		}

		// match trade
		size := order.Amount
		var fee float64

		// trade fee
		if immediate {
			fee = size / order.Price * b.takerFeeRate
		} else {
			fee = size / order.Price * b.makerFeeRate
		}

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, size, order.Price)

		order.AvgPrice = order.Price
		order.FilledAmount = order.Amount
		order.UpdateTime = ob.Time
		order.Status = OrderStatusFilled
		changed = true
	} else { // Ask order
		if order.Price > ob.BidPrice() {
			return
		}

		if immediate && order.PostOnly {
			order.UpdateTime = ob.Time
			order.Status = OrderStatusRejected
			changed = true
			return
		}

		// match trade
		size := order.Amount
		var fee float64

		// trade fee
		if immediate {
			fee = size / order.Price * b.takerFeeRate
		} else {
			fee = size / order.Price * b.makerFeeRate
		}

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, -size, order.Price)

		order.AvgPrice = order.Price
		order.FilledAmount = order.Amount
		order.UpdateTime = ob.Time
		order.Status = OrderStatusFilled
		changed = true
	}
	return
}

// 更新持仓
func (b *DeribitSim) updatePosition(symbol string, size float64, price float64) {
	position := b.getPosition(symbol)
	if position == nil {
		log.Fatalf("position error symbol=%v", symbol)
	}

	if position.Size > 0 && size < 0 || position.Size < 0 && size > 0 {
		b.closePosition(position, size, price)
	} else {
		b.addPosition(position, size, price)
	}
}

// 增加持仓
func (b *DeribitSim) addPosition(position *Position, size float64, price float64) (err error) {
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
	return
}

// 平仓，超过数量，则开立新仓
func (b *DeribitSim) closePosition(position *Position, size float64, price float64) (err error) {
	if position.Size == 0 {
		err = errors.New("当前无持仓")
		return
	}
	if position.Size > 0 && size > 0 || position.Size < 0 && size < 0 {
		err = errors.New("方向错误")
		return
	}
	remaining := math.Abs(size) - math.Abs(position.Size)
	if remaining > 0 {
		// 先平掉原有持仓
		// 计算盈利
		pnl, _ := CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		b.addPnl(pnl)
		position.AvgPrice = price
		position.Size = position.Size + size
	} else if remaining == 0 {
		// 完全平仓
		pnl, _ := CalcPnl(position.Side(), math.Abs(size), position.AvgPrice, price)
		b.addPnl(pnl)
		position.AvgPrice = 0
		position.Size = 0
	} else {
		// 部分平仓
		pnl, _ := CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		b.addPnl(pnl)
		//position.AvgPrice = position.AvgPrice
		position.Size = position.Size + size
	}
	return
}

// 增加Balance
func (b *DeribitSim) addBalance(value float64) {
	b.balance += value
}

// 增加P/L
func (b *DeribitSim) addPnl(pnl float64) {
	b.balance += pnl
}

// 获取持仓
func (b *DeribitSim) getPosition(symbol string) *Position {
	if position, ok := b.positions[symbol]; ok {
		return position
	} else {
		position = &Position{
			Symbol:    symbol,
			OpenTime:  time.Time{},
			OpenPrice: 0,
			Size:      0,
			AvgPrice:  0,
		}
		b.positions[symbol] = position
		return position
	}
}

func (b *DeribitSim) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	for _, v := range b.openOrders {
		if v.Symbol == symbol {
			result = append(result, v)
		}
	}
	return
}

func (b *DeribitSim) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	order, ok := b.orders[id]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = order
	return
}

func (b *DeribitSim) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	if order, ok := b.orders[id]; ok {
		if !order.IsOpen() {
			err = errors.New("status error")
			return
		}
		switch order.Status {
		case OrderStatusCreated, OrderStatusNew, OrderStatusPartiallyFilled:
			order.Status = OrderStatusCancelled
			result = order
			delete(b.openOrders, id)
		default:
			err = errors.New("error")
		}
	} else {
		err = errors.New("not found")
	}
	return
}

func (b *DeribitSim) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
	var idsToBeRemoved []string

	for _, order := range b.openOrders {
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
		delete(b.openOrders, id)
	}
	return
}

func (b *DeribitSim) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (b *DeribitSim) GetPositions(symbol string) (result []*Position, err error) {
	position, ok := b.positions[symbol]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = []*Position{position}
	return
}

func (b *DeribitSim) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	return nil
}

func (b *DeribitSim) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	return nil
}

func (b *DeribitSim) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	return nil
}

func (b *DeribitSim) SubscribePositions(market Market, callback func(positions []*Position)) error {
	return nil
}

func (b *DeribitSim) SetExchangeLogger(l ExchangeLogger) {
	b.eLog = l
}

func (b *DeribitSim) RunEventLoopOnce() (err error) {
	var changed bool
	for _, order := range b.openOrders {
		changed, err = b.matchOrder(order, false)
		if changed {
			b.eLog.Warnw("Match order",
				SimEventKey, SimEventDeal,
				"order", order,
				"orderbook", b.data.GetOrderBook())
		}
	}
	return
}

func NewDeribitSim(data *dataloader.Data, cash float64, makerFeeRate float64, takerFeeRate float64) *DeribitSim {
	return &DeribitSim{
		data:          data,
		balance:       cash,
		makerFeeRate:  makerFeeRate, // -0.00025 // Maker 费率
		takerFeeRate:  takerFeeRate, // 0.00075	// Taker 费率
		orders:        make(map[string]*Order),
		openOrders:    make(map[string]*Order),
		historyOrders: make(map[string]*Order),
		positions:     make(map[string]*Position),
	}
}
