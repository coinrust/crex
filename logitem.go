package gotrader

import (
	"time"
)

type LogStats struct {
	Balance float64 `json:"balance"`
	Equity  float64 `json:"equity"`
}

type LogItem struct {
	Time  time.Time  `json:"time"`
	Ask   float64    `json:"ask"`
	Bid   float64    `json:"bid"`
	Stats []LogStats `json:"stats"`
}

func (i *LogItem) Price() float64 {
	return (i.Ask + i.Bid) / 2.0
}

func (i *LogItem) TotalEquity() float64 {
	var total float64
	for _, v := range i.Stats {
		total += v.Equity
	}
	return total
}

type LogItems []*LogItem
