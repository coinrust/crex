package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/coinrust/crex/dataloader"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, `mongodb-to-csv version: v1.0.0
Usage: mongodb-to-csv [-h] 
Options:
`)
	flag.PrintDefaults()
}

func main() {
	var exchangeName string
	var symbol string
	var start string
	var end string
	var help bool
	flag.StringVar(&exchangeName, "e", "huobi", "exchange name")
	flag.StringVar(&symbol, "s", "BTC-USDT", "symbol")
	flag.StringVar(&start, "st", "2020-05-01 00:00:00", "start time, 2020-05-01 00:00:00")
	flag.StringVar(&end, "et", "2020-05-01 00:00:00", "end time, 2020-05-01 00:00:00")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	fmt.Printf("exchange name: %s\nsymbol: %s\nstart: %s\nend: %s\n", exchangeName, symbol, start, end)
	st, _ := time.Parse("2006-01-02 15:04:05", start)
	et, _ := time.Parse("2006-01-02 15:04:05", end)

	loader0 := dataloader.NewMongoDBDataLoader("mongodb://localhost:27017",
		"tick_db", exchangeName, symbol)
	data0 := dataloader.NewData(loader0)
	data0.Reset(st, et)

	file, err := os.Create(fmt.Sprintf("%s_%s_%s_%d%d.csv", exchangeName, symbol, st.Format("2006-01-02 15_04_05"), time.Now().Minute(), time.Now().Second()))

	if err != nil {
		panic(err)
	}
	targetCsv := csv.NewWriter(file)
	data := []string{"t"}
	for k := 0; k < 20; k++ {
		data = append(data, fmt.Sprintf("asks[%d].price", k))
		data = append(data, fmt.Sprintf("asks[%d].amount", k))
	}
	for k := 0; k < 20; k++ {
		data = append(data, fmt.Sprintf("bids[%d].price", k))
		data = append(data, fmt.Sprintf("bids[%d].amount", k))
	}
	targetCsv.Write(data)
	targetCsv.Flush()

	for data0.Next() {
		ob := data0.GetOrderBook()

		data := []string{strconv.Itoa(int(ob.Time.UnixNano() / int64(time.Millisecond)))}

		for kk, ask := range ob.Asks {
			if kk >= 20 {
				break
			}
			data = append(data, strconv.FormatFloat(ask.Price, 'f', -1, 64))
			data = append(data, strconv.FormatFloat(ask.Amount, 'f', -1, 64))
		}
		for kk, bid := range ob.Bids {
			if kk >= 20 {
				break
			}
			data = append(data, strconv.FormatFloat(bid.Price, 'f', -1, 64))
			data = append(data, strconv.FormatFloat(bid.Amount, 'f', -1, 64))
		}
		targetCsv.Write(data)
		targetCsv.Flush()
	}
}
