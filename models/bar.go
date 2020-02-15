package models

type Bar struct {
	Event
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}
