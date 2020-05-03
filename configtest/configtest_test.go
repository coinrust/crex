package configtest

import "testing"

func TestLoadTestConfig(t *testing.T) {
	testDataDir = "../testdata/"
	tCfg := LoadTestConfig("bybit")
	t.Logf("%#v", tCfg)
}
