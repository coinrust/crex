package models

import "time"

// Position 持仓
type Position struct {
	Symbol    string    `json:"symbol"`     // 标
	OpenI     time.Time `json:"open_i"`     // 开仓时间
	OpenPrice float64   `json:"open_price"` // 开仓价
	Size      float64   `json:"size"`       // 仓位大小
	AvgPrice  float64   `json:"avg_price"`  // 平均价
}

func (p *Position) Side() Direction {
	if p.Size > 0 {
		return Buy
	} else if p.Size < 0 {
		return Sell
	}
	return Buy
}

// Amount 持仓量
func (p *Position) Amount() float64 {
	if p.IsLong() {
		return p.Size
	}
	return -p.Size
}

// IsOpen 是否持仓
func (p *Position) IsOpen() bool {
	return p.Size != 0
}

// IsLong 是否多仓
func (p *Position) IsLong() bool {
	return p.Size > 0
}

// IsShort 是否空仓
func (p *Position) IsShort() bool {
	return p.Size < 0
}
