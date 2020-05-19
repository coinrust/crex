package dataloader

import (
	. "github.com/coinrust/crex"
	"time"
)

type Data struct {
	index      int
	offset     int // 数据偏移量，使用过的数据被清理。这个值记录被清理的数据量
	maxIndex   int
	data       []*OrderBook
	dataLoader DataLoader
}

func (d *Data) Len() int {
	return len(d.data)
}

func (d *Data) GetIndex() int {
	return d.index + d.offset
}

func (d *Data) GetMaxIndex() int {
	return d.maxIndex + d.offset
}

func (d *Data) Reset(start time.Time, end time.Time) {
	d.dataLoader.Setup(start, end)
	d.readMore()
	d.index = 0
	d.offset = 0
	d.maxIndex = len(d.data) - 1
}

func (d *Data) GetOrderBook() *OrderBook {
	return d.data[d.index]
}

func (d *Data) GetRecords(size int) []*Record {
	return nil
}

func (d *Data) Next() bool {
	if d.index < d.maxIndex {
		d.index++
		return true
	}
	if n := d.readMore(); n > 0 {
		d.index = 0
		d.maxIndex = n - 1
		return true
	}
	return false
}

func (d *Data) readMore() int {
	if !d.dataLoader.HasMoreData() {
		return 0
	}
	data := d.dataLoader.ReadOrderBooks()
	if len(data) == 0 {
		return 0
	}
	d.offset += len(d.data)
	d.data = data
	return len(data)
}

func NewData(loader DataLoader) *Data {
	return &Data{
		index:      0,
		maxIndex:   0,
		data:       nil,
		dataLoader: loader,
	}
}
