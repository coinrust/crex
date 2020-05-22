//go:generate statik -f -src=./ -include=*.html
package backtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/coinrust/crex"
	_ "github.com/coinrust/crex/backtest/statik"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/log"
	"github.com/coinrust/crex/utils"
	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/datatypes"
	"github.com/rakyll/statik/fs"
	"github.com/tidwall/gjson"
	"io/ioutil"
	slog "log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	OriginEChartsJs = "https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"
	MyEChartsJs     = "https://cdnjs.cloudflare.com/ajax/libs/echarts/4.7.0/echarts.min.js"

	OriginEChartsBulmaCss = "https://go-echarts.github.io/go-echarts-assets/assets/bulma.min.css"
	MyEChartsBulmaCss     = "https://cdnjs.cloudflare.com/ajax/libs/bulma/0.8.2/css/bulma.min.css"
)

type PlotData struct {
	NameItems []string
	Prices    []float64
	Equities  []float64
}

type DataState struct {
	Time  int64 // ns
	Index int   // datas 中的索引
}

type Backtest struct {
	datas            []*dataloader.Data
	symbol           string
	strategy         Strategy
	exchanges        []ExchangeSim
	baseOutputDir    string
	outputDir        string
	exchangeLogFiles []string // 撮合日志记录文件

	logs LogItems

	eLoggers []ExchangeLogger

	start         time.Time // 开始时间
	end           time.Time // 结束时间
	currentTimeNS int64     // ns
	sortedDatas   []*DataState

	startedAt time.Time // 运行开始时间
	endedAt   time.Time // 运行结束时间
}

const SimpleDateTimeFormat = "2006-01-02 15:04:05.000"

var (
	reportHistoryTemplate string
)

func init() {
	statikFS, err := fs.New()
	if err != nil {
		slog.Fatal(err)
	}
	f, err := statikFS.Open("/ReportHistoryTemplate.html")
	if err != nil {
		slog.Fatal(err)
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		slog.Fatal(err)
	}
	reportHistoryTemplate = string(d)
}

// NewBacktest Create backtest
// data: The data
// outputDir: 日志输出目录
func NewBacktest(datas []*dataloader.Data, symbol string, start time.Time, end time.Time, strategy Strategy, exchanges []ExchangeSim, outputDir string) *Backtest {
	b := &Backtest{
		datas:         datas,
		symbol:        symbol,
		start:         start,
		end:           end,
		strategy:      strategy,
		baseOutputDir: outputDir,
	}
	b.exchanges = exchanges
	var exs []Exchange
	for _, v := range exchanges {
		exs = append(exs, v)
	}

	if strategy != nil {
		strategy.Setup(TradeModeBacktest, exs...)
	}

	return b
}

// SetData Set data for backtest
func (b *Backtest) SetDatas(datas []*dataloader.Data) {
	b.datas = datas
}

// GetTime get current time
func (b *Backtest) GetTime() time.Time {
	return time.Unix(0, b.currentTimeNS)
}

