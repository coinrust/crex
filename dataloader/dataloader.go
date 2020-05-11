package dataloader

import (
	. "github.com/coinrust/crex"
	"time"
)

type DataLoader interface {
	Setup(start time.Time, end time.Time) error
	ReadData() []*OrderBook
	HasMoreData() bool
}
