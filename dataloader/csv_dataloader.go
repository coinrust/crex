package dataloader

import (
	"bufio"
	. "github.com/coinrust/crex"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CsvDataLoader struct {
	file        *os.File
	reader      *bufio.Reader
	filename    string
	hasMoreData bool
}

func (l *CsvDataLoader) Setup(start time.Time, end time.Time) error {
	return nil
}

func (l *CsvDataLoader) ReadData() (result []*OrderBook) {
	if !l.hasMoreData {
		return nil
	}

	var count int
	for {
		if count >= 10000 {
			break
		}
		rawLine, _, err := l.reader.ReadLine()
		if err != nil {
			l.close()
			if err == io.EOF {
				return
			}
			return
		}

		line := strings.TrimSpace(string(rawLine))
		tick, ok := l.readLine(line)
		if !ok || tick == nil {
			continue
		}
		result = append(result, tick)
		count++
	}

	return
}

func (l *CsvDataLoader) HasMoreData() bool {
	return l.hasMoreData
}

func (l *CsvDataLoader) open() {
	var err error
	l.file, err = os.Open(l.filename)
	if err != nil {
		log.Fatal(err)
	}

	l.reader = bufio.NewReader(l.file)
}

func (l *CsvDataLoader) close() {
	l.file.Close()
	l.hasMoreData = false
}

func (l *CsvDataLoader) readLine(line string) (result *OrderBook, ok bool) {
	ss := strings.Split(line, ",")
	n := len(ss)
	if n < 5 {
		//log.Printf("End [line: %v]", line)
		return
	}
	if (n-1)%4 != 0 {
		return
	}

	// 忽略标题行
	if ss[0] == "t" {
		ok = true
		return
	}

	t, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	timestamp := time.Unix(0, t*int64(time.Millisecond))

	nDepth := (n - 1) / 4

	var asks []Item
	var bids []Item

	bidOffset := nDepth * 2

	for i := 0; i < nDepth; i++ {
		ask, _ := strconv.ParseFloat(ss[1+2*i], 64)
		askAmount, _ := strconv.ParseFloat(ss[2+2*i], 64)
		bid, _ := strconv.ParseFloat(ss[1+2*i+bidOffset], 64)
		bidAmount, _ := strconv.ParseFloat(ss[2+2*i+bidOffset], 64)
		asks = append(asks, Item{
			Price:  ask,
			Amount: askAmount,
		})
		bids = append(bids, Item{
			Price:  bid,
			Amount: bidAmount,
		})
	}

	if asks[0].Price > asks[1].Price {
		sort.Slice(asks, func(i, j int) bool {
			return asks[i].Price < asks[j].Price
		})
	}
	if bids[0].Price < bids[1].Price {
		sort.Slice(bids, func(i, j int) bool {
			return bids[i].Price > bids[j].Price
		})
	}

	result = &OrderBook{
		Time:   timestamp,
		Symbol: "",
		Asks:   asks,
		Bids:   bids,
	}
	ok = true
	return
}

func NewCsvDataLoader(filename string) *CsvDataLoader {
	loader := &CsvDataLoader{filename: filename, hasMoreData: true}
	loader.open()
	return loader
}

func NewCsvData(filename string) *Data {
	loader := NewCsvDataLoader(filename)
	return &Data{
		index:      0,
		maxIndex:   0,
		data:       nil,
		dataLoader: loader,
	}
}
