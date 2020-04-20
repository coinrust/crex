package binancefutures

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func testExchange() Exchange {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	proxyURL := viper.GetString("proxy_url")
	log.Printf("accessKey: %v", accessKey)
	log.Printf("secretKey: %v", secretKey)
	params := &Parameters{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	ex := NewBinanceFutures(params)
	if proxyURL != "" {
		ex.SetProxy(proxyURL)
	}
	return ex
}

func TestBinanceFutures_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("USDT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestBinanceFutures_GetOrderBook(t *testing.T) {
	ex := testExchange()
	ob, err := ex.GetOrderBook("BTCUSDT", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestBinanceFutures_GetRecords(t *testing.T) {
	ex := testExchange()
	now := time.Now()
	start := now.Add(-300 * time.Minute)
	end := now
	records, err := ex.GetRecords("BTCUSDT",
		PERIOD_1MIN, start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("Timestamp: %v %#v", v.Timestamp, v)
	}
}

func TestBinanceFutures_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	orders, err := ex.GetOpenOrders("BTCUSDT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orders)
}
