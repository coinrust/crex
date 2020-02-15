package data

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coinrust/gotrader/models"
)

func NewDeribitData(filename string) *Data {
	loader := NewDeribitCsvDataLoader(filename)
	return &Data{
		index:      0,
		maxIndex:   0,
		data:       nil,
		dataLoader: loader,
	}
}

func LoadDeribitTickByTickBookSnapshots(filename string) (result *Data, err error) {
	result = &Data{}
	ReadFile(filename, func(s string) bool {
		ss := strings.Split(s, ",")
		if len(ss) != 42 {
			log.Printf("End [line: %v]", s)
			return false
		}

		// 或略非BTC永续合约
		if ss[0] != "BTC-PERPETUAL" {
			return true
		}

		// 或略标题行
		if ss[1] == "timestamp" {
			return true
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

		result.data = append(result.data, &models.Tick{
			Timestamp: timestamp,
			Bid:       bid,
			Ask:       ask,
			BidVolume: int64(bidAmount),
			AskVolume: int64(askAmount),
		})

		return true
	})
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
