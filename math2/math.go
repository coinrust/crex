package math2

import "math"

// ToFixed
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}

// ToFixedE5 类似四舍五入法，规整到 0,0.5,1.0,1.5...
// 用于 BitMEX/Deribit 等平台价格规整
func ToFixedE5(x float64) float64 {
	t := math.Trunc(x)
	if x > t+0.5 {
		t += 0.5
	}
	if d := math.Abs(x - t); d > 0.25 {
		return t + math.Copysign(0.5, x)
	}
	return t
}

// ToFixedE5P 类似四舍五入法
// XBT: precision=0 0,0.5,1.0,1.5...
// ETH: precision=1 0.05,0.10,0.15...
func ToFixedE5P(x float64, precision int) float64 {
	if precision == 0 {
		return ToFixedE5(x)
	}
	p := math.Pow(10, float64(precision))
	y := ToFixedE5(x*p) / p
	return y
}
