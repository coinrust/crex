package data

import (
	"bufio"
	"github.com/coinrust/gotrader/models"
	"io"
	"log"
	"os"
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

func (l *CsvDataLoader) ReadData() (result []*models.Tick) {
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
		if !ok {
			l.close()
			break
		}
		if tick == nil {
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

func (l *CsvDataLoader) readLine(line string) (result *models.Tick, ok bool) {
	ss := strings.Split(line, ",")
	if len(ss) != 41 {
		//log.Printf("End [line: %v]", line)
		return
	}

	// 或略标题行
	if ss[0] == "t" {
		ok = true
		return
	}

	t, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	timestamp := time.Unix(0, t*1e6)

	ask, _ := strconv.ParseFloat(ss[1], 64)
	askAmount, _ := strconv.ParseFloat(ss[2], 64)
	bid, _ := strconv.ParseFloat(ss[1+20], 64)
	bidAmount, _ := strconv.ParseFloat(ss[2+20], 64)

	// log.Printf("Ask,AskAmount,Bid,BidAmount=%v/%v/%v/%v", ask, askAmount, bid, bidAmount)

	result = &models.Tick{
		Timestamp: timestamp,
		Bid:       bid,
		Ask:       ask,
		BidVolume: int64(bidAmount),
		AskVolume: int64(askAmount),
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
