package dataloader

import (
	. "github.com/coinrust/crex"
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

func (d *Data) Reset() {
	d.readMore()
	d.index = 0
	d.offset = 0
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
		//d.maxIndex += n
		d.maxIndex = n - 1
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
	//d.data = append(d.data, data...)
	d.offset += len(d.data)
	d.data = data
	return len(data)
}
