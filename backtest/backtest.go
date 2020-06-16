//go:generate statik -f -src=./ -include=*.html
package backtest

import (
	. "github.com/coinrust/crex"
	_ "github.com/coinrust/crex/backtest/statik"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/log"
	"github.com/coinrust/crex/utils"
	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/datatypes"
	"github.com/json-iterator/go"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	slog "log"
	"os"
	"path/filepath"
	"time"
)

const (
	OriginEChartsJs = "https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"
	MyEChartsJs     = "https://cdnjs.cloudflare.com/ajax/libs/echarts/4.7.0/echarts.min.js"

	OriginEChartsBulmaCss = "https://go-echarts.github.io/go-echarts-assets/assets/bulma.min.css"
	MyEChartsBulmaCss     = "https://cdnjs.cloudflare.com/ajax/libs/bulma/0.8.2/css/bulma.min.css"

	SimpleDateTimeFormat = "2006-01-02 15:04:05.000"
)

var (
	reportHistoryTemplate string
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

type PlotData struct {
	NameItems []string
	Prices    []float64
	Equities  []float64
}

type DataState struct {
	PrevTime int64 // ns
	Time     int64 // ns
	Index    int   // datas 中的索引
}

type Backtest struct {
	datas          []*dataloader.Data
	symbol         string
	strategyTester *StrategyTester
	baseOutputDir  string
	outputDir      string

	start         time.Time // 开始时间
	end           time.Time // 结束时间
	currentTimeNS int64     // ns
	timeNsDatas   []int64

	startedAt time.Time // 运行开始时间
	endedAt   time.Time // 运行结束时间
}

// NewBacktest 创建回测
// datas: 数据
// symbol: 标
// start: 起始时间
// end: 结束时间
// strategyHold: 策略和交易所
// outputDir: 回测输出目录
func NewBacktestFromParams(datas []*dataloader.Data, symbol string, start time.Time, end time.Time, strategyParams *StrategyTesterParams, outputDir string) *Backtest {
	strategyTester := &StrategyTester{
		StrategyTesterParams: strategyParams,
	}
	b := &Backtest{
		datas:          datas,
		symbol:         symbol,
		start:          start,
		end:            end,
		strategyTester: strategyTester,
		baseOutputDir:  outputDir,
	}

	strategyTester.backtest = b

	if err := strategyTester.Setup(); err != nil {
		panic(err)
	}

	return b
}

// NewBacktest 创建回测
// datas: 数据
// symbol: 标
// start: 起始时间
// end: 结束时间
// strategy: 策略
// exchanges: 交易所对象
// outputDir: 回测输出目录
func NewBacktest(datas []*dataloader.Data, symbol string, start time.Time, end time.Time, strategy Strategy, exchanges []ExchangeSim, outputDir string) *Backtest {
	b := &Backtest{
		datas:         datas,
		symbol:        symbol,
		start:         start,
		end:           end,
		baseOutputDir: outputDir,
	}

	strategyTester := &StrategyTester{
		StrategyTesterParams: &StrategyTesterParams{
			strategy:  strategy,
			exchanges: exchanges,
		},
		backtest: b,
	}

	if err := strategyTester.Setup(); err != nil {
		panic(err)
	}

	b.strategyTester = strategyTester

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
	strategyTester := b.strategyTester

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

	strategyTester.Init()

	b.startedAt = time.Now()

	for _, data := range b.datas {
		data.Reset(b.start, b.end)
	}

	if !b.next() {
		log.Error("error")
		return
	}

	// 初始净值
	strategyTester.addInitItemStats()
	strategyTester.OnInit()

	for {
		strategyTester.OnTick()
		strategyTester.RunEventLoopOnce()
		strategyTester.addItemStats()
		if !b.next() {
			break
		}
	}

	// Exit
	strategyTester.OnExit()

	// Sync logs
	strategyTester.Sync()

	log.Sync()

	b.endedAt = time.Now()
}

// 新的 next 方法
func (b *Backtest) next() bool {
	if len(b.datas) == 1 {
		return b.nextOne()
	}

	if b.currentTimeNS == 0 {
		for _, data := range b.datas {
			if !data.Next() {
				return false
			}
		}

		b.resetSortedDatas()

		// 取时间最大项
		b.currentTimeNS = b.timeNsDatas[len(b.timeNsDatas)-1]
		n := len(b.datas)
		for i := 0; i < n; i++ {
			data := b.datas[i]
			for {
				if data.GetOrderBook().Time.UnixNano() >= b.currentTimeNS {
					break
				}
				if !data.Next() {
					return false
				}
			}
		}
		return true
	}

	for {
		for _, timeNs := range b.timeNsDatas {
			if b.currentTimeNS < timeNs {
				b.currentTimeNS = timeNs
				if !b.ensureMoveNext(b.currentTimeNS) {
					return false
				}
				return true
			}
		}

		for _, data := range b.datas {
			if !data.Next() {
				return false
			}
		}

		b.resetSortedDatas()
	}
}

func (b *Backtest) ensureMoveNext(ns int64) bool {
	n := len(b.datas)
	count := 0
	for i := 0; i < n; i++ {
		data := b.datas[i]
		for {
			if data.GetOrderBook().Time.UnixNano() >= ns {
				break
			}
			if !data.Next() {
				return false
			}
			count++
		}
	}
	if count > 0 {
		// 重新排序
		b.resetSortedDatas()
	}
	return true
}

func (b *Backtest) resetSortedDatas() {
	nDatas := len(b.datas)
	if len(b.timeNsDatas) != nDatas*2 {
		b.timeNsDatas = make([]int64, nDatas*2)
	}

	for i := 0; i < nDatas; i++ {
		index := i * 2
		b.timeNsDatas[index] = b.datas[i].GetOrderBookRaw(1).Time.UnixNano()
		b.timeNsDatas[index+1] = b.datas[i].GetOrderBook().Time.UnixNano()
	}

	utils.SortInt64(b.timeNsDatas)
}

func (b *Backtest) nextOne() bool {
	ret := b.datas[0].Next()
	if ret {
		b.currentTimeNS = b.datas[0].GetOrderBook().Time.UnixNano()
	}
	return ret
}

//func (b *Backtest) next() bool {
//	if b.currentTimeNS == 0 {
//		return b.nextInternal()
//	}
//
//	for _, data := range b.sortedDatas {
//		if b.currentTimeNS < data.Time {
//			b.currentTimeNS = data.Time
//			return true
//		}
//	}
//
//	return b.nextInternal()
//}

//func (b *Backtest) nextInternal() bool {
//	if len(b.datas) == 1 {
//		ret := b.datas[0].Next()
//		if ret {
//			b.currentTimeNS = b.datas[0].GetOrderBook().Time.UnixNano()
//		}
//		return ret
//	}
//
//	for _, data := range b.datas {
//		if !data.Next() {
//			return false
//		}
//	}
//
//	b.setDataCache()
//
//	return true
//}

//func (b *Backtest) setDataCache() {
//	// 数据对齐，提前排序
//	n := len(b.datas)
//	if n == 0 {
//		return
//	}
//
//	for i := 0; i < n; i++ {
//		b.sortedDatas[i].Time = b.datas[i].GetOrderBook().Time.UnixNano()
//		b.sortedDatas[i].Index = i
//	}
//
//	sort.Slice(b.sortedDatas, func(i, j int) bool {
//		return b.sortedDatas[i].Time < b.sortedDatas[j].Time
//	})
//
//	if b.currentTimeNS == 0 {
//		// 第一次按右对齐
//		b.currentTimeNS = b.sortedDatas[len(b.sortedDatas)-1].Time
//	} else {
//		b.currentTimeNS = b.sortedDatas[0].Time
//	}
//}

func (b *Backtest) GetPrices() (result []float64) {
	n := len(b.datas)
	result = make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = b.datas[i].GetOrderBook().Price()
	}
	return
}

func (b *Backtest) GetLogs() LogItems {
	return b.strategyTester.GetLogs()
}

// ComputeStats Calculating Backtest Statistics
func (b *Backtest) ComputeStats() (result *Stats) {
	return b.strategyTester.ComputeStats()
}

// HTMLReport 创建Html报告文件
func (b *Backtest) HtmlReport() {
	b.strategyTester.HtmlReport()
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

	for _, v := range b.strategyTester.logs {
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
