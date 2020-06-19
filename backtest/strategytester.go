package backtest

import (
	"bytes"
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/log"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type StrategyTesterParams struct {
	strategy  Strategy
	exchanges []ExchangeSim
}

func NewStrategyTesterParams(strategy Strategy, exchanges []ExchangeSim) *StrategyTesterParams {
	return &StrategyTesterParams{
		strategy:  strategy,
		exchanges: exchanges,
	}
}

type StrategyTester struct {
	*StrategyTesterParams

	backtest  *Backtest
	logs      LogItems
	eLogFiles []string // 撮合日志记录文件
	eLoggers  []ExchangeLogger
}

func (t *StrategyTester) Setup() error {
	if t.strategy == nil {
		return nil
	}

	var exs []interface{}
	for _, v := range t.exchanges {
		exs = append(exs, v)
	}

	err := t.strategy.Setup(TradeModeBacktest, exs...)
	return err
}

func (t *StrategyTester) Init() {
	for i := 0; i < len(t.exchanges); i++ {
		t.exchanges[i].SetBacktest(t.backtest)

		path := filepath.Join(t.backtest.outputDir, fmt.Sprintf("trade_%v.log", i))
		t.eLogFiles = append(t.eLogFiles, path)
		eLogger := NewBtLogger(t.backtest,
			path,
			log.DebugLevel,
			true,
			false)
		t.exchanges[i].SetExchangeLogger(eLogger)
		t.eLoggers = append(t.eLoggers, eLogger)
	}
}

func (t *StrategyTester) addInitItemStats() {
	b := t.backtest
	item := &LogItem{
		Time:    time.Unix(0, b.currentTimeNS).Local(),
		RawTime: time.Unix(0, b.currentTimeNS).Local(),
		Prices:  b.GetPrices(),
		Stats:   nil,
	}
	t.fetchItemStats(item)
	t.logs = append(t.logs, item)
}

func (t *StrategyTester) addItemStats() {
	b := t.backtest
	tm := b.GetTime().Local()
	update := false
	timestamp := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute()+1, 0, 0, time.Local)
	var lastItem *LogItem

	if len(t.logs) > 0 {
		lastItem = t.logs[len(t.logs)-1]
		if timestamp.Unix() == lastItem.Time.Unix() {
			update = true
			return
		}
	}

	var item *LogItem
	if update {
		item = lastItem
		item.RawTime = tm
		item.Prices = b.GetPrices()
		item.Stats = nil
		t.fetchItemStats(item)
	} else {
		item = &LogItem{
			Time:    timestamp,
			RawTime: tm,
			Prices:  b.GetPrices(),
			Stats:   nil,
		}
		t.fetchItemStats(item)
		t.logs = append(t.logs, item)
		//log.Printf("%v / %v", tick.Timestamp, timestamp)
	}
}

func (t *StrategyTester) fetchItemStats(item *LogItem) {
	b := t.backtest
	n := len(t.exchanges)

	for i := 0; i < n; i++ {
		// 期货
		ex, ok := t.exchanges[i].(Exchange)
		if ok {
			balance, err := ex.GetBalance(b.symbol)
			if err != nil {
				panic(err)
			}
			item.Stats = append(item.Stats, LogStats{
				Balance: balance.Available,
				Equity:  balance.Equity,
			})
		}

		// 现货
		spotEx, ok := t.exchanges[i].(SpotExchange)
		if ok {
			balance, err := spotEx.GetBalance(b.symbol)
			if err != nil {
				panic(err)
			}
			total := balance.Base.Available + balance.Base.Frozen - balance.Base.Borrow
			totalQuote := balance.Quote.Available + balance.Quote.Frozen - balance.Quote.Borrow
			price := item.Prices[i]
			item.Stats = append(item.Stats, LogStats{
				Balance: totalQuote,
				Equity:  total*price + totalQuote,
			})
		}
	}
}

func (t *StrategyTester) GetLogs() LogItems {
	return t.logs
}

func (t *StrategyTester) OnInit() error {
	return t.strategy.OnInit()
}

func (t *StrategyTester) OnTick() error {
	return t.strategy.OnTick()
}

func (t *StrategyTester) OnExit() error {
	return t.strategy.OnExit()
}

func (t *StrategyTester) RunEventLoopOnce() {
	for _, v := range t.exchanges {
		v.RunEventLoopOnce()
	}
}

func (t *StrategyTester) Sync() {
	for _, v := range t.eLoggers {
		v.Sync()
	}
}

