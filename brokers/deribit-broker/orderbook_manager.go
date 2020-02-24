package deribit_broker

import (
	. "github.com/coinrust/gotrader/models"
	"github.com/frankrap/deribit-api/models"
	"sync"
)

type OrderBookManager struct {
	data sync.Map // key: InstrumentName value: *OrderBookLocal
}

func NewOrderBookManager() *OrderBookManager {
	return &OrderBookManager{}
}

func (m *OrderBookManager) GetOrderBook(instrumentName string) (ob OrderBook, ok bool) {
	v, okL := m.data.Load(instrumentName)
	if !okL {
		return
	}
	if v == nil {
		return
	}
	obLocal, okC := v.(*OrderBookLocal)
	if !okC {
		return
	}
	ob = obLocal.GetOrderbook()
	ok = true
	return
}

func (m *OrderBookManager) Update(newOrderBook *models.OrderBookNotification) {
	v, ok := m.data.Load(newOrderBook.InstrumentName)
	if !ok {
		v, _ = m.data.LoadOrStore(newOrderBook.InstrumentName, NewOrderBookLocal(newOrderBook.InstrumentName))
	}
	if v == nil {
		return
	}
	obLocal, ok := v.(*OrderBookLocal)
	if !ok {
		return
	}
	obLocal.Update(newOrderBook)
}
