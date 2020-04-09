package crex

import "time"

type Record struct {
	Symbol    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}