// ComputeStats Calculating Backtest Statistics
func (t *StrategyTester) ComputeStats() (result *Stats) {
	result = &Stats{}

	logs := t.logs
	n := len(logs)

	if n == 0 {
		return
	}

	b := t.backtest

	result.Start = logs[0].Time
	result.End = logs[n-1].Time
	result.Duration = result.End.Sub(result.Start)
	result.RunDuration = b.endedAt.Sub(b.startedAt)
	result.EntryPrice = logs[0].Prices[0]
	result.ExitPrice = logs[n-1].Prices[0]
	result.EntryEquity = logs[0].TotalEquity()
	result.ExitEquity = logs[n-1].TotalEquity()
	result.BaHReturn = (result.ExitPrice - result.EntryPrice) / result.EntryPrice * result.EntryEquity
	result.BaHReturnPnt = (result.ExitPrice - result.EntryPrice) / result.EntryPrice
	result.EquityReturn = result.ExitEquity - result.EntryEquity
	result.EquityReturnPnt = result.EquityReturn / result.EntryEquity
	result.AnnReturn = t.CalAnnReturn(result)
	result.MaxDrawDown = t.CalMaxDrawDown()

	return
}

// 计算年化收益
func (t *StrategyTester) CalAnnReturn(s *Stats) float64 {
	days := s.Duration.Hours() / 24.0
	if days < 2.0 { // 小于2天直接返回0
		return 0
	}
	return math.Pow(s.EquityReturnPnt+1.0, 365.0/days) - 1
}

// 计算最大回撤
func (t *StrategyTester) CalMaxDrawDown() (result float64) {
	n := len(t.logs)
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = t.logs[i].TotalEquity()
	}

	maxValueUntil := func(untilIdx int) float64 {
		maxVal := 0.0
		for i := 0; i <= untilIdx; i++ {
			if values[i] > maxVal {
				maxVal = values[i]
			}
		}
		return maxVal
	}

	maxDrawDown := 0.0

	for i := 0; i < n; i++ {
		maxVal := maxValueUntil(i - 1)
		drawDown := 1.0 - (values[i] / maxVal)
		if drawDown > 0 && drawDown > maxDrawDown {
			maxDrawDown = drawDown
		}
	}

	return maxDrawDown
}

// HTMLReport 创建Html报告文件
func (t *StrategyTester) HtmlReport() {
	for _, v := range t.eLogFiles {
		t.htmlReport(v)
	}
}

func (t *StrategyTester) htmlReport(path string) (err error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	name = name[:len(name)-len(ext)]
	//slog.Printf("%v", name)
	htmlPath := filepath.Join(dir, name+".html")
	//slog.Printf("htmlPath: %v", htmlPath)

	var sOrders []*SOrder
	var dealOrders []*SOrder
	sOrders, dealOrders, err = t.readTradeLog(path)
	if err != nil {
		return
	}

	var html string
	html, err = t.buildReportHtml(sOrders, dealOrders)
	err = ioutil.WriteFile(htmlPath, []byte(html), os.ModePerm)
	return
}

func (t *StrategyTester) buildReportHtml(sOrders []*SOrder, dealOrders []*SOrder) (html string, err error) {
	// <!--{order-row}-->
	html = strings.ReplaceAll(reportHistoryTemplate, "<!--{Symbol}-->", t.backtest.symbol)
	html = strings.ReplaceAll(html, "<!--{Period}-->", fmt.Sprintf("%v - %v", t.backtest.start.String(), t.backtest.end.String())) // 2018.11.01 - 2018.12.01
	// Parameters: A=1
	stats := t.ComputeStats()
	html = strings.ReplaceAll(html, "<!--Initial Equity-->", fmt.Sprint(stats.EntryEquity))
	html = strings.ReplaceAll(html, "<!--Exit Equity-->", fmt.Sprint(stats.ExitEquity))
	html = strings.ReplaceAll(html, "<!--Duration-->", stats.Duration.String())
	html = strings.ReplaceAll(html, "<!--Return-->", fmt.Sprintf("%.8f", stats.EquityReturn))
	html = strings.ReplaceAll(html, "<!--Return [%]-->", fmt.Sprintf("%.4f", stats.EquityReturnPnt*100))
	html = strings.ReplaceAll(html, "<!--Run Duration-->", stats.RunDuration.String())
	html = strings.ReplaceAll(html, "<!--Buy & Hold Return-->", fmt.Sprintf("%.8f", stats.BaHReturn))
	html = strings.ReplaceAll(html, "<!--Buy & Hold Return [%]-->", fmt.Sprintf("%.4f", stats.BaHReturnPnt*100))
	s := t.buildSOrders(sOrders)
	html = strings.Replace(html, `<!--{order-rows}-->`, s, -1)
	s = t.buildSOrders(dealOrders)
	html = strings.Replace(html, `<!--{deal-order-rows}-->`, s, -1)

	var orderCommissionTotal float64
	var dealCommissionTotal float64
	for _, v := range sOrders {
		orderCommissionTotal += v.Order.Commission
	}
	for _, v := range dealOrders {
		dealCommissionTotal += v.Order.Commission
	}

	var orderCommissionTotalString string
	if orderCommissionTotal != 0 {
		orderCommissionTotalString = fmt.Sprintf("%.8f", orderCommissionTotal)
	}
	var dealCommissionTotalString string
	if dealCommissionTotal != 0 {
		dealCommissionTotalString = fmt.Sprintf("%.8f", dealCommissionTotal)
	}
	html = strings.Replace(html, `<!--{order-commission-total}-->`, orderCommissionTotalString, -1)
	html = strings.Replace(html, `<!--{deal-commission-total}-->`, dealCommissionTotalString, -1)
	return
}

