package data

import (
	. "github.com/coinrust/gotrader"
)

type DataLoader interface {
	ReadData() []*OrderBook
	HasMoreData() bool
}
