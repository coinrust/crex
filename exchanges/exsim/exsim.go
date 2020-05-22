package exsim

import (
	"errors"
	"fmt"
	"github.com/chuckpreslar/emission"
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

// 持仓，用于多仓
type Positions []*Position // 单向持仓只有一项; 双向持仓 Index: 0-Long Index: 1-Short

// ExSim the exchange for backtest
type ExSim struct {
	data           *dataloader.Data
	makerFeeRate   float64 // -0.00025	// Maker fee rate
	takerFeeRate   float64 // 0.00075	// Taker fee rate
	hedgedPosition bool    // 双向持仓
	balance        float64
	orders         map[string]*Order     // All orders key: OrderID value: Order
	openOrders     map[string]*Order     // Open orders
	historyOrders  map[string]*Order     // History orders
	positions      map[string]*Positions // Position key: symbol

	emitter  *emission.Emitter
	backtest IBacktest
	eLog     ExchangeLogger
}

func (b *ExSim) GetName() (name string) {
	return "exsim"
}

func (b *ExSim) GetTime() (tm int64, err error) {
	tm = b.backtest.GetTime().UnixNano() / int64(time.Millisecond)
	return
}

func (b *ExSim) GetBalance(symbol string) (result *Balance, err error) {
	result = &Balance{}
	result.Available = b.balance
	positions := b.getPositions(symbol)
	ob := b.data.GetOrderBook()

	result.Equity = result.Available
	for _, position := range *positions {
		var price float64
		side := position.Side()
		if side == Buy {
			price = ob.AskPrice()
		} else if side == Sell {
			price = ob.BidPrice()
		}
		pnl, _ := CalcPnl(side, math.Abs(position.Size), position.AvgPrice, price)
		result.Equity += pnl
	}
	return
}

func (b *ExSim) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
	result = b.data.GetOrderBook()
	return
}

func (b *ExSim) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []*Record, err error) {
	return
}

func (b *ExSim) SetContractType(pair string, contractType string) (err error) {
	return
}

func (b *ExSim) GetContractID() (symbol string, err error) {
	return
}

func (b *ExSim) SetLeverRate(value float64) (err error) {
	return
}

func (b *ExSim) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *ExSim) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *ExSim) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *ExSim) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result *Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *ExSim) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
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

	// 如果是减仓单，判断是否足够
	if params.ReduceOnly {
		side := b.getOrderSide(order)
		position := b.getPosition(symbol, side)
		if direction == Buy &&
			(position.Size >= 0 || math.Abs(position.Size) < size) { // 平空
			err = ErrInvalidAmount
			return
		} else if direction == Sell &&
			(position.Size <= 0 || math.Abs(position.Size) < size) { // 平多
			err = ErrInvalidAmount
			return
		}
	}

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

	b.logOrderInfo("Place order", SimEventOrder, order)

	var orders = []*Order{order}
	b.emitter.Emit(WSEventOrder, orders)

	return
}

// 撮合成交
func (b *ExSim) matchOrder(order *Order, immediate bool) (match bool, err error) {
	switch order.Type {
	case OrderTypeMarket:
		match, err = b.matchMarketOrder(order)
	case OrderTypeLimit:
		match, err = b.matchLimitOrder(order, immediate)
	}
	return
}

func (b *ExSim) matchMarketOrder(order *Order) (changed bool, err error) {
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

	side := b.getOrderSide(order)
	position := b.getPosition(order.Symbol, side)

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
	var filledAmount float64
	var avgPrice float64

	// 市价成交
	if order.Direction == Buy {
		maxSize = margin * 100 * ob.AskPrice()
		if order.Amount > maxSize {
			err = errors.New(fmt.Sprintf("rejected, maximum size of future position is %v", maxSize))
			return
		}

		filledAmount, avgPrice = b.matchBid(order.Amount, ob.Asks...)

		// trade fee
		fee := filledAmount / avgPrice * b.takerFeeRate

		// Update balance
		b.addBalance(-fee)

		// Update position
		pnl := b.updatePosition(order.Symbol, filledAmount, avgPrice, side)
		order.Pnl += pnl
		order.Commission += fee
		order.AvgPrice = avgPrice
	} else if order.Direction == Sell {
		maxSize = margin * 100 * ob.BidPrice()
		if order.Amount > maxSize {
			err = errors.New(fmt.Sprintf("rejected, maximum size of future position is %v", maxSize))
			return
		}

		filledAmount, avgPrice = b.matchBid(order.Amount, ob.Bids...)

		// trade fee
		fee := filledAmount / avgPrice * b.takerFeeRate

		// Update balance
		b.addBalance(-fee)

		// Update position
		pnl := b.updatePosition(order.Symbol, -filledAmount, avgPrice, side)
		order.Pnl += pnl
		order.Commission += fee
		order.AvgPrice = avgPrice
	}
	order.FilledAmount = filledAmount
	order.UpdateTime = ob.Time
	order.Status = OrderStatusFilled
	changed = true
	return
}

