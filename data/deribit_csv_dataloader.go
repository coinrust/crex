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

type DeribitCsvDataLoader struct {
	file        *os.File
	reader      *bufio.Reader
	filename    string
	hasMoreData bool
}

func (l *DeribitCsvDataLoader) ReadData() (result []*models.Tick) {
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

func (l *DeribitCsvDataLoader) HasMoreData() bool {
	return l.hasMoreData
}

func (l *DeribitCsvDataLoader) open() {
	var err error
	l.file, err = os.Open(l.filename)
	if err != nil {
		log.Fatal(err)
	}

	l.reader = bufio.NewReader(l.file)
}

func (l *DeribitCsvDataLoader) close() {
	l.file.Close()
	l.hasMoreData = false
}

func (l *DeribitCsvDataLoader) readLine(line string) (result *models.Tick, ok bool) {
	ss := strings.Split(line, ",")
	if len(ss) != 42 {
		log.Printf("End [line: %v]", line)
		return
	}

	// 或略标题行
	if ss[1] == "timestamp" {
		ok = true
		return
	}

	//log.Printf(ss[1])
	timestamp, err := time.Parse("2006-01-02T15:04:05.000Z", ss[1])
	if err != nil {
		log.Fatal(err)
	}

	ask, _ := strconv.ParseFloat(ss[2], 64)
	askAmount, _ := strconv.ParseFloat(ss[3], 64)
	bid, _ := strconv.ParseFloat(ss[2+20], 64)
	bidAmount, _ := strconv.ParseFloat(ss[3+20], 64)

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

func NewDeribitCsvDataLoader(filename string) *DeribitCsvDataLoader {
	loader := &DeribitCsvDataLoader{filename: filename, hasMoreData: true}
	loader.open()
	return loader
}