// Run Run backtest
func (b *Backtest) Run() {
	b.logs = LogItems{}

	SetIdGenerate(utils.NewIdGenerate(b.start))

	err := os.MkdirAll(b.baseOutputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	b.outputDir = filepath.Join(b.baseOutputDir, time.Now().Format("20060102150405"))

	logger := NewBtLogger(b,
		filepath.Join(b.outputDir, "result.log"),
		log.DebugLevel,
		false,
		true)
	log.SetLogger(logger)

	for i := 0; i < len(b.exchanges); i++ {
		b.exchanges[i].SetBacktest(b)

		path := filepath.Join(b.outputDir, fmt.Sprintf("trade_%v.log", i))
		b.exchangeLogFiles = append(b.exchangeLogFiles, path)
		eLogger := NewBtLogger(b,
			path,
			log.DebugLevel,
			true,
			false)
		b.exchanges[i].SetExchangeLogger(eLogger)
		b.eLoggers = append(b.eLoggers, eLogger)
	}

	b.startedAt = time.Now()

	for _, data := range b.datas {
		data.Reset(b.start, b.end)
	}

	// 初始化数据
	b.sortedDatas = make([]*DataState, 0, len(b.datas))
	for i := 0; i < len(b.datas); i++ {
		b.sortedDatas = append(b.sortedDatas, &DataState{
			Time:  0,
			Index: i,
		})
	}

	// 设置数据缓存
	b.setCacheData()

	// 初始净值
	item := &LogItem{
		Time:    time.Unix(0, b.currentTimeNS),
		RawTime: time.Unix(0, b.currentTimeNS),
		Prices:  b.getPrices(),
		Stats:   nil,
	}
	b.fetchItemStats(item)
	b.logs = append(b.logs, item)

	// Init
	b.strategy.OnInit()

	for {
		b.strategy.OnTick()
		b.runEventLoopOnce()
		b.addItemStats()
		if !b.next() {
			break
		}
	}

	// Exit
	b.strategy.OnExit()

	// Sync logs
	log.Sync()

	for _, v := range b.eLoggers {
		v.Sync()
	}

	b.endedAt = time.Now()
}

func (b *Backtest) next() bool {
	if b.currentTimeNS == 0 {
		return b.nextInternal()
	}

	for _, data := range b.sortedDatas {
		if b.currentTimeNS < data.Time {
			b.currentTimeNS = data.Time
			return true
		}
	}

	return b.nextInternal()
}

func (b *Backtest) nextInternal() bool {
	if len(b.datas) == 1 {
		ret := b.datas[0].Next()
		if ret {
			b.currentTimeNS = b.datas[0].GetOrderBook().Time.UnixNano()
		}
		return ret
	}

	for _, data := range b.datas {
		if !data.Next() {
			return false
		}
	}

	b.setCacheData()

	return true
}

func (b *Backtest) setCacheData() {
	// 数据对齐，提前排序
	n := len(b.datas)

	for i := 0; i < n; i++ {
		b.sortedDatas[i].Time = b.datas[i].GetOrderBook().Time.UnixNano()
		b.sortedDatas[i].Index = i
	}

	sort.Slice(b.sortedDatas, func(i, j int) bool {
		return b.sortedDatas[i].Time < b.sortedDatas[j].Time
	})

	b.currentTimeNS = b.sortedDatas[0].Time
}

func (b *Backtest) getPrices() (result []float64) {
	for _, data := range b.datas {
		result = append(result, data.GetOrderBook().Price())
	}
	return
}

func (b *Backtest) runEventLoopOnce() {
	for _, exchange := range b.exchanges {
		exchange.RunEventLoopOnce()
	}
}

func (b *Backtest) addItemStats() {
	tm := b.GetTime()
	update := false
	timestamp := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute()+1, 0, 0, time.UTC)
	var lastItem *LogItem

	if len(b.logs) > 0 {
		lastItem = b.logs[len(b.logs)-1]
		if timestamp.Unix() == lastItem.Time.Unix() {
			update = true
			return
		}
	}
	var item *LogItem
	if update {
		item = lastItem
		item.RawTime = tm
		item.Prices = b.getPrices()
		item.Stats = nil
		b.fetchItemStats(item)
	} else {
		item = &LogItem{
			Time:    timestamp,
			RawTime: tm,
			Prices:  b.getPrices(),
			Stats:   nil,
		}
		b.fetchItemStats(item)
		b.logs = append(b.logs, item)
		//log.Printf("%v / %v", tick.Timestamp, timestamp)
	}
}

func (b *Backtest) fetchItemStats(item *LogItem) {
	n := len(b.exchanges)
	for i := 0; i < n; i++ {
		balance, err := b.exchanges[i].GetBalance(b.symbol)
		if err != nil {
			panic(err)
		}
		item.Stats = append(item.Stats, LogStats{
			Balance: balance.Available,
			Equity:  balance.Equity,
		})
	}
}

func (b *Backtest) GetLogs() LogItems {
	return b.logs
}

// ComputeStats Calculating Backtest Statistics
func (b *Backtest) ComputeStats() (result *Stats) {
	result = &Stats{}

	if len(b.logs) == 0 {
		return
	}

	logs := b.logs

	n := len(logs)

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

	return
}

// HTMLReport 创建Html报告文件
func (b *Backtest) HtmlReport() {
	for _, v := range b.exchangeLogFiles {
		b.htmlReport(v)
	}
}

