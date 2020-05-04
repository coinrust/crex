package bitmexsim

import (
	. "github.com/coinrust/crex"
)

func testExchange() Exchange {
	return NewBitMEXSim(nil, 10000, 0, 0)
}
