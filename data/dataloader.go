package data

import "github.com/coinrust/gotrader/models"

type DataLoader interface {
	ReadData() []*models.Tick
	HasMoreData() bool
}
