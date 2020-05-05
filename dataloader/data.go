package dataloader

import (
	. "github.com/coinrust/crex"
)

type Data struct {
	index      int
	maxIndex   int
	data       []*OrderBook
	dataLoader DataLoader
}

func (d *Data) Len() int {
	return len(d.data)
}

func (d *Data) GetIndex() int {
	return d.index
}

func (d *Data) GetMaxIndex() int {
	return d.maxIndex
}

func (d *Data) Reset() {
	d.readMore()
	d.index = 0
	d.maxIndex = len(d.data) - 1
}

func (d *Data) GetOrderBook() *OrderBook {
	return d.data[d.index]
}

func (d *Data) Next() bool {
	if d.index < d.maxIndex {
		d.index++
		return true
	}
	if n := d.readMore(); n > 0 {
		d.maxIndex += n
		return true
	}
	return false
}

func (d *Data) readMore() int {
	if !d.dataLoader.HasMoreData() {
		return 0
	}
	data := d.dataLoader.ReadData()
	if len(data) == 0 {
		return 0
	}
	d.data = append(d.data, data...)
	return len(data)
}
