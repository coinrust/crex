package spotsim

import (
	"github.com/coinrust/crex"
	"testing"
)

var ss = New(
	nil,
	crex.SpotBalance{
		Base:  crex.SpotAsset{"BTC", 1, 0},
		Quote: crex.SpotAsset{"USDT", 10000, 0},
	},
	0.0001,
	0.0003,
)

func TestSpotSim_GetName(t *testing.T) {
	t.Log(ss.GetName())
}
