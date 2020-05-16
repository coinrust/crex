package main

import (
	"bufio"
	"flag"
	"github.com/spf13/cast"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	input  string
	output string
	sample int

	depth    int // 10档
	mode     int // 0-42列 1-41列
	database string
	exchange string
	symbol   string

	fileCache = map[string]*os.File{}

	store *Store
	n     int
)

// --m 1

func main() {
	flag.StringVar(&input, "i", `../../data-samples/deribit/deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv`, "")
	flag.StringVar(&output, "o", `./BTC-PERPETUAL`, "")
	flag.IntVar(&sample, "s", 0, "")
	flag.IntVar(&depth, "d", 10, "")
	flag.IntVar(&mode, "m", 1, "")
	flag.StringVar(&database, "db", "tick_db", "")
	flag.StringVar(&exchange, "e", "deribit", "")
	flag.StringVar(&symbol, "symbol", "BTC-PERPETUAL", "")
	flag.Parse()

	store = NewStore("mongodb://localhost:27017",
		database)

	defer store.Close()

	store.GetCollection(exchange, symbol, true)

	// headerString = "symbol,timestamp,asks[0].price,asks[0].amount,asks[1].price,asks[1].amount,asks[2].price,asks[2].amount,asks[3].price,asks[3].amount,asks[4].price,asks[4].amount,asks[5].price,asks[5].amount,asks[6].price,asks[6].amount,asks[7].price,asks[7].amount,asks[8].price,asks[8].amount,asks[9].price,asks[9].amount,bids[0].price,bids[0].amount,bids[1].price,bids[1].amount,bids[2].price,bids[2].amount,bids[3].price,bids[3].amount,bids[4].price,bids[4].amount,bids[5].price,bids[5].amount,bids[6].price,bids[6].amount,bids[7].price,bids[7].amount,bids[8].price,bids[8].amount,bids[9].price,bids[9].amount"
	//timestamp(ms),ask_price_0,ask_size_0,bid_price_0,bid_size_0,ask_price_1,ask_size_1,bid_price_1,bid_size_1,...
	n = 0
	err := ReadFile(input, handleLine)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func ReadFile(filePath string, handle func(string) bool) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		rawLine, _, err := buf.ReadLine()
		line := strings.TrimSpace(string(rawLine))
		if !handle(line) {
			break
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func handleLine(s string) bool {
	ss := strings.Split(s, ",")
	if mode == 0 && len(ss) != 42 || mode == 1 && len(ss) != 41 {
		log.Printf("[line: %v]", s)
		return false
	}

	// 或略标题行
	if mode == 0 && ss[1] == "timestamp" {
		return true
	} else if mode == 1 && ss[0] == "t" {
		return true
	}

	writeTo(symbol, ss...)

	if sample != 0 {
		n++
		if n >= sample {
			return false
		}
	}
	return true
}

func writeTo(symbol string, ss ...string) {
	// BTC-PERPETUAL,2019-10-01T00:00:00.531Z,8304.5,10270,8305,60,8305.5,1220,8306,80,8307,200,8307.5,20370,8308,65760,8308.5,120000,8309,38400,8309.5,8400,8304,185010,8303.5,52200,8303,20600,8302.5,4500,8302,2000,8301.5,18200,83 01,18000,8300.5,5090,8300,71320,8299.5,310
	//log.Printf("line=%v", ss)
	var time2 time.Time
	var err error
	if mode == 0 {
		time2, err = time.Parse("2006-01-02T15:04:05.000Z", ss[1])
		if err != nil {
			log.Fatal(err)
		}
		symbol = ss[0]
	} else if mode == 1 {
		timestamp, err := strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		time2 = time.Unix(0, timestamp*int64(time.Millisecond))
	}

	if err != nil {
		log.Fatal(err)
	}

	// 10档
	var asks []Item
	var bids []Item

	var offset = 2
	if mode == 1 {
		offset = 1
	}

	for i := offset; i < offset+20; i += 2 {
		price := cast.ToFloat64(ss[i])
		amount := cast.ToFloat64(ss[i+1])
		asks = append(asks, Item{
			Price:  price,
			Amount: amount,
		})
	}

	offset += 20

	for i := offset; i < offset+20; i += 2 {
		price := cast.ToFloat64(ss[i])
		amount := cast.ToFloat64(ss[i+1])
		bids = append(bids, Item{
			Price:  price,
			Amount: amount,
		})
	}

	ob := &OrderBook{
		Symbol:    symbol,
		Timestamp: time2,
		Asks:      asks,
		Bids:      bids,
	}
	store.Insert(ob)
}
