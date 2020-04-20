package hbdmswap

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"testing"
	"time"
)

func testExchange() Exchange {
	testConfig := configtest.LoadTestConfig("hbdmswap")

	params := &Parameters{}
	params.AccessKey = testConfig.AccessKey
	params.SecretKey = testConfig.SecretKey
	params.ProxyURL = testConfig.ProxyURL
	params.Testnet = testConfig.Testnet
	return NewHbdmSwap(params)
}

func TestHbdmSwap_GetRecords(t *testing.T) {
	ex := testExchange()
	symbol := "BTC-USD"
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := ex.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
