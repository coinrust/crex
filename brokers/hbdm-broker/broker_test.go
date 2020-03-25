package hbdm_broker

import (
	. "github.com/coinrust/crex"
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

func TestHBDMBroker_GetAccountSummary(t *testing.T) {
	b := newTestBroker()
	accountSummary, err := b.GetAccountSummary("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
}

func TestHBDMBroker_GetOrderBook(t *testing.T) {
	b := newTestBroker()
	b.SetContractType("BTC", "W1")
	ob, err := b.GetOrderBook("BTC200327", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestHBDMBroker_GetContractID(t *testing.T) {
	b := newTestBroker()
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
}
