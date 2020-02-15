package models

import "time"

type Tick struct {
	// Event
	Timestamp time.Time
	Bid       float64
	Ask       float64
	BidVolume int64
	AskVolume int64
}

// Price returns the middle of Bid and Ask.
func (t Tick) Price() float64 {
	latest := (t.Bid + t.Ask) / float64(2)
	return latest
}

// Spread returns the difference or spread of Bid and Ask.
func (t Tick) Spread() float64 {
	return t.Bid - t.Ask
}