func (t *StrategyTester) buildSOrders(sOrders []*SOrder) string {
	s := bytes.Buffer{}
	for i := 0; i < len(sOrders); i++ {
		sOrder := sOrders[i]
		order := sOrders[i].Order
		bgColor := "#FFFFFF"
		if i%2 != 0 {
			bgColor = "#F7F7F7"
		}
		price := fmt.Sprintf("%v", order.Price)
		orderType := strings.ToLower(order.Direction.String())
		if order.Type == OrderTypeMarket {
			price = "market"
		} else {
			orderType += " " + strings.ToLower(order.Type.String())
		}
		if order.PostOnly {
			orderType += " postOnly"
		}
		if order.ReduceOnly {
			orderType += " reduceOnly"
		}
		positions := ""
		sort.Slice(sOrder.Positions, func(i, j int) bool {
			return sOrder.Positions[i].Size > sOrder.Positions[j].Size
		})
		for _, v := range sOrder.Positions {
			if positions != "" {
				positions += " / "
			}
			positions += fmt.Sprintf("%v", v.Size)
		}
		s.WriteString(fmt.Sprintf(`<tr bgcolor="%v" align="right">`, bgColor))                  // #FFFFFF
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Time.Format("2006-01-02 15:04:05.000"))) // 2018.07.06 11:08:44
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.ID))                                     // 11573668
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Symbol))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, orderType))                             // buy limit/buy
		s.WriteString(fmt.Sprintf(`<td>%v / %v</td>`, order.Amount, order.FilledAmount)) // 0.20 / 0.00
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, price))                                 // 1.16673
		var avgPriceString string
		if order.AvgPrice > 0 {
			avgPriceString = fmt.Sprintf(`%v`, order.AvgPrice)
		}
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, avgPriceString))
		var pnlString string
		if order.Pnl != 0 {
			pnlString = fmt.Sprintf("%.8f", order.Pnl)
		}
		s.WriteString(fmt.Sprintf(`<td>%s</td>`, pnlString))
		var commissionString string
		if order.Commission != 0 {
			commissionString = fmt.Sprintf(`%.8f`, order.Commission)
		}
		s.WriteString(fmt.Sprintf(`<td>%s</td>`, commissionString))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, sOrder.BalancesString()))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.UpdateTime.Format("2006-01-02 15:04:05.000")))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Status.String())) // canceled
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, positions))
		s.WriteString(`</tr>`)
	}
	return s.String()
}

func (t *StrategyTester) readTradeLog(path string) (orders []*SOrder, dealOrders []*SOrder, err error) {
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	ss := strings.Split(string(data), "\n")

	for _, s := range ss {
		if s == "" {
			continue
		}
		var event string
		var so *SOrder
		event, so, err = t.parseSOrder(s)
		if err != nil {
			return
		}
		switch event {
		case SimEventOrder:
			orders = append(orders, so)
		case SimEventDeal:
			dealOrders = append(dealOrders, so)
		}
	}

	return
}

func (t *StrategyTester) parseSOrder(s string) (event string, so *SOrder, err error) {
	ret := gjson.Parse(s)
	if eventValue := ret.Get("event"); eventValue.Exists() {
		var order Order
		var orderbook OrderBook
		var positions []*Position

		event = eventValue.String()
		tsString := ret.Get("ts").String() // 2019-10-01T08:00:00.143+0800
		msg := ret.Get("msg").String()
		orderJson := ret.Get("order").String()
		orderbookJson := ret.Get("orderbook").String()
		positionsJson := ret.Get("positions").String()
		var balances []float64
		if v := ret.Get("balances"); v.Exists() {
			values := v.Value().([]interface{})
			for _, v := range values {
				balances = append(balances, cast.ToFloat64(v))
			}
		} else if v := ret.Get("balance"); v.Exists() {
			balances = []float64{v.Float()}
		}

		err = json.Unmarshal([]byte(orderJson), &order)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(orderbookJson), &orderbook)
		if err != nil {
			return
		}
		if positionsJson != "" {
			err = json.Unmarshal([]byte(positionsJson), &positions)
			if err != nil {
				return
			}
		}
		var ts time.Time
		ts, err = time.Parse("2006-01-02T15:04:05.000Z0700", tsString)
		if err != nil {
			return
		}
		so = &SOrder{
			Ts:        ts,
			Order:     &order,
			OrderBook: &orderbook,
			Positions: positions,
			Balances:  balances,
			Comment:   msg,
		}
	}
	return
}
