package dataloader

import (
	. "github.com/coinrust/crex"
)

type DataLoader interface {
	ReadData() []*OrderBook
	HasMoreData() bool
}
