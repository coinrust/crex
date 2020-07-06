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

	relData *Data // 关联 Data
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

func (d *Data) GetOrderBookByNS(symbol string, ns int64) *OrderBook {
	if ob := d.getOrderBookByNS(ns); ob != nil && (symbol == "" || ob.Symbol == symbol) {
		return ob
	}
	if d.relData != nil {
		return d.relData.GetOrderBookByNS(symbol, ns)
	}
	return nil
}

func (d *Data) getOrderBookByNS(ns int64) *OrderBook {
	if d.data == nil {
		return nil
	}
	ob := d.data[d.index]
	if ob.Time.UnixNano() <= ns {
		return ob
	} else if ob = d.GetOrderBookRaw(1); ob != nil && ob.Time.UnixNano() <= ns {
		return ob
	} else {
		return nil
	}
}

func (d *Data) GetOrderBook() *OrderBook {
	if d.data == nil {
		return nil
	}
	return d.data[d.index]
}

func (d *Data) GetOrderBookRaw(offset int) *OrderBook {
	if d.data == nil {
		return nil
	}
	index := d.index - offset
	if index < 0 {
		return nil
	}
	return d.data[index]
}

func (d *Data) GetRecords(size int) []*Record {
	return nil
}

func (d *Data) Next() bool {
	if d.index < d.maxIndex {
		d.index++
		return true
	}
	if o, n := d.readMore(); n > 0 {
		d.index = o
		d.maxIndex = n - 1
		return true
	}
	return false
}

func (d *Data) readMore() (offset int, count int) {
	if !d.dataLoader.HasMoreData() {
		return 0, 0
	}
	data := d.dataLoader.ReadOrderBooks()
	if len(data) == 0 {
		return 0, 0
	}
	d.offset += len(data)
	dataCount := len(d.data)
	if dataCount > 0 {
		n := 5 // 需要保留的周期
		if dataCount < n {
			n = dataCount
		}
		d.data = append(d.data[dataCount-n:], data...)
		offset = n
	} else {
		d.data = data
	}
	count = len(d.data)
	return
}

func (d *Data) GetDataRel() *Data {
	return d.relData
}

func (d *Data) SetDataRel(relData *Data) {
	d.relData = relData
}

func NewData(loader DataLoader) *Data {
	return &Data{
		index:      0,
		maxIndex:   0,
		data:       nil,
		dataLoader: loader,
	}
}
