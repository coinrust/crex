package huobi_broker

import (
	. "github.com/coinrust/gotrader"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newTestBroker() Broker {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	baseURL := "https://api.btcgateway.pro"
	return NewBroker(baseURL, accessKey, secretKey)
}

func TestHuobiBroker_GetAccountSummary(t *testing.T) {
	b := newTestBroker()
	accountSummary, err := b.GetAccountSummary("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
}
