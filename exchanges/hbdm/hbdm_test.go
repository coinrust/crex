package hbdm

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func testExchange(websocket bool) *Hbdm {
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
	params.WebSocket = websocket
	return NewHbdm(params)
}

func TestHbdm_GetBalance(t *testing.T) {
	ex := testExchange(false)
	balance, err := ex.GetBalance("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", balance)
}

func TestHbdm_GetContractInfo(t *testing.T) {
	ex := testExchange(false)
	symbol, contractType, err := ex.GetContractInfo("BTC200424")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("symbol: %v", symbol)
	t.Logf("contractType: %v", contractType)
}

func TestHbdm_GetOrderBook(t *testing.T) {
	ex := testExchange(false)
	ex.SetContractType("BTC", "W1")
	ob, err := ex.GetOrderBook("BTC200327", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestHbdm_GetRecords(t *testing.T) {
	ex := testExchange(false)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := ex.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}

func TestHbdm_GetContractID(t *testing.T) {
	ex := testExchange(false)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
}

func TestHbdm_GetOpenOrders(t *testing.T) {
	ex := testExchange(false)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("symbol: %v", symbol)

	orders, err := ex.GetOpenOrders(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestHbdm_GetOrder(t *testing.T) {
	ex := testExchange(false)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := ex.GetOrder(symbol, "694901372910391296")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestHbdm_PlaceOrder(t *testing.T) {
	ex := testExchange(false)
	ex.SetLeverRate(10)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := ex.PlaceOrder(symbol,
		Buy,
		OrderTypeLimit,
		3000,
		0,
		1,
		false,
		false,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestHbdm_PlaceOrder2(t *testing.T) {
	ex := testExchange(false)
	ex.SetLeverRate(10)
	ex.SetContractType("BTC", ContractTypeW1)
	symbol, err := ex.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := ex.PlaceOrder(symbol,
		Sell,
		OrderTypeMarket,
		3000,
		0,
		1,
		false,
		true,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestHbdm_WebSocket(t *testing.T) {
	ex := testExchange(true)
	err := ex.SubscribeLevel2Snapshots(Market{
		Symbol: "BTC200424",
	}, func(ob *OrderBook) {
		t.Logf("%#v", ob)
	})
	if err != nil {
		t.Error(err)
	}

	select {}
}
