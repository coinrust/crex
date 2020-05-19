package dataloader

import (
	. "github.com/coinrust/crex"
	"time"
)

type DataLoader interface {
	Setup(start time.Time, end time.Time) error
	ReadOrderBooks() []*OrderBook
	ReadRecords(limit int) []*Record
	HasMoreData() bool
}