func (b *Backtest) htmlReport(path string) (err error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	name = name[:len(name)-len(ext)]
	//slog.Printf("%v", name)
	htmlPath := filepath.Join(dir, name+".html")
	//slog.Printf("htmlPath: %v", htmlPath)

	var sOrders []*SOrder
	var dealOrders []*SOrder
	sOrders, dealOrders, err = b.readTradeLog(path)
	if err != nil {
		return
	}

	//for _, v := range sOrders {
	//	slog.Printf("sOrders Ts: %v Order: %v OrderBook: %v Comment: %v",
	//		v.Ts, v.Order, v.OrderBook, v.Comment)
	//}

	var html string
	html, err = b.buildReportHtml(sOrders, dealOrders)
	err = ioutil.WriteFile(htmlPath, []byte(html), os.ModePerm)
	return
}

func (b *Backtest) buildReportHtml(sOrders []*SOrder, dealOrders []*SOrder) (html string, err error) {
	//var templateData []byte
	//templateData, err = ioutil.ReadFile("./ReportHistory-template.html")
	//if err != nil {
	//	slog.Printf("%v", err)
	//	return
	//}
	//reportHistoryTemplate := string(templateData)
	//slog.Printf("%v", reportHistoryTemplate)
	// <tr bgcolor="#FFFFFF" align="right"><td>2018.07.06 11:08:44</td><td>11573668</td><td>EURUSD</td><td>buy limit</td><td colspan="2">0.20 / 0.00</td><td>1.16673</td><td></td><td></td><td colspan="2">2018.07.06 11:17:24</td><td>canceled</td><td></td></tr>
	// <tr bgcolor="#F7F7F7" align="right"><td>2018.07.06 11:08:57</td><td>11573671</td><td>EURUSD</td><td>sell limit</td><td colspan="2">0.20 / 0.00</td><td>1.17106</td><td></td><td></td><td colspan="2">2018.07.06 11:13:03</td><td>canceled</td><td></td></tr>
	// <!--{order-row}-->
	html = strings.ReplaceAll(reportHistoryTemplate, "<!--{Symbol}-->", b.symbol)
	html = strings.ReplaceAll(html, "<!--{Period}-->", fmt.Sprintf("%v - %v", b.start.String(), b.end.String())) // 2018.11.01 - 2018.12.01
	// Parameters: A=1
	stats := b.ComputeStats()
	html = strings.ReplaceAll(html, "<!--Initial Equity-->", fmt.Sprint(stats.EntryEquity))
	html = strings.ReplaceAll(html, "<!--Exit Equity-->", fmt.Sprint(stats.ExitEquity))
	html = strings.ReplaceAll(html, "<!--Duration-->", stats.Duration.String())
	html = strings.ReplaceAll(html, "<!--Return-->", fmt.Sprintf("%.8f", stats.EquityReturn))
	html = strings.ReplaceAll(html, "<!--Return [%]-->", fmt.Sprintf("%.4f", stats.EquityReturnPnt*100))
	html = strings.ReplaceAll(html, "<!--Run Duration-->", stats.RunDuration.String())
	html = strings.ReplaceAll(html, "<!--Buy & Hold Return-->", fmt.Sprintf("%.8f", stats.BaHReturn))
	html = strings.ReplaceAll(html, "<!--Buy & Hold Return [%]-->", fmt.Sprintf("%.4f", stats.BaHReturnPnt*100))
	s := b.buildSOrders(sOrders)
	html = strings.Replace(html, `<!--{order-rows}-->`, s, -1)
	s = b.buildSOrders(dealOrders)
	html = strings.Replace(html, `<!--{deal-order-rows}-->`, s, -1)
	return
}

func (b *Backtest) buildSOrders(sOrders []*SOrder) string {
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
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Time.Format("2006.01.02 15:04:05.000"))) // 2018.07.06 11:08:44
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.ID))                                     // 11573668
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Symbol))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, orderType))                             // buy limit/buy
		s.WriteString(fmt.Sprintf(`<td>%v / %v</td>`, order.Amount, order.FilledAmount)) // 0.20 / 0.00
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, price))                                 // 1.16673
		if order.AvgPrice > 0 {
			s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.AvgPrice))
		} else {
			s.WriteString(`<td></td>`)
		}
		s.WriteString(fmt.Sprintf(`<td>%.8f</td>`, order.Pnl))
		s.WriteString(fmt.Sprintf(`<td>%.8f</td>`, order.Commission))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, sOrder.Balance))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.UpdateTime.Format("2006.01.02 15:04:05.000")))
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, order.Status.String())) // canceled
		s.WriteString(fmt.Sprintf(`<td>%v</td>`, positions))
		s.WriteString(`</tr>`)
	}
	return s.String()
}