func (b *ExSim) getOrderSide(order *Order) (side int) {
	if b.hedgedPosition {
		switch order.Direction {
		case Buy:
			if order.ReduceOnly {
				side = 1
			} else {
				side = 0
			}
		case Sell:
			if order.ReduceOnly {
				side = 0
			} else {
				side = 1
			}
		}
	}
	return
}

func (b *ExSim) matchLimitOrder(order *Order, immediate bool) (match bool, err error) {
	if !order.IsOpen() {
		return
	}

	side := b.getOrderSide(order)

	ob := b.data.GetOrderBook()
	if order.Direction == Buy { // Bid order
		filledAmount, avgPrice := b.matchBid(order.Amount, ob.Asks...)
		//if order.Price < ob.AskPrice() {
		if filledAmount == 0 {
			return
		}

		if immediate && order.PostOnly {
			order.UpdateTime = ob.Time
			order.Status = OrderStatusRejected
			match = true
			return
		}

		// match trade
		var fee float64

		// trade fee
		if immediate {
			fee = filledAmount / avgPrice * b.takerFeeRate
		} else {
			fee = filledAmount / avgPrice * b.makerFeeRate
		}

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, filledAmount, avgPrice, side)

		order.Commission -= fee
		order.AvgPrice = avgPrice
		order.FilledAmount = filledAmount
		order.UpdateTime = ob.Time
		order.Status = OrderStatusFilled
		match = true
	} else { // Ask order
		filledAmount, avgPrice := b.matchBid(order.Amount, ob.Asks...)
		//if order.Price > ob.BidPrice() {
		if filledAmount == 0 {
			return
		}

		if immediate && order.PostOnly {
			order.UpdateTime = ob.Time
			order.Status = OrderStatusRejected
			match = true
			return
		}

		// match trade
		var fee float64

		// trade fee
		if immediate {
			fee = filledAmount / avgPrice * b.takerFeeRate
		} else {
			fee = filledAmount / avgPrice * b.makerFeeRate
		}

		// Update balance
		b.addBalance(-fee)

		// Update position
		b.updatePosition(order.Symbol, -filledAmount, avgPrice, side)

		order.Commission -= fee
		order.AvgPrice = avgPrice
		order.FilledAmount = filledAmount
		order.UpdateTime = ob.Time
		order.Status = OrderStatusFilled
		match = true
	}
	return
}

