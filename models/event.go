package models

import "time"

type Event struct {
	timestamp time.Time
	symbol    string
}

// Time returns the timestamp of an event
func (e Event) Time() time.Time {
	return e.timestamp
}

// SetTime returns the timestamp of an event
func (e *Event) SetTime(t time.Time) {
	e.timestamp = t
}

// Symbol returns the symbol string of the event
func (e Event) Symbol() string {
	return e.symbol
}

// SetSymbol returns the symbol string of the event
func (e *Event) SetSymbol(s string) {
	e.symbol = s
}
