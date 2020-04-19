package hbdmswap

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func newForTest() Broker {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	params := &Parameters{}
	params.AccessKey = viper.GetString("access_key")
	params.SecretKey = viper.GetString("secret_key")
	params.ProxyURL = viper.GetString("proxy_url")
	params.Testnet = true
	return New(params)
}

func TestHBDMSwap_GetRecords(t *testing.T) {
	b := newForTest()
	symbol := "BTC-USD"
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := b.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