func (b *ExSim) matchBid(size float64, asks ...Item) (filledSize float64, avgPrice float64) {
	type item = struct {
		Amount float64
		Price  float64
	}

	var items []item
	lSize := size
	for i := 0; i < len(asks); i++ {
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

func (b *ExSim) matchAsk(size float64, bids ...Item) (filledSize float64, avgPrice float64) {
	type item = struct {
		Amount float64
		Price  float64
	}

	var items []item
	lSize := size
	for i := 0; i < len(bids); i++ {
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

// 更新持仓
// side: 0-Long 1-Short
func (b *ExSim) updatePosition(symbol string, size float64, price float64, side int) (pnl float64) {
	positions := b.getPositions(symbol)
	if positions == nil {
		log.Fatalf("position error symbol=%v", symbol)
	}

	position := (*positions)[side]
	if position.Size > 0 && size < 0 || position.Size < 0 && size > 0 {
		pnl, _ = b.closePosition(position, size, price)
	} else {
		b.addPosition(position, size, price)
	}
	return
}

// 增加持仓
func (b *ExSim) addPosition(position *Position, size float64, price float64) (err error) {
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
func (b *ExSim) closePosition(position *Position, size float64, price float64) (pnl float64, err error) {
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
		pnl, _ = CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		b.addPnl(pnl)
		position.AvgPrice = price
		position.Size = position.Size + size
	} else if remaining == 0 {
		// 完全平仓
		pnl, _ = CalcPnl(position.Side(), math.Abs(size), position.AvgPrice, price)
		b.addPnl(pnl)
		position.AvgPrice = 0
		position.Size = 0
	} else {
		// 部分平仓
		pnl, _ = CalcPnl(position.Side(), math.Abs(position.Size), position.AvgPrice, price)
		b.addPnl(pnl)
		//position.AvgPrice = position.AvgPrice
		position.Size = position.Size + size
	}
	return
}

// 增加Balance
func (b *ExSim) addBalance(value float64) {
	b.balance += value
}

// 增加P/L
func (b *ExSim) addPnl(pnl float64) {
	b.balance += pnl
}

func (b *ExSim) getPosition(symbol string, side int) *Position {
	if !b.hedgedPosition {
		side = 0
	}
	positions := b.getPositions(symbol)
	return (*positions)[side]
}

// 获取持仓
func (b *ExSim) getPositions(symbol string) *Positions {
	if positions, ok := b.positions[symbol]; ok {
		return positions
	} else {
		if b.hedgedPosition {
			positions = &Positions{
				&Position{
					Symbol:    symbol,
					OpenTime:  time.Time{},
					OpenPrice: 0,
					Size:      0,
					AvgPrice:  0,
				},
				&Position{
					Symbol:    symbol,
					OpenTime:  time.Time{},
					OpenPrice: 0,
					Size:      0,
					AvgPrice:  0,
				},
			}
		} else {
			positions = &Positions{
				&Position{
					Symbol:    symbol,
					OpenTime:  time.Time{},
					OpenPrice: 0,
					Size:      0,
					AvgPrice:  0,
				},
			}
		}
		b.positions[symbol] = positions
		return positions
	}
}

func (b *ExSim) GetOpenOrders(symbol string, opts ...OrderOption) (result []*Order, err error) {
	for _, v := range b.openOrders {
		if v.Symbol == symbol {
			result = append(result, v)
		}
	}
	return
}

func (b *ExSim) GetOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
	order, ok := b.orders[id]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = order
	return
}

func (b *ExSim) CancelOrder(symbol string, id string, opts ...OrderOption) (result *Order, err error) {
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

func (b *ExSim) CancelAllOrders(symbol string, opts ...OrderOption) (err error) {
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

func (b *ExSim) AmendOrder(symbol string, id string, price float64, size float64, opts ...OrderOption) (result *Order, err error) {
	return
}

func (b *ExSim) GetPositions(symbol string) (result []*Position, err error) {
	positions, ok := b.positions[symbol]
	if !ok {
		err = errors.New("not found")
		return
	}
	result = *positions
	return
}

func (b *ExSim) SubscribeTrades(market Market, callback func(trades []*Trade)) error {
	return nil
}

func (b *ExSim) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	return nil
}

func (b *ExSim) SubscribeOrders(market Market, callback func(orders []*Order)) error {
	return nil
}

func (b *ExSim) SubscribePositions(market Market, callback func(positions []*Position)) error {
	return nil
}

func (b *ExSim) SetBacktest(backtest IBacktest) {
	b.backtest = backtest
}

func (b *ExSim) SetExchangeLogger(l ExchangeLogger) {
	b.eLog = l
}

func (b *ExSim) RunEventLoopOnce() (err error) {
	var match bool
	for _, order := range b.openOrders {
		match, err = b.matchOrder(order, false)
		if match {
			b.logOrderInfo("Match order", SimEventDeal, order)
			var orders = []*Order{order}
			b.emitter.Emit(WSEventOrder, orders)
		}
	}
	return
}

func (b *ExSim) logOrderInfo(msg string, event string, order *Order) {
	if b.eLog == nil {
		return
	}

	ob := b.data.GetOrderBook()
	positions := b.getPositions(order.Symbol)
	b.eLog.Infow(msg,
		SimEventKey, event,
		"order", order,
		"orderbook", ob,
		"balance", b.balance,
		"positions", *positions)
}

// NewExSim 创建模拟交易所
// cash: 初始资金
// makerFeeRate: Maker 费率
// takerFeeRate: Taker 费率
// hedgedPosition: 双向持仓
func NewExSim(data *dataloader.Data, cash float64, makerFeeRate float64, takerFeeRate float64, hedgedPosition bool) *ExSim {
	return &ExSim{
		data:           data,
		balance:        cash,
		makerFeeRate:   makerFeeRate, // -0.00025 // Maker 费率
		takerFeeRate:   takerFeeRate, // 0.00075	// Taker 费率
		hedgedPosition: hedgedPosition,
		orders:         make(map[string]*Order),
		openOrders:     make(map[string]*Order),
		historyOrders:  make(map[string]*Order),
		positions:      make(map[string]*Positions),
		emitter:        emission.NewEmitter(),
	}
}
