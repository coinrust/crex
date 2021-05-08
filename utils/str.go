package utils

import "strconv"

func ParseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ParseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func ParseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