func (b *Backtest) readTradeLog(path string) (orders []*SOrder, dealOrders []*SOrder, err error) {
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
		event, so, err = b.parseSOrder(s)
		if err != nil {
			return
		}
		switch event {
		case "order":
			orders = append(orders, so)
		case "deal":
			dealOrders = append(dealOrders, so)
		}
	}

	return
}

func (b *Backtest) parseSOrder(s string) (event string, so *SOrder, err error) {
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
		balance := ret.Get("balance").Float()

		err = json.Unmarshal([]byte(orderJson), &order)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(orderbookJson), &orderbook)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(positionsJson), &positions)
		if err != nil {
			return
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
			Balance:   balance,
			Comment:   msg,
		}
	}
	return
}

func (b *Backtest) priceLine(plotData *PlotData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "价格", Width: "1270px", Height: "500px"},
		charts.ToolboxOpts{Show: true},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true},
	)

	line.AddXAxis(plotData.NameItems)
	line.AddYAxis("price", plotData.Prices,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
		//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	return line
}

func (b *Backtest) equityLine(plotData *PlotData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "净值", Width: "1270px", Height: "400px"},
		charts.ToolboxOpts{Show: true},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true},
	)

	line.AddXAxis(plotData.NameItems)

	line.AddYAxis("equity", plotData.Equities,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
		//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	return line
}

// Plot Output backtest results
func (b *Backtest) Plot() {
	var plotData PlotData

	for _, v := range b.logs {
		plotData.NameItems = append(plotData.NameItems, v.Time.Format(SimpleDateTimeFormat))
		plotData.Prices = append(plotData.Prices, v.Prices[0])
		plotData.Equities = append(plotData.Equities, v.TotalEquity())
	}

	p := charts.NewPage()
	p.Add(b.priceLine(&plotData), b.equityLine(&plotData))

	filename := filepath.Join(b.outputDir, "result.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}

	replaceJSAssets(&p.JSAssets)
	replaceCssAssets(&p.CSSAssets)

	p.Render(f)
}

// 替换Js资源，使用cdn加速资源，查看网页更快
func replaceJSAssets(jsAssets *datatypes.OrderedSet) {
	for i := 0; i < len(jsAssets.Values); i++ {
		if jsAssets.Values[i] == OriginEChartsJs {
			jsAssets.Values[i] = MyEChartsJs
		}
	}
}

func replaceCssAssets(cssAssets *datatypes.OrderedSet) {
	for i := 0; i < len(cssAssets.Values); i++ {
		if cssAssets.Values[i] == OriginEChartsBulmaCss {
			cssAssets.Values[i] = MyEChartsBulmaCss
		}
	}
}

func (b *Backtest) PlotOld() {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "回测", Width: "1270px", Height: "600px"},
		charts.ToolboxOpts{Show: true},
		charts.ToolboxOpts{Show: true},
		charts.TitleOpts{Title: "回测"},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		//charts.LegendOpts{Right: "80%"},
		//charts.SplitLineOpts{Show: true},
		//charts.SplitAreaOpts{Show: true},
	)
	nameItems := make([]string, 0)
	prices := make([]float64, 0)
	equities := make([]float64, 0)

	for _, v := range b.logs {
		nameItems = append(nameItems, v.Time.Format(SimpleDateTimeFormat))
		prices = append(prices, v.Prices[0])
		equities = append(equities, v.TotalEquity())
	}

	line.AddXAxis(nameItems)
	line.AddYAxis("price", prices,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
	//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	line.AddYAxis("equity", equities,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
	//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	line.SetGlobalOptions(charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true})

	filename := filepath.Join(b.outputDir, "result.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}
	line.Render(f)
}
